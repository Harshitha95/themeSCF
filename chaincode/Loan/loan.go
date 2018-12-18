package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"errors"
	"strings"
	"time"
	"html"
	"bytes"
/* 	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/hyperledger/fabric/protos/msp" */
	"math"
	"math/rand"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type SimpleChaincode struct {
}

type loanInfo struct {
	InstNum                     string    `json:"InstrumentNo"`          //[1]//Instrument Number
	ExposureBusinessID          string    `json:"ExposureBusinessID"`    //[2]//buyer for now
	ProgramID                   string    `json:"ProgramID"`             //[3]
	SanctionAmt                 int64     `json:"SanctionAmt"`           //[4]
	SanctionDate                time.Time `json:"SanctionDate"`          //auto generated as created
	SanctionAuthority           string    `json:"SanctionAuthority"`     //[5]
	ROI                         float64   `json:"ROI"`                   //[6]
	DueDate                     time.Time `json:"DueDate"`               //[7]
	ValueDate                   time.Time `json:"ValueDate"`             //[8]//with time
	LoanStatus                  string    `json:"LoanStatus"`            //[]
	LoanDisbursedWalletID       string    `json:"DisbursementWallet"`    //[9]
	LoanChargesWalletID         string    `json:"ChargesWallet"`         //[10]
	LoanAccruedInterestWalletID string    `json:"AccruedInterestWallet"` //[11]
	BuyerBusinessID             string    `json:"BuyerID"`               //[12]
	SellerBusinessID            string    `json:"SellerID"`              //[13]
//	InterestInAdvanceCollection string	  `json:"InterestInAdvanceCollection"`
}

type loanInfoPPR struct {
	InstNum                     string    `json:"InstrumentNo"`          //[1]//Instrument Number
	ExposureBusinessID          string    `json:"ExposureBusinessID"`    //[2]//buyer for now
	ProgramID                   string    `json:"ProgramID"`             //[3]
	SanctionAmt                 int64     `json:"SanctionAmt"`           //[4]
	SanctionDate                time.Time `json:"SanctionDate"`          //auto generated as created
	SanctionAuthority           string    `json:"SanctionAuthority"`     //[5]
	ROI                         float64   `json:"ROI"`                   //[6]
	DueDate                     time.Time `json:"DueDate"`               //[7]
	ValueDate                   time.Time `json:"ValueDate"`             //[8]//with time
	LoanStatus                  string    `json:"LoanStatus"`            //[]
	LoanDisbursedWalletID       string    `json:"DisbursementWallet"`    //[9]
	LoanChargesWalletID         string    `json:"ChargesWallet"`         //[10]
	LoanAccruedInterestWalletID string    `json:"AccruedInterestWallet"` //[11]
	BuyerBusinessID             string    `json:"BuyerID"`               //[12]
	SellerBusinessID            string    `json:"SellerID"`              //[13]
	PPRId						string	  `json:PPRId"`
//	InterestInAdvanceCollection string	  `json:"InterestInAdvanceCollection"`
}

type loanInfoVal struct {
	
	LoanDisbursedWalletID       string    `json:"DisbursementWallet"`    //[1]
	LoanChargesWalletID         string    `json:"ChargesWallet"`         //[2]
	LoanAccruedInterestWalletID string    `json:"AccruedInterestWallet"` //[3]
	LoanDisbursedWalletBal      string    `json:"DisbursementWalletBal"`    //[4]
	LoanChargesWalletBal        string	  `json:"ChargesWalletBal"`         //[5]
	LoanAccruedInterestWalletBal string   `json:"AccruedInterestWalletBal"` //[6]
}

