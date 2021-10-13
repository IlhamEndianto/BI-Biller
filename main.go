package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/j03hanafi/bankiso/iso20022/head"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Found error in log ", err)
	}
	log.SetOutput(file)
	log.Println("Log setup")

	path := pathHandler()

	address := ":6066"
	log.Printf("Biller started at %v", address)
	err = http.ListenAndServe(address, path)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func pathHandler() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/biller", biller).Methods("POST")

	return router
}

func biller(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Simulator return response based on msgID
	// Mapping JSON config untuk request-response (ex: req->pacs008, res->pacs002). Referensi: BI-FAST Participant Guide - Format Interface v1.0
	// Lama waktu close connection BI-Biller

	log.Println("New Request from BIFast Connector")
	fmt.Println("New Request from BIFast Connector")

	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))

	request := ChannelInput{}
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Error unmarshal JSON: %s", err.Error())
	}
	log.Println("request: ", request)
	var response interface{}

	// var msgID string
	var fileName string
	bzMsgID := fmt.Sprintf("%v", *request.BusMsg.AppHdr.BusinessMessageIdentifier)
	// msgID = fmt.Sprintf("%v", *request.BusMsg.AppHdr.MessageDefinitionIdentifier)
	// log.Println("MsgDefIdn:", msgID)
	businessCode := bzMsgID[16:19]
	fmt.Println(businessCode)

	switch businessCode {
	case "010": // Credit Transfer
		fileName = "sampleCreditTransferResponse.json"
		fmt.Println("010")
	case "019":
		fileName = "sampleFItoFICreditTransfer.json"
		fmt.Println("019")
	case "011":
		fileName = "sampleReverseCreditTransfer.json"
		fmt.Println("011")
	case "110":
		fileName = "sampleCreditTransferResponsewithProxy.json"
		fmt.Println("110")
	case "510":
		fileName = "sampleAccountEnquiry.json"
		fmt.Println("510")
	case "610":
		fileName = "sampleProxyResolution.json"
		fmt.Println("610")
	case "620":
		fileName = "sampleProxyRegistrationInquiry.json"
		fmt.Println("620")
	case "710":
		fileName = "sampleRegisterNewProxy.json"
		fmt.Println("710")
	case "720":
		fileName = "sampleProxyMaintenance.json"
		fmt.Println("720")
	}

	//fmt.Println("Enter file name: ")

	//fmt.Scanln(&fileName)
	//
	fileName = "samples/" + fileName
	fmt.Println("filename:", fileName)

	file, _ := os.Open(fileName)
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(b, &response)
	fmt.Println("response:", response)

	responseFormatter(w, response, 200)
}

func responseFormatter(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// type tempDocumentXML interface{}
type ChannelInput struct {
	BusMsg BusMsg `xml:"BusMsg" json:"BusMsg"`
}

type BusMsg struct {
	AppHdr   *head.BusinessApplicationHeaderV01 `xml:"AppHdr" json:"AppHdr"`
	Document json.RawMessage                    `xml:"Document" json:"Document"`
}

func getFileName(BzID string) (fileName string) {
	switch BzID {
	case "20210301INDOIDJA610ORB11111111":
		fileName = "sample_prxy.002_dummyDeregist.json"
	case "20210301INDOIDJA610ORB22222222":
		fileName = "sample_prxy.002_dummyUpdate.json"
	case "20210301INDOIDJA610ORB33333333":
		fileName = "sample_prxy.002_dummySuspend.json"
	case "20210301INDOIDJA610ORB44444444":
		fileName = "sample_prxy.002_dummyActivation.json"
	case "20210301INDOIDJA610ORB55555555":
		fileName = "sample_prxy.002_dummyPorting.json"
	default:
		fileName = "sample_prxy.002_response_alias_mgmt_NEWR_CIHUB_to_OFI.xml.json"
	}
	return fileName
}
