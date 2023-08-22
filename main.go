package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/j03hanafi/bankiso/iso20022/pacs"
	"github.com/j03hanafi/bankiso/iso20022/prxy"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Found error in log ", err)
	}
	log.SetOutput(file)
	log.Println("Log setup")

	path := pathHandler()

	address := ":5000"
	log.Printf("Biller started at %v", address)
	err = http.ListenAndServe(address, path)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func pathHandler() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", biller).Methods("POST")

	return router
}

func biller(w http.ResponseWriter, r *http.Request) {
	log.Println("New Request from BIFast Connector")
	fmt.Println("New Request from BIFast Connector")

	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))

	request := ChannelInput{}
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Error unmarshal JSON: %s", err.Error())
	}
	log.Println("request: ", string(request.BusMsg.Document))
	var response interface{}

	// var msgID string
	var fileName string
	bzMsgID := fmt.Sprintf("%v", *request.BusMsg.AppHdr.BusinessMessageIdentifier)
	DocumentValue := request.BusMsg.Document
	fmt.Println(string(DocumentValue))
	trxType := bzMsgID[16:19]
	fmt.Println("trxType:", trxType)

	switch trxType {
	// ##################### Account Enquiry ##################################
	case "510":
		document := pacs.Document00800108{}
		err := json.Unmarshal(DocumentValue, &document)
		if err != nil {
			fmt.Println("Error unmarshal: ", err)
		}
		CrAccId := *document.Message.CreditTransferTransactionInformation[0].CdtrAcct.Id.Other.Identification
		fmt.Println(CrAccId)
		switch CrAccId {
		case "510654300": //AE U000 CT U170
			fileName = "sampleAccountEnquiry.json"
		case "510654301": //AE U000 CT U000
			// timeoutVal := rand.Intn(5)
			// time.Sleep(time.Duration(timeoutVal) * time.Second)
			fileName = "sampleAccountEnquiry3.json"
		case "510654302": //AE U000 CT PEND
			fileName = "sampleAccountEnquiry4.json"
		case "510654303":
			fileName = "sampleAccountEnquiry5.json"
		case "510654305":
			fileName = "sampleAccountEnquiry6.json"
		case "510654306":
			fileName = "sampleAccountEnquiry7.json"
		case "510654307":
			fileName = "sampleAccountEnquiry8.json"
		case "510654308":
			fileName = "sampleAccountEnquiry9.json"
		case "510654309":
			fileName = "sampleAccountEnquiry10.json"
		case "510654310":
			fileName = "sampleAccountEnquiry11.json"
		case "511654182":
			fileName = "sampleAccountEnquiry2.json"
		case "511654999":
			fileName = "sampleAccountEnquiry5.json"
		case "0000000000":
			fileName = "rejectMessage.json"
		case "511654129":
			fileName = "sampleAccountEnquiry4.json"
		case "020600022014769":
			fileName = "sampleAccountEnquiryMega.json"
		}

	//##################### Credit Transfer ###################################
	case "010": // Credit Transfer
		document := pacs.Document00800108{}
		err := json.Unmarshal(DocumentValue, &document)
		if err != nil {
			fmt.Println("Error unmarshal: ", err)
		}
		CrAccId := *document.Message.CreditTransferTransactionInformation[0].CdtrAcct.Id.Other.Identification
		trigger := *document.Message.CreditTransferTransactionInformation[0].PmtId.EndToEndId
		switch CrAccId {
		case "510654300":
			fileName = "sampleCreditTransferResponse.json"
		case "510654301":
			fileName = "sampleCreditTransferResponse3.json"
		case "510654302":
			time.Sleep(100 * time.Second)
			fileName = "sampleCreditTransferResponse4.json"
		case "510654303":
			time.Sleep(5 * time.Second)
			log.Panic("Panic Triggered!")
			fileName = "sampleCreditTransferResponse5.json"
		case "510654305":
			fileName = "sampleCreditTransferResponse6.json"
		case "510654306":
			fileName = "sampleCreditTransferResponse7.json"
		case "510654307":
			fileName = "sampleCreditTransferResponse8.json"
		case "510654308":
			fileName = "sampleCreditTransferResponse9.json"
		case "510654309":
			fileName = "sampleCreditTransferResponse10.json"
		case "0102345600":
			fileName = "sampleCreditTransferResponse.json"
		case "0102345184":
			fileName = "sampleCreditTransferResponse2.json"
		case "0102345129":
			fileName = "sampleCreditTransferResponse4.json"
		case "0102345999":
			fileName = "sampleCreditTransferResponse5.json"
		case "0000000000":
			fileName = "rejectMessage.json"
		case "0000000001":
			switch trigger {
			case "20210301INDOIDJA000ORB00000000":
				time.Sleep(80 * time.Second)
				fileName = "sampleCreditTransferResponse.json"
			case "20210301INDOIDJA000ORB11111111":
				time.Sleep(80 * time.Second)
				fileName = "sampleCreditTransferResponse.json"
			case "20210301INDOIDJA000ORB22222222":
				time.Sleep(80 * time.Second)
				fileName = "sampleCreditTransferResponse.json"
			}
		}
	case "012":
		fileName = "sampleCreditTransferResponse012.json"
		fmt.Println("012")
	case "110":
		fileName = "sampleCreditTransferResponsewithProxy.json"
		fmt.Println("110")
	//==========================================================================

	case "019":
		document := pacs.Document00900109{}
		err := json.Unmarshal(DocumentValue, &document)
		if err != nil {
			fmt.Println("Error Unmarshal: ", err)
		}
		trigger := *document.Message.CreditTransferTransactionInformation[0].PmtId.EndToEndId
		switch trigger {
		case "20210301INDOIDJA000ORB99999999":
			time.Sleep(80 * time.Second)
			fileName = "sampleFItoFICreditTransfer.json"
		default:
			fileName = "sampleFItoFICreditTransfer.json"
			fmt.Println("019")
		}

	// case "119":
	// 	time.Sleep(80 * time.Second)
	// 	fileName = "sampleFItoFICreditTransfer.json"
	// 	fmt.Println("019")
	case "011":
		fileName = "sampleReverseCreditTransfer.json"
		fmt.Println("011")

	// ################# Proxy Resolution #####################################
	case "610":
		document := prxy.Document00300101{}
		err := json.Unmarshal(DocumentValue, &document)
		if err != nil {
			fmt.Println("Error unmarshal: ", err)
		}
		PxValue := *document.Message.LookUp.PrxyOnly.PrxyRtrvl.Val
		switch PxValue {
		case "086102345000":
			fileName = "sampleProxyResolution.json"
		case "086112345101":
			fileName = "sampleProxyResolution2.json"
		case "086112345804":
			fileName = "sampleProxyResolution3.json"
		case "086132345600":
			fileName = "sampleProxyResolution4.json"
		case "086142345804":
			fileName = "sampleProxyResolution5.json"
		case "08615234804":
			fileName = "sampleProxyResolution6.json"
		case "08616234811":
			fileName = "sampleProxyResolution7.json"
		case "08617234805":
			fileName = "sampleProxyResolution8.json"
		case "08617234129":
			fileName = "sampleProxyResolution9.json"
		case "08617234999":
			fileName = "sampleProxyResolution10.json"
		case "0000000000":
			fileName = "rejectMessage.json"
		}
	case "611":
		fileName = "sampleProxyResolution611.json"
		fmt.Println("611")
	case "612":
		fileName = "sampleProxyResolution612.json"
		fmt.Println(("612"))
	// =========================================================================

	// ################# Proxy Registration Inquiry ############################
	case "620":
		document := prxy.Document00500101{}
		err := json.Unmarshal(DocumentValue, &document)
		if err != nil {
			fmt.Println("Error unmarshal: ", err)
		}
		PxRegId := *document.Message.Nqry.RegnId
		fmt.Println(PxRegId)
		switch PxRegId {
		case "6202345600":
			fileName = "sampleProxyRegistrationInquiry.json"
		case "6212345101":
			fileName = "sampleProxyRegistrationInquiry2.json"
		case "6222345808":
			fileName = "sampleProxyRegistrationInquiry3.json"
		case "6232345600":
			fileName = "sampleProxyRegistrationInquiry4.json"
		case "6242345600":
			fileName = "sampleProxyRegistrationInquiry5.json"
		case "6252345808":
			fileName = "sampleProxyRegistrationInquiry6.json"
		case "6262345806":
			fileName = "sampleProxyRegistrationInquiry7.json"
		case "6262345129":
			fileName = "sampleProxyRegistrationInquiry8.json"
		case "6262345999":
			fileName = "sampleProxyRegistrationInquiry9.json"
		case "0000000000":
			fileName = "rejectMessage.json"
		}
	case "621":
		fileName = "sampleProxyRegistrationInquiry621.json"
		fmt.Println("621")
	case "622":
		fileName = "sampleProxyRegistrationInquiry622.json"
		fmt.Println("622")
	// =========================================================================

	case "710":
		fileName = "sampleRegisterNewProxy.json"
		fmt.Println("710")

	//#################### Proxy Maintenance ###################################
	case "720":
		document := prxy.Document00100101{}
		err := json.Unmarshal(DocumentValue, &document)
		if err != nil {
			fmt.Println("Error unmarshal: ", err)
		}
		PxValue := *document.Message.Regn.Prxy.Val
		switch PxValue {
		case "7202345600":
			fileName = "sampleProxyMaintenance.json"
		case "7212345101":
			fileName = "sampleProxyMaintenance2.json"
		case "7212345129":
			fileName = "sampleProxyMaintenance3.json"
		case "7212345999":
			fileName = "sampleProxyMaintenance4.json"
		case "0000000000":
			fileName = "rejectMessage.json"
		}
	case "721":
		fileName = "sampleProxyMaintenance721.json"
		fmt.Println("721")
		//============================================================================

	case "000":
		document := PSRDocument{}
		err := json.Unmarshal(DocumentValue, &document)
		if err != nil {
			fmt.Println("error unmarshal: ", err)
		}
		trigger := document.FItoFIPmtStsReq.TxInf[0].OrgnlEndToEndID
		// switch trigger {
		// case "20210301INDOIDJA000ORB00000000":
		// 	fileName = "PaymentStatusReqResponse.json"
		// case "20210301INDOIDJA000ORB11111111":
		// 	fileName = "PaymentStatusReqResponse2.json"
		// case "20210301INDOIDJA000ORB22222222":
		// 	fileName = "PaymentStatusReqResponse3.json"
		// case "20210301INDOIDJA000ORB99999999":
		// 	fileName = "PaymentStatusReqResponse.json"
		switch {
		case strings.Contains(trigger, "ROYBIDJ1010O02"):
			fileName = "PaymentStatusReqResponse4.json"
		case strings.Contains(trigger, "ROYBIDJ1010O01"):
			fileName = "PaymentStatusReqResponse5.json"
		case strings.Contains(trigger, "LFIBIDJ1010"):
			fileName = "PaymentStatusReqResponse.json"
		}

	// case "001":
	// 	fileName = "PaymentStatusReqResponse2.json"
	// 	fmt.Println("001")
	// case "002":
	// 	fileName = "PaymentStatusReqResponse3.json"
	// 	fmt.Println("001")

	default:
		fileName = "rejectMessage.json"
		fmt.Println("Default")
	}

	//fmt.Println("Enter file name: ")

	//fmt.Scanln(&fileName)
	//
	fileName = "samples/" + fileName
	fmt.Println("filename:", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(b, &response)
	// fmt.Println("response:", response)

	responseFormatter(w, response, 200)
}

func responseFormatter(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
