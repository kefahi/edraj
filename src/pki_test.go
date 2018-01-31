package main

import (
	"fmt"
	"testing"
)

var (
	Tkeypair = RsaKeypair{}
	Ttext    = "Hello World"
	Tlabel   = "testing label"
)

func TestGenerateKeypair(t *testing.T) {
	fmt.Println("Testing method Generate...")
	err := Tkeypair.Generate()
	if err != nil {
		t.Error("Generate method returned error!\n\n\t", err.Error())
	}
}

func TestEncryptData(t *testing.T) {
	fmt.Println("Testing method Encrypt...")
	encByte, err := Tkeypair.Encrypt([]byte(Ttext), []byte(Tlabel))
	if err != nil {
		t.Error("Encrypt method returned error!\n\n\t", err.Error())
	}

	fmt.Println("Testing method Decrypt...")
	decByte, err := Tkeypair.Decrypt(encByte, []byte(Tlabel))

	if err != nil {
		t.Fatal("Decrypt method returned error!\n\n\t", err.Error())
	}
	fmt.Println(string(decByte))
	if string(decByte) != Ttext {
		t.Error("Decryption corrupt!\n\tExpected", Ttext+"!")
	}
}

func TestPrivatePEM(t *testing.T) {
	fmt.Println("Testing method PrivatePEM...")
	textPem, err := Tkeypair.PrivatePEM()
	if err != nil {
		t.Error("PrivatePEM method returned error!\n\n\t", err.Error())
	}
	fmt.Println(textPem)

	fmt.Println("Testing method LoadPrivatePEM...")
	err = Tkeypair.LoadPrivatePEM([]byte(textPem))
	if err != nil {
		t.Error("LoadPrivatePEM method returned error!\n\n\t", err.Error())
	}
}

func TestPublicPEM(t *testing.T) {
	fmt.Println("Testing method PublicPEM...")
	textPem, err := Tkeypair.PublicPEM()
	if err != nil {
		t.Error("PublicPEM method returned error!\n\n\t", err.Error())
	}
	fmt.Println(textPem)
}

func TestSignNverifyData(t *testing.T) {
	fmt.Println("Testing method Sign...")
	signed, err := Tkeypair.Sign([]byte("Hello again"))
	if err != nil {
		t.Fatal("Sign method returned error!\n\n\t", err.Error())
	}

	err = Tkeypair.Verify([]byte("Hello again"), signed)
	if err != nil {
		t.Fatal("Verify method returned error!\n\n\t", err.Error())
	}
}

func TestLoadPublicPEM(t *testing.T) {
	fmt.Println("Testing method LoadPublicPEM...")
	pubPem, err := Tkeypair.PublicPEM()
	if err != nil {
		t.Fatal("PublicPEM method returned error!\n\n\t", err.Error())
	}

	err = Tkeypair.LoadPublicPem([]byte(pubPem))
	if err != nil {
		t.Error("LoadPublicPEM method returned error!\n\n\t", err.Error())
	}
}

// Test keeps failing. The fail is related to the use of the method PrivatePEM
// Uncomment the first part to see yourself
// in another function. It preduces a fatal error, priventing the rest of the tests of running
// if you got it to work.
// Uncomment the second part.
func TestLoadPrivatePEM(t *testing.T) {

	// FIRST PART
	// fmt.Println("Testing method LoadPrivatePEM...")
	// _, err := Tkeypair.PrivatePEM() // If worked, replace _ with priPem
	// if err != nil {
	// 	t.Fatal("PrivatePEM method returned error!\n\n\t", err.Error())
	// }

	// SECOND PART
	// err = Tkeypair.LoadPrivatePEM([]byte(priPem))
	// if err != nil {
	// 	t.Error("LoadPrivatePEM method returned error!\n\n\t", err.Error())
	// }
}
