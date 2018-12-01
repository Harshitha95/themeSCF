package main

import (
	"math"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"bytes"
	"html"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type instrumentInfo struct {
	//Instrument ID for storing is auto generated
	ObjectType 		string	  `json:"docType"`  //Referring the Objecttype InstrumentInfo as doctype named 
	InstrumentRefNo string    `json:"InstrumentRefNo"`         //[0]
	InstrumentDate  time.Time `json:"Date"`          //[1]
	SellBusinessID  string    `json:"SellerID"`      //[2]
	BuyBusinsessID  string    `json:"BuyerID"`       //[3]
	InsAmount       string    `json:"Amount"`        //[4]// use int64 for convertion
	InsStatus       string    `json:"Status"`        // not required
	InsDueDate      time.Time `json:"DueDate"`       //[5]
	ProgramID       string    `json:"ProgramID"`     //[6]
	PPRid           string    `json:"PPRID"`         //[7]
	UploadBatchNo   string    `json:"UploadBatchNo"` //[8]
	ValueDate       time.Time `json:"ValueDate"`     //[9]
}

type instrumentInfoReturn struct {
	//Instrument ID for storing is auto generated
	Key             string    `json:"Key"`
	Record          instrumentInfo 
/* 	ObjectType 		string	  `json:"docType"`  //Referring the Objecttype InstrumentInfo as doctype named 
	InstrumentRefNo string    `json:"RefNo"`         //[0]
	InstrumentDate  time.Time `json:"Date"`          //[1]
	SellBusinessID  string    `json:"SellerID"`      //[2]
	BuyBusinsessID  string    `json:"BuyerID"`       //[3]
	InsAmount       string    `json:"Amount"`        //[4]// use int64 for convertion
	InsStatus       string    `json:"Status"`        // not required
	InsDueDate      time.Time `json:"DueDate"`       //[5]
	ProgramID       string    `json:"ProgramID"`     //[6]
	PPRid           string    `json:"PPRID"`         //[7]
	UploadBatchNo   string    `json:"UploadBatchNo"` //[8]
	ValueDate       time.Time `json:"ValueDate"`     //[9] */
}


func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	indexName := "InstrumentRefNo~SellBusinessID~InsAmount"
	inst := instrumentInfo{}

	refNoSellIDkey, err := stub.CreateCompositeKey(indexName, []string{inst.InstrumentRefNo, inst.SellBusinessID, inst.InsAmount})
	if err != nil {
		return shim.Error("instrumetcc: " + "Composite key InstrumentRefNo~SellBusinessID~InsAmount can not be created (instrument)")
	}
	value := []byte{0x00}
	stub.PutState(refNoSellIDkey, value)
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	var jsonresp string
	var err1 error
	fmt.Println("**********Invoke Instrument *********")
	if function == "enterInstrument" {
		//Used to enter new instrument data
		return enterInstrument(stub, args)
	} else if function == "getInstrument" {
		//used to retrieve the instrument data
		return getInstrument(stub, args)
	} else if function == "getInstrumentDetails" {
		//used to retrieve the instrument data
		return getInstrumentDetails(stub, args)
	} else if function == "updateInstrumentStatus" { //instrumentStatus
		//Updates instrument status accordingly
		return updateInstrumentStatus(stub, args)
	} else if function == "getInstrumentAmt" {
		return getInstrumentAmt(stub, args)
	} else if function == "queryInstrument" {
		return queryInstrument(stub, args)
	} else if function == "instrumentIDexists" {
		jsonresp,err1 = instrumentIDexists(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			return shim.Success([]byte(jsonresp))
		}
		//return instrumentIDexists(stub, args)
	} /** else if function == "getBusinessInfoQuery" {
		//Return the Response payload JSON
		//return getBusinessInfoQuery(stub, args)
		jsonresp,err1 = getBusinessInfoQuery(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			return shim.Success([]byte(jsonresp))
		}
	}*/
	
	
	return shim.Error("instrumetcc: " + "No function named " + function + " in Instrumentsssss")

}

