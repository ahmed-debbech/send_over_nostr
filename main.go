package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"

	"github.com/ahmed-debbech/nostr_talk/config"
	"github.com/ahmed-debbech/nostr_talk/datasource"
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

		d_source := datasource.FileSource{}
		d_source.Init()
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

		go func() {
			file, err := os.OpenFile(
				"test1.jpg",
				os.O_CREATE|os.O_WRONLY|os.O_APPEND,
				0644,
			)
			if err != nil {
				log.Fatal(err)
			}
			writer := bufio.NewWriterSize(file, 4*1024*1024)
			defer file.Close()

			for {
				if _, err := writer.Write([]byte(bufferer.Pop().Content)); err != nil {
					log.Fatal(err)
				}
				writer.Flush()
			}
		}()

		for {
			ev := <-chEvents
			//log.Println(string(ev))
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
