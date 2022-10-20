package superstratum

type Decoder interface {
	Decode(data []byte, session *TcpSession) (Request, error)
}
