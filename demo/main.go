package main

import (
	"log"

	"github.com/twigaeng/mpesa-api-go"
)

const (
	appKey    = "GvzjNnYgNJtwgwfLBkZh65VPwfuKvs0V" // sandbox --> change to yours
	appSecret = "oOpJICRVlyrGSAkM"                 // sandbox --> change to yours
)

func main() {
	// These examples are taken from the mpesa-java-sdk examples
	// at https://github.com/safaricom/mpesa-java-sdk

	m, err := mpesa.New(appKey, appSecret, mpesa.PRODUCTION)
	if err != nil {
		panic(err)
	}

	response, err := m.B2BRequest(mpesa.B2B{})
	if err != nil {
		log.Println(err)
	}
	log.Println(response)

	c2b := mpesa.C2B{
		ShortCode:     "600576",
		CommandID:     "CustomerPayBillOnline",
		Amount:        "2",
		Msisdn:        "254708374149",
		BillRefNumber: "hkjhjkhjkh"}
	response, err = m.C2BSimulation(c2b)
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
}
