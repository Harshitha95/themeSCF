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

	if function == "newRepayInfo" {
		return newRepayInfo(stub, args)
	}
	return shim.Error("repaymentcc: " + "no function named " + function + " found in Repayment")
}
func newRepayInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start Repayment *************************")
	
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("repaymentcc: " + "Invalid number of arguments in newRepayInfo(repayment) (required:10) given:" + xLenStr)
	}
	fmt.Print("args[0]", args[0])
	fmt.Print("args[1]", args[1])
	fmt.Print("args[2]", args[2])
	fmt.Print("args[3]", args[3])
	fmt.Print("args[4]", args[5])
	fmt.Print("args[5]", args[5])
	fmt.Print("args[6]", args[6])
	fmt.Print("args[7]" args[7])
	fmt.Print("args[8] ",,args[8])
	fmt.Print("args[9] ",,args[9])
	/*
	 *TxnType string    //args[1]
	 *TxnDate time.Time //args[2]
	 *LoanID  string    //args[3]
	 *InsID   string    //args[4]
	 *Amt     int64     //args[5]
	 *BankID  string    //args[6]  Bank
	 *SellID  string    //args[7]  Seller
	 *BuyID	  string	//args[8]  Buyer
	 *By      string    //args[9]
	 */

	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// The transaction object has been created and written into the ledger
	// The JSON file is 'transaction'function
	// Now to create a TXN_Bal_Update obj for 10 times
	// Calling TXN_Balance CC based on TXN_Type {ex: Disbursement}
	/*
			    a. Debiting (decreasing) Business Wallet (Buyer)
		            i. Txn amt
		        b. Crediting (Increasing) Bank Wallet
		            i. Txn amt
		        c. Debiting (decreasing) Bank Asset Wallet
		            i. Loan Disbursed Wallet Balance + Loan Charges Wallet Balance
		        d. Crediting (Increasing) Bank Refund Wallet (if applicable)
		            i. Txn Amt – Loan Disbursed Wallet balance – Loan Charges Wallet Balance
		        e. Debiting (decreasing) Business Loan Wallet (Seller)
		            i. If Txn Amt is >/= (Loan Charges Wallet Balance + Loan Disbursed Wallet Balance)
		                1. Loan disbursed Wallet Balance + Loan Charges Wallet Balance
		            ii. If Txn Amt is  < (Loan Charges Wallet Balance + Loan Disbursed Wallet Balance)
		                1. Txn Amt – Loan Charges Wallet Balance
		        f. Debiting (decreasing) Business Charges O/s Wallet
		            i. Loan Charges Wallet Balance
		        g. Debiting (decreasing) Business Principal O/s Wallet
		            i. Loan Disbursed Wallet Balance
		        h. Debiting (Decreasing) Loan Charges Wallet
		            i. Loan Charges Wallet Balance
		        i. Debiting (Decreasing) Loan Disbursed Wallet
		            i. If Txn Amt is >/= (Loan Charges Wallet Balance + Loan Disbursed Wallet Balance)
		                1. Loan Disbursed Wallet is reduced to Zero
		                2. Loan Status is updated to Collected
		            ii. If Txn Amt is  < (Loan Charges Wallet Balance + Loan Disbursed Wallet Balance)
		                1. Txn Amt – Loan Charges Wallet Balance
		                2. Loan Status is updated to Part Collected
		        j. Debiting (Decreasing) Business Liability Wallet (Buyer)
		            i. Txn Amt
	*/

	amt, _ := strconv.ParseInt(args[5], 10, 64)

	//#####################################################################################################################
	//Calling for updating Business Main_Wallet
	//####################################################################################################################

	cAmtString := "0"
	dAmtString := args[5]

	walletID, err := getWalletID(stub, "businesscc", args[8], "main")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Main WalletID " + err.Error())
	}

	openBalance, err := getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Main Wallet Value " + err.Error())
	}
	openBalString := strconv.FormatInt(openBalance, 10)
	bal := openBalance - amt

	response := walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + response.Message)
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
		return shim.Error("Invalid transaction Amount for Business Main Wallet Wallet")
	}
	argsList := []string{StringUUID1, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	txnResponse := putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Bank Main_Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, err = getWalletID(stub, "bankcc", args[6], "main")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Bank Main WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Bank Main WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance + amt
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + response.Message)
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
		return shim.Error("Invalid transaction Amount for Bank Main  Wallet")
	}
	argsList = []string{StringUUID2, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Loan Disbursed  Wallet
	//####################################################################################################################

	cAmtString = "0"

	//Loan Disbursed Wallet balance
	loanDisbursedWalletID, err := getWalletID(stub, "loancc", args[3], "disbursed")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment loanDisbursedWalletID " + err.Error())
	}
	loanDisbursedWalletValue, err := getWalletValue(stub, loanDisbursedWalletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment loanDisbursedWalletValue " + err.Error())
	}

	//Loan Charges Wallet Balance
	loanChargesWalletID, err := getWalletID(stub, "loancc", args[3], "charges")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment loanChargesWalletID " + err.Error())
	}
	loanChargesWalletValue, err := getWalletValue(stub, loanChargesWalletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment loanChargesWalletValue " + err.Error())
	}

	//Bank Asset Wallet
	walletID, err = getWalletID(stub, "bankcc", args[6], "asset")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Bank Asset WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Bank Asset WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	bal = openBalance - loanChargesWalletValue - loanDisbursedWalletValue
	fmt.Println("openBalance  Bank Asset  Update" ,openBalance);
	fmt.Println("loanChargesWalletValue  Bank Asset Update" ,loanChargesWalletValue);
	fmt.Println("loanDisbursedWalletValue  Bank Asset Update" ,loanDisbursedWalletValue);
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + "Repayment Bank Asset Wallet " + response.Message)
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
		return shim.Error("Invalid transaction Amount for Bank Asset Wallet")
	}
	dAmt := loanChargesWalletValue + loanDisbursedWalletValue
	dAmtString = strconv.FormatInt(dAmt, 10)
	argsList = []string{StringUUID3, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Bank liability Wallet
	//####################################################################################################################

	dAmtString = "0"
	walletID, err = getWalletID(stub, "bankcc", args[6], "liability")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Bank Liability WalletID " + err.Error())
	}
	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Bank Liability WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	fmt.Println(" openBalance in bank liability ",openBalance)
	fmt.Println(" amt ",amt)
	fmt.Println(" loanChargesWalletValue ",loanChargesWalletValue)
	fmt.Println(" loanDisbursedWalletValue ",loanDisbursedWalletValue)
		cAmt := amt - loanChargesWalletValue - loanDisbursedWalletValue

	if cAmt > 0 {
		bal = openBalance + cAmt
	} else {
		bal = openBalance
	}
	fmt.Println("bal for updating Bank liability  ",bal)
	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + "Repayment Bank Liability Wallet " + response.Message)
	}
	txnBalString = strconv.FormatInt(bal, 10)
	fmt.Println("txnBalString Ammount  Bank liability ",txnBalString)
	fmt.Println(" Bank liability walletID",walletID)
	cAmtString = strconv.FormatInt(cAmt, 10)
	u4 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v Bank liability", u4)
	StringUUID4 := u4.String();
	fmt.Print("StringUUID4 ",StringUUID4);
	i3, err3 := strconv.ParseFloat(txnBalString, 64)
	if err3!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i3) {
		return shim.Error("Invalid transaction Amount for Bank liability  Wallet")
	}
	argsList = []string{StringUUID4, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Business Loan_Wallet (seller)
	//####################################################################################################################

	// geting seller's ID using loan ID

	cAmtString = "0"
	chaincodeArgs := toChaincodeArgs("getSellerID", args[3])
	response = stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + response.Message)
	}
	sellerID := string(response.Payload)
	SellerLoanwalletID, err := getWalletID(stub, "businesscc", sellerID, "loan")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Loan_WalletID (seller) " + err.Error())
	}
	SellerLoanwalletopenBalance, err := getWalletValue(stub, SellerLoanwalletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Loan_WalletValue (seller) " + err.Error())
	}
	openBalString = strconv.FormatInt(SellerLoanwalletopenBalance, 10)

	if (amt > loanChargesWalletValue+loanDisbursedWalletValue) || (amt == loanChargesWalletValue+loanDisbursedWalletValue) {
		dAmt = loanChargesWalletValue + loanDisbursedWalletValue
	} else if amt < loanChargesWalletValue+loanDisbursedWalletValue {
		dAmt = amt - loanChargesWalletValue
	}
	fmt.Println(" openBalance in Business Loan_Wallet  ",SellerLoanwalletopenBalance)
	fmt.Println(" amt ",amt)
	fmt.Println(" dAmt ",dAmt)
	fmt.Println(" SellerLoanwalletID  ",SellerLoanwalletID)
	
	bal = SellerLoanwalletopenBalance - dAmt
	fmt.Println(" updatinf Balance in Business Loan_Wallet ",bal)
	txnBalString = strconv.FormatInt(bal, 10)
	response = walletUpdation(stub, SellerLoanwalletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + "Repayment Business Loan_Wallet (seller) " + response.Message)
	}
	dAmtString = strconv.FormatInt(dAmt, 10)
	u5 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u5)
	StringUUID5 := u5.String();
	fmt.Print("StringUUID5 ",StringUUID5);
	i4, err4 := strconv.ParseFloat(txnBalString, 64)
	if err4!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i4) {
		return shim.Error("Invalid transaction Amount for Business Loan  Wallet")
	}
	argsList = []string{StringUUID5, args[0], args[2], args[3], args[4], SellerLoanwalletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Business Charges/Interest O/s Wallet
	//####################################################################################################################

	cAmtString = "0"
	walletID, err = getWalletID(stub, "businesscc", args[7], "chargesOut")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Charges/Interest O/s WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Charges/Interest O/s WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)
	dAmt = loanChargesWalletValue
	bal = openBalance - loanChargesWalletValue
	txnBalString = strconv.FormatInt(bal, 10)
