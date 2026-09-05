package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"

	"github.com/ahmed-debbech/nostr_talk/config"
	"github.com/ahmed-debbech/nostr_talk/protocol"
)

func main() {
	log.Println("Hello Nostr!")

	peerName := flag.String("n", "", "Name")
	mode := flag.String("m", "", "Mode")
	flag.Parse()

	if *peerName == "" {
		log.Fatal("can not preceed without peer name, use -n to set it")
	}
	if *mode != "t" && *mode != "r" {
		log.Fatal("specify a mode using -m t/r")
	}

	config.PeerName = *peerName
	if err := CreateAndStoreIfNotDone(); err != nil {
		log.Fatal(err)
	}

	if *mode == "t" {
		log.Println("Using t (Transmitter mode)")

		relay := Relay{
			Host:   "nos.lol",
			Scheme: "wss",
			Path:   "",
		}
		if err := relay.Connect(); err != nil {
			log.Fatal(err)
		}

		var sequenceNumber int64 = 0
		for {
			file, err := os.Open("input")
			if err != nil {
				panic(err)
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				//log.Println(scanner.Text())
				keys, err := GeyKeys()
				if err != nil {
					log.Fatal(err)
				}
				relay.Send(protocol.EVENTevent(
					protocol.BuildNoteToBytes(
						scanner.Text(),
						sequenceNumber,
						keys,
					),
				))
				sequenceNumber++
				time.Sleep(time.Second * 2)
			}
		}
	}

	if *mode == "r" {
		log.Println("Using r (Receiving mode)")

		relay := Relay{
			Host:   "nos.lol",
			Scheme: "wss",
			Path:   "",
		}
		if err := relay.Connect(); err != nil {
			log.Fatal(err)
		}
		chEvents := relay.Subscribe()
		relay.Send(protocol.REQEvent())

		bufferer := NewBufferer(100)
		for {
			ev := <-chEvents
			//log.Println(string(ev))
			nostrNote := protocol.ParseEvent(ev)
			if nostrNote.Id == "" {
				continue
			}

			if !nostrNote.VerifyNote() {
				log.Println("[", nostrNote.Id[:7], "]", "Note is NOT verified.")
				continue
			}
			log.Println("[", nostrNote.Id[:7], "]", "Verified Note.")
			bufferer.Buffer(nostrNote)
		}
	}
}
