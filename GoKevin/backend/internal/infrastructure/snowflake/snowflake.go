package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	// Epoch is the custom epoch (2024-01-01 00:00:00 UTC)
	Epoch int64 = 1704067200000

	// Bit allocations
	WorkerIDBits    = 5
	DatacenterIDBits = 5
	SequenceBits    = 12

	// Max values
	MaxWorkerID    = -1 ^ (-1 << WorkerIDBits)    // 31
	MaxDatacenterID = -1 ^ (-1 << DatacenterIDBits) // 31
	MaxSequence    = -1 ^ (-1 << SequenceBits)     // 4095

	// Bit shifts
	WorkerIDShift     = SequenceBits
	DatacenterIDShift = SequenceBits + WorkerIDBits
	TimestampShift    = SequenceBits + WorkerIDBits + DatacenterIDBits
)

// Snowflake represents a snowflake ID generator
type Snowflake struct {
	mu           sync.Mutex
	workerID     int64
	datacenterID int64
	sequence     int64
	lastTimestamp int64
}

// NewSnowflake creates a new snowflake generator
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > MaxWorkerID {
		return nil, errors.New("worker ID out of range")
	}
	if datacenterID < 0 || datacenterID > MaxDatacenterID {
		return nil, errors.New("datacenter ID out of range")
	}

	return &Snowflake{
		workerID:     workerID,
		datacenterID: datacenterID,
		sequence:     0,
		lastTimestamp: -1,
	}, nil
}

// NextID generates the next unique ID
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixMilli()

	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & MaxSequence
		if s.sequence == 0 {
			// Sequence exhausted, wait for next millisecond
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	id := ((timestamp - Epoch) << TimestampShift) |
		(s.datacenterID << DatacenterIDShift) |
		(s.workerID << WorkerIDShift) |
		s.sequence

	return id, nil
}

// ParseID parses a snowflake ID into its components
func ParseID(id int64) map[string]int64 {
	timestamp := (id >> TimestampShift) + Epoch
	datacenterID := (id >> DatacenterIDShift) & MaxDatacenterID
	workerID := (id >> WorkerIDShift) & MaxWorkerID
	sequence := id & MaxSequence

	return map[string]int64{
		"timestamp":    timestamp,
		"datacenterID": datacenterID,
		"workerID":     workerID,
		"sequence":     sequence,
	}
}