fmt.Println("openBalance Business Charges/Interest O/s Wallet ",openBalance)

fmt.Println(" Business Charges/Interest O/s Wallet",walletID)

fmt.Println("Updating Balance Business Charges/Interest O/s Wallet",bal);

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + "Repayment Business Charges/Interest O/s Wallet " + response.Message)
	}
	dAmtString = strconv.FormatInt(dAmt, 10)
	u6 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u6)
	StringUUID6 := u6.String();
	fmt.Print("StringUUID6 ",StringUUID6);
	i5, err5 := strconv.ParseFloat(txnBalString, 64)
	if err5!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i5) {
		return shim.Error("Invalid transaction Amount for Business Charges  Wallet")
	}
	argsList = []string{StringUUID6, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Business Principal O/s Wallet
	//####################################################################################################################

	cAmtString = "0"

	walletID, err = getWalletID(stub, "businesscc", args[7], "principalOut")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Principal O/s WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Business Principal O/s WalletValue " + err.Error())
	}
	if (amt >openBalance ) || (amt == openBalance) {
		dAmt = openBalance
	} else if amt < openBalance {
		dAmt = openBalance - amt
	}
	//openBalString = strconv.FormatInt(openBalance, 10)
	
	/* fmt.Println(" **** Added get getInterestCollectedStatus **** ")
	chaincodeArgs1 := toChaincodeArgs("getInterestCollectedStatus", args[3])
	response1 := stub.InvokeChaincode("loancc", chaincodeArgs1, "myc")
	if response1.Status != shim.OK {
		return shim.Error("interestAdvcc: can't get InterestCollectedStatus" + response1.Message)
	}
	InterestCollectedStatus := string(response1.Payload)
	fmt.Println( " **** End get getInterestCollectedStatus **** " ) */
