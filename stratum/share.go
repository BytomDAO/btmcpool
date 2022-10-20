package superstratum

// ShareState codes
type ShareState int

const (
	ShareStateUnVerified ShareState = iota // share is unverified
	ShareStateRejected                     // share is rejected
	ShareStateAccepted                     // share meets job's difficulty and is accepted
	ShareStateBlock                        // share meets block requirement and is ready to be submitted to node
)

func (state ShareState) String() string {
	switch state {
	case ShareStateUnVerified:
		return "unverified"
	case ShareStateRejected:
		return "rejected"
	case ShareStateAccepted:
		return "accepted"
	case ShareStateBlock:
		return "block found"
	default:
		return "unknown state"
	}
}

// RejectReason codes
type RejectReason int

const (
	RejectReasonUndefined     RejectReason = iota
	RejectReasonPass                       // accepted share
	RejectReasonInvalidJob                 // no job found
	RejectReasonInvalidWorker              // no worker authorized
	RejectReasonDuplicate                  // duplicate share
	RejectReasonStale                      // stale share
	RejectReasonLowDiff                    // lower than required difficulty
	RejectReasonInvalidSol                 // incorrect solution
)

func (reason RejectReason) String() string {
	switch reason {
	case RejectReasonUndefined:
		return "undefined"
	case RejectReasonPass:
		return "no error"
	case RejectReasonInvalidJob:
		return "invalid job"
	case RejectReasonInvalidWorker:
		return "invalid worker"
	case RejectReasonDuplicate:
		return "duplicate share"
	case RejectReasonStale:
		return "stale share"
	case RejectReasonLowDiff:
		return "low diff share"
	case RejectReasonInvalidSol:
		return "invalid solution"
	default:
		return "unknown reason"
	}
}

// RejectReason => ErrorType
func (reason RejectReason) Error() ErrorType {
	switch reason {
	case RejectReasonPass:
		return ErrorNone
	case RejectReasonInvalidJob, RejectReasonStale:
		return ErrorJobNotFound
	case RejectReasonInvalidWorker:
		return ErrorUnauthorized
	case RejectReasonDuplicate:
		return ErrorDuplicateShare
	case RejectReasonLowDiff:
		return ErrorLowDiffShare
	case RejectReasonInvalidSol:
		return ErrorInvalidSolShare
	default:
		return ErrorUnknown
	}
}

// Share is one submitted result from worker
type Share interface {
	// build pb sharelog from the share for logging
	BuildLog(port uint64) ([]byte, error)

	// update share state
	UpdateState(state ShareState, reason RejectReason) error

	// get share state
	GetState() ShareState

	// get Reject Reason
	GetReason() RejectReason

	// get the worker submitted the share
	GetWorker() *Worker
}
