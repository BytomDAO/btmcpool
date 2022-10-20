package superstratum

import (
	"math/rand"
	"strconv"
)

type JobId uint64

func AllocJobId() JobId {
	return JobId(rand.Uint64())
}

func AllocJobId32() uint32 {
	return rand.Uint32()
}

func (id JobId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func StringToJobId(id string) (JobId, error) {
	u, err := strconv.ParseUint(id, 10, 64)
	return JobId(u), err
}

func HexToJobId(id string) (JobId, error) {
	u, err := strconv.ParseUint(id, 16, 64)
	return JobId(u), err
}

// Job defines common job interface
type Job interface {
	// get unique job id
	GetId() JobId

	// get job difficulty
	GetDiff() uint64

	// return an encoded JSON byte array with job info and ready to be sent to client
	Encode() (interface{}, error)

	// return a job target and two bool which means whether to notify target or diff
	GetTarget() (string, bool, bool)
}
