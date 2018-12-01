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

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "newDisbInfo" {
		//Creates new disbursement info
		return newDisbInfo(stub, args)
	}
	return shim.Error("disbursementcc: " + "no function named " + function + " found in Disbursement")
}

func newDisbInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start Disbursement *************************")
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("disbursementcc: " + "Invalid number of arguments in newDisbInfo(disbursement) (required:10) given:" + xLenStr)
	}
	fmt.Println("args []",args [0], " ",args [1], " ",args [2], " ",args [3], " ",args [4], " ",args [5])
	fmt.Print("args[6] &[7]  &[8] &[9]", args[6] ," -----> ",args[7]," -----> ",args[8]," -----> ",args[9])
	/*
	 *TxnType string    //args[1]
	 *TxnDate time.Time //args[2]
	 *LoanID  string    //args[3]
	 *InsID   string    //args[4]
	 *Amt     int64     //args[5]
	 *FromID  string    //args[6]	bank
	 *ToID    string    //args[7]	seller
	 *By      string    //args[8]
	 *PprID   string    //args[9]
	 */
	//Validations
	//Getting the sanction amount and the status
	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	status := string(response.Payload)
	fmt.Println(status)
	if (status != "sanctioned") && (status != "part disbursed") {
		return shim.Error("disbursementcc: " + "loan status for loanID " + args[3] + " is not sanctioned / part disbursed")
	}

	//sancAmt, _ := strconv.ParseInt(statusNamt[1], 10, 64)
	chaincodeArgs = toChaincodeArgs("getLoanSancAmt", args[3])
	response = stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	sancAmt, err := strconv.ParseInt(string(response.Payload), 10, 64)
	if err != nil {
		return shim.Error("loancc: error parsing sancAmt")
	}
	fmt.Println(sancAmt)

	//Getting the disbursed wallet
	chaincodeArgs = toChaincodeArgs("getWalletID", args[3], "disbursed")
	response = stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	walletid := string(response.Payload)
	disbAmt, err := getWalletValues(stub, walletid)
	if err != nil {
		return shim.Error("disbursementcc: " + err.Error())
	}
	//converting dibAmt to integer for testing
	amtToBeDisburesed := sancAmt - disbAmt

	amt, _ := strconv.ParseInt(args[5], 10, 64)
	fmt.Println("sancAmt %s",sancAmt);
	fmt.Println("amt ,disbAmt %s",amt," ",disbAmt);
	fmt.Println("amtToBeDisburesed %s",amtToBeDisburesed);
	
	if amt > amtToBeDisburesed {
		return shim.Error("disbursementcc: " + "Amount is greater than Amount to be disbursed")
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// The transaction object has been created and written into the ledger
	// The JSON file is 'transaction'function
	// Now to create a TXN_Bal_Update obj for 6 times
	// Calling TXN_Balance CC based on TXN_Type
	/*
	   a. Debiting (Reducing) Bank Wallet
	   b. Crediting (Increasing) Business Wallet
	   c. Crediting (Increasing) Bank Asset Wallet
	   d. Crediting (Increasing) Business Loan Wallet
	   e. Crediting (Increasing) Business Principal O/s Wallet
	   f. Crediting (Increasing) Loan Disbursed Wallet
	*/

	//####################################################################################################################
	//Calling for updating Bank Main_Wallet
	//####################################################################################################################

	cAmtString := "0"
	dAmtString := args[5]

	walletID, openBalString, txnBalString, err := getWalletInfo(stub, args[6], "main", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("disbursementcc: " + "Bank Main Wallet(Disbursement):" + err.Error())
	}
	u1 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u1)
	StringUUID1 := u1.String();
	
	fmt.Print("StringUUID1 ",StringUUID1);
	i, err := strconv.ParseFloat(txnBalString, 64) //String to float 
	if err!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i) {
		return shim.Error("Invalid tranaction Amount for Bank Main Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList := []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	
	fmt.Println("calling the putTxnBalInfo chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//#####################################################################################################################
	//Calling for updating Business Main Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[7], "main", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("disbursementcc: " + " Business Main  Wallet(Disbursement):" + err.Error())
	}
	u2 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u2)
	StringUUID2 := u2.String();
	
	fmt.Print("StringUUID2 ",StringUUID2);
	
//	i := float64(bal) // int to float
	i1, err1 := strconv.ParseFloat(txnBalString, 64) //String to float 
	if err1!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i1) {
		return shim.Error("Invalid tranaction Amount for Business Main Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the putTxnBalInfo chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Business Loan_Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[7], "loan", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("disbursementcc: " + "Business Loan Wallet(Disbursement)" + err.Error())
	}
	u3 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u3)
	StringUUID3 := u3.String();
	
	fmt.Print("StringUUID3 ",StringUUID3);
	i2, err2 := strconv.ParseFloat(txnBalString, 64) //String to float 
	if err2!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i2) {
		return shim.Error("Invalid tranaction Amount for Business Loan Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID3, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the putTxnBalInfo chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Bank Asset_Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[6], "asset", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("disbursementcc: " + "Bank Asset Wallet(Disbursement)" + err.Error())
	}

	u4 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u4)
	StringUUID4 := u4.String();
	
	fmt.Print("StringUUID4 ",StringUUID4);
	i3, err3 := strconv.ParseFloat(txnBalString, 64) //String to float 
	if err3!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i3) {
		return shim.Error("Invalid tranaction Amount for Bank Asset Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID4, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the putTxnBalInfo chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Business principal O/S Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[7], "principalOut", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("disbursementcc: " + "Business principal O/S Wallet(Disbursement)" + err.Error())
	}
	u5 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u5)
	StringUUID5 := u5.String();
	
	fmt.Print("StringUUID5 ",StringUUID5);
	i4, err4 := strconv.ParseFloat(txnBalString, 64) //String to float 
	if err4!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i4) {
		return shim.Error("Invalid tranaction Amount for Business principal O/S Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID5, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the putTxnBalInfo chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Loan Disbursed Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[3], "disbursed", "loancc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("disbursementcc: " + "Loan Disbursed Wallet(Disbursement)" + err.Error())
	}
	/* hash5 := sha256.New()
	txnId5 := rand.Intn(10000);
	hash5.Write([]byte(strconv.Itoa(txnId5)))
	md5 := hash5.Sum(nil)
	txnIDsha5 := hex.EncodeToString(md5)
	fmt.Println("Loan Disbursed Wallet txnIDsha5",txnIDsha5) */
	u6 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u6)
	StringUUID6 := u6.String();
	
	fmt.Print("StringUUID6 ",StringUUID6);
	i5, err5 := strconv.ParseFloat(txnBalString, 64) //String to float 
	if err5!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i5) {
		return shim.Error("Invalid tranaction Amount for Loan Disbursed Wallet")
	}
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID6, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the putTxnBalInfo chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################

	//####################################################################################################################
	//Calling Loan to change the status
	//####################################################################################################################

	if amtToBeDisburesed == 0 {
		status = "disbursed"
	} else if amtToBeDisburesed > 0 {
		status = "part disbursed"
	}

	//calling to change loan status
	chaincodeArgs = toChaincodeArgs("updateLoanInfo", args[3], status, "disbursement")
	response = stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("disbursementcc: " + response.Message)
	}
	fmt.Println("******************** End Disbursement *************************")
	return shim.Success(nil)
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func getWalletInfo(stub shim.ChaincodeStubInterface, participantID string, walletType string, ccName string, cAmtStr string, dAmtStr string) (string, string, string, error) {

	// STEP-1
	// using FromID, get a walletID from bank structure
	// bankID = bankID
	fmt.Println(" ******************** Start getWalletInfo with value ",participantID ," ",walletType," ",ccName," ",cAmtStr," ",dAmtStr)
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
	fmt.Println(" ******************** End getWalletInfo ************************* ")
	return walletID, openBalString, txnBalString, nil
}

func getWalletValues(stub shim.ChaincodeStubInterface, walletID string) (int64, error) {
	fmt.Println(" ******************** Start getWalletValues *********** ")
	fmt.Println("args " ,walletID )
	walletArgs := toChaincodeArgs("getWallet", walletID)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return 0, errors.New(walletResponse.Message)
	}
	openBalString := string(walletResponse.Payload)
	openBal, err := strconv.ParseInt(openBalString, 10, 64)
	if err != nil {
		return 0, errors.New("Error in converting the openBalance in getWalletValues(disbursement)")
	}
	fmt.Println(" ******************** End getWalletValues *********** ")
	return openBal, nil
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("disbursementcc: " + "Unable to start the chaincode")
	}
}
