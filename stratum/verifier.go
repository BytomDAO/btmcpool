package superstratum

// Verifier defines common interface for share verification
type Verifier interface {
	// verify share and update its state
	Verify(Share) error
}
