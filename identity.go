package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/ahmed-debbech/nostr_talk/config"
	"github.com/ahmed-debbech/nostr_talk/protocol"
)

var cachedPubKey = ""
var cachedPrvKey = ""

func CreateAndStoreIfNotDone() error {
	if _, err := getKeysOrCreate(config.PeerName + ".keys"); err != nil {
		return err
	}
	return nil
}

func GetPubKey() string {

	if cachedPubKey != "" {
		return cachedPubKey
	}

	data, err := os.ReadFile(config.PeerName + ".keys")
	if err != nil {
		log.Println("could not read .keys file", err)
		return ""
	}

	var keys protocol.Keys

	if err := json.Unmarshal(data, &keys); err == nil &&
		keys.Pub != "" &&
		keys.Prv != "" {
		cachedPubKey = keys.Pub
		return keys.Pub
	}
	log.Println(".keys file doesn't seem to have valid keys format")
	return ""
}

func GeyKeys() (protocol.Keys, error) {
	if cachedPubKey != "" && cachedPrvKey != "" {
		return protocol.Keys{
			Pub: cachedPubKey,
			Prv: cachedPrvKey,
		}, nil
	}

	data, err := os.ReadFile(config.PeerName + ".keys")
	if err != nil {
		log.Println("could not read .keys file", err)
		return protocol.Keys{}, errors.New("could not read .keys file: " + err.Error())
	}

	var keys protocol.Keys

	if err := json.Unmarshal(data, &keys); err == nil &&
		keys.Pub != "" &&
		keys.Prv != "" {
		cachedPubKey = keys.Pub
		cachedPrvKey = keys.Prv
		return protocol.Keys{
			Pub: keys.Pub,
			Prv: keys.Prv,
		}, nil
	}
	log.Println(".keys file doesn't seem to have valid keys format")
	return protocol.Keys{}, errors.New(".keys file doesn't seem to have valid keys format")
}

func getKeysOrCreate(path string) (protocol.Keys, error) {
	// 1. Try to read the existing file
	data, err := os.ReadFile(path)
	if err == nil {
		// 2. File exists: parse the JSON
		var keys protocol.Keys

		if err := json.Unmarshal(data, &keys); err == nil &&
			keys.Pub != "" &&
			keys.Prv != "" {
			// Existing file contains valid keys
			return keys, nil
		}
		// Existing file is invalid
		errors.New("Invalid keys file: the keys don't seem to be in valid format")
	} else {
		// 3. File does not exist
		newKeys, err := protocol.GenerateKeyPair()
		if err != nil {
			return protocol.Keys{}, err
		}
		data, err = json.MarshalIndent(newKeys, "", "  ")
		if err != nil {
			return protocol.Keys{}, err
		}

		if err := os.WriteFile(path, data, 0600); err != nil {
			return protocol.Keys{}, err
		}
		return newKeys, nil
	}
	return protocol.Keys{}, errors.New("could not create keys because: " + err.Error())
}
