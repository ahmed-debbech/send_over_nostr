package main

import (
	"log"
	"strconv"

	"github.com/ahmed-debbech/nostr_talk/protocol"
)

type Bufferer struct {
	queue           chan protocol.NostrNote
	bufferSize      int
	nextExpectedSeq int64
	seenSequences   map[int64]protocol.NostrNote
}

func NewBufferer(bufferSize int) *Bufferer {
	return &Bufferer{
		queue:           make(chan protocol.NostrNote, bufferSize),
		bufferSize:      bufferSize,
		nextExpectedSeq: 0,
		seenSequences:   make(map[int64]protocol.NostrNote),
	}
}

func (b *Bufferer) Buffer(note protocol.NostrNote) {
	// USED for debugging purposes
	/*log.Println("[INFO]: sequence", note.Tags)
	log.Println("[INFO]: next expected sequence", b.nextExpectedSeq)
	keys := make([]string, 0, len(b.seenSequences))
	for key := range b.seenSequences {
		keys = append(keys, strconv.FormatInt(key, 10))
	}
	log.Println("[INFO]: seen sequences", keys)

	if len(note.Tags) == 0 || len(note.Tags[0]) < 2 {
		log.Println("[ERROR]: sequence number not found in note tags")
		return
	}*/

	if note.Tags[0][0] != "seq" {
		log.Println("[ERROR]: sequence number not found in note tags")
		return
	}

	seq, err := strconv.Atoi(note.Tags[0][1])
	if err != nil {
		log.Println("[ERROR]: sequence number received is not a number")
		return
	}

	if b.nextExpectedSeq < int64(seq) {
		b.seenSequences[int64(seq)] = note
		return
	}
	if b.nextExpectedSeq > int64(seq) {
		return
	}

	if v, ok := b.seenSequences[int64(seq)]; ok {
		b.queue <- v
		delete(b.seenSequences, int64(seq))
	} else {
		b.queue <- note
	}
	b.nextExpectedSeq++
}

func (b *Bufferer) Pop() protocol.NostrNote {
	return <-b.queue
}
