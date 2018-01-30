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

	keypair1 := RsaKeypair{}
	var err error
	err = keypair1.generate()
	if err != nil {
		log.Fatal("While generating", err.Error())
	}

	private, err := keypair1.privatePEM()
	if err != nil {
		log.Fatal("asString:", "private:", private, "/err:", err.Error())
	}
	public, err := keypair1.publicPEM()
	if err != nil {
		log.Fatal("asString:", "public:", public, "/err:", err.Error())
	}

	fmt.Println("public: \n"+public, "private: \n"+private)
	keypair := RsaKeypair{}
	//err = keypair.loadPrivatePEM([]byte(private))
	err = keypair.loadPublicPem([]byte(public))
	if err != nil {
		log.Fatal("Load:", err.Error())
	}

	fmt.Println(keypair)

	data := "Hello World"
	signature, err := keypair1.sign([]byte(data))
	if err != nil {
		log.Fatal("Sign:", err.Error())
	}

	sigb64 := base64.StdEncoding.EncodeToString(signature)
	fmt.Printf("Signature: %v\n", sigb64)

	err = keypair.verify([]byte(data), signature)
	if err != nil {
		log.Fatal("Verification failed:", err.Error())
	}

	enc, err := keypair.encrypt([]byte("Hello there"), []byte("some label"))
	if err != nil {
		log.Fatalf("encrypt: %s", err)
	}

	// Decrypt the data
	dec, err := keypair1.decrypt(enc, []byte("some label"))
	if err != nil {
		log.Fatalf("decrypt: %s", err)
	}
	fmt.Println(string(dec))
	fmt.Println("All good")
}

// RsaKeypair the pub/private
type RsaKeypair struct {
	public  *rsa.PublicKey // This is essentially Private.PublicKey
	private *rsa.PrivateKey
}

func (keypair *RsaKeypair) generate() (err error) {
	keypair.private, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	keypair.public = &keypair.private.PublicKey
	return
}

func (keypair *RsaKeypair) encrypt(clear []byte, label []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, keypair.public, clear, label)
}

func (keypair *RsaKeypair) decrypt(encrypted []byte, label []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, keypair.private, encrypted, label)
}

func (keypair *RsaKeypair) privatePEM() (string, error) {
	var err error
	var buffer bytes.Buffer
	key := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(keypair.private)}
	err = pem.Encode(&buffer, key)
	if err != nil {
		return "", errors.New("Print privateKey:" + err.Error())
	}
	return buffer.String(), nil
}
func (keypair *RsaKeypair) publicPEM() (string, error) {
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

func (keypair *RsaKeypair) sign(data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	dataHash := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, keypair.private, crypto.SHA256, dataHash)
	return signature, err
}

func (keypair *RsaKeypair) verify(data []byte, signature []byte) error {
	hash := sha256.New()
	hash.Write(data)
	dataHash := hash.Sum(nil)
	err := rsa.VerifyPKCS1v15(keypair.public, crypto.SHA256, dataHash, signature)
	return err
}

func (keypair *RsaKeypair) loadPublicPem(publicPEM []byte) error {
	// From now this is public-key only
	// Hence only encrypt and verify can be used
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

func (keypair *RsaKeypair) loadPrivatePEM(privatePEM []byte /*, public []byte*/) error {
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