func enterInstrument(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************While Writing Instrument *************************")
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("instrumetcc: " + "Invalid number of arguments in enterInstrument (required:10) given:" + xLenStr)
	}
	fmt.Print("args[0]", args[0])
	fmt.Print("args[1]", args[1])
	fmt.Print("args[2]", args[2])
	fmt.Print("args[3]", args[3])
	fmt.Print("args[4]", args[5])
	fmt.Print("args[5]", args[5])
	fmt.Print("args[6]", args[6])
	fmt.Print("args[7]" ,args[7])
	fmt.Print("args[8] ",args[8])
	fmt.Print("args[9] ",args[9])
	// Checking existence of Instrument Reference No. – Supplier ID pair
	indexName := "InstrumentRefNo~SellBusinessID~InsAmount"
	inst := instrumentInfo{}
	refNoSellIDiterator, _ := stub.GetStateByPartialCompositeKey("InstrumentRefNo~SellBusinessID~InsAmount", []string{args[0],args[2]})
	refNoSellIDdata, _ := refNoSellIDiterator.Next()
	if refNoSellIDdata != nil {
		return shim.Error("instrumetcc: " + "Instrument Reference No. – Supplier ID pair already exists")
	}
	defer refNoSellIDiterator.Close() 
	
	//create composite primary key for each new entry,check duplicate prior this
	refNoSellIDkey, err := stub.CreateCompositeKey(indexName, []string{inst.InstrumentRefNo, inst.SellBusinessID})
	if err != nil {
		return shim.Error("instrumetcc: " + "Composite key InstrumentRefNo~SellBusinessID~InsAmount can not be created (instrument)")
	}
	value := []byte{0x00}
	stub.PutState(refNoSellIDkey, value)
	
	//Checking existence of ProgramID
	chaincodeArgs := toChaincodeArgs("programIDexists", args[6])
	response := stub.InvokeChaincode("programcc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "ProgramId " + args[6] + " does not exits")
	}

	//Checking existence of pprID
	chaincodeArgs = toChaincodeArgs("pprIDexists", args[7])
	response = stub.InvokeChaincode("pprcc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "PprId " + args[7] + " does not exits")
	}

	//Checking existence of SellerBusinessID
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[2])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "BusinessId " + args[2] + " does not exits")
	}

	//Checking existence of BuyerBusinessID
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[3])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "BusinessId " + args[3] + " does not exits")
	}
	i, err := strconv.ParseFloat(args[4], 64)
	if err!= nil {
		return shim.Error("Error while converting String to Int ");
	}
	if math.Signbit(i) {
		return shim.Error("Invalid Instrument Amount")
	}
	//InstrumentDate -> instDate
	instDate, err := time.Parse("02/01/2006", args[1])
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}

	_, err = strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}

	//InsDueDate -> insDate
	insDueDate, err := time.Parse("02/01/2006", args[5])
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}
	if insDueDate.Weekday().String() == "Sunday" {
		fmt.Println("Since the due date falls on sunday, due date is extended to Monday(instrument) : ", insDueDate.AddDate(0, 0, 1))
	}
	insDueDate = insDueDate.AddDate(0, 0, 1)
	//Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing
	vString := args[9][:10] + "T" + args[9][11:] //removing the ":" part from the string

	//ValueDate -> vDate
	vDate, err := time.Parse("02/01/2006T15:04:05", vString)
	if err != nil {
		return shim.Error("instrumetcc: " + "error in parsing the date and time (instrument)" + err.Error())
	}

	// Hashing for key to store in ledger
	hash := sha256.New()
	instID := strings.ToLower(args[0] + args[2])
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md)
    //Doctype is passed as 
	inst1 := instrumentInfo{"InstruObject",args[0], instDate, args[2], args[3], args[4], "open", insDueDate, args[6], args[7], args[8], vDate}
	instBytes, err := json.Marshal(inst1)
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}
	stub.PutState(instIDsha, instBytes)
	fmt.Println("********************End  Writing Instrument  *************************")
	return shim.Success([]byte("Successfully added instrument to the ledger"))
}

func updateInstrumentStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start updateInstrumentStatus *************************")
	/*
		args[0] -> instrument reference number
		args[1] -> seller ID
		args[2] -> status
	*/
 key := strings.ToLower(args[0] + args[1])
	hash := sha256.New()
	instID := strings.ToLower(key)
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md) 

instBytes, err := stub.GetState(instIDsha)
		//instBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("instrumetcc: " + "Unable to fetch instrument info for status updation")
	}
	fmt.Println("instBytes ",instBytes)
	inst := instrumentInfo{}
	err = json.Unmarshal(instBytes, &inst)
	if err != nil {
		return shim.Error("instrumetcc: " + "Error in unmarshaling the instrument (updateInsStatus)")
	}
	/*
	 updated sequentially Open > Sanctioned > Overdue>Settled or Open > Sanctioned > Settled
	*/
	if (args[2] == "sanctioned") && (inst.InsStatus != "open") {
		return shim.Error("instrumetcc: " + "Instrument status cannot be sanctioned as it is not open")
	} else if (args[2] == "overdue") && (inst.InsStatus != "sanctioned") {
		return shim.Error("instrumetcc: " + "Instrument status cannot be overdue as it is not sanctioned")
	} else if (args[2] == "settled") && ((inst.InsStatus != "overdue") && (inst.InsStatus != "sanctioned")) {
		return shim.Error("instrumetcc: " + "Instrument status cannot be settled as it is not overdue or sanctioned")
	}
	inst.InsStatus = args[2]
	fmt.Println("inst b4 marshalling" ,inst)
	instBytes, _ = json.Marshal(inst)
	fmt.Println("inst after marshalling" ,instBytes)
	stub.PutState(instIDsha, instBytes)
	err = json.Unmarshal(instBytes, &inst)
	fmt.Println("After Updating ",inst)
	fmt.Println("******************** end updateInstrumentStatus *************************")
	return shim.Success([]byte("Instrument status updated successfully"))

}

func getInstrument(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Inside get instrument with arg" ,args [0], " ",args[1])
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("instrumetcc: " + "Invalid number of arguments in getInstrument (required:2) given:" + xLenStr)
	}
	/*
		args[0] -> InstrumentRefNo
		args[1] -> SellBusinessID
	*/
	hash := sha256.New()
	instID := strings.ToLower(args[0] + args[1])
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md)

	insBytes, err := stub.GetState(instIDsha)
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	} else if insBytes == nil {
		return shim.Error("instrumetcc: " + "No data exists on this InstrumentID: " + args[0])
	}

	ins := instrumentInfo{}
	err1 := json.Unmarshal(insBytes, &ins)
	if err1 != nil {
		return shim.Error("Error while unmarshalling in getInstrument()" +err1.Error())
	}
//	insString := fmt.Sprintf("%+v", ins)
	fmt.Println(" insBytes ",insBytes)
	
	fmt.Println(" ins ",ins)
	data, err1 := json.Marshal(ins)
	if err1 != nil {
		return shim.Error("after adding error while marshaling  instrument" + err1.Error() )
	} 	
	fmt.Println("End  getInstrument")
	return shim.Success(data)
}

func getInstrumentDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("******************** Inside getInstrumentDetails with arg" ,args[0])
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("instrumetcc: " + "Invalid number of arguments in getInstrument (required:1) given:" + xLenStr)
	}
	/*
		args[0] -> InstrumentRefNo
		args[1] -> SellBusinessID
	*/

	//bArrayMemberAlreadyWritten := false
	fmt.Println("args[0] ",args[0])                      
	insBytes, err := stub.GetStateByPartialCompositeKey("InstrumentRefNo~SellBusinessID~InsAmount", []string{args[0]})
/* 	insData, err := insBytes.Next()
	if err != nil {
		return shim.Error(" No Instrument Exist on this ID " + args[0])
	} */
	if err != nil {
		return shim.Error(" Error with GetStateByPartialCompositeKey " + err.Error())
	} 
	defer insBytes.Close()
	fmt.Println("insBytes ",insBytes)
	buffer, err1 := constructQueryResponseFromIterator(insBytes)
	if err1 != nil {
		return shim.Error("Error while contructing Query Response ")
	} 
	fmt.Println("buffer& ",&buffer)
	fmt.Println("buffer ",buffer)

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

