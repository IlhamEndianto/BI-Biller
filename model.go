package main

import (
	"encoding/json"

	"github.com/j03hanafi/bankiso/iso20022/head"
)

type ChannelInput struct {
	BusMsg BusMsg `xml:"BusMsg" json:"BusMsg"`
}

type BusMsg struct {
	AppHdr   *head.BusinessApplicationHeaderV01 `xml:"AppHdr" json:"AppHdr"`
	Document json.RawMessage                    `xml:"Document" json:"Document"`
}

type PSRDocument struct {
	FItoFIPmtStsReq struct {
		GrpHdr struct {
			MsgID   string `json:"MsgId"`
			CreDtTm string `json:"CreDtTm"`
		} `json:"GrpHdr"`
		TxInf []struct {
			OrgnlEndToEndID string `json:"OrgnlEndToEndId"`
		} `json:"TxInf"`
	} `json:"FItoFIPmtStsReq"`
}
