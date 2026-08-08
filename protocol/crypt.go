package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

type Keys struct {
	Pub string `json:"pub"`
	Prv string `json:"prv"`
}

func GenerateKeyPair() (Keys, error) {
	// Conceptual: Generate a new private key using a secp256k1 library
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Println("Error generating secp256k1 key:", err)
		return Keys{}, errors.New("not able to generate keys because: " + err.Error())
	}

	// Conceptual: Derive the public key
	publicKey := privateKey.PubKey()
	log.Println("Conceptual: secp256k1 keypair generated.")
	log.Printf("Private Key (Hex): %x\n", privateKey.Serialize())
	log.Printf("Public Key (Hex): %x\n", publicKey.SerializeCompressed()[1:])
	return Keys{
		Pub: fmt.Sprintf("%x", publicKey.SerializeCompressed()[1:]),
		Prv: fmt.Sprintf("%x", privateKey.Serialize()),
	}, nil
}

func GenerateIdFromSerializedEvent(ser_event string) ([32]byte, string) {
	hash := sha256.Sum256([]byte(ser_event))
	hash_id := hex.EncodeToString(hash[:])
	return hash, hash_id
}

func GenerateSignature(privateKey string, hex [32]byte) (string, error) {

	privKey, err := GetPrivateKeyFromHex(privateKey)
	if err != nil {
		log.Println("not able to get private key from hex when signing the event because:", err)
		return "", errors.New("not able to get private key from hex when signing the event because: " + err.Error())
	}

	sig, err := schnorr.Sign(privKey, hex[:])
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Signature", string(fmt.Sprintf("%x", sig.Serialize())))
	return fmt.Sprintf("%x", sig.Serialize()), nil
}

func GetPrivateKeyFromHex(hexKey string) (*btcec.PrivateKey, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Println("Error decoding hex key:", err)
		return nil, errors.New("not able to decode hex key because: " + err.Error())
	}
	privKey, _ := btcec.PrivKeyFromBytes(keyBytes)

	return privKey, nil
}

func VerifySignature(publicKey string, signature string, h [32]byte) bool {
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Println("Error decoding public key:", err)
		return false
	}

	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		log.Println("Error decoding signature:", err)
		return false
	}

	pubKey, err := schnorr.ParsePubKey(pubKeyBytes)
	if err != nil {
		log.Println("Error parsing public key:", err)
		return false
	}

	sig, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		log.Println("Error parsing signature:", err)
		return false
	}
	return sig.Verify(h[:], pubKey)
}
