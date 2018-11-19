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

	if function == "newMarginInfo" {
		return newMarginInfo(stub, args)
	}
	return shim.Error("marginrefundcc: " + "no function named " + function + " found in Margin Refund")
}

func newMarginInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	fmt.Println("******************** Start newMarginInfo *************************")
	
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("marginrefundcc: " + "Invalid number of arguments in newMarginInfo(Margin Refund) (required:10) given:" + xLenStr)
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
	 *FromID  string    //args[6]  Bank
	 *ToID    string    //args[7]  Business
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
	   a. Crediting (Increasing) Business Wallet
	   b. Debiting (Decreasing) Bank Wallet
	   c. Debiting (Decreasing) Bank Refund Wallet
	*/

	//Validations

	// Must be Existing Loan with Status as Collected
 	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("marginrefundcc: can't get loanStatus" + response.Message)
	}
	status := string(response.Payload)
	if status != "collected" {
		return shim.Error("marginrefundcc: " + "loan status for loanID " + args[3] + " is not collected")
	} 

	//TXN Amt must be > Zero
	if amt <= 0 {
		return shim.Error("marginrefundcc: " + "Transaction Amount in margin refund is less than or equal to zero")
	}

	//Loan disbursed Wallet balance must be Zero
	loanDisbursedWalletID, err := getWalletID(stub, "loancc", args[3], "disbursed")
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund loanDisbursedWalletID " + err.Error())
	}
	loanDisbursedWalletValue, err := getWalletValue(stub, loanDisbursedWalletID)
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund loanDisbursedWalletValue " + err.Error())
	}

	//Loan Charges Wallet balance must be Zero
	loanChargesWalletID, err := getWalletID(stub, "loancc", args[3], "charges")
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund loanChargesWalletID " + err.Error())
	}
	loanChargesWalletValue, err := getWalletValue(stub, loanChargesWalletID)
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund loanChargesWalletValue " + err.Error())
	}

	// Loan Accrued Wallet balance must be Zero
	loanAccruedWalletID, err := getWalletID(stub, "loancc", args[3], "accrued")
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund loanAccruedWalletID " + err.Error())
	}
	loanAccruedWalletValue, err := getWalletValue(stub, loanAccruedWalletID)
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund loanAccruedWalletValue " + err.Error())
	}

	if (loanDisbursedWalletValue + loanChargesWalletValue + loanAccruedWalletValue) != 0 {

		errString := fmt.Sprintf("The wallet values are not zero loanDisbursedWalletValue: %d; loanChargesWalletValue:%d ;loanAccruedWalletValue:%d", loanDisbursedWalletValue, loanChargesWalletValue, loanAccruedWalletValue)
		return shim.Error(errString)
	}
	//####################################################################################################################

	//#####################################################################################################################
	//Calling for updating Business Main_Wallet
	//####################################################################################################################

	cAmtString := args[5]
	dAmtString := "0"

	
	walletID, err := getWalletID(stub, "businesscc", args[7], "main")
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund Business Main WalletID " + err.Error())
	}

	openBalance, err := getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund Business Main WalletValue " + err.Error())
	}
	openBalString := strconv.FormatInt(openBalance, 10)
	bal := openBalance + amt

	response1 := walletUpdation(stub, walletID, bal)
	if response1.Status != shim.OK {
		return shim.Error("marginrefundcc: " + response1.Message)
	}
	txnBalString := strconv.FormatInt(bal, 10)
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger

	u1 := uuid.New()

	fmt.Printf("generated Version 4 UUID Business Main_Wallet %v", u1)
	StringUUID1 := u1.String();
	fmt.Print("StringUUID1 ",StringUUID1);

	argsList := []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	txnResponse := putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("marginrefundcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Bank Main_Wallet
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, err = getWalletID(stub, "bankcc", args[6], "main")
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund Bank Main WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund Bank Main WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance - amt
	txnBalString = strconv.FormatInt(bal, 10)

	response2 := walletUpdation(stub, walletID, bal)
	if response2.Status != shim.OK {
		return shim.Error("marginrefundcc: " + response2.Message)
	}
	u2 := uuid.New()
	fmt.Printf("generated Version 4 UUID  Bank Main_Wallet %v", u2)
	StringUUID2 := u2.String();
	fmt.Print("StringUUID2 ",StringUUID2);
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("marginrefundcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Bank Refund_Wallet
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, err = getWalletID(stub, "bankcc", args[6], "liability")
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund Bank Refund_WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("marginrefundcc: " + "Margin Refund Bank Refund_WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance - amt
	txnBalString = strconv.FormatInt(bal, 10)

	response3 := walletUpdation(stub, walletID, bal)
	if response3.Status != shim.OK {
		return shim.Error("marginrefundcc: " + response3.Message)
	}
	u3 := uuid.New()
	fmt.Printf("generated Version 4 UUID Bank Refund_Wallet %v", u3)
	StringUUID3 := u3.String();
	fmt.Print("StringUUID3 ",StringUUID3);
	argsList = []string{StringUUID3, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("marginrefundcc: " + txnResponse.Message)
	}

	//####################################################################################################################

	//Getting sellerId using loanID
	chaincodeArgs1 := toChaincodeArgs("getSellerID", args[3])
	response = stub.InvokeChaincode("loancc", chaincodeArgs1, "myc")
	if response.Status != shim.OK {
		return shim.Error("marginrefundcc: " + response.Message)
	}

	sellerID := string(response.Payload)
	walletID, err = getWalletID(stub, "businesscc", sellerID, "loan")
	//Calling instrument chaincode to update the status
	if args[2] == "collected" {
		argsList := []string{args[4], sellerID, "settled"}
		argsListStr := strings.Join(argsList, ",")
		chaincodeArgs := toChaincodeArgs("updateInsStatus", argsListStr)
		response := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("marginrefundcc: " + response.Message)
		}
	}
	fmt.Println("******************** End marginRefund *************************")
	return shim.Success(nil)
}

func putInTxnBal(stub shim.ChaincodeStubInterface, argsListStr string) pb.Response {

	chaincodeArgs := toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the txnbalcc chaincode from Margin Refund")
	response := stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("marginrefundcc: " + response.Message)
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
		return "0", errors.New("marginrefundcc: " + response.Message)
	}
	walletID := string(response.GetPayload())
	return walletID, nil

}

func getWalletValue(stub shim.ChaincodeStubInterface, walletID string) (int64, error) {

	walletArgs := toChaincodeArgs("getWallet", walletID)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return 0, errors.New("marginrefundcc: " + walletResponse.Message)
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
		return shim.Error("marginrefundcc: " + walletResponse.Message)
	}
	return shim.Success(nil)

}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start Margin Refund chaincode:", err)
	}
}
