package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnowflake_NextID(t *testing.T) {
	sf, err := NewSnowflake(1, 1)
	assert.NoError(t, err)

	id1, err := sf.NextID()
	assert.NoError(t, err)
	assert.Greater(t, id1, int64(0))

	id2, err := sf.NextID()
	assert.NoError(t, err)
	assert.NotEqual(t, id1, id2)
}

func TestSnowflake_InvalidWorkerID(t *testing.T) {
	_, err := NewSnowflake(-1, 1)
	assert.Error(t, err)

	_, err = NewSnowflake(32, 1)
	assert.Error(t, err)
}

func TestSnowflake_InvalidDatacenterID(t *testing.T) {
	_, err := NewSnowflake(1, -1)
	assert.Error(t, err)

	_, err = NewSnowflake(1, 32)
	assert.Error(t, err)
}

func TestSnowflake_ValidBoundaryIDs(t *testing.T) {
	sf, err := NewSnowflake(0, 0)
	assert.NoError(t, err)
	id, err := sf.NextID()
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	sf, err = NewSnowflake(MaxWorkerID, MaxDatacenterID)
	assert.NoError(t, err)
	id, err = sf.NextID()
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))
}

func TestSnowflake_UniqueIDs(t *testing.T) {
	sf, err := NewSnowflake(1, 1)
	assert.NoError(t, err)

	ids := make(map[int64]bool)
	for i := 0; i < 10000; i++ {
		id, err := sf.NextID()
		assert.NoError(t, err)
		assert.False(t, ids[id], "duplicate ID generated: %d", id)
		ids[id] = true
	}
	assert.Equal(t, 10000, len(ids))
}

func TestParseID(t *testing.T) {
	sf, _ := NewSnowflake(1, 2)
	id, _ := sf.NextID()

	parsed := ParseID(id)
	assert.Equal(t, int64(2), parsed["datacenterID"])
	assert.Equal(t, int64(1), parsed["workerID"])
	assert.Greater(t, parsed["timestamp"], Epoch)
	assert.GreaterOrEqual(t, parsed["sequence"], int64(0))
}

func TestParseID_RoundTrip(t *testing.T) {
	sf, err := NewSnowflake(15, 20)
	assert.NoError(t, err)

	id, err := sf.NextID()
	assert.NoError(t, err)

	parsed := ParseID(id)
	assert.Equal(t, int64(20), parsed["datacenterID"])
	assert.Equal(t, int64(15), parsed["workerID"])
}