type loanstruct struct {

	Key             string    `json:"Key"`
	Record          loanInfoPPR 

}
// Converts the args to byte of args
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

	if function == "newLoanInfo" {
		//Creates a new Loan Data
		return newLoanInfo(stub, args)
	} else if function == "getLoanInfo" {
		//Retrieves the existing data
		return getLoanInfo(stub, args)
	} else if function == "getLoanwithPPR" {
		//Retrieves the existing data
		return getLoanwithPPR(stub, args)
	} else if function == "updateLoanInfo" {
		//Updates variables for loan structure
		return updateLoanInfo(stub, args)
	} else if function == "loanIDexists" {
		//Checks the existence of loan ID
		return loanIDexists(stub, args[0])
	} else if function == "getLoanStatus" {
		//Returns the Loan status
		return getLoanStatus(stub, args[0])
	} else if function == "getLoanSancAmt" {
		//Returns the Sanc Amt
		return getLoanSancAmt(stub, args[0])
	} else if function == "getWalletID" {
		//Returns the walletID for the required wallet type
		return getWalletID(stub, args)
	} else if function == "getSellerID" {
		//Returns the Seller Id
		return getSellerID(stub, args[0])
	} else if function == "getBuyerID" {
		return getBuyerID(stub, args[0])
	} else if function == "queryLoan" {
		return queryLoan(stub, args)
	} else if function == "getWalletsofLoan" {
		return getWalletsofLoan(stub, args)
	}  else if function == "updateLoanStatusToOverDue" {
		return updateLoanStatusToOverDue(stub, args)
	} else if function == "getWLoan" {
		//return getWLoan(stub, args)
		jsonresp,err1 := getWLoan(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			return shim.Success([]byte(jsonresp))
		}
	} /* else if function == "getCreator" {
		return getCreator(stub)
	} else if function == "getLoanInfoBalTemp" {
		//Returns the walletID for the required Loan wallet type
		jsonrespasstruct,err1 := getLoanInfoBalTemp(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			//return shim.Success([]byte(strings.Join(jsonresp,"")))
			return shim.Success([]byte(jsonrespasstruct))
		}
	}  else if function == "updateInterestInAdvanceCollection" {
		//Updates variables for loan interest in advance flag
		return updateInterestInAdvanceCollection(stub, args)
	}	

 else if function == "getInterestCollectedStatus" {
		//Checks the existence of loan ID
		return getInterestCollectedStatus(stub, args[])
	}  */
	return shim.Error("loancc: " + "No function named " + function + " in Loan")
}


/* func (t *chainCode) getCreator(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Printf("\nBegin*** getCreator \n")
	creator, err := stub.GetCreator()
	if err != nil {
		fmt.Printf("GetCreator Error")
		return shim.Error(err.Error())
	}

	si := &msp.SerializedIdentity{}
	err2 := proto.Unmarshal(creator, si)
	if err2 != nil {
		fmt.Printf("Proto Unmarshal Error")
		return shim.Error(err2.Error())
	}
	buf := &bytes.Buffer{}
	protolator.DeepMarshalJSON(buf, si)
	fmt.Printf("End*** getCreator \n")
	fmt.Printf(string(buf.Bytes()))

	return shim.Success([]byte(buf.Bytes()))
} */

func newLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************While Writing Loan *************************")
	//Added instrument Status Update code 
	if len(args) != 14 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("loancc: " + "Invalid number of arguments in newLoanInfo(loan) (required:14) given: " + xLenStr)
	}
	fmt.Println("args[0]", args[0])
	fmt.Println("args[1]", args[1])
	fmt.Println("args[2]", args[2])
	fmt.Println("args[3]", args[3])
	fmt.Println("args[4]", args[5])
	fmt.Println("args[5]", args[5])
	fmt.Println("args[6]", args[6])
	fmt.Println("args[7]" ,args[7])
	fmt.Println("args[8] ",args[8])
	fmt.Println("args[9] ",args[9])
	
	//Checking existence of loanID
	println("Checking existence of loanID")
