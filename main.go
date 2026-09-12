package main

import (
	"flag"
	"log"
	"time"

	"github.com/ahmed-debbech/nostr_talk/config"
	"github.com/ahmed-debbech/nostr_talk/datasource"
	"github.com/ahmed-debbech/nostr_talk/datatarget"
	"github.com/ahmed-debbech/nostr_talk/protocol"
)

const sendId = "jdjdd"

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
			Host:   "relay.pocketstr.com",
			Scheme: "wss",
			Path:   "",
		}
		if err := relay.Connect(); err != nil {
			log.Fatal(err)
		}

		d_source := datasource.FileSource{}
		if err := d_source.Init(); err != nil {
			log.Fatal(err)
		}

		datasourceCh := d_source.StartFetch()

		var sequenceNumber int64 = 0
		for dataBin := range datasourceCh {

			keys, err := GeyKeys()
			if err != nil {
				log.Fatal(err)
			}
			relay.Send(protocol.EVENTevent(
				protocol.BuildNoteToBytes(
					dataBin,
					sendId,
					sequenceNumber,
					keys,
				),
			))
			sequenceNumber++
			time.Sleep(time.Microsecond * 100)
		}

		d_source.Destroy()
	}

	if *mode == "r" {
		log.Println("Using r (Receiving mode)")

		relay := Relay{
			Host:   "relay.pocketstr.com",
			Scheme: "wss",
			Path:   "",
		}
		if err := relay.Connect(); err != nil {
			log.Fatal(err)
		}
		chEvents := relay.Subscribe()
		relay.Send(protocol.REQEvent(sendId))

		bufferer := NewBufferer(10000)

		fileTarget := datatarget.FileTarget{}
		fileTarget.Init()
		fileTargetCh := fileTarget.ListenForData()

		go func() {
			for {
				fileTargetCh <- []byte(bufferer.Pop().Content)
			}
		}()

		for {
			ev := <-chEvents
			nostrNote := protocol.ParseEvent(ev)
			if nostrNote.Id == "" {
				log.Println("could not parse the event into a NostrNote, skipping this one...")
				continue
			}

			if !nostrNote.VerifyNote() {
				log.Fatal("[FATAL] [", nostrNote.Id[:7], "]", "Note is NOT verified.")
			}
			log.Println("[", nostrNote.Id[:7], "]", "Verified Note.")

			if err := protocol.NormalizeNote(&nostrNote); err != nil {
				log.Fatal("[FATAL] [", nostrNote.Id[:7], "]", "Could not normalize Note because:", err)
			}
			bufferer.Buffer(nostrNote)

		}
	}
}
