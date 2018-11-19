package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	/* "github.com/gofrs/uuid" */
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/satori/go.uuid"
	"github.com/google/uuid"
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

	if function == "newChargesInfo" {
		//Creates new charges info
		return newChargesInfo(stub, args)
	}
	return shim.Error("chargescc: " + "no function named " + function + " found in charges")
}

func newChargesInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start Charges *************************")
/* 	fmt.Println("args []",args[0], " ",args[1], " ",args[2], " ",args[3], " ",args[4], " ",args[5])
	fmt.Print("args[6] &[7]  &[8]", args[6] ," -----> ",args[7]," -----> ",args[8]) */
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("chargescc: " + "Invalid number of arguments in newChargesInfo(charges) (required:10) given:" + xLenStr)
	}
	/*
	 *TxnType string    //args[1]
	 *TxnDate time.Time //args[2]
	 *LoanID  string    //args[3]
	 *InsID   string    //args[4]
	 *Amt     int64     //args[5]
	 *BankID  string    //args[6]
	 *SellID  string    //args[7]
	 *By      string    //args[8]
	 *PprID   string    //args[9]
	 */

	//Validations
	//Getting the sanction amount and the status
	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("chargescc: can't get loanStatus" + response.Message)
	}
	status := string(response.Payload)
	if status != "sanctioned" && status != "part disbursed" && status != "disbursed" {
		return shim.Error("chargescc: " + "loan status for loanID " + args[3] + " is not Sanctioned / part disbursed / disbursed")
	}

	//sancAmt, _ := strconv.ParseInt(statusNamt[0], 10, 64)
	//txnAmt > 0
	txnAmt, _ := strconv.ParseInt(args[5], 10, 64)
	if txnAmt <= 0 {
		return shim.Error("chargescc: txnAmt is zero or less")
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// Now to create a TXN_Bal_Update obj for 3 times
	// Calling TXN_Balance CC based on TXN_Type
	/*
	   a. Crediting (Increasing) Bank Revenue Wallet
	   b. Crediting (Increasing) Business Charges O/s Wallet
	   c. Crediting (Increasing) Loan Charges Wallet
	*/

	//####################################################################################################################
	//Calling for updating Bank Revenue Wallet
	//####################################################################################################################

	cAmtString := args[5]
	dAmtString := "0"

	walletID, openBalString, txnBalString, err := getWalletInfo(stub, args[6], "charges", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("chargescc: " + "Bank Revenue Wallet(charges):" + err.Error())
	}
	u2 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u2)
	StringUUID := u2.String();
	
	fmt.Print("StringUUID ",StringUUID);
	
	fmt.Println("u2 ",u2)
	fmt.Println("StringUUID ",StringUUID)
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList := []string{StringUUID, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	fmt.Println("argsListStr 1 ",argsListStr)
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the bank asset chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("chargescc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Business Charges O/s Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[7], "chargesOut", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("chargescc: " + "Business Charges O/s Wallet(charges):" + err.Error())
	}
	u3 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u3)
	StringUUID1 := u3.String();
	
	fmt.Print("StringUUID1 ",StringUUID1);
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger

	argsList = []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the Business COS  chaincode")
	fmt.Println("argsListStr 2",argsListStr)
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("chargescc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Loan Charges Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[3], "charges", "loancc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("chargescc: " + "Loan charges Wallet(charges):" + err.Error())
	}
	u4 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u2)
	StringUUID2 := u4.String();
	
	fmt.Print("StringUUID2 ",StringUUID2);
	
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	fmt.Println("argsListStr 3",argsListStr)
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the Loan Charges Wallet chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("chargescc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))
	fmt.Println("******************** End Charges *************************")
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
/* func genUUIDv4() uuid.UUID{
   /*  id := uuid.NewV4()
	fmt.Printf("github.com/satori/go.uuid:  %s\n", id)
	return id;
	id, err:= uuid.NewV4()
	if err != nil {
		fmt.Errorf("failed to generate UUID: %v", err)
	}
	return id
} */
func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
}
