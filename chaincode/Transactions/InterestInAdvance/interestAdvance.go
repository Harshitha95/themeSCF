package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	if function == "newInterestInAdvInfo" {
		//Creates new newInterestInAdvInfo
		return newInterestInAdvInfo(stub, args)
	}
	return shim.Error("interestAdv.cc: " + "no function named " + function + " found in InterestAdv")
}

func newInterestInAdvInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start newInterestInAdvInfo *************************")
	
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("interestAdv.cc: " + "Invalid number of arguments in newInterestAdvInfo(interestAdv) (required:10) given:" + xLenStr)
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
	 *FromID  string    //args[6]
	 *ToID    string    //args[7]
	 *By      string    //args[8]
	 *PprID   string    //args[9]
	 */

	//Validations
	//Getting the sanction amount and the status
	//Validations
	//Getting the sanction amount and the status
	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("interestAdvcc: can't get loanStatus" + response.Message)
	}
	status := string(response.Payload)
	if status != "part disbursed" && status != "disbursed" {
		return shim.Error("interestAdv.cc: " + "loan status for loanID " + args[3] + " is not Sanctioned / part disbursed / disbursed")
	}
	txnAmt, _ := strconv.ParseInt(args[5], 10, 64)
	if txnAmt <= 0 {
		return shim.Error("interestAdv.cc: txnAmt is zero or less")
	}
	
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// Now to create a TXN_Bal_Update obj for 3 times
	// Calling TXN_Balance CC based on TXN_Type
	/*
	   a. Crediting (Increasing) Business Loan Wallet
	   b. Crediting (Increasing) Loan Disbursed Wallet
	   c. Debiting (Decreasing) Business Charges O/s Wallet
	   d. Debiting (Decreasing) Loan Charges O/s Wallet
	   e. Crediting (Increasing) Bank Asset Wallet
	*/

	//####################################################################################################################
	//Calling for updating Business Loan Wallet
	//####################################################################################################################

	cAmtString := args[5]
	dAmtString := "0"

	walletID, openBalString, txnBalString, err := getWalletInfo(stub, args[7], "loan", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("interestAdv.cc: " + "Business Loan Wallet(interestAdv):" + err.Error())
	}
	//ADDED SHA256

	u2 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u2)
	StringUUID := u2.String();
	
	fmt.Print("StringUUID ",StringUUID);

	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList := []string{StringUUID, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the other chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("interestAdv.cc: " + response.Message)
	}
	fmt.Println("interestAdvcc: " + string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Loan Disbursed Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[3], "disbursed", "loancc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("insterestAdv.cc: " + "Loan Disbursed Wallet(interestAdv):" + err.Error())
	}
	u3 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u3)
	StringUUID1 := u3.String();
	fmt.Print("StringUUID1 "+StringUUID1 );
	//TXiDHASH1 := hex.EncodeToString(hash1)
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the other chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("interestAdv.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))
	//####################################################################################################################
	//Calling for updating Business Charges O/s Wallet
	//####################################################################################################################
	cAmtString = "0"
	dAmtString = args[5]

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[7], "chargesOut", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("interestAdv.cc: " + "Business charges Wallet(interestAdv):" + err.Error())
	}
	u4 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u4)
	StringUUID2 := u4.String();
	fmt.Print("StringUUID2 "+StringUUID2 );

	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the other chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("interestAdv.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Loan Charges O/s Wallet
	//####################################################################################################################
	cAmtString = "0"
	dAmtString = args[5]
	fmt.Println("Loan charges Wallet args[3] ",args[3])
	fmt.Println("cAmtString ",cAmtString)
	fmt.Println("dAmtString ",dAmtString)

	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[3], "charges", "loancc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("interestAdv.cc: " + "Loan charges Wallet(interestAdv):" + err.Error())
	}
	
	u5 := uuid.New()
	fmt.Println("generated Version 4 UUID %v", u5)
	StringUUID3 := u5.String();
		fmt.Println("StringUUID3 "+StringUUID3 );
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID3, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the other chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("interestAdv.cc: " + response.Message)
	}
	fmt.Println(string(response.GetPayload()))

	//####################################################################################################################
	//Calling for updating Bank Asset Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"
	fmt.Println("Bank Asset Wallet args[3] ",args[3])
	fmt.Println("cAmtString ",cAmtString)
	fmt.Println("dAmtString ",dAmtString)
	walletID, openBalString, txnBalString, err = getWalletInfo(stub, args[8], "asset", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("interestAdv.cc: " + "Bank Asset Wallet(interestAdv):" + err.Error())
	}

	u6 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u6)
	StringUUID4 := u6.String();
		
	fmt.Print("StringUUID4 "+StringUUID4 )
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{StringUUID4, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the other chaincode")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("interestAdv.cc: " + response.Message)
	}
	// Call updateInterestInAdvanceCollection()
	chaincodeArgs1 := toChaincodeArgs("updateInterestInAdvanceCollection", args[3], "true")
		response1 := stub.InvokeChaincode("loancc", chaincodeArgs1, "myc")
		if response.Status != shim.OK {
			return shim.Error("interestinadvance: " + response1.Message)
		}
	fmt.Println(" updateInterestInAdvanceCollection Response payload " ,string(response1.GetPayload()))

	fmt.Println("******************** End newInterestInAdvInfo *************************")
	return shim.Success(nil)
}

func getWalletInfo(stub shim.ChaincodeStubInterface, participantID string, walletType string, ccName string, cAmtStr string, dAmtStr string) (string, string, string, error) {

	// STEP-1
	// using FromID, get a walletID from participant / loan
	// bankID = bankID

	fmt.Println("participantID ",participantID);
	fmt.Println("walletType ",walletType)
	chaincodeArgs := toChaincodeArgs("getWalletID", participantID, walletType)
	response := stub.InvokeChaincode(ccName, chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "", "", "", errors.New(response.Message)
	}
	walletID := string(response.GetPayload())
	fmt.Println(" Wallet Id from getWalletInfo ",walletID)
	// STEP-2
	// getting Balance from walletID
	// walletFcn := "getWallet"
	walletArgs := toChaincodeArgs("getWallet", walletID)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return "", "", "", errors.New(walletResponse.Message)
	}
	openBalString := string(walletResponse.Payload)
	fmt.Println(" Balance in Wallet ",openBalString)

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
		fmt.Println("interestAdv.cc: " + "Unable to start the chaincode")
	}
}
