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

	if function == "newPICinfo" {
		return newPICinfo(stub, args)
	}
	return shim.Error("penalinterestcc: " + "no function named " + function + " found in Interest Refund")
}

func newPICinfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	fmt.Println("******************** Start Penal Interst Collection *************************")

	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("penalCharges.cc: " + "Invalid number of arguments in newPICinfo(Interest Refund) (required:10) given:" + xLenStr)
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
	 *LoanID  string    //args[3] loan
	 *InsID   string    //args[4]
	 *Amt     int64     //args[5]
	 *FromID  string    //args[6]  Bank
	 *ToID    string    //args[7]  Seller
	 *By      string    //args[8]
	 *PprID   string    //args[9]
	 */

	amt, _ := strconv.ParseInt(args[5], 10, 64)
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// The transaction object has been created and written into the ledger
	// The JSON file is 'transaction'function
	// Now to create a TXN_Bal_Update obj for 4 times
	// Calling TXN_Balance CC based on TXN_Type
	/*
			    a. Debiting (Decreasing) Business Wallet
		        b. Crediting (Incresing) Bank Wallet
		        c. Debiting (Decreasing) Business Charges O/s Wallet
		        d. Debiting (Decreasing) Loan Charges Wallet
	*/

	//Validations

	// Must be Existing Loan with Status as Overdue
	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalChargescc: can't get loanStatus" + response.Message)
	}
	status := string(response.Payload)
	if status != "overdue" {
		return shim.Error("penalCharges.cc: " + "loan status for loanID " + args[3] + " is not overdue")
	}

	//TXN Amt must be > Zero
	if amt <= 0 {
		return shim.Error("penalCharges.cc: " + "Transaction Amount in Penal Interest Collectionis less than or equal to zero")
	}

	//####################################################################################################################

	//#####################################################################################################################
	//Calling for updating Business Main_Wallet
	//####################################################################################################################

	cAmtString := "0"
	dAmtString := args[5]

	walletID, err := getWalletID(stub, "businesscc", args[7], "main")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Business Main WalletID " + err.Error())
	}

	openBalance, err := getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Business Main WalletValue " + err.Error())
	}
	openBalString := strconv.FormatInt(openBalance, 10)
	bal := openBalance - amt

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	txnBalString := strconv.FormatInt(bal, 10)
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	u1 := uuid.New()

	fmt.Printf("generated Version 4 UUID %v", u1)
	StringUUID1 := u1.String();
	fmt.Print("StringUUID1 ",StringUUID1);
	i, err := strconv.ParseFloat(txnBalString, 64)
	if err!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i) {
		return shim.Error("Invalid transaction Amount for Business Main Wallet")
	}
	argsList := []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	txnResponse := putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Bank Main_Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, err = getWalletID(stub, "bankcc", args[6], "main")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Bank Main WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Bank Main WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance + amt
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	u2 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u2)
	StringUUID2 := u2.String();
	fmt.Print("StringUUID2 ",StringUUID2);
	i1, err1 := strconv.ParseFloat(txnBalString, 64)
	if err1!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i1) {
		return shim.Error("Invalid transaction Amount for Bank Main Wallet")
	}
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Business Charges O/s Wallet
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, err = getWalletID(stub, "businesscc", args[7], "chargesOut")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Business Charges O/s WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Business Charges O/s WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance - amt
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	u3 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u3)
	StringUUID3 := u3.String();
	fmt.Print("StringUUID3 ",StringUUID3);
	i2, err2 := strconv.ParseFloat(txnBalString, 64)
	if err2!= nil {
		return shim.Error("Error while converting String to Int");
	}
	if math.Signbit(i2) {
		return shim.Error("Invalid transaction Amount for Business Charges O/s Wallet")
	}
	argsList = []string{StringUUID3, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Loan Charges Wallet
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, err = getWalletID(stub, "loancc", args[3], "charges")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Loan Charges Wallet WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Loan Charges WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance - amt
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
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
		return shim.Error("Invalid transaction Amount for Loan Charges Wallet")
	}
	argsList = []string{StringUUID4, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################
	fmt.Println("******************** End Penal Interst collection *************************")
	return shim.Success(nil)
}

func putInTxnBal(stub shim.ChaincodeStubInterface, argsListStr string) pb.Response {

	chaincodeArgs := toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the txnbalcc chaincode from Interest Refund")
	response := stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	fmt.Println(string(response.Payload))
	return shim.Success(nil)
}

func getWalletID(stub shim.ChaincodeStubInterface, ccName string, id string, walletType string) (string, error) {

	// STEP-1
	// using FromID, get a walletID from bank structure

	chaincodeArgs := toChaincodeArgs("getWalletID", id, walletType)
	response := stub.InvokeChaincode(ccName, chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "0", errors.New(response.Message)
	}
	walletID := string(response.GetPayload())
	return walletID, nil

}

func getWalletValue(stub shim.ChaincodeStubInterface, walletID string) (int64, error) {

	walletArgs := toChaincodeArgs("getWallet", walletID)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return 0, errors.New(walletResponse.Message)
	}
	balString := string(walletResponse.Payload)
	balance, _ := strconv.ParseInt(balString, 10, 64)
	return balance, nil
}

func walletUpdation(stub shim.ChaincodeStubInterface, walletID string, amt int64) pb.Response {

	txnBalString := strconv.FormatInt(amt, 10)
	walletArgs := toChaincodeArgs("updateWallet", walletID, txnBalString)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return shim.Error(walletResponse.Message)
	}
	return shim.Success(nil)

}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start Penal Interest Collectionchaincode:", err)
	}
}
