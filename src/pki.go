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

func pki_examples() {

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
	err = keypair.load([]byte(private), []byte(public))
	if err != nil {
		log.Fatal("Load:", err.Error())
	}

	fmt.Println(keypair)

	data := "Hello World"
	signature, err := keypair.sign([]byte(data))
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
	dec, err := keypair.decrypt(enc, []byte("some label"))
	if err != nil {
		log.Fatalf("decrypt: %s", err)
	}
	fmt.Println(string(dec))
	fmt.Println("All good")
}

// RsaKeypair the pub/private
type RsaKeypair struct {
	Public  *rsa.PublicKey // This is essentially Private.PublicKey
	Private *rsa.PrivateKey
}

func (keypair *RsaKeypair) generate() (err error) {
	keypair.Private, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	keypair.Public = &keypair.Private.PublicKey
	return
}

func (keypair *RsaKeypair) encrypt(clear []byte, label []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, keypair.Public, clear, label)
}

func (keypair *RsaKeypair) decrypt(encrypted []byte, label []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, keypair.Private, encrypted, label)
}

func (keypair *RsaKeypair) privatePEM() (string, error) {
	var err error
	var private bytes.Buffer
	key := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(keypair.Private)}
	err = pem.Encode(&private, key)
	if err != nil {
		return "", errors.New("Print privateKey:" + err.Error())
	}
	return private.String(), nil
}
func (keypair *RsaKeypair) publicPEM() (string, error) {
	var err error

	data, err := x509.MarshalPKIXPublicKey(keypair.Public)
	if err != nil {
		return "", errors.New("public key marshal:" + err.Error())
	}

	var public bytes.Buffer
	key := &pem.Block{Type: "PUBLIC KEY", Bytes: data}
	err = pem.Encode(&public, key)
	if err != nil {
		return "", err
	}

	return public.String(), nil
}

func (keypair *RsaKeypair) sign(data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	dataHash := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, keypair.Private, crypto.SHA256, dataHash)
	return signature, err
}

func (keypair *RsaKeypair) verify(data []byte, signature []byte) error {
	hash := sha256.New()
	hash.Write(data)
	dataHash := hash.Sum(nil)
	err := rsa.VerifyPKCS1v15(keypair.Public, crypto.SHA256, dataHash, signature)
	return err
}

func (keypair *RsaKeypair) load(privateKey []byte, publicKey []byte) error {
	var block *pem.Block
	var err error

	block, _ = pem.Decode(privateKey)
	if block == nil {
		return errors.New("private key not found")
	} else if block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("Unsupported private key type %q", block.Type)
	}
	keypair.Private, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("While parsing private key: " + err.Error())
	}

	block, _ = pem.Decode(publicKey)
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
		keypair.Public = t
	default:
		fmt.Println("pub is of type:", t)
		return fmt.Errorf("Unsupported public key type %T", key)
	}

	return nil
}
