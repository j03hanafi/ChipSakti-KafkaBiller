package main

import "testing"

func TestGetIso(t *testing.T) {

	request := map[int]string{
		3:  "380001",
		48: "2021                     USER01          WOM             2                        KIOS01                   2018-05-15 15:10:052020",
	}
	mti := "0200"
	isoStruct := getIso(request, mti)

	expected := "0200a00000000001000000000000000000003800011302021                     USER01          WOM             2                        KIOS01                   2018-05-15 15:10:052020"
	result, _ := isoStruct.ToString()

	if result != expected {
		t.Errorf("getIso() failed, \nexpected\t: %v, \ngot\t\t\t: %v", expected, result)
	} else {
		t.Log("getIso() success")
	}
}

func TestGetIsoPPOBInquiry(t *testing.T) {

	var jsonRequest PPOBInquiryResponse

	jsonRequest = PPOBInquiryResponse{
		Rc:           "00",
		Msg:          "approve",
		Produk:       "WOM",
		Nopel:        "2",
		Nama:         "HANAFI",
		Tagihan:      1500,
		Admin:        3300,
		TotalTagihan: 4800,
		Reffid:       "12345",
		Data:         "2020",
		Restime:      "2021-03-18 08:03:23",
	}

	isoRequest := getIsoPPOBInquiry(jsonRequest)

	expected := "0210bc0000000a21000400000000000001c038000100000000150000000000330000000000480012345       00200HANAFI                                  0192021-03-18 08:03:230042020007approve003WOM0012"
	result, _ := isoRequest.ToString()

	if result != expected {
		t.Errorf("getIsoPPOBInquiry() failed, \nexpected\t: %v, \ngot\t\t\t: %v", expected, result)
	} else {
		t.Log("getIsoPPOBInquiry() success")
	}
}

func TestGetIsoPPOBPayment(t *testing.T) {

	var jsonRequest PPOBPaymentResponse

	jsonRequest = PPOBPaymentResponse{
		Rc:           "00",
		Msg:          "approve",
		Produk:       "WOM",
		Nopel:        "2",
		Nama:         "HANAFI",
		Tagihan:      870000,
		Admin:        3300,
		TotalTagihan: 873300,
		Reffid:       "12345",
		TglLunas:     "2021-03-18 08:03:35",
		Struk: []string{
			"pembayaranWOM",
			"",
			"ID PEL :2",
			"NAMA :HANAFI",
			"REF : 5/4-3-2-1",
			"ANGSURAN KE: 5",
			"TAGIHAN : Rp 870000",
			"BIAYA ADMIN : Rp 3300",
			"TTL TAGIHAN : Rp 873300",
			"",
			"STRUK INI ADALAH BUKTI PEMBAYARAN YANG SAH",
			"TERIMA KASIH",
		},
		ReffNo:  "54321",
		Restime: "",
	}

	isoRequest := getIsoPPOBPayment(jsonRequest)

	expected := "0210bc0000000a21000400000000000001e081000100000087000000000000330000000087330012345       00200HANAFI                                  0192021-03-18 08:03:35191pembayaranWOM,,ID PEL :2,NAMA :HANAFI,REF : 5/4-3-2-1,ANGSURAN KE: 5,TAGIHAN : Rp 870000,BIAYA ADMIN : Rp 3300,TTL TAGIHAN : Rp 873300,,STRUK INI ADALAH BUKTI PEMBAYARAN YANG SAH,TERIMA KASIH007approve003WOM001200554321"
	result, _ := isoRequest.ToString()

	if result != expected {
		t.Errorf("getIsoPPOBPayment() failed, \nexpected\t: %v, \ngot\t\t\t: %v", expected, result)
	} else {
		t.Log("getIsoPPOBPayment() success")
	}
}

func TestGetIsoPPOBStatus(t *testing.T) {

	var jsonRequest PPOBStatusResponse

	jsonRequest = PPOBStatusResponse{
		Rc:           "00",
		Msg:          "approve",
		Produk:       "WOM",
		Nopel:        "2",
		Nama:         "HANAFI",
		Tagihan:      10000,
		Admin:        0,
		TotalTagihan: 10000,
		Reffid:       "12345",
		TglLunas:     "2021-03-18 08:03:38",
		Struk: []string{
			"<b>PT. MULTI ACCESS INDONESIA - CHIPSAKTI</b>",
			"",
			"LOKET : ZONATIK",
			"TGL BAYAR : 02/07/2018 / 14:16:44",
			"",
			"STRUK PEMBAYARAN LANGGANANWOM",
			"",
			"IDPEL 2",
			"NAMA : HANAFI",
			"TTL TAGIHAN : Rp 10000",
			"",
			"STRUK INI ADALAH BUKTI PEMBAYARAN YANG SAH",
			"TERIMA KASIH",
		},
		ReffNo: "0123",
		Status: "payment Successfull",
	}

	isoRequest := getIsoPPOBStatus(jsonRequest)

	expected := "0210bc0000000a21000400000000000001f038000200000001000000000000000000000001000012345       00200HANAFI                                  0192021-03-18 08:03:38230<b>PT. MULTI ACCESS INDONESIA - CHIPSAKTI</b>,,LOKET : ZONATIK,TGL BAYAR : 02/07/2018 / 14:16:44,,STRUK PEMBAYARAN LANGGANANWOM,,IDPEL 2,NAMA : HANAFI,TTL TAGIHAN : Rp 10000,,STRUK INI ADALAH BUKTI PEMBAYARAN YANG SAH,TERIMA KASIH007approve003WOM00120040123019payment Successfull"
	result, _ := isoRequest.ToString()

	if result != expected {
		t.Errorf("getIsoPPOBStatus() failed, \nexpected\t: %v, \ngot\t\t\t: %v", expected, result)
	} else {
		t.Log("getIsoPPOBStatus() success")
	}
}

func TestGetIsoTopupBuy(t *testing.T) {

	var jsonRequest TopupBuyResponse

	jsonRequest = TopupBuyResponse{
		Rc:      "00",
		Msg:     "PembelianWOMberhasil. Harga Rp. 1000",
		Restime: "2021-03-18 08:03:42",
		SN:      "12345678",
		Price:   "1000",
	}

	isoRequest := getIsoTopupBuy(jsonRequest)

	expected := "0210a00000000201000000000000000001c0810002002000192021-03-18 08:03:42036PembelianWOMberhasil. Harga Rp. 1000008123456780041000"
	result, _ := isoRequest.ToString()

	if result != expected {
		t.Errorf("getIsoTopupBuy() failed, \nexpected\t: %v, \ngot\t\t\t: %v", expected, result)
	} else {
		t.Log("getIsoTopupBuy() success")
	}
}

func TestGetIsoTopupCheck(t *testing.T) {

	var jsonRequest TopupCheckResponse

	jsonRequest = TopupCheckResponse{
		Rc:      "00",
		Msg:     "PembelianWOMberhasil. Harga Rp. 1000",
		Restime: "2021-03-18 08:03:45",
		SN:      "12345678",
		Price:   "1000",
	}

	isoRequest := getIsoTopupCheck(jsonRequest)

	expected := "0210a00000000201000000000000000001c0380003002000192021-03-18 08:03:45036PembelianWOMberhasil. Harga Rp. 1000008123456780041000"
	result, _ := isoRequest.ToString()

	if result != expected {
		t.Errorf("getIsoTopupCheck() failed, \nexpected\t: %v, \ngot\t\t\t: %v", expected, result)
	} else {
		t.Log("getIsoTopupCheck() success")
	}
}
