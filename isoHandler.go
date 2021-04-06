package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mofax/iso8583"
	"github.com/rivo/uniseg"
)

// Handler to new consumed request in consumerChan and send new response to billerChan
func requestHandler() {

	// loop for checking if there is any new request from Consumer (Kafka) that has been sent to consumerChan
	for {
		select {
		// execute if there is a new request in consumerChan
		case newRequest := <-consumerChan:

			start := time.Now()
			// Send new request to `Biller` and get response that ready to produce
			msg := newRequest
			log.Printf("[Time: %v. Elapsed: %.6fs] Received new Request\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())
			isoParsed := getResponse(msg, start)

			// Send new response to billerChan
			billerChan <- isoParsed

			// Done with requestHandler
			log.Printf("[Time: %v. Elapsed: %.6fs] Send response to consumer\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())
			log.Println("New request handled")

		// keep looping if there is none new request
		default:
			continue
		}
	}

}

// Return response from `Biller` in ISO8583 Format
func getResponse(message string, start time.Time) (isoResponse string) {

	var response Iso8583
	data := message[4:]

	// Parse new ISO8583 message to ISO Struct
	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)
	msg, err := isoStruct.Parse(data)
	if err != nil {
		log.Println(err)
	}

	var isoParsed iso8583.IsoStruct

	// Check processing code and send request to appropriate `Biller` endpoints
	pcode := msg.Elements.GetElements()[3]
	switch pcode {
	// Process PPOB Inquiry request
	case "380001":
		// Convert ISO message to JSON format
		jsonIso := getJsonPPOBInquiry(msg)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert ISO message to JSON format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Send JSON data to Biller
		serverResp := responseJsonPPOBInquiry(jsonIso)
		log.Printf("[Time: %v. Elapsed: %.6fs] Send JSON data to Biller\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Convert response from JSON data to ISO8583 format
		isoParsed = getIsoPPOBInquiry(serverResp)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert response from JSON data to ISO8583 format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

	// Process PPOB Payment request
	case "810001":
		// Convert ISO message to JSON format
		jsonIso := getJsonPPOBPayment(msg)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert ISO message to JSON format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Send JSON data to Biller
		serverResp := responsePPOBPayment(jsonIso)
		log.Printf("[Time: %v. Elapsed: %.6fs] Send JSON data to Biller\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Convert response from JSON data to ISO8583 format
		isoParsed = getIsoPPOBPayment(serverResp)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert response from JSON data to ISO8583 format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

	// Process PPOB Status request
	case "380002":
		// Convert ISO message to JSON format
		jsonIso := getJsonPPOBStatus(msg)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert ISO message to JSON format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Send JSON data to Biller
		serverResp := responsePPOBStatus(jsonIso)
		log.Printf("[Time: %v. Elapsed: %.6fs] Send JSON data to Biller\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Convert response from JSON data to ISO8583 format
		isoParsed = getIsoPPOBStatus(serverResp)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert response from JSON data to ISO8583 format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

	// Process Topup Buy
	case "810002":
		// Convert ISO message to JSON format
		jsonIso := getJsonTopupBuy(msg)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert ISO message to JSON format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Send JSON data to Biller
		serverResp := responseTopupBuy(jsonIso)
		log.Printf("[Time: %v. Elapsed: %.6fs] Send JSON data to Biller\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Convert response from JSON data to ISO8583 format
		isoParsed = getIsoTopupBuy(serverResp)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert response from JSON data to ISO8583 format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

	// Process Topup Check
	case "380003":
		// Convert ISO message to JSON format
		jsonIso := getJsonTopupCheck(msg)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert ISO message to JSON format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Send JSON data to Biller
		serverResp := responseTopupCheck(jsonIso)
		log.Printf("[Time: %v. Elapsed: %.6fs] Send JSON data to Biller\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())

		// Convert response from JSON data to ISO8583 format
		isoParsed = getIsoTopupCheck(serverResp)
		log.Printf("[Time: %v. Elapsed: %.6fs] Convert response from JSON data to ISO8583 format\n", time.Now().Format("15:04:05"), time.Since(start).Seconds())
	}

	isoMessage, _ := isoParsed.ToString()
	isoHeader := fmt.Sprintf("%04d", uniseg.GraphemeClusterCount(isoMessage))

	response.Header, _ = strconv.Atoi(isoHeader)
	response.MTI = isoParsed.Mti.String()
	response.Hex, _ = iso8583.BitMapArrayToHex(isoParsed.Bitmap)
	response.Message = isoMessage

	isoResponse = isoHeader + isoMessage
	log.Printf("\n\nResponse: \n\tHeader: %v\n\tMTI: %v\n\tHex: %v\n\tIso Message: %v\n\tFull Message: %v\n\n",
		response.Header,
		response.MTI,
		response.Hex,
		response.Message,
		isoResponse)

	// create file from response
	filename := "Response_to_" + isoParsed.Elements.GetElements()[3] + "@" + fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	file := CreateFile("storage/response/"+filename, isoResponse)
	log.Println("File created: ", file)

	return isoResponse

}

// Return ISO Message by converting data from map[int]string
func getIso(data map[int]string, mti string) (iso iso8583.IsoStruct) {
	log.Println("Converting to ISO8583...")

	isoStruct := iso8583.NewISOStruct("spec1987.yml", true)
	spec, _ := specFromFile("spec1987.yml")

	if isoStruct.Mti.String() != "" {
		log.Printf("Empty generates invalid MTI")
	}

	// Compare request data length and spec data length, add padding if different
	for field, data := range data {

		fieldSpec := spec.fields[field]

		// Check length for field with Length Type "fixed"
		if fieldSpec.LenType == "fixed" {
			lengthValidate, _ := iso8583.FixedLengthIntegerValidator(field, fieldSpec.MaxLen, data)

			if lengthValidate == false {
				if fieldSpec.ContentType == "n" {
					// Add padding for numeric field
					data = leftPad(data, fieldSpec.MaxLen, "0")
				} else {
					// Add padding for non-numeric field
					data = rightPad(data, fieldSpec.MaxLen, " ")
				}
			}
		}

		// Add data to isoStruct
		isoStruct.AddField(int64(field), data)
	}

	// Add MTI to isoStruct
	isoStruct.AddMTI(mti)

	// Logging isoStruct field and value
	printSortedDE(isoStruct)

	return isoStruct
}

// Return ISO message for PPOB Inquiry JSON response
func getIsoPPOBInquiry(jsonResponse PPOBInquiryResponse) iso8583.IsoStruct {

	log.Println("Converting PPOB Inquiry JSON Response to ISO8583")
	log.Printf("PPOB Inquiry Response (JSON): %v\n", jsonResponse)

	// Assign data to map and add MTI
	var response map[int]string
	if jsonResponse.Rc == "00" {
		response = map[int]string{
			4:   strconv.Itoa(jsonResponse.Tagihan),
			5:   strconv.Itoa(jsonResponse.Admin),
			6:   strconv.Itoa(jsonResponse.TotalTagihan),
			37:  jsonResponse.Reffid,
			39:  jsonResponse.Rc,
			43:  jsonResponse.Nama,
			48:  jsonResponse.Restime,
			62:  jsonResponse.Data,
			120: jsonResponse.Msg,
			121: jsonResponse.Produk,
			122: jsonResponse.Nopel,
		}
	} else {
		response = map[int]string{
			39:  jsonResponse.Rc,
			48:  jsonResponse.Restime,
			120: jsonResponse.Msg,
		}
	}
	mti := "0210"

	// Converting request map to isoStruct
	isoStruct := getIso(response, mti)

	// Adding PAN for PPOB Inquiry Response
	isoStruct.AddField(3, "380001")
	isoMessage, _ := isoStruct.ToString()

	log.Println("Convert Success")
	log.Printf("PPOB Inquiry Response (ISO8583): %s\n", isoMessage)
	return isoStruct

}

// Return ISO message for PPOB Payment JSON response
func getIsoPPOBPayment(jsonResponse PPOBPaymentResponse) iso8583.IsoStruct {

	log.Println("Converting PPOB Payment JSON Response to ISO8583")
	log.Printf("PPOB Payment Response (JSON): %v\n", jsonResponse)

	// Assign data to map and add MTI
	struk := strings.Join(jsonResponse.Struk, ",")
	var response map[int]string
	if jsonResponse.Rc == "00" {
		response = map[int]string{
			4:   strconv.Itoa(jsonResponse.Tagihan),
			5:   strconv.Itoa(jsonResponse.Admin),
			6:   strconv.Itoa(jsonResponse.TotalTagihan),
			37:  jsonResponse.Reffid,
			39:  jsonResponse.Rc,
			43:  jsonResponse.Nama,
			48:  jsonResponse.TglLunas,
			62:  struk,
			120: jsonResponse.Msg,
			121: jsonResponse.Produk,
			122: jsonResponse.Nopel,
			123: jsonResponse.ReffNo,
		}
	} else {
		response = map[int]string{
			39:  jsonResponse.Rc,
			48:  jsonResponse.Restime,
			120: jsonResponse.Msg,
		}
	}
	mti := "0210"

	// Converting request map to isoStruct
	isoStruct := getIso(response, mti)

	// Adding PAN for PPOB Payment Response
	isoStruct.AddField(3, "810001")
	isoMessage, _ := isoStruct.ToString()

	log.Println("Convert Success")
	log.Printf("PPOB Payment Response (ISO8583): %s\n", isoMessage)
	return isoStruct

}

// Return ISO message for PPOB Status JSON response
func getIsoPPOBStatus(jsonResponse PPOBStatusResponse) iso8583.IsoStruct {

	log.Println("Converting PPOB Status JSON Response to ISO8583")
	log.Printf("PPOB Status Response (JSON): %v\n", jsonResponse)

	// Assign data to map and add MTI
	struk := strings.Join(jsonResponse.Struk, ",")
	var response map[int]string
	if jsonResponse.Rc == "00" {
		response = map[int]string{
			4:   strconv.Itoa(jsonResponse.Tagihan),
			5:   strconv.Itoa(jsonResponse.Admin),
			6:   strconv.Itoa(jsonResponse.TotalTagihan),
			37:  jsonResponse.Reffid,
			39:  jsonResponse.Rc,
			43:  jsonResponse.Nama,
			48:  jsonResponse.TglLunas,
			62:  struk,
			120: jsonResponse.Msg,
			121: jsonResponse.Produk,
			122: jsonResponse.Nopel,
			123: jsonResponse.ReffNo,
			124: jsonResponse.Status,
		}
	} else {
		response = map[int]string{
			39:  jsonResponse.Rc,
			48:  jsonResponse.Restime,
			120: jsonResponse.Msg,
		}
	}
	mti := "0210"

	// Converting request map to isoStruct
	isoStruct := getIso(response, mti)

	// Adding PAN for PPOB Status Response
	isoStruct.AddField(3, "380002")
	isoMessage, _ := isoStruct.ToString()

	log.Println("Convert Success")
	log.Printf("PPOB Status Response (ISO8583): %s\n", isoMessage)
	return isoStruct

}

// Return ISO message for Topup Buy JSON response
func getIsoTopupBuy(jsonResponse TopupBuyResponse) iso8583.IsoStruct {

	log.Println("Converting Topup Buy JSON Response to ISO8583")
	log.Printf("Topup Buy Response (JSON): %v\n", jsonResponse)

	// Assign data to map and add MTI
	var response map[int]string
	if jsonResponse.Rc == "00" {
		response = map[int]string{
			39:  jsonResponse.Rc,
			48:  jsonResponse.Restime,
			120: jsonResponse.Msg,
			121: jsonResponse.SN,
			122: jsonResponse.Price,
		}
	} else {
		response = map[int]string{
			39:  jsonResponse.Rc,
			48:  jsonResponse.Restime,
			120: jsonResponse.Msg,
		}
	}
	mti := "0210"

	// Converting request map to isoStruct
	isoStruct := getIso(response, mti)

	// Adding PAN for Topup Buy Response
	isoStruct.AddField(3, "810002")
	isoMessage, _ := isoStruct.ToString()

	log.Println("Convert Success")
	log.Printf("Topup Buy Response (ISO8583): %s\n", isoMessage)
	return isoStruct

}

// Return ISO message for Topup Check JSON response
func getIsoTopupCheck(jsonResponse TopupCheckResponse) iso8583.IsoStruct {

	log.Println("Converting Topup Check JSON Response to ISO8583")
	log.Printf("Topup Check Response (JSON): %v\n", jsonResponse)

	// Assign data to map and add MTI
	var response map[int]string
	if jsonResponse.Rc == "00" {
		response = map[int]string{
			39:  jsonResponse.Rc,
			48:  jsonResponse.Restime,
			120: jsonResponse.Msg,
			121: jsonResponse.SN,
			122: jsonResponse.Price,
		}
	} else {
		response = map[int]string{
			39:  jsonResponse.Rc,
			48:  jsonResponse.Restime,
			120: jsonResponse.Msg,
		}
	}
	mti := "0210"

	// Converting request map to isoStruct
	isoStruct := getIso(response, mti)

	// Adding PAN for Topup Check Response
	isoStruct.AddField(3, "380003")
	isoMessage, _ := isoStruct.ToString()

	log.Println("Convert Success")
	log.Printf("Topup Check Response (ISO8583): %s\n", isoMessage)
	return isoStruct

}

// Log sorted converted ISO Message
func printSortedDE(parsedMessage iso8583.IsoStruct) {
	dataElement := parsedMessage.Elements.GetElements()
	int64toSort := make([]int, 0, len(dataElement))
	for key := range dataElement {
		int64toSort = append(int64toSort, int(key))
	}
	sort.Ints(int64toSort)
	for _, key := range int64toSort {
		log.Printf("[%v] : %v\n", int64(key), dataElement[int64(key)])
	}
}

// Spec contains a structured description of an iso8583 spec
// properly defined by a spec file
type Spec struct {
	fields map[int]fieldDescription
}

// readFromFile reads a yaml specfile and loads
// and iso8583 spec from it
func (s *Spec) readFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(content, &s.fields) // expecting content to be valid yaml
	return nil
}

// SpecFromFile returns a brand new empty spec
func specFromFile(filename string) (Spec, error) {
	s := Spec{}
	err := s.readFromFile(filename)
	if err != nil {
		return s, err
	}
	return s, nil
}

// Add pad on left of data,
// Used to format number by adding "0" in front of number data
func leftPad(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return padding + s
}

// Add pad on right of data,
// Used to format string by adding " " at the end of string data
func rightPad(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return s + padding
}