/* 	response := loanIDexists(stub, args[0])
	if response.Status != shim.OK {
		return shim.Error("loancc: " + response.Message)
	} */
	//Checking existence of ExposureBusinessID
	println("Checking existence of ExposureBusinessID")
	chaincodeArgs := toChaincodeArgs("bisIDexists", args[2])
	response1 := stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response1.Status == shim.OK {
		return shim.Error("loancc: " + "ExposureBusinessID " + args[2] + " does not exits")
	}
	//Checking if Instrument ID is Instrument Ref. No.
	println("Checking if Instrument ID is Instrument Ref. No.")
	chaincodeArgs = toChaincodeArgs("getInstrument", args[1], args[13])
	response2 := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
	if response2.Status != shim.OK {
		return shim.Error("loancc: " + "Instrument refrence no " + args[1] + " does not exits")
	}
	// getting the sanction amount from the instrument
	chaincodeArgs = toChaincodeArgs("getInstrumentAmt", args[1], args[13])
	response3 := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
	if response3.Status != shim.OK {
		return shim.Error("loancc: " + response3.Message)
	}
	instAmtStr := string(response3.Payload)
	fmt.Println("instAmtStr: " + instAmtStr)
	instAmt, err := strconv.ParseInt(instAmtStr, 10, 64)
	if err != nil {
		return shim.Error("loancc: " + "Unable to parse instAmt(loan): " + err.Error())
	}
	//Getting the discount percentage
	println("Getting the discount percentage")
	chaincodeArgs = toChaincodeArgs("discountPercentage", args[3], args[2])
	response4 := stub.InvokeChaincode("pprcc", chaincodeArgs, "myc")
	if response4.Status == shim.OK {
		return shim.Error("loancc: " + "PprId " + args[8] + " does not exits")
	}
	discountPercentStr := string(response4.Payload)
	discountPercent, _ := strconv.ParseInt(discountPercentStr, 10, 64)
	amt := instAmt - ((discountPercent * instAmt) / 100)
	//SanctionAmt -> sAmt
	println("SanctionAmt -> sAmt")
	sAmt, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	if sAmt > amt && sAmt == 0 {
		return shim.Error("loancc: " + "Sanction amount exceeds the required value or it is zero : " + args[4])
	}
	//SanctionDate ->sDate
	println("SanctionDate ->sDate")
	sDate := time.Now()
	roi, err := strconv.ParseFloat(args[6], 32)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	//Parsing into date for storage but hh:mm:ss will also be stored as
	println("Parsing into date for storage")
	//00:00:00 .000Z with the date
	//DueDate -> dDate 02-01-2006T15:04:05
	dDate, err := time.Parse("02/01/2006", args[7])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	if dDate.Weekday().String() == "Sunday" {
		fmt.Println("Since the due date falls on sunday, due date is extended to Monday(loan) : ", dDate.AddDate(0, 0, 1))
	}
	dDate = dDate.AddDate(0, 0, 1)
	//Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing
	println("Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing")
	vDateStr := args[8][:10]
	vTime := args[8][11:]
	vStr := vDateStr + "T" + vTime

	//ValueDate ->vDate
	println("ValueDate ->vDate")
	vDate, err := time.Parse("02/01/2006T15:04:05", vStr)
	//vDate, err := time.Parse("02/01/2006", vStr)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	hash := sha256.New()
	println("Hashing wallets")
	// Hashing LoanDisbursedWalletID
	LoanDisbursedWalletStr := args[9] + "LoanDisbursedWallet"
	hash.Write([]byte(LoanDisbursedWalletStr))
	md := hash.Sum(nil)
	LoanDisbursedWalletIDsha := hex.EncodeToString(md)
	fmt.Print("LoanDisbursedWalletIDsha ",LoanDisbursedWalletIDsha ,"args[9] " ,args[9])
	createWallet(stub, LoanDisbursedWalletIDsha, args[9])

	// Hashing LoanChargesWalletID
	LoanChargesWalletStr := args[10] + "LoanChargesWallet"
	hash.Write([]byte(LoanChargesWalletStr))
	md = hash.Sum(nil)
	LoanChargesWalletIDsha := hex.EncodeToString(md)
	fmt.Print("LoanChargesWalletIDsha ",LoanChargesWalletIDsha ,"args[10] " ,args[10])
	createWallet(stub, LoanChargesWalletIDsha, args[10])

	// Hashing LoanAccruedInterestWalletID
	LoanAccruedInterestWalletStr := args[11] + "LoanAccruedInterestWallet"
	hash.Write([]byte(LoanAccruedInterestWalletStr))
	md = hash.Sum(nil)
	LoanAccruedInterestWalletIDsha := hex.EncodeToString(md)
	fmt.Print("LoanAccruedInterestWalletIDsha ",LoanAccruedInterestWalletIDsha ,"args[11] " ,args[11])
	createWallet(stub, LoanAccruedInterestWalletIDsha, args[11])

	//Checking existence of BuyerBusinessID
	println("Checking existence of BuyerBusinessID")
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[13])
	response5 := stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response5.Status == shim.OK {
		return shim.Error("loancc: " + "BuyerBusinessID " + args[13] + " does not exits")
	}
	//Checking existence of SellerBusinessID
	println("Checking existence of SellerBusinessID")
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[13])
	response6 := stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response6.Status == shim.OK {
		return shim.Error("loancc: " + "SellerBusinessID " + args[13] + " does not exits")
	}

	fmt.Println("marshalling loaninfo")
	loan := loanInfo{args[1], args[2], args[3], sAmt, sDate, args[6], roi, dDate, vDate, "sanctioned", LoanDisbursedWalletIDsha, LoanChargesWalletIDsha, LoanAccruedInterestWalletIDsha, args[12], args[13]}
	loanBytes, err := json.Marshal(loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	stub.PutState(args[0], loanBytes)
//********************************
fmt.Println("Business Liability Wallet Update code ")
cAmtString := args[4]
dAmtString := "0"

walletID, err := getWalletIDForBus(stub, "businesscc", args[12], "liability")
if err != nil {
	return shim.Error("businesscc: " + "Sanction Business Liability WalletID " + err.Error())
}
println("Got WalletId" ,walletID)
openBalance, err := getWalletValue(stub, walletID)
if err != nil {
	return shim.Error("businesscc: " + "Sanction Business Liability WalletValue " + err.Error())
}
openBalString := strconv.FormatInt(openBalance, 10)
println("Got Opening Balance",openBalString)
//amt, _ = strconv.ParseInt(args[5], 10, 64)

bal := openBalance + amt
i := float64(bal)
	
	if math.Signbit(i) {
		return shim.Error("Invalid tranaction Amount for Business Liability Wallet")
	}
txnBalString := strconv.FormatInt(bal, 10)
println("openBalString ",openBalString)
println("Updating Balance ",bal)
response7 := walletUpdation(stub, walletID, bal)
if response7.Status != shim.OK {
	return shim.Error("loancc: " + response7.Message)
}
hash1 := sha256.New()
txnId1 := rand.Intn(10000);
txnId12 := rand.Intn(10000);
hash1.Write([]byte(strconv.Itoa(txnId1)))
hash1.Write([]byte(strconv.Itoa(txnId12)))
md1 := hash1.Sum(nil)
txnIDsha1 := hex.EncodeToString(md1)
md12 := hash1.Sum(nil)
txnIDsha12 := hex.EncodeToString(md12)

argsList := []string{txnIDsha1, txnIDsha12,args[7], args[0], args[1], walletID, openBalString, "loan_sanction",args[4], cAmtString, dAmtString, txnBalString, args[5]}
argsListStr := strings.Join(argsList, ",")
txnResponse := putInTxnBal(stub, argsListStr)
if txnResponse.Status != shim.OK {
	println("Error in PutTxnBal");
	return shim.Error("loancc: " + txnResponse.Message)
	}
println("End of Transaction")
//********************
	println("changing inst status")
	//argsList := []string{args[1], args[13], "sanctioned"}
	//argsListStr := strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("updateInstrumentStatus", args[1], args[13], "sanctioned")
	response8 := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
	if response8.Status != shim.OK {
		return shim.Error("loancc: " + response8.Message)
	}
	println("End of Loan ")
	return shim.Success([]byte("Successfully added loan info into ledger"))
}

