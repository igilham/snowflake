package snowflake

import (
	"fmt"
	"time"
)

const (
	// MAX_ID is 2^10
	MAX_ID = 1024

	// MAX_SEQUENCE is 2^12
	MAX_SEQUENCE = 4096

	// MAX_TIME = 2^41
	MAX_TIME = 2199023255552
)

type Snowflake struct {
	Timestamp int64
	Worker    int16
	Sequence  int16
}

// String converts a snowflake to a string for printing to the console.
func (s Snowflake) String() string {
	return fmt.Sprintf("%X-%d-%d", s.Timestamp, s.Worker, s.Sequence)
}

// TODO: implement translation to binary
// func (s *Snowflake) Binary() int64 {}

// TODO: implement translation to bytes
// func (s *Snowflake) Bytes() byte[] {}

// TODO: implement Parse
// func Parse(snowflake int64) Snowflake {}

type Generator struct {
	epoch    int64
	id       int16
	sequence int16
	prev     Snowflake
}

// NewGenerator creates a new generator. It does not check the configuration for errors.
func NewGenerator(epoch int64, id int16) *Generator {
	return &Generator{epoch, id, 0, Snowflake{}}
}

// Next gets the next Snowflake value. The sequence number resets when the timestamp is updated
// but can roll over if called more than MAX_SEQUENCE times per second.
// This error is not checked.
func (g *Generator) Next() Snowflake {
	now := time.Now().UTC().Unix() - g.epoch

	if g.prev.Timestamp == 0 {
		g.updatePrev(now)
		return g.prev
	}

	if now == g.prev.Timestamp {
		g.updatePrev(now)
		return g.prev
	}

	g.sequence = 0
	g.updatePrev(now)
	return g.prev
}

func (g *Generator) updatePrev(now int64) {
	g.prev = Snowflake{
		now,
		g.id,
		g.sequence,
	}
	g.sequence++

	if g.sequence > MAX_SEQUENCE {
		g.sequence = 0
		// We could raise an error instead here if we need IDs to be unique
	}
}
