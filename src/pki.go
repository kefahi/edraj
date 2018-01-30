package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

/*

	RsaKeypair struct with:
		generate,
		encrypt/decrypt,
		sign/verify,
		privatePEM/publicPEM (for exporting the keys as text)
		loadPrivatePEM/loadPublicPEM to parse and load the respective keys.
		Note: loadPublicPEM nullifies RsaKeypair.Private so only encrypt and verify are possible
*/

func pkiExamples() {

	keypair := RsaKeypair{}
	var err error
	err = keypair.Generate()
	if err != nil {
		log.Fatal("While generating", err.Error())
	}

	private, err := keypair.PrivatePEM()
	if err != nil {
		log.Fatal("asString:", "private:", private, "/err:", err.Error())
	}
	public, err := keypair.PublicPEM()
	if err != nil {
		log.Fatal("asString:", "public:", public, "/err:", err.Error())
	}

	fmt.Println("public: \n"+public, "private: \n"+private)
	publicOnlyKeypair := RsaKeypair{}
	//err = keypair.loadPrivatePEM([]byte(private))
	err = publicOnlyKeypair.LoadPublicPem([]byte(public))
	if err != nil {
		log.Fatal("Load:", err.Error())
	}

	data := "Hello World"
	signature, err := keypair.Sign([]byte(data))
	if err != nil {
		log.Fatal("Sign:", err.Error())
	}

	sigb64 := base64.StdEncoding.EncodeToString(signature)
	fmt.Printf("Signature: %v\n", sigb64)

	err = publicOnlyKeypair.Verify([]byte(data), signature)
	if err != nil {
		log.Fatal("Verification failed:", err.Error())
	}

	enc, err := publicOnlyKeypair.Encrypt([]byte("Hello there"), []byte("some label"))
	if err != nil {
		log.Fatalf("encrypt: %s", err)
	}

	dec, err := keypair.Decrypt(enc, []byte("some label"))
	if err != nil {
		log.Fatalf("decrypt: %s", err)
	}
	fmt.Println(string(dec))
	fmt.Println("All good")
}

// RsaKeypair PKI
type RsaKeypair struct {
	public  *rsa.PublicKey
	private *rsa.PrivateKey
}

// Generate new keypair
func (keypair *RsaKeypair) Generate() (err error) {
	keypair.private, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	keypair.public = &keypair.private.PublicKey
	return
}

// Encrypt data
func (keypair *RsaKeypair) Encrypt(clear []byte, label []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, keypair.public, clear, label)
}

// Decrypt data
func (keypair *RsaKeypair) Decrypt(encrypted []byte, label []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, keypair.private, encrypted, label)
}

// PrivatePEM produce text-pem formated string for the private key
func (keypair *RsaKeypair) PrivatePEM() (string, error) {
	var err error
	var buffer bytes.Buffer
	key := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(keypair.private)}
	err = pem.Encode(&buffer, key)
	if err != nil {
		return "", errors.New("Print privateKey:" + err.Error())
	}
	return buffer.String(), nil
}

// PublicPEM produce text-pem formatted string for the public key
func (keypair *RsaKeypair) PublicPEM() (string, error) {
	var err error

	data, err := x509.MarshalPKIXPublicKey(keypair.public)
	if err != nil {
		return "", errors.New("public key marshal:" + err.Error())
	}

	var buffer bytes.Buffer
	key := &pem.Block{Type: "PUBLIC KEY", Bytes: data}
	err = pem.Encode(&buffer, key)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// Sign data
func (keypair *RsaKeypair) Sign(data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	dataHash := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, keypair.private, crypto.SHA256, dataHash)
}

// Verify previously signed data
func (keypair *RsaKeypair) Verify(data []byte, signature []byte) error {
	hash := sha256.New()
	hash.Write(data)
	dataHash := hash.Sum(nil)
	err := rsa.VerifyPKCS1v15(keypair.public, crypto.SHA256, dataHash, signature)
	return err
}

// LoadPublicPem loads publickey from a pem string.
// Note: this disables all the private-key based functions
// Only Encrypt and Verify are possible
func (keypair *RsaKeypair) LoadPublicPem(publicPEM []byte) error {
	keypair.private = nil
	var block *pem.Block
	var err error

	block, _ = pem.Decode(publicPEM)
	if block == nil {
		return errors.New("public key not found")
	} else if block.Type != "PUBLIC KEY" {
		return fmt.Errorf("Unsupported public key  %q", block.Type)
	}

	var key interface{}

	key, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return errors.New("While parsing public key: " + err.Error())
	}
	fmt.Printf("key type %T\n", key)

	switch t := key.(type) {
	case *rsa.PublicKey:
		keypair.public = t
	default:
		fmt.Println("pub is of type:", t)
		return fmt.Errorf("Unsupported public key type %T", key)
	}

	return nil
}

// LoadPrivatePEM Parse and load text-pem formatted private key string
func (keypair *RsaKeypair) LoadPrivatePEM(privatePEM []byte) error {
	var block *pem.Block
	var err error

	block, _ = pem.Decode(privatePEM)
	if block == nil {
		return errors.New("private key not found")
	} else if block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("Unsupported private key type %q", block.Type)
	}
	keypair.private, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("While parsing private key: " + err.Error())
	}

	keypair.public = &keypair.private.PublicKey

	return nil
}
