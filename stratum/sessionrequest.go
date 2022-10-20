package superstratum

type Request interface {
	Name() string
	Handle(session *TcpSession) error
	Forward(session *TcpSession) error

	// first param return true if miner is in banned list
	// second param return true if we need to close the session
	CheckMiner(session *TcpSession) (bool, bool)
}
