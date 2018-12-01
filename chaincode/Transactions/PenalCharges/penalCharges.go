package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"math"
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

	if function == "newPenalChargesInfo" {
		//Creates new penalCharges info
		return newPenalChargesInfo(stub, args)
	}
	return shim.Error("penalCharges.cc: " + "no function named " + function + " found in penalCharges")
}

func newPenalChargesInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start newPenalChargesInfo *************************")

	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("penalCharges.cc: " + "Invalid number of arguments in newPenalChargesInfo(penalCharges) (required:10) given:" + xLenStr)
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
	 *By      string    //args[8]
	 *PprID   string    //args[9]
	 */

	//Validations
	//Getting the sanction amount and the status
	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalChargescc: can't get loanStatus" + response.Message)
	}

	status := string(response.Payload)
	if status != "overdue" {
		return shim.Error("penalCharges.cc: " + "loan status for loanID " + args[3] + " is not Sanctioned / part disbursed / disbursed")
	}
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// Now to create a TXN_Bal_Update obj for 3 times
	// Calling TXN_Balance CC based on TXN_Type
	/*
	   a. Crediting (Increasing) Bank Liability  Wallet
	   b. Crediting (Increasing) Business Charges O/s Wallet
	   c. Crediting (Increasing) Loan Charges Wallet
	*/

	//####################################################################################################################
	//Calling for updating Bank Revenue Wallet
	//####################################################################################################################

	cAmtString := args[5]
	dAmtString := "0"

	walletID, openBalString, txnBalString, err := getWalletInfo(stub, args[6], "liability", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Bank Liability Wallet(penalCharges):" + err.Error())
	}
	u1 := uuid.New()

	fmt.Printf("Generated Version 4 UUID Bank Liability Wallet %v", u1)
	StringUUID1 := u1.String();
	fmt.Print("StringUUID1 ",StringUUID1);
	i, err := strconv.ParseFloat(txnBalString, 64)
	if err!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i) {
		return shim.Error("Invalid transaction Amount for Bank liability Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList := []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the bankcc chaincode Bank Liability Wallet")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Business Charges O/s Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[7], "chargesOut", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Business Charges O/s Wallet(penalCharges):" + err.Error())
	}
	u2 := uuid.New()
	fmt.Printf("generated Version 4 UUID Business Charges O/s Wallet %v", u2)
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
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the  businesscc chaincode Business Charges O/s Wallet")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Loan Charges Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[3], "charges", "loancc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Loan charges Wallet(penalCharges):" + err.Error())
	}
	u3 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u3)
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
	argsList = []string{StringUUID3, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the loancc  chaincode Loan Charges Wallet")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))
	fmt.Println("******************** End Penal Charges *************************")
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
		fmt.Println("penalCharges.cc: " + "Unable to start the chaincode")
	}
}
