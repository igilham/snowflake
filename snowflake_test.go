package snowflake

import (
	"testing"
	"time"
)

var (
	epoch             int64 = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	now               int64 = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedTimestamp int64 = now - epoch
)

func TestSnowflakeTimestamp(t *testing.T) {
	s := NewSnowflake(expectedTimestamp, 0, 0)
	actual := s.Timestamp()
	if expectedTimestamp != actual {
		t.Errorf("expected timestamp to be %x but was %x", expectedTimestamp, actual)
	}
}

func TestSnowflakeWorker(t *testing.T) {
	s := NewSnowflake(0, 1001, 0)
	actual := s.Worker()
	if actual != 1001 {
		t.Errorf("expected worker ID to be %d but was %d", 1001, actual)
	}
}

func TestSnowflakeSequence(t *testing.T) {
	s := NewSnowflake(0, 0, 17)
	actual := s.Sequence()
	if actual != 17 {
		t.Errorf("expected sequence number to be %d but was %d", 17, actual)
	}
}

func TestConversionToInfo(t *testing.T) {
	s := NewSnowflake(expectedTimestamp, 1001, 17)
	i := s.SnowflakeInfo()
	if s.Timestamp() != i.Timestamp {
		t.Errorf("timestamp is different: expeceted %d but was %d", expectedTimestamp, i.Timestamp)
	}
}

// Sequence numbering tests might be flaky due to the advancement of time
func TestSequenceIncrements(t *testing.T) {
	w := NewWorker(epoch, 1001)
	s1 := w.Next()
	s2 := w.Next()
	if s2.Sequence() <= s1.Sequence() {
		t.Errorf("Sequence numbering not incrementing")
	}
}