//	fmt.Println("InterestCollectedStatus ",InterestCollectedStatus)
	//IsInterestCollectedStatus := false

 /*    if InterestCollectedStatus {
		fmt.Println("InterestCollectedStatus is true")
		bal = openBalance - SellerLoanwalletopenBalance
		txnBalString = strconv.FormatInt(bal, 10)
    }
    if !InterestCollectedStatus {
		fmt.Println("InterestCollectedStatus is false")
		bal = openBalance - loanDisbursedWalletValue
	txnBalString = strconv.FormatInt(bal, 10)
    } 
		if InterestCollectedStatus == "true" {
			fmt.Println("InterestCollectedStatus is true")
		bal = openBalance - SellerLoanwalletopenBalance
		txnBalString = strconv.FormatInt(bal, 10)
		}
		if  InterestCollectedStatus == "false"  {
			fmt.Println("InterestCollectedStatus is false")
			bal = openBalance - loanDisbursedWalletValue
			
		txnBalString = strconv.FormatInt(bal, 10)
		} */
		
	response2 := walletUpdation(stub, walletID, bal)
	if response2.Status != shim.OK {
		return shim.Error("repaymentcc: " + "Repayment Business Principal O/s Wallet " + response2.Message)
	}
	dAmt = loanDisbursedWalletValue
	dAmtString = strconv.FormatInt(dAmt, 10)
	u7 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u7)
	StringUUID7 := u7.String();
	fmt.Print("StringUUID7 ",StringUUID7);
	i6, err6 := strconv.ParseFloat(txnBalString, 64)
	if err6!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i6) {
		return shim.Error("Invalid transaction Amount for Business Principal  Wallet")
	}
	argsList = []string{StringUUID7, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}
	fmt.Println("0")
	//####################################################################################################################
	//Calling for updating Loan Charges Wallet
	//####################################################################################################################

	cAmtString = "0"
	walletID, err = getWalletID(stub, "loancc", args[3], "charges")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Loan Charges WalletID " + err.Error())
	}
	fmt.Println("1")

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Loan Charges WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)
	fmt.Println("2")

	bal = openBalance - loanChargesWalletValue
	txnBalString = strconv.FormatInt(bal, 10)
	fmt.Println("3")

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + "Repayment Loan Charges Wallet " + response.Message)
	}
	fmt.Println("4")

	dAmt = loanChargesWalletValue
	dAmtString = strconv.FormatInt(dAmt, 10)
	fmt.Println("5")
	u8 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u8)
	StringUUID8 := u8.String();
	fmt.Print("StringUUID8 ",StringUUID8);
	i7, err7 := strconv.ParseFloat(txnBalString, 64)
	if err7!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i7) {
		return shim.Error("Invalid transaction Amount for Loan Charges  Wallet")
	}
	argsList = []string{StringUUID8, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Loan Disbursed Wallet
	//####################################################################################################################

	cAmtString = "0"
	walletID, err = getWalletID(stub, "loancc", args[3], "disbursed")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Loan Disbursed WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment Loan Disbursed WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	if (amt > loanChargesWalletValue+loanDisbursedWalletValue) || (amt == loanChargesWalletValue+loanDisbursedWalletValue) {
		bal = 0
		dAmt = openBalance
		chaincodeArgs := toChaincodeArgs("updateLoanInfo", args[3], "repayment", "collected")
		response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("repaymentcc: " + response.Message)
		}
	} else if amt < loanChargesWalletValue+loanDisbursedWalletValue {
		dAmt = amt - loanChargesWalletValue
		bal = openBalance - dAmt
		chaincodeArgs := toChaincodeArgs("updateLoanInfo", args[3], "repayment", "part collected")
		response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("repaymentcc: " + response.Message)
		}
	}
	txnBalString = strconv.FormatInt(bal, 10)
	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + "Repayment Loan Charges Wallet " + response.Message)
	}
	dAmtString = strconv.FormatInt(dAmt, 10)
	u9 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u9)
	StringUUID9 := u9.String();
	fmt.Print("StringUUID9 ",StringUUID9);
	i8, err8 := strconv.ParseFloat(txnBalString, 64)
	if err8!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i8) {
		return shim.Error("Invalid transaction Amount for Loan Charges  Wallet")
	}
	argsList = []string{StringUUID9, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Business Liability Wallet (Buyer)
	//####################################################################################################################
	
	cAmtString = "0"
	dAmtString = args[5]

	walletID, err = getWalletID(stub, "businesscc", args[8], "liability")
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment businesscc Liability WalletID " + err.Error())
	}
	fmt.Print("walletID ",walletID);
	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("repaymentcc: " + "Repayment businesscc Liability WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)
	//reducing loan amt instead of  repayment amount 
	if (openBalance < amt) || (openBalance == amt) {
	bal = openBalance - amt
	} else { 
	bal = openBalance
	}
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + response.Message)
	}
	u10 := uuid.New()
	fmt.Printf("generated Version 4 UUID %v", u10)
	StringUUID10 := u10.String();
	fmt.Print("StringUUID10 ",StringUUID10)
	i9, err9 := strconv.ParseFloat(txnBalString, 64)
	if err9!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i9) {
		return shim.Error("Invalid transaction Amount for businesscc Liability  Wallet")
	};
	argsList = []string{StringUUID10, args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("repaymentcc: " + txnResponse.Message)
	}

	//####################################################################################################################
	fmt.Println("******************** End Repayment *************************")
	return shim.Success(nil)
}

func putInTxnBal(stub shim.ChaincodeStubInterface, argsListStr string) pb.Response {

	chaincodeArgs := toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the txnbalcc chaincode from repayment")
	response := stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("repaymentcc: " + response.Message)
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
	fmt.Println("walletUpdation  WalletId",walletID)
	fmt.Println("walletUpdation  Balance",txnBalString)
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
		fmt.Println("Unable to start Repayment chaincode:", err)
	}
}
