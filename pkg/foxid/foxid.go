package foxid

import (
	"github.com/carlmjohnson/crockford"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	// Increment stores global incrementing counter.
	Increment uint32
	// Epoch stores milliseconds from start of FOxID Epoch.
	Epoch uint64 = 1691159571000
	// Datacenter stores current datacenter ID of FOxID generator.
	Datacenter uint16
	// Worker stores current worker id of FOxID generator.
	Worker uint16 = uint16(os.Getpid())
)

// FOxID is a unique identifier with lexicographic sorting capabilities.
type FOxID [16]byte

// Empty returns the new empty FOxID.
//
//goland:noinspection GoUnusedExportedFunction
func Empty() (id FOxID) {
	return id
}

// Generate returns the new generated FOxID.
func Generate() (id FOxID) {
	id.SetTime(time.Now())
	id.SetDatacenter(Datacenter)
	id.SetWorker(Worker)
	id.IncrementCounter()
	id.GenerateRandom()
	return id
}

// Parse parses FOxID from its string representation.
//
//goland:noinspection GoUnusedExportedFunction
func Parse(s string) (id *FOxID, err error) {
	upperString := strings.ToUpper(s)
	decodeString, err := crockford.Upper.DecodeString(upperString)
	if err != nil {
		return nil, err
	} else {
		return (*FOxID)(decodeString), nil
	}
}

// Timestamp returns the timestamp encoded in the FOxID.
func (id *FOxID) Timestamp() uint64 {
	return (uint64(id[5]) | uint64(id[4])<<8 | uint64(id[3])<<16 | uint64(id[2])<<24 | uint64(id[1])<<32 | uint64(id[0])<<40) + Epoch
}

// SetTimestamp sets the timestamp component to the given amount of milliseconds since Unix epoch.
func (id *FOxID) SetTimestamp(timestamp uint64) {
	timestamp -= Epoch
	id[0] = byte(timestamp >> 40)
	id[1] = byte(timestamp >> 32)
	id[2] = byte(timestamp >> 24)
	id[3] = byte(timestamp >> 16)
	id[4] = byte(timestamp >> 8)
	id[5] = byte(timestamp)
}

// Time returns the timestamp component encoded in the FOxID as time.Time.
func (id *FOxID) Time() time.Time {
	return time.UnixMilli(int64(id.Timestamp())).UTC()
}

// SetTime sets the time component of the FOxID to the given in time.Time amount of milliseconds since Unix epoch.
func (id *FOxID) SetTime(t time.Time) {
	id.SetTimestamp(uint64(t.UTC().UnixMilli()))
}

// Generator returns the generator id encoded in the FOxID.
func (id *FOxID) Generator() uint32 {
	return uint32(id[9]) | uint32(id[8])<<8 | uint32(id[7])<<16 | uint32(id[6])<<24
}

// SetGenerator sets the generator component of FOxID to the given value.
func (id *FOxID) SetGenerator(value uint32) {
	id[6] = byte(value >> 24)
	id[7] = byte(value >> 16)
	id[8] = byte(value >> 8)
	id[9] = byte(value)
}

// Datacenter returns the datacenter id encoded in the FOxID.
func (id *FOxID) Datacenter() uint16 {
	return uint16(id[7]) | uint16(id[6])<<8
}

// SetDatacenter sets the datacenter component of FOxID to the given value.
func (id *FOxID) SetDatacenter(value uint16) {
	id[6] = byte(value >> 8)
	id[7] = byte(value)
}

// Worker returns the worker id encoded in the FOxID.
func (id *FOxID) Worker() uint16 {
	return uint16(id[9]) | uint16(id[8])<<8
}

// SetWorker sets the worker component of FOxID to the given value.
func (id *FOxID) SetWorker(value uint16) {
	id[8] = byte(value >> 8)
	id[9] = byte(value)
}

// Counter returns the number of generated ids encoded in the FOxID.
func (id *FOxID) Counter() uint32 {
	return uint32(id[12]) | uint32(id[11])<<8 | uint32(id[10])<<16
}

// SetCounter sets the counter component of FOxID to the given value.
func (id *FOxID) SetCounter(value uint32) {
	id[10] = byte(value >> 16)
	id[11] = byte(value >> 8)
	id[12] = byte(value)
}

// IncrementCounter sets the FOxID counter component to the global incrementing counter.
func (id *FOxID) IncrementCounter() {
	Increment++
	id.SetCounter(Increment)
}

// Random returns the random value used to create the FOxID.
func (id *FOxID) Random() uint32 {
	return uint32(id[15]) | uint32(id[14])<<8 | uint32(id[13])<<16
}

// SetRandom sets the random component of FOxID to the given value.
func (id *FOxID) SetRandom(value uint32) {
	id[13] = byte(value >> 16)
	id[14] = byte(value >> 8)
	id[15] = byte(value)
}

// GenerateRandom sets the random component of FOxID to random value.
func (id *FOxID) GenerateRandom() {
	id.SetRandom(rand.Uint32())
}

// Bytes returns a bytes slice representation of FOxID.
func (id *FOxID) Bytes() []byte {
	return id[:]
}

// String returns a lexicographically sortable string encoded FOxID.
func (id *FOxID) String() string {
	return crockford.Upper.EncodeToString(id.Bytes())
}