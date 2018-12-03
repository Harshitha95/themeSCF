package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"math"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type transactionInfo struct {
	TxnType string    `json:"TxnType"`      //args[1]
	TxnDate time.Time `json:"TxnDate"`      //args[2]
	LoanID  string    `json:"LoanID"`       //args[3]
	InsID   string    `json:"InstrumentID"` //args[4]
	Amt     int64     `json:"TxnAmount"`    //args[5]
	FromID  string    `json:"From"`         //args[6]
	ToID    string    `json:"To"`           //args[7]
	By      string    `json:"By"`           //args[8]
	//PprID   string    `json:"PPR_ID"`
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

	if function == "newTxnInfo" {
		//Creates new Transaction Information
		return newTxnInfo(stub, args)
	} else if function == "getTxnInfo" {
		//Retrieves an existing transcation information
		return getTxnInfo(stub, args)
	}
	return shim.Error("transactioncc: " + "No function named " + function + " in Transactionsssss")
}

func newTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start newTxnInfo *************************")
	//fmt.Println("Transaction Type",args[1])
	if len(args) != 9 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("transactioncc: " + "Invalid number of arguments in newTxnInfo(transactions) (required:9) given: " + xLenStr)
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

	tTypeValues := map[string]bool{
		"Disbursement":              true,
		"Repayment":                 true,
		"Marginrefund":             true,
		"Interestrefund":           true,
		"Penalinterestcollection": true,
		"loan_sanction":             true,
		"Charges":                   true,
		"Interestinadvance":       true,
		"Accrual":                   true,
		"Interestaccruedcharges":  true,
		"Penalcharges":             true,
		"TDS":                       true,
	}

	//Converting into lower case for comparison
	tTypeLower := strings.ToLower(args[1])
	if !tTypeValues[tTypeLower] {
		return shim.Error("transactioncc: " + "Invalid transaction type " + args[1])
	}
	//TxnDate -> tDate
	tDate, err := time.Parse("02/01/2006", args[2])
	if err != nil {
		return shim.Error("transactioncc: " + err.Error())
	}
	i, err := strconv.ParseFloat(args[5], 64)
	if err!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i) {
		return shim.Error("Invalid Amount")
	}
	amt, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error("transactioncc: " + err.Error())
	}

	//TODO: put it at last for redability

	var sellerID string
	var buyerID string

	switch tTypeLower {

	case "disbursement":
		//bank -> seller
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newDisbInfo", argsStr)
		fmt.Println("calling the disbursement chaincode")
		response := stub.InvokeChaincode("disbursementcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

		transaction := transactionInfo{tTypeLower, tDate, args[3], args[4], amt, args[6], args[7], args[8]}
		fmt.Println(transaction)

		txnBytes, err := json.Marshal(transaction)
		err = stub.PutState(args[0], txnBytes)
		if err != nil {
			return shim.Error("transactioncc: " + "Cannot write into ledger the transactino details")
		}
		fmt.Println("Successfully inserted disbursement transaction into the ledger")

	//#######################################################################################################

	case "repayment":
		//seller -> bank
		sellerID = getSellerID(stub, args[3])
		buyerID = getBuyerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[7], sellerID, buyerID, args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newRepayInfo", argsStr)
		fmt.Println("calling the repayment chaincode")
		response := stub.InvokeChaincode("repaymentcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

		transaction := transactionInfo{tTypeLower, tDate, args[3], args[4], amt, args[6], args[7], args[8]}
		fmt.Println(transaction)
		txnBytes, err := json.Marshal(transaction)
		err = stub.PutState(args[0], txnBytes)
		if err != nil {
			return shim.Error("transactioncc: " + "Cannot write into ledger the transaction details")
		}
		fmt.Println("Successfully inserted repayment transaction into the ledger")

	//#######################################################################################################

	case "marginrefund":
		//bank -> seller
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newMarginInfo", argsStr)
		fmt.Println("calling the margin_refundcc chaincode")
		response := stub.InvokeChaincode("margin_refundcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

		transaction := transactionInfo{tTypeLower, tDate, args[3], args[4], amt, args[6], args[7], args[8]}
		fmt.Println(transaction)
		txnBytes, err := json.Marshal(transaction)
		err = stub.PutState(args[0], txnBytes)
		if err != nil {
			return shim.Error("transactioncc: " + "Cannot write into ledger the transaction details")
		}
		fmt.Println("Successfully inserted margin refund transaction into the ledger")

	//#######################################################################################################

	case "interestrefund":
		//bank -> seller
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newInterestInfo", argsStr)
		fmt.Println("calling the interest_refundcc chaincode")
		response := stub.InvokeChaincode("interest_refundcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

		transaction := transactionInfo{tTypeLower, tDate, args[3], args[4], amt, args[6], args[7], args[8]}
		fmt.Println(transaction)
		txnBytes, err := json.Marshal(transaction)
		err = stub.PutState(args[0], txnBytes)
		if err != nil {
			return shim.Error("transactioncc: " + "Cannot write into ledger the transaction details")
		}
		fmt.Println("Successfully inserted interest refund transaction into the ledger")

	//#######################################################################################################

	case "penalinterestcollection":
		//seller -> bank
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newPICinfo", argsStr)
		fmt.Println("calling the penal_interest_collectioncc chaincode")
		response := stub.InvokeChaincode("penal_interest_collectioncc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

		transaction := transactionInfo{tTypeLower, tDate, args[3], args[4], amt, args[6], args[7], args[8]}
		fmt.Println(transaction)
		txnBytes, err := json.Marshal(transaction)
		err = stub.PutState(args[0], txnBytes)
		if err != nil {
			return shim.Error("transactioncc: " + "Cannot write into ledger the transaction details")
		}
		fmt.Println("Successfully inserted penal interest collection transaction into the ledger")

	//#######################################################################################################
	/*
		case "loan_sanction":
			argsStr := strings.Join(args, ",")
			chaincodeArgs := toChaincodeArgs("newLoanSancInfo", argsStr)
			fmt.Println("calling the loan_sanctioncc chaincode")
			response := stub.InvokeChaincode("loan_sanctioncc", chaincodeArgs, "myc")
			if response.Status != shim.OK {
				return shim.Error("transactioncc: " + response.Message)
			}
	*/
	//#######################################################################################################

	case "charges":

		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newChargesInfo", argsStr)
		fmt.Println("calling the chargescc chaincode")
		response := stub.InvokeChaincode("chargescc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

	//#######################################################################################################
	case "interestinadvance":
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newInterestInAdvInfo", argsStr)
		fmt.Println("calling the interest_in_advancecc chaincode")
		response := stub.InvokeChaincode("interest_in_advancecc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

	//#######################################################################################################

	case "accrual":
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newAccrualInfo", argsStr)
		fmt.Println("calling the accrualcc chaincode")
		response := stub.InvokeChaincode("accrualcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

	//#######################################################################################################

	case "interestaccruedcharges":
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newInterestAccruedInfo", argsStr)
		fmt.Println("calling the interest_accrued_chargescc chaincode")
		response := stub.InvokeChaincode("interest_accrued_chargescc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

	//#######################################################################################################

	case "penalcharges":
		sellerID = getSellerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, args[7], args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newPenalChargesInfo", argsStr)
		fmt.Println("calling the penal_chargescc chaincode")
		response := stub.InvokeChaincode("penal_chargescc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

	//#######################################################################################################

	case "tds":
		sellerID = getSellerID(stub, args[3])
		buyerID = getBuyerID(stub, args[3])
		txnArgs := []string{args[0], args[1], args[2], args[3], args[4], args[5], args[6], sellerID, buyerID, args[8]}
		argsStr := strings.Join(txnArgs, ",")
		chaincodeArgs := toChaincodeArgs("newTDSInfo", argsStr)
		fmt.Println("calling the tdscc chaincode")
		response := stub.InvokeChaincode("tdscc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("transactioncc: " + response.Message)
		}

	//#######################################################################################################

	default:
		fmt.Println("incorrect txnType")
		return shim.Error("incorrect txnType from txncc")
	}
	fmt.Println("******************** End newTxnInfo *************************")
	return shim.Success(nil)
}

func getTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("transactioncc: " + "Invalid number of arguments in getTxnInfo (required:1) given: " + xLenStr)
	}

	txnBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if txnBytes == nil {
		return shim.Error("transactioncc: " + "No data exists on this txnID: " + args[0])
	}

	transaction := transactionInfo{}
	err = json.Unmarshal(txnBytes, &transaction)
	if err != nil {
		return shim.Error("transactioncc: " + "error while unmarshaling:" + err.Error())
	}

	tString := fmt.Sprintf("%+v", transaction)
	return shim.Success([]byte(tString))

}

func getSellerID(stub shim.ChaincodeStubInterface, loanID string) string {

	chaincodeArgs := toChaincodeArgs("getSellerID", loanID)
	fmt.Println("transactioncc: " + "calling the loan chaincode")
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "not_found"
	}
	return string(response.GetPayload())
}

func getBuyerID(stub shim.ChaincodeStubInterface, loanID string) string {

	chaincodeArgs := toChaincodeArgs("getBuyerID", loanID)
	fmt.Println("transactioncc: " + "calling the loan chaincode")
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "not_found"
	}
	return string(response.GetPayload())
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("transactioncc: "+"Error starting Transaction chaincode: %s\n", err)
	}
}
