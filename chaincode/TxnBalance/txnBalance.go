package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type txnBalanceInfo struct {
	TxnID      string    `json:"TxnID"`
	TxnDate    time.Time `json:"TxnDate"`
	LoanID     string    `json:"LoanID"`
	InsID      string    `json:"InsID"`
	WalletID   string    `json:"WalletID"`
	OpeningBal int64     `json:"OpeningBalance"`
	TxnType    string    `json:"TxnType"`
	Amt        int64     `json:"Amount"`
	CAmt       int64     `json:"CreditAmount"`
	DAmt       int64     `json:"DebitAmount"`
	TxnBal     int64     `json:"TxnBalance"`
	By         string    `json:"By"`
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "putTxnBalInfo" { //Inserting a New Business information
		return c.putTxnBalInfo(stub, args)
	} else if function == "getTxnBalInfo" { // To view a Business information
		return c.getTxnBalInfo(stub, args)
	}
	return shim.Error("txnbalcc: " + "Inside txnBalcc:Invoke(), Function does not exit" + function)
}

func (c *chainCode) putTxnBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start putTxnBalInfo *************************")
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
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
	if len(args) != 13 {
		return shim.Error("txnbalcc: " + "Invalid number of arguments for txnBal. Needed 13 arguments")
	}

	ifExists, err := stub.GetState(args[0])
	if ifExists != nil {
		return shim.Error("txnbalcc: " + "TxnBalanceId " + args[0] + " exits. Cannot create new ID")
	}

	//TxnDate ->txnDate
	txnDate, err := time.Parse(time.RFC3339, args[2])
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in txndate " + err.Error())
	}

	openBal, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in openbal " + err.Error())
	}

	txnTypeValues := map[string]bool{
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

	txnTypeLower := strings.ToLower(args[7])
	if !txnTypeValues[txnTypeLower] {
		return shim.Error("txnbalcc: " + "txnbal Invalid Transaction type"+args[7])
	}

	amt, err := strconv.ParseInt(args[8], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in amt" + err.Error())
	}

	cAmt, err := strconv.ParseInt(args[9], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in camt" + err.Error())
	}

	dAmt, err := strconv.ParseInt(args[10], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in damt " + err.Error())
	}

	txnBal, err := strconv.ParseInt(args[11], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in txnbal " + err.Error())
	}

	txnBalance := txnBalanceInfo{args[1], txnDate, args[3], args[4], args[5], openBal, txnTypeLower, amt, cAmt, dAmt, txnBal, args[12]}
	txnBalanceBytes, err := json.Marshal(txnBalance)
	if err != nil {
		return shim.Error("txnbalcc: " + err.Error())
	}
	err = stub.PutState(args[0], txnBalanceBytes)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal cannot write to ledger: " + err.Error())
	}
	fmt.Println("Succefully wrote txnID " + args[0] + " into the ledger")
	fmt.Println("******************** end putTxnBalInfo *************************")
	return shim.Success(nil)

}

func (c *chainCode) getTxnBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("txnbalcc: " + "Required only one argument")
	}

	txnBalance := txnBalanceInfo{}
	txnBalanceBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("txnbalcc: " + "Failed to get the business information: " + err.Error())
	} else if txnBalanceBytes == nil {
		return shim.Error("txnbalcc: " + "No information is avalilable on this businessID " + args[0])
	}

	err = json.Unmarshal(txnBalanceBytes, &txnBalance)
	if err != nil {
		return shim.Error("txnbalcc: " + "Unable to parse into the structure " + err.Error())
	}
	jsonString := fmt.Sprintf("%+v", txnBalance)
	fmt.Printf("Transaction info %s : %s", args[0], jsonString)
	return shim.Success([]byte(jsonString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("txnbalcc: "+"Error starting Simple chaincode: %s\n", err)
	}
}