func createWallet(stub shim.ChaincodeStubInterface, walletID string, amt string) pb.Response {
	chaincodeArgs := toChaincodeArgs("newWallet", walletID, amt)
	response9 := stub.InvokeChaincode("walletcc", chaincodeArgs, "myc")
	if response9.Status != shim.OK {
		return shim.Error("loancc: " + "Unable to create new wallet from business")
	}
	return shim.Success([]byte("created new wallet from business"))
}

func loanIDexists(stub shim.ChaincodeStubInterface, loanID string) pb.Response {
	fmt.Println("******** Start loanIDexists with arg ",loanID)
	ifExists, _ := stub.GetState(loanID)
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("loancc: " + "LoanId " + loanID + " exits. Cannot create new ID")
	}
	fmt.Println("******** end loanIDexists")
	return shim.Success(nil)
}

func getLoanStatus(stub shim.ChaincodeStubInterface, loanID string) pb.Response {
	fmt.Println(" ********* loancc: inside getLoanStatus")
	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "Error unmarshiling in loanstatus(loan):" + err.Error())
	}
	fmt.Println(" ********* End  getLoanStatus")
	return shim.Success([]byte(loan.LoanStatus))
}

/* func getInterestCollectedStatus(stub shim.ChaincodeStubInterface, loanID string) pb.Response {
	fmt.Println(" ********* loancc: inside getInterestCollectedStatus")
	loanBytes1, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes1 == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + loanID)
	}

	loan1 := loanInfo{}
	err1 := json.Unmarshal(loanBytes1, &loan1)
	if err1 != nil {
		return shim.Error("loancc: " + "Error unmarshiling in InterestCollectedStatus(loan):" + err1.Error())
	}
	fmt.Println(" ********* End  getInterestCollectedStatus")
	return shim.Success([]byte(loan1.InterestInAdvanceCollection))
} */
func getLoanSancAmt(stub shim.ChaincodeStubInterface, loanID string) pb.Response {
	fmt.Println("loancc: inside getLoanSancAmt")
	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "Error unmarshiling in loanstatus(loan):" + err.Error())
	}
	fmt.Println(loan.SanctionAmt)
	sancAmtString := strconv.FormatInt(loan.SanctionAmt, 10)
	fmt.Println(sancAmtString)
	fmt.Println(" ********* End  getLoanSancAmt")
	return shim.Success([]byte(sancAmtString))
	
}

func getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(" ********* Inside getWalletID with arg" ,args[0])
	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "Unable to parse into loan the structure (loanWalletValues)" + err.Error())
	}

	walletID := ""

	switch args[1] {
	case "accrued":
		walletID = loan.LoanAccruedInterestWalletID
	case "charges":
		walletID = loan.LoanChargesWalletID
	case "disbursed":
		walletID = loan.LoanDisbursedWalletID
	default:
		return shim.Error("loancc: " + "There is no wallet of this type in Loan :" + args[1])
	}
	fmt.Println(" ********* End getWalletID ********* ")
	return shim.Success([]byte(walletID))
}

func getLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(" ********* Inside getLoanInfo with arg * ",args[0])
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("loancc: " + "Invalid number of arguments in getLoanInfo (required:1) given:" + xLenStr)

	}

	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	loanString := fmt.Sprintf("%+v", loan)
	fmt.Printf("Loan Info:%s\n ", loanString)
	
	fmt.Println(" End getLoanInfo")
	return shim.Success(loanBytes)
}

func getLoanwithPPR(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(" ********* Inside getLoanwithPPR with arg 456 * ",args[0])
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("loancc: " + "Invalid number of arguments in getLoanInfo (required:1) given:" + xLenStr)

	}
	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}
	loan1 := loanInfoPPR{}
	err = json.Unmarshal(loanBytes, &loan1)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	//add code to get PPRId from instrument Id 
	//instrumentNumber =: loan.InstNum;
	chaincodeArgs := toChaincodeArgs("getPPRID", loan1.InstNum, loan1.SellerBusinessID)
		response := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			fmt.Println("instrumentcc :",response.Payload)
			return shim.Error("instrumentcc: " + response.Message)
		}
		pprid := string(response.GetPayload())

	loan1.PPRId=pprid;
	loanString := fmt.Sprintf("%+v", loan1)
	fmt.Printf("Loan Info:%s\n ", loanString)
	loanBytes1, _ := json.Marshal(loan1)
	fmt.Println(" End getLoanInfo")
	return shim.Success(loanBytes1)
}