func getInstrumentAmt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Start getInstrumentAmt ")
	if len(args) != 2 {
		return shim.Error("invalid no. of arguments (requiered 2)")
	}

	/*fmt.Println("before refNoSellIDiterator")
	refNoSellIDiterator, err := stub.GetStateByPartialCompositeKey("InstrumentRefNo~SellBusinessID~InsAmount", args)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("before refNoSellIDiterator")
	refNoSellIDdata, err := refNoSellIDiterator.Next()
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Println(refNoSellIDdata.Key)
	_, insAmt, err := stub.SplitCompositeKey(refNoSellIDdata.Key)
	if err != nil {
		return shim.Error("Error splitting the composite key " + err.Error())
	}
	fmt.Println(insAmt)
	*/
	hash := sha256.New()
	instID := strings.ToLower(args[0] + args[1])
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md)

	insBytes, err := stub.GetState(instIDsha)
	if err != nil {
		return shim.Error(err.Error())
	} else if insBytes == nil {
		return shim.Error("No data exists on this InstrumentID: " + instIDsha)
	}

	ins := instrumentInfo{}
	err = json.Unmarshal(insBytes, &ins)
	if err != nil {
		return shim.Error("Cannot unmarshal json")
	}
fmt.Println("End getInstrumentAmt ")
	return shim.Success([]byte(ins.InsAmount))
}

/** func instrumentIDexists(stub shim.ChaincodeStubInterface, instrumentId string) pb.Response {
	insBytes, err  := stub.GetState(instrumentId)

	if err != nil {
		return shim.Error(err.Error())
	} else if insBytes == nil {
		return shim.Error("No data exists on this InstrumentID: " + instrumentId)
	}
	fmt.Println("instrumentId" +instrumentId)
	instrAsJson := instrumentInfo{}
	err = json.Unmarshal(insBytes, &instrAsJson)
	if err != nil {
		return shim.Error("Cannot unmarshal json")
	}
	return shim.Success([]byte(instrAsJson.InsAmount))
	//return shim.Success([]byte(instrAsJson.InsStatus))
}*/

func instrumentIDexists(stub shim.ChaincodeStubInterface, instrumentId []string) (string , error){
	fmt.Println("******** Start instrumentIDexists with arg ",instrumentId[0])
	value, err := stub.GetState(instrumentId[0])
    if err != nil {
            return "", fmt.Errorf("Failed to get asset: %s with error: %s", instrumentId[0], err)
    }
    if value == nil {
            return "", fmt.Errorf("Asset not found: %s", instrumentId[0])
	}
	fmt.Println(" ******** End instrumentIDexists ******** ")
    return string(value), nil
}
func queryInstrument(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	instrumentvalasstruct := []instrumentInfoReturn{}
	//var instrumentvalasstruct instrumentInfo;
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
//	queryString := args;
	fmt.Println(" Hello .. Inside queryInstrument1 ")
	fmt.Println("args " ,html.UnescapeString(args[0]))
	//fmt.Println("args ",args);
	queryResults, err := getQueryResultForQueryString(stub, html.UnescapeString(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	} 	
		fmt.Println("queryInstrument ",string(queryResults));
	///	businessValues,err4 = getBusinessInfoVal(stub, args)
	//var instrumentvalasstruct1;
	  err2 := json.Unmarshal(queryResults, &instrumentvalasstruct)
	if err2 != nil {
		return shim.Error("Error while unmarshalling instrumentvalasstruct "+ err2.Error())
	}   
	fmt.Println(" instrumentvalasstruct  RefNo",instrumentvalasstruct)
//	fmt.Println(" Key  ",instrumentvalasstruct.Key)
	data, err1 := json.Marshal(instrumentvalasstruct)
	if err1 != nil {
		return shim.Error("error while marshaling  businessvalasstruct" + err1.Error() )
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
		fmt.Printf("instrumetcc: "+"Error starting Instrument chaincode: %s\n", err)
	}
}
