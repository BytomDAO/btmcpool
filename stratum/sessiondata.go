package superstratum

// SessionData is generic session data type
type SessionData interface {
	GetWorker() *Worker
	SetWorker(*Worker)
}

// SessionDataBuilder is a helper for building session data
type SessionDataBuilder interface {
	// Build builds a new session data
	Build(id uint) SessionData
}