func getSellerID(stub shim.ChaincodeStubInterface, loanID string) pb.Response {

	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID (getSellerID): " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	return shim.Success([]byte(loan.SellerBusinessID))
}

func getBuyerID(stub shim.ChaincodeStubInterface, loanID string) pb.Response {

	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID (getSellerID): " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	return shim.Success([]byte(loan.BuyerBusinessID))
}


func updateLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/*
		Updating the variables for loan structure
	*/
	//args = strings.Split(args[0], ",")
	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "error in unmarshiling loan: in updateLoanInfo" + err.Error())
	}

	// To change the LoanStatus from "sanction" to "disbursed"
	if args[2] == "disbursement" {
		if (loan.LoanStatus != "sanctioned") && (loan.LoanStatus != "part disbursed") {
			return shim.Error("loancc: " + "Loan is not Sanctioned, so cannot be disbursed/ part Disbursed : " + loan.LoanStatus)
		}
		//Updating Loan status for disbursement
		loan.LoanStatus = args[1]
		loanBytes, _ := json.Marshal(loan)
		err = stub.PutState(args[0], loanBytes)
		if err != nil {
			return shim.Error("loancc: " + "Error in loan updation " + err.Error())
		}

		//Calling instrument chaincode to update the status
		//argsList := []string{loan.InstNum, loan.SellerBusinessID, "disbursed"}
		//argsListStr := strings.Join(argsList, ",")
		chaincodeArgs := toChaincodeArgs("updateInstrumentStatus", loan.InstNum, loan.SellerBusinessID, "disbursed")
		response := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("loancc: " + response.Message)
		}
		return shim.Success([]byte("sanction updated succesfully"))

	} else if (args[1] == "repayment") && ((args[2] == "collected") || (args[2] == "part collected")) {
		if (loan.LoanStatus != "disbursed") && (loan.LoanStatus != "part disbursed") {
			return shim.Error("loancc: " + "Loan is not disbursed or part disbursed, so cannot be repayed")
		}
		//Updating Loan status for repayment
		loan.LoanStatus = args[2]
		loanBytes, _ = json.Marshal(loan)
		err = stub.PutState(args[0], loanBytes)
		if err != nil {
			return shim.Error("loancc: " + "Error in loan status updation " + err.Error())
		}

		return shim.Success([]byte("Successfully updated loan status with data from repayment"))
	}
	return shim.Error("loancc: " + "Invalid info for update loan")
}

func updateLoanStatusToOverDue(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println(" ********** Inside Update Status ****************** ")
	loanBytes1, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes1 == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}
	loan := loanInfo{}
	err = json.Unmarshal(loanBytes1, &loan)
	if err != nil {
		return shim.Error("loancc: " + "error in unmarshiling loan: in updateLoanInfo" + err.Error())
	}
	loan.LoanStatus = "overdue"
	loanBytes2, _ := json.Marshal(loan)
	err = stub.PutState(args[0], loanBytes2)
	if err != nil {
		return shim.Error("loancc: " + "Error in loan status updation " + err.Error())
	}
	fmt.Println(" ********** End Update Status ****************** ")
	return shim.Success([]byte(" Updated loan status successfully "))
}
/* func updateInterestInAdvanceCollection(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		Updating the InterestInAdvanceCollection
	
	
	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}
	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "error in unmarshiling loan: in updateLoanInfo" + err.Error())
	}
		loan.InterestInAdvanceCollection = args[1];
		loanBytes, _ = json.Marshal(loan)
		err = stub.PutState(args[0], loanBytes)
		if err != nil {
			return shim.Error("loancc: " + "Error in InterestInAdvanceCollection(loan)  updation " + err.Error())
		}
		return shim.Success([]byte("Successfully updated loan InterestInAdvanceCollection with data from repayment"))
	} */
