package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Return PPOB Inquiry response in JSON
func responseJsonPPOBInquiry(jsonIso PPOBInquiryRequest) PPOBInquiryResponse {
	var response PPOBInquiryResponse

	// Client setup for custom http request
	client := &http.Client{}
	var baseURL = "https://chipsakti-mock.herokuapp.com"

	// Set data to be encoded
	var param = url.Values{}
	param.Set("transaction_id", jsonIso.TransactionID)
	param.Set("partner_id", jsonIso.PartnerID)
	param.Set("product_code", jsonIso.ProductCode)
	param.Set("customer_no", jsonIso.CustomerNo)
	param.Set("periode", jsonIso.Periode)
	param.Set("merchant_code", jsonIso.MerchantCode)
	param.Set("request_time", jsonIso.RequestTime)
	param.Set("signature", jsonIso.Signature)

	log.Printf("Send request to https://chipsakti-mock.herokuapp.com/inquiry\n")

	// Request to Biller
	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest("POST", baseURL+"/inquiry", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatalf("Failed to sent request to https://chipsakti-mock.herokuapp.com/inquiry. Error: %v\n", err)
	}

	// Check response from Biller
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response from https://chipsakti-mock.herokuapp.com/inquiry. Error: %v\n", err)
	}

	defer resp.Body.Close()

	log.Printf("Receive response from https://chipsakti-mock.herokuapp.com/inquiry\n")

	// Read response from Biller
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	return response
}

// Return PPOB Payment response in JSON
func responsePPOBPayment(jsonIso PPOBPaymentRequest) PPOBPaymentResponse {
	var response PPOBPaymentResponse

	// Client setup for custom http request
	client := &http.Client{}
	var baseURL = "https://chipsakti-mock.herokuapp.com"
	amount := strconv.Itoa(jsonIso.Amount)

	// Set data to be encoded
	var param = url.Values{}
	param.Set("transaction_id", jsonIso.TransactionID)
	param.Set("partner_id", jsonIso.PartnerID)
	param.Set("product_code", jsonIso.ProductCode)
	param.Set("customer_no", jsonIso.CustomerNo)
	param.Set("reff_id", jsonIso.ReffID)
	param.Set("amount", amount)
	param.Set("merchant_code", jsonIso.MerchantCode)
	param.Set("request_time", jsonIso.RequestTime)
	param.Set("signature", jsonIso.Signature)

	log.Printf("Send request to https://chipsakti-mock.herokuapp.com/payment\n")

	// Request to Biller
	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest("POST", baseURL+"/payment", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatalf("Failed to sent request to https://chipsakti-mock.herokuapp.com/payment. Error: %v\n", err)
	}

	// Check response from Biller
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response from https://chipsakti-mock.herokuapp.com/payment. Error: %v\n", err)
	}

	defer resp.Body.Close()

	log.Printf("Receive response from https://chipsakti-mock.herokuapp.com/payment\n")

	// Read response from Biller
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	return response
}

// Return PPOB Status response in JSON
func responsePPOBStatus(jsonIso PPOBStatusRequest) PPOBStatusResponse {
	var response PPOBStatusResponse

	// Client setup for custom http request
	client := &http.Client{}
	var baseURL = "https://chipsakti-mock.herokuapp.com"
	amount := strconv.Itoa(jsonIso.Amount)

	// Set data to be encoded
	var param = url.Values{}
	param.Set("transaction_id", jsonIso.TransactionID)
	param.Set("partner_id", jsonIso.PartnerID)
	param.Set("product_code", jsonIso.ProductCode)
	param.Set("customer_no", jsonIso.CustomerNo)
	param.Set("reff_id", jsonIso.ReffID)
	param.Set("amount", amount)
	param.Set("merchant_code", jsonIso.MerchantCode)
	param.Set("request_time", jsonIso.RequestTime)
	param.Set("signature", jsonIso.Signature)

	log.Printf("Send request to https://chipsakti-mock.herokuapp.com/status\n")

	// Request to Biller
	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest("POST", baseURL+"/status", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatalf("Failed to sent request to https://chipsakti-mock.herokuapp.com/status. Error: %v\n", err)
	}

	// Check response from Biller
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response from https://chipsakti-mock.herokuapp.com/status. Error: %v\n", err)
	}

	defer resp.Body.Close()

	log.Printf("Receive response from https://chipsakti-mock.herokuapp.com/status\n")

	// Read response from Biller
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	return response
}

// Return Topup Buy response in JSON
func responseTopupBuy(jsonIso TopupBuyRequest) TopupBuyResponse {
	var response TopupBuyResponse

	// Client setup for custom http request
	client := &http.Client{}
	var baseURL = "https://chipsakti-mock.herokuapp.com"

	// Set data to be encoded
	var param = url.Values{}
	param.Set("transaction_id", jsonIso.TransactionID)
	param.Set("partner_id", jsonIso.PartnerID)
	param.Set("product_code", jsonIso.ProductCode)
	param.Set("customer_no", jsonIso.CustomerNo)
	param.Set("merchant_code", jsonIso.MerchantCode)
	param.Set("request_time", jsonIso.RequestTime)
	param.Set("signature", jsonIso.Signature)

	log.Printf("Send request to https://chipsakti-mock.herokuapp.com/buy\n")

	// Request to Biller
	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest("POST", baseURL+"/buy", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatalf("Failed to sent request to https://chipsakti-mock.herokuapp.com/buy. Error: %v\n", err)
	}

	// Check response from Biller
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response from https://chipsakti-mock.herokuapp.com/buy. Error: %v\n", err)
	}

	defer resp.Body.Close()

	log.Printf("Receive response from https://chipsakti-mock.herokuapp.com/buy\n")

	// Read response from Biller
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	return response
}

// Return Topup Check response in JSON
func responseTopupCheck(jsonIso TopupCheckRequest) TopupCheckResponse {
	var response TopupCheckResponse

	// Client setup for custom http request
	client := &http.Client{}
	var baseURL = "https://chipsakti-mock.herokuapp.com"

	// Set data to be encoded
	var param = url.Values{}
	param.Set("transaction_id", jsonIso.TransactionID)
	param.Set("partner_id", jsonIso.PartnerID)
	param.Set("product_code", jsonIso.ProductCode)
	param.Set("customer_no", jsonIso.CustomerNo)
	param.Set("merchant_code", jsonIso.MerchantCode)
	param.Set("request_time", jsonIso.RequestTime)
	param.Set("signature", jsonIso.Signature)

	log.Printf("Request to https://chipsakti-mock.herokuapp.com/check\n")

	// Request to Biller
	var payload = bytes.NewBufferString(param.Encode())
	req, err := http.NewRequest("POST", baseURL+"/check", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatalf("Failed to sent request to https://chipsakti-mock.herokuapp.com/check. Error: %v\n", err)
	}

	// Check response from Biller
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response from https://chipsakti-mock.herokuapp.com/check. Error: %v\n", err)
	}

	defer resp.Body.Close()

	log.Printf("Receive response from https://chipsakti-mock.herokuapp.com/check\n")

	// Read response from Biller
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	return response
}
