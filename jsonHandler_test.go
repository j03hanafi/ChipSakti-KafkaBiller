package main

import (
	"github.com/mofax/iso8583"
	"testing"
)

func TestGetJsonPPOBInquiry(t *testing.T) {
	var result PPOBInquiryRequest
	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)
	parsedIso := "0200a00000000001000000000000000000003800011302021                     USER01          WOM             2                        KIOS01                   2018-05-15 15:10:052020"
	iso, err := isoStruct.Parse(parsedIso)
	if err != nil {
		t.Errorf("Error parsing iso message. Error: %v", err)
	}
	result = getJsonPPOBInquiry(iso)
	var expected PPOBInquiryRequest

	expected.TransactionID = "2021"
	expected.PartnerID = "USER01"
	expected.ProductCode = "WOM"
	expected.CustomerNo = "2"
	expected.MerchantCode = "KIOS01"
	expected.RequestTime = "2018-05-15 15:10:05"
	expected.Periode = "2020"
	expected.Signature = "41146b35700a4f6c05bd9dfd9da52b63f308050baa61d395e4f78a7d1625d70a"

	if result != expected {
		t.Errorf("getJsonPPOBInquiry() failed. \nExpected\t: %v. Got\t: %v", expected, result)
	} else {
		t.Log("getJsonPPOBInquiry() success")
	}
}

func TestGetsonPPOBPayment(t *testing.T) {
	var result PPOBPaymentRequest
	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)
	parsedIso := "0200b000000008010000000000000000000081000100000087330012345       1262015                     USER01          WOM             2                        KIOS01                   2018-05-15 15:10:05"
	iso, err := isoStruct.Parse(parsedIso)
	if err != nil {
		t.Errorf("Error parsing iso message. Error: %v", err)
	}
	result = getJsonPPOBPayment(iso)
	var expected PPOBPaymentRequest

	expected.TransactionID = "2015"
	expected.PartnerID = "USER01"
	expected.ProductCode = "WOM"
	expected.CustomerNo = "2"
	expected.MerchantCode = "KIOS01"
	expected.RequestTime = "2018-05-15 15:10:05"
	expected.ReffID = "12345"
	expected.Amount = 873300
	expected.Signature = "2ba6d0fba94ea4d189af60ed83f286ada47f3d442d232458ab6ea1ff76ef93fb"

	if result != expected {
		t.Errorf("getJsonPPOBPayment() failed. \nExpected\t: %v. Got\t: %v", expected, result)
	} else {
		t.Log("getJsonPPOBPayment() success")
	}
}

func TestGetJsonPPOBStatus(t *testing.T) {
	var result PPOBStatusRequest
	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)
	parsedIso := "0200b000000008010000000000000000000038000200000001000012345       1262021                     USER01          WOM             2                        KIOS01                   2018-05-15 15:10:05"
	iso, err := isoStruct.Parse(parsedIso)
	if err != nil {
		t.Errorf("Error parsing iso message. Error: %v", err)
	}
	result = getJsonPPOBStatus(iso)
	var expected PPOBStatusRequest

	expected.TransactionID = "2021"
	expected.PartnerID = "USER01"
	expected.ProductCode = "WOM"
	expected.CustomerNo = "2"
	expected.MerchantCode = "KIOS01"
	expected.RequestTime = "2018-05-15 15:10:05"
	expected.ReffID = "12345"
	expected.Amount = 10000
	expected.Signature = "8aa59cf131fa5022ed2c2cec39faa7f3a2a4990159d0bfb004d3af97b9265254"

	if result != expected {
		t.Errorf("getJsonPPOBStatus() failed. \nExpected\t: %v. Got\t: %v", expected, result)
	} else {
		t.Log("getJsonPPOBStatus() success")
	}
}

func TestGetJsonTopupBuy(t *testing.T) {
	var result TopupBuyRequest
	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)
	parsedIso := "0200a00000000001000000000000000000008100021262021                     USER01          WOM             1                        KIOS01                   2018-05-15 15:10:05"
	iso, err := isoStruct.Parse(parsedIso)
	if err != nil {
		t.Errorf("Error parsing iso message. Error: %v", err)
	}
	result = getJsonTopupBuy(iso)
	var expected TopupBuyRequest

	expected.TransactionID = "2021"
	expected.PartnerID = "USER01"
	expected.ProductCode = "WOM"
	expected.CustomerNo = "1"
	expected.MerchantCode = "KIOS01"
	expected.RequestTime = "2018-05-15 15:10:05"
	expected.Signature = "44b514050e89d90ab8e43ca9716670b04d5f65b7ef963112bd87e7a866c6c0f7"

	if result != expected {
		t.Errorf("getJsonTopupBuy() failed. \nExpected\t: %v. Got\t: %v", expected, result)
	} else {
		t.Log("getJsonTopupBuy() success")
	}
}

func TestGetJsonTopupCheck(t *testing.T) {
	var result TopupCheckRequest
	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)
	parsedIso := "0200a00000000001000000000000000000008100021262021                     USER01          WOM             1                        KIOS01                   2018-05-15 15:10:05"
	iso, err := isoStruct.Parse(parsedIso)
	if err != nil {
		t.Errorf("Error parsing iso message. Error: %v", err)
	}
	result = getJsonTopupCheck(iso)
	var expected TopupCheckRequest

	expected.TransactionID = "2021"
	expected.PartnerID = "USER01"
	expected.ProductCode = "WOM"
	expected.CustomerNo = "1"
	expected.MerchantCode = "KIOS01"
	expected.RequestTime = "2018-05-15 15:10:05"
	expected.Signature = "3912621a2a919f95a55986310c7ebc4625ed31fd6d1248fb1716e66dac0979fe"

	if result != expected {
		t.Errorf("getJsonTopupCheck() failed. \nExpected\t: %v. Got\t: %v", expected, result)
	} else {
		t.Log("getJsonTopupCheck() success")
	}
}
