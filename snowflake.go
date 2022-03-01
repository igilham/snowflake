package snowflake

import (
	"fmt"
	"time"
)

const (
	SEQUENCE_BITS = 12
	MAX_SEQUENCE  = (1 << SEQUENCE_BITS) - 1
	SEQUENCE_MASK = MAX_SEQUENCE

	WORKER_ID_BITS  = 10
	MAX_WORKER_ID   = (1 << WORKER_ID_BITS) - 1
	WORKER_ID_SHIFT = SEQUENCE_BITS
	WORKER_MASK     = MAX_WORKER_ID << WORKER_ID_SHIFT

	TIMESTAMP_BITS  = 41
	MAX_TIMESTAMP   = (1 << TIMESTAMP_BITS) - 1
	TIMESTAMP_SHIFT = SEQUENCE_BITS + WORKER_ID_BITS
	TIMESTAMP_MASK  = MAX_TIMESTAMP << TIMESTAMP_SHIFT
)

type SnowflakeInfoer interface {
	SnowflakeInfo() SnowflakeInfo
}

type Snowflaker interface {
	Snowflake() Snowflake
}

// Snowflake is a number consisting of the packed bits that make up a snowflake identifier.
// It consists of a timestamp, worker ID and sequence number.
type Snowflake uint64

// NewSnowflake creates a new snowflake from the given separate values.
func NewSnowflake(timestamp int64, worker uint16, sequence uint16) Snowflake {
	return Snowflake((uint64(timestamp) << TIMESTAMP_SHIFT) | (uint64(worker) << WORKER_ID_SHIFT) | (uint64(sequence)))
}

// String converts a Snowflake to a human-readable string.
func (s *Snowflake) String() string {
	return fmt.Sprintf("%d-%d-%d", s.Timestamp(), s.Worker(), s.Sequence())
}

// SnowflakeInfo converts a Snowflake to a struct containing the constituent values.
func (s *Snowflake) SnowflakeInfo() SnowflakeInfo {
	return SnowflakeInfo{s.Timestamp(), s.Worker(), s.Sequence()}
}

// Timestamp extracts the timestamp.
func (s *Snowflake) Timestamp() int64 {
	return int64((*s & TIMESTAMP_MASK) >> TIMESTAMP_SHIFT)
}

// Worker extracts the worker ID
func (s *Snowflake) Worker() uint16 {
	return uint16((*s & WORKER_MASK) >> WORKER_ID_SHIFT)
}

// Sequence extracts the sequence number
func (s *Snowflake) Sequence() uint16 {
	return uint16((*s & SEQUENCE_MASK))
}

// SnowflakeInfo is a compound struct representing the separate values held in a Snowflake.
type SnowflakeInfo struct {
	// Timestamp represents the number of seconds since some epoch.
	Timestamp int64
	// Worker represents the ID of the worker that created this Snowflake value.
	Worker uint16
	// Sequence represents the sequence number used by the worker to create this Snowflake value
	Sequence uint16
}

// String converts a snowflake info object to a human readable string
func (s *SnowflakeInfo) String() string {
	return fmt.Sprintf("%d-%d-%d", s.Timestamp, s.Worker, s.Sequence)
}

// Snowflake converts the struct to a Snowflake.
func (s *SnowflakeInfo) Snowflake() Snowflake {
	return NewSnowflake(s.Timestamp, s.Worker, s.Sequence)
}

// Worker exposes a Next function to generate new Snowflakes.
type Worker struct {
	epoch    int64
	id       uint16
	sequence uint16
	prev     int64
}

// NewWorker creates a new worker. It does not check the configuration for errors.
func NewWorker(epoch int64, id uint16) *Worker {
	return &Worker{epoch, id, 0, 0}
}

// Next gets the next Snowflake value. The sequence number resets when the timestamp is updated
// but can roll over if called more than MAX_SEQUENCE times per second.
// This error is not checked.
func (g *Worker) Next() Snowflake {
	now := time.Now().UTC().Unix()

	// TODO: handle clock going backwards
	// if now < g.prev {
	// 	// Error: clock going backwards
	// }

	if now == g.prev {
		g.sequence = (g.sequence + 1) & SEQUENCE_MASK
		if g.sequence == 0 {
			// TODO: handle sequence rollover
			fmt.Printf("Sequence rolled over\n")
		}
	} else {
		g.sequence = 0
	}

	g.prev = now
	return NewSnowflake(now-g.epoch, g.id, g.sequence)
}
