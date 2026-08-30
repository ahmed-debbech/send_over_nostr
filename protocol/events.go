package protocol

import (
	"encoding/json"
	"fmt"
	"log"
)

func REQEvent() []byte {
	data := []interface{}{
		"REQ",
		"myssss",
		map[string]interface{}{
			"kinds": []int{1},
			"authors": []string{
				"807e8dfa713d4f32b78732b939e27fb3bfaf0d93432e8ed56676e334f4425e8b",
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("could not encode json when building REQ event:", err)
		return nil
	}
	log.Println("final REQ structure: ", string(jsonData))
	return jsonData
}

func EVENTevent(event []byte) []byte {
	data := []interface{}{
		"EVENT",
		json.RawMessage(event),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("could not encode json when building EVENT event:", err)
		return nil
	}
	return jsonData
}

func ParseEvent(event []byte) NostrNote {

	var deserialized []json.RawMessage
	json.Unmarshal(event, &deserialized)

	if len(deserialized) != 3 {
		log.Println("Not correct length of EVENT type, skipping...")
		return NostrNote{}
	}

	if string(deserialized[0]) != "\"EVENT\"" {
		log.Println("Not an EVENT type, skipping...")
		return NostrNote{}
	}

	deserializedNote := NostrNote{}
	if err := json.Unmarshal([]byte(deserialized[2]), &deserializedNote); err != nil {
		log.Println("Not able to parse content of EVENT type into NostrNote because:", err, ", skipping...")
		return NostrNote{}
	}
	log.Println("[", deserializedNote.Id[:7], "]", "Recevied and parsed EVENT content into NostrNote successfully.")
	return deserializedNote
}
