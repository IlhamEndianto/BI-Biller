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

	address := "localhost:6066"
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

	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))

	request := ChannelInput{}
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Error unmarshal JSON: %s", err.Error())
	}
	log.Println("request: ", request)
	var response interface{}

	var msgID string
	var fileName string
	var BzID string
	msgID = fmt.Sprintf("%v", request.BusMsg.AppHdr.MessageDefinitionIdentifier)
	BzID = fmt.Sprintf("%v", request.BusMsg.AppHdr.BusinessMessageIdentifier)
	log.Println("MsgDefIdn:", msgID)

	switch msgID {
	case "pacs.008.001.08":
		switch BzID {
		case "20210301INDOIDJA010ORB12356789":
			fileName = "sample_AccountEnquiry_AccountID_pacs.002_response_CIHUB_to_OFI.xml.json"
		case "20210301INDOIDJA010ORB12357890":
			fileName = "sample_AccountEnquiry_AccountID_pacs.002_response_CIHUB_to_OFI.xml.json"
		default:
			fileName = "sample_AccountEnquiry_AccountID_pacs.002_response_CIHUB_to_OFI.xml.json"
		}
		fileName = "sample_CreditTransfer_AccountID_pacs.002_response_CIHUB_to_OFI.xml.json"
	case "pacs.009.001.09":
		fileName = "sample_FItoFICreditTransfer_pacs.002_response_CIHUB_to_OFI.xml.json"
	case "pacs.028.001.04":
		fileName = "PaymentStatusReqResponse.json"
	case "prxy.001.001.01":
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
	case "prxy.003.001.01":
		fileName = "sample_prxy.004_response_alias_resolution_CIHUB_to_OFI.xml.json"
	case "prxy.005.001.01":
		fileName = "sample_prxy.006_response_alias_enquiry_CIHUB_to_OFI.xml.json"
	default:
		fileName = "sample_AccountEnquiry_AccountID_pacs.002_response_CIHUB_to_OFI.xml.json"
	}

	//fmt.Println("Enter file name: ")

	//fmt.Scanln(&fileName)
	//
	fileName = "samples/" + fileName
	log.Println("filename:", fileName)

	file, _ := os.Open(fileName)
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(b, &response)
	log.Println("response:", response)

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