func getLoanInfoBalTemp(stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {

	loanvalasstruct := loanInfoVal{}
	loan, err := stub.GetState(args[0])
	if err != nil {
		return nil,fmt.Errorf("loancc : " + "Unable to fetch the state" + err.Error())
	}
	if loan == nil {
		
		return nil,fmt.Errorf("loancc : " + "Data does not exist for " + args[0])
	}
	err = json.Unmarshal(loan, &loanvalasstruct)
	if err != nil {

		return nil,fmt.Errorf("loancc : " + "Uable to paser into the json format")
	}	
	data, err := json.Marshal(loanvalasstruct)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("data as percentage s %s",data)

return data,nil 
}

func getWLoan(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
	
	return "", fmt.Errorf("bankcc : " + "Invalid number of arguments in getWalletID(bank) (required:2) given:" + xLenStr)
	}
	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return "",fmt.Errorf("loancc: " + err.Error())
	} else if loanBytes == nil {
		return "", fmt.Errorf("loancc: No data exists on this loanID: " , args[0])
	}
	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return "",fmt.Errorf("loancc: Unable to parse into loan the structure (loanWalletValues)" ,err.Error())
	}
	walletID := ""

	switch args[1] {
	case "accrued":
		walletID = loan.LoanAccruedInterestWalletID
	case "charges":
		walletID = loan.LoanChargesWalletID
	case "disbursed":
		walletID = loan.LoanDisbursedWalletID
	default:
		return "",fmt.Errorf("loancc: " + "There is no wallet of this type in Loan :" + args[1])
	}
	fmt.Printf("walletID %s " ,walletID)
	
	return string(walletID),nil;
	
}
func getWalletsofLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	var response1 pb.Response  ;
	var response2 pb.Response  ;
	var response3 pb.Response  ;


	var intarg1 []string;
	var intarg2 []string;
	var intarg3 []string;
	
	intarg1=append(intarg1,args[0])
	intarg1=append(intarg1,args[1])

	intarg2=append(intarg2,args[0])
	intarg2=append(intarg2,args[2])

	intarg3=append(intarg3,args[0])
	intarg3=append(intarg3,args[3])
	
	var comres []string;  
	var loanValues []byte;

	var err4 error;

	disbursementWalletId,err1 := getWLoan(stub ,intarg1)
	chargesWalletId,err1 := getWLoan(stub ,intarg2)
	accruedWalletId,err1 := getWLoan(stub ,intarg3)
	
	if err1 != nil {
//
	}
	loanvalasstruct := loanInfoVal{}
	loanValues,err4 = getLoanInfoBalTemp(stub, args)
	fmt.Printf("loanValues %s",loanValues);
	 err := json.Unmarshal(loanValues, &loanvalasstruct)
	if err != nil {
		return shim.Error("Error while unmarshalling loanValues"+ err.Error())
    }  
	data, err1 := json.Marshal(loanvalasstruct)
	if err1 != nil {
		return shim.Error("error while marshaling  loanvalasstruct" + err1.Error() )
	} 
	fmt.Printf("Data as percentage at calling point %s",data)

	fmt.Printf("LoanValues from getBankInfoBalTemp %s" ,loanValues)
	
	if err4 != nil {
		return shim.Error(" Unable to get Loan Info >>>> "+err4.Error())
	} else {
			fmt.Print(" loanValues " ,loanValues)
	}
	
	chaincodeArgs1 := toChaincodeArgs("getWallet", disbursementWalletId)
	response1 = stub.InvokeChaincode("walletcc", chaincodeArgs1, "myc")
	comres=append(comres,string(response1.Payload))
	loanvalasstruct.LoanDisbursedWalletBal = string(response1.Payload);
	fmt.Println("LoanDisbursedWalletBal balance %s",string(response1.Payload))
	
	chaincodeArgs2 := toChaincodeArgs("getWallet", chargesWalletId)
	response2 = stub.InvokeChaincode("walletcc", chaincodeArgs2, "myc")
	comres=append(comres,string(response2.Payload))
	loanvalasstruct.LoanChargesWalletBal=string(response2.Payload);
	fmt.Println("LoanChargesWalletBal balance %s",string(response2.Payload))
	
	chaincodeArgs3 := toChaincodeArgs("getWallet", accruedWalletId)
	response3 = stub.InvokeChaincode("walletcc", chaincodeArgs3, "myc")
	comres=append(comres,string(response3.Payload))
	loanvalasstruct.LoanAccruedInterestWalletBal=string(response2.Payload);
	fmt.Println("LoanAccruedInterestWalletBal balance %s",string(response3.Payload))		
	data, err1 = json.Marshal(loanvalasstruct)
	if err1 != nil {
		return shim.Error("After adding error while marshaling  loanvalasstruct" + err1.Error() )
	} 	
	fmt.Printf("after asset data as percentage at calling point %s",data)
	fmt.Println("comres as string %s",comres);
	return shim.Success(data)
}
func putInTxnBal(stub shim.ChaincodeStubInterface, argsListStr string) pb.Response {

	chaincodeArgs := toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the txnbalcc chaincode from repayment")
	response := stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("loancc: " + response.Message)
	}
	fmt.Println(string(response.Payload))
	return shim.Success(nil)
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
func getWalletIDForBus(stub shim.ChaincodeStubInterface, ccName string, id string, walletType string) (string, error) {

	// STEP-1
	// using FromID, get a walletID from Bus structure

	chaincodeArgs := toChaincodeArgs("getWalletID", id, walletType)
	response := stub.InvokeChaincode(ccName, chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "0", errors.New(response.Message)
	}
	walletID := string(response.GetPayload())
	return walletID, nil

}

func queryLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	loanstructvarname := []loanstruct{}
	//var instrumentvalasstruct instrumentInfo;
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
//	queryString := args;
	fmt.Println(" Hello .. Inside queryLoan ")
	fmt.Println("args " ,html.UnescapeString(args[0]))
	//fmt.Println("args ",args);
	queryResults, err := getQueryResultForQueryString(stub, html.UnescapeString(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	} 	
	fmt.Println("queryLoan ",string(queryResults));
	///	businessValues,err4 = getBusinessInfoVal(stub, args)
	//var instrumentvalasstruct1;
	  err2 := json.Unmarshal(queryResults, &loanstructvarname)
	if err2 != nil {
		return shim.Error("Error while unmarshalling instrumentvalasstruct "+ err2.Error())
	}   
	//fmt.Println(" loanstruct  RefNo",loanstruct)
	//loanstructvarname.Record.InstNum;
//	fmt.Println("loanstructvarname.Record.InstNum ",loanstructvarname[0].Record.InstNum)
	//inst :=loanstruct.Record.InstNum;
	//sellID :=loanstruct.Record.SellerBusinessID
//	fmt.Println(" Key  ",instrumentvalasstruct.Key)
   // loanstruct.Record.PPRId:= getPPRID(inst,sellID);
   for i:=0; i < len(loanstructvarname); i++ {
    chaincodeArgs := toChaincodeArgs("getPPRID",  loanstructvarname[i].Record.InstNum,  loanstructvarname[i].Record.SellerBusinessID)
		response := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			fmt.Println("instrumentcc : getPPRID",response.Payload)
			return shim.Error("instrumentcc: " + response.Message)
		}
		loanstructvarname[i].Record.PPRId =string(response.Payload); 
	}
	data, err1 := json.Marshal(loanstructvarname)
	if err1 != nil {
		return shim.Error("error while marshaling  loanstruct" + err1.Error() )
	} 
	fmt.Println(" Before sending %s ",data)
	return shim.Success(data)
}
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
//	buff := buffer.String();
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	//fmt.Println("Buffer RefNo",buff.Record)
	return buffer.Bytes(), nil
}
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	var i int
	fmt.Println("Inside constructQueryResponseFromIterator")
	for i = 0; resultsIterator.HasNext();  i++ {
		queryResponse, err := resultsIterator.Next()
		fmt.Println("Inside resultsIterator ",queryResponse.Value)
		if err != nil {
			return nil, err
		}
		fmt.Println("queryResponse from constructQueryResponseFromIterator ",queryResponse.Value)
		// Add a comma before array members, suppress it for the first array member
		
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}


func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("loancc: "+"Error starting Loan chaincode: %s\n", err)
	}
}
