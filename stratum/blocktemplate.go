package superstratum

// BlockTemplate defines common interface for block template
type BlockTemplate interface {
	// create a new job from the block template
	CreateJob(*TcpSession) (Job, error)

	// compare with another block template
	// 1 : newer than the other
	// 0 : same as the other
	// -1 : older than the other
	Compare(BlockTemplate) int
}
