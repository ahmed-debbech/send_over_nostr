package protocol

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type NostrNote struct {
	Id        string     `json:"id"`
	PubKey    string     `json:"pubkey"`
	CreatedAt int64      `json:"created_at"`
	Kind      int        `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Signature string     `json:"sig"`
}

func BuildNoteToBytes(content []byte, sendId string, sequence int64, keys Keys) []byte {
	contentBase64 := base64.StdEncoding.EncodeToString([]byte(content))

	ser_event := serialize(
		keys.Pub,
		time.Now().Unix(),
		1,
		[][]string{
			{"seq", strconv.FormatInt(sequence, 10)},
			{"t", sendId},
		},
		contentBase64,
	)
	hashBytes, hash_id := GenerateIdFromSerializedEvent(ser_event)

	sig, err := GenerateSignature(keys.Prv, hashBytes)
	if err != nil {
		return nil
	}

	data := NostrNote{
		Id:        hash_id,
		PubKey:    keys.Pub,
		CreatedAt: time.Now().Unix(),
		Kind:      1,
		Tags: [][]string{
			{"seq", strconv.FormatInt(sequence, 10)},
			{"t", sendId},
		},
		Content:   contentBase64,
		Signature: sig,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("could not encode json when building REQ event:", err)
		return nil
	}
	log.Println("[", data.Id[:7], "]", "Build and now sending EVENT Note sequence: ", sequence, "...")
	return jsonData
}

func (n *NostrNote) VerifyNote() bool {
	IdTo32 := [32]byte{}
	idBytes, err := hex.DecodeString(n.Id)
	if err != nil {
		fmt.Println("could not decode id hex string to bytes when verifying the note:", err)
		return false
	}
	copy(IdTo32[:], idBytes)
	return VerifySignature(n.PubKey, n.Signature, IdTo32)
}

func serialize(pubKey string, created_at int64, kind int, tags [][]string, content string) string {

	arr := []any{
		0,
		pubKey,
		created_at,
		kind,
		tags,
		content,
	}

	jsonBytes, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)
	//log.Println(jsonString)
	return jsonString
}

func NormalizeNote(note *NostrNote) error {
	c, err := base64.StdEncoding.DecodeString(string(note.Content))
	if err != nil {
		log.Println("could not decode base64 content of the note:", err)
		return errors.New("could not decode base64 content of the note: " + err.Error())
	}
	note.Content = string(c)
	return nil
}
