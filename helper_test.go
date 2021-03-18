package main

import "testing"

func TestCreateFile(t *testing.T) {
	createdFile := CreateFile("testFile", "testContent")
	expectedFileName := "testFile.txt"

	if createdFile != expectedFileName {
		t.Errorf("CreateFile() failed. Expected: %v. Got: %v", expectedFileName, createdFile)
	} else {
		t.Log("CreateFile() success")
	}
}

func TestSignatureSHA256(t *testing.T) {
	data := "$inquiry$2021$USER01$KIOS01$2018-05-15 15:10:05$unand$"
	expected := "41146b35700a4f6c05bd9dfd9da52b63f308050baa61d395e4f78a7d1625d70a"
	result := signatureSHA256(data)

	if result != expected {
		t.Errorf("signatureSHA256() failed. Expected: %v. Got: %v", expected, result)
	} else {
		t.Log("signatureSHA256() success")
	}
}
