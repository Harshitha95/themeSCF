package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"math"
/* 	"crypto/sha256"
	"encoding/hex"
	"math/rand" */
	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "newTDSInfo" {
		//Creates new tds info
		return newTDSInfo(stub, args)
	}
	return shim.Error("tds.cc: " + "no function named " + function + " found in tds")
}

func newTDSInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start TDS *************************")
	
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("tds.cc: " + "Invalid number of arguments in newTDSInfo(tds) (required:10) given:" + xLenStr)
	}
	fmt.Println("args[0]", args[0] )
	fmt.Println("args[1]", args[1] )
	fmt.Println("args[2]", args[2] )
	fmt.Println("args[3]", args[3] )
	fmt.Println("args[4]", args[4] )
	fmt.Println("args[5]", args[5] )
	fmt.Println("args[6]", args[6] )
	fmt.Println("args[7]", args[7] )
	fmt.Println("args[8]", args[8] ) 
	/*
	 *TxnType string    //args[1]
	 *TxnDate time.Time //args[2]
	 *LoanID  string    //args[3]
	 *InsID   string    //args[4]
	 *Amt     int64     //args[5]
	 *bank    string    //args[6]
	 *seller  string    //args[7]
	 *buyer	  string	//args[8]
	 *By      string    //args[9]
	 */

	//Validations
	//Getting the sanction amount and the status
	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("interestrefundcc: can't get loanStatus" + response.Message)
	}
	status := string(response.Payload)
	if status != "overdue" {
		return shim.Error("tds.cc: " + "loan status for loanID " + args[3] + " is not Overdue")
	}

	//sancAmt, _ := strconv.ParseInt(statusNamt[0], 10, 64)
	//txnAmt > 0
	txnAmt, _ := strconv.ParseInt(args[5], 10, 64)
	if txnAmt <= 0 {
		return shim.Error("tds.cc: txnAmt is zero or less")
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// Now to create a TXN_Bal_Update obj for 3 times
	// Calling TXN_Balance CC based on TXN_Type
	/*
	   a. Debiting (Decreasing) Business Loan Wallet (Seller)
	   b. Debiting (Decreasing) Business Liability Wallet (buyer)
	   c. Debiting (Decreasing) Bank Asset Wallet
	   d. Crediting (Increasing) TDS Receivable Wallet
	*/

	//####################################################################################################################
	//Calling for updating Business Loan Wallet (Seller)
	//####################################################################################################################

	cAmtString := "0"
	dAmtString := args[5]

	walletID, openBalString, txnBalString, err := getWalletInfo(stub, args[7], "loan", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("tds.cc: " + "Business Loan Wallet (Seller)(tds):" + err.Error())
	}
	u1 := uuid.New()

	fmt.Printf("generated Version 4 UUID Business Loan Wallet %v", u1)
	StringUUID1 := u1.String();
	fmt.Print("StringUUID1 ",StringUUID1);
	i, err := strconv.ParseFloat(txnBalString, 64)
	if err!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i) {
		return shim.Error("Invalid transaction Amount for Business Loan Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList := []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[9]}
	argsListStr := strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the businesscc chaincode to update  Business Loan Wallet")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("tds.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Debiting (Decreasing) Business Liability Wallet (buyer)
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[8], "liability", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("tds.cc: " + "Business Charges O/s Wallet(tds):" + err.Error())
	}
	u2 := uuid.New()
	fmt.Printf("generated Version 4 UUID  Business Liability Wallet  %v", u2)
	StringUUID2 := u2.String();
	fmt.Print("StringUUID2 ",StringUUID2);
	i1, err1 := strconv.ParseFloat(txnBalString, 64)
	if err1!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i1) {
		return shim.Error("Invalid transaction Amount for Business Charges Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[9]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the businesscc  chaincode  Business Liability Wallet ")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("tds.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Debiting (Decreasing) Bank Asset Wallet
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[6], "asset", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("tds.cc: " + "Loan charges Wallet(tds):" + err.Error())
	}
	u3 := uuid.New()
	fmt.Printf("generated Version 4 UUID Bank Asset Wallet %v", u3)
	StringUUID3 := u3.String();
	fmt.Print("StringUUID3 ",StringUUID3);
	i2, err2 := strconv.ParseFloat(txnBalString, 64)
	if err2!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i2) {
		return shim.Error("Invalid transaction Amount for Loan Charges Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID3, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[9]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the Bank  chaincode Bank Asset Wallet ")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("tds.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Crediting (Increasing) TDS Receivable Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[6], "tds", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("tds.cc: " + "Business Charges O/s Wallet(tds):" + err.Error())
	}
	u4 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u4)
	StringUUID4 := u4.String();
	fmt.Print("StringUUID4 ",StringUUID4);
	i3, err3 := strconv.ParseFloat(txnBalString, 64)
	if err3!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i3) {
		return shim.Error("Invalid transaction Amount for Business Charges Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID4, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[9]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the bankcc  chaincode Business Charges O/s Wallet(tds) ")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("tds.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))
	fmt.Println("******************** End TDS *************************")
	return shim.Success(nil)
}

func getWalletInfo(stub shim.ChaincodeStubInterface, participantID string, walletType string, ccName string, cAmtStr string, dAmtStr string) (string, string, string, error) {

	// STEP-1
	// using FromID, get a walletID from participant / loan
	// bankID = bankID

	chaincodeArgs := toChaincodeArgs("getWalletID", participantID, walletType)
	response := stub.InvokeChaincode(ccName, chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "", "", "", errors.New(response.Message)
	}
	walletID := string(response.GetPayload())

	// STEP-2
	// getting Balance from walletID
	// walletFcn := "getWallet"
	walletArgs := toChaincodeArgs("getWallet", walletID)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return "", "", "", errors.New(walletResponse.Message)
	}
	openBalString := string(walletResponse.Payload)

	openBal, err := strconv.ParseInt(openBalString, 10, 64)
	if err != nil {
		return "", "", "", errors.New("Error in converting the openBalance")
	}
	cAmt, err := strconv.ParseInt(cAmtStr, 10, 64)
	if err != nil {
		return "", "", "", errors.New("Error in converting the cAmt")
	}
	dAmt, err := strconv.ParseInt(dAmtStr, 10, 64)
	if err != nil {
		return "", "", "", errors.New("Error in converting the dAmt")
	}

	txnBal := openBal - dAmt + cAmt
	txnBalString := strconv.FormatInt(txnBal, 10)

	// STEP-3
	// update wallet of ID walletID here, and write it to the wallet_ledger
	// walletFcn := "updateWallet"

	walletArgs = toChaincodeArgs("updateWallet", walletID, txnBalString)
	walletResponse = stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return "", "", "", errors.New(walletResponse.Message)
	}

	return walletID, openBalString, txnBalString, nil
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("tdscc: " + "Unable to start the chaincode")
	}
}
