package superstratum

import "github.com/segmentio/encoding/json"

// Generic JSON RPC messages
// These are generic format, including neccessary header info
// individual coin should define its only parameter structs

// JSONRpcReq defines generic JSON RPC client request
type JSONRpcReq struct {
	Id     *json.RawMessage `json:"id"`     // message id from client. server response should use the same id
	Method string           `json:"method"` // protocol method name
	Params *json.RawMessage `json:"params"` // parameters, per coin definition
}

// JSONRpcResp defines generic JSON RPC server response
type JSONRpcResp struct {
	Id      *json.RawMessage `json:"id"`      // message id. should match with the initial request from client
	Version string           `json:"jsonrpc"` // fixed version number.
	Result  interface{}      `json:"result"`  // result string, per coin definition
	Error   *ErrorReply      `json:"error"`   // error message
}

// JSONRpcOmitErrorResp defines generic JSON RPC server response omitting empty error
type JSONRpcOmitErrorResp struct {
	Id      *json.RawMessage `json:"id"`              // message id. should match with the initial request from client
	Version string           `json:"jsonrpc"`         // fixed version number.
	Result  interface{}      `json:"result"`          // result string, per coin definition
	Error   *ErrorReply      `json:"error,omitempty"` // error message
}

// StratumJSONRpcNotify defines generic JSON RPC server notification
type StratumJSONRpcNotify struct {
	Version string      `json:"jsonrpc"` // fixed version number
	Method  string      `json:"method"`  // protocol method name
	Params  interface{} `json:"params"`  // paramters, per coin definition
}

// ErrorReply defines standard stratum error message
type ErrorReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorType int

const (
	ErrorNone ErrorType = 0 // no error
	// defined in stratum protocol (DO NOT change value)
	ErrorNotConnected    ErrorType = 19
	ErrorUnknown         ErrorType = 20
	ErrorJobNotFound     ErrorType = 21
	ErrorDuplicateShare  ErrorType = 22
	ErrorLowDiffShare    ErrorType = 23
	ErrorUnauthorized    ErrorType = 24
	ErrorUnsubscribed    ErrorType = 25
	ErrorInvalidSolShare ErrorType = 26

	// custom errors
	ErrorFormatVersion   ErrorType = 28
	ErrorFormatConfigure ErrorType = 29
	ErrorFormatSubscribe ErrorType = 30
	ErrorFormatAuthorize ErrorType = 31
	ErrorFormatSubmit    ErrorType = 32
	ErrorFormatShare     ErrorType = 33
	ErrorUnsupported     ErrorType = 34

	ErrorMultipleAuth ErrorType = 40
)

func (t ErrorType) genError() *ErrorReply {
	var message string
	switch t {
	case ErrorNone:
		message = "No Error"
	case ErrorJobNotFound:
		message = "Job Not Found"
	case ErrorDuplicateShare:
		message = "Duplicate Share"
	case ErrorLowDiffShare:
		message = "Low Difficulty Share"
	case ErrorInvalidSolShare:
		message = "Invalid solution"
	case ErrorFormatVersion:
		message = "Invalid Multion Version Format"
	case ErrorUnauthorized:
		message = "Unauthorized Worker"
	case ErrorUnsubscribed:
		message = "Not Subscribed"
	case ErrorFormatConfigure:
		message = "Invalid Configure Format"
	case ErrorFormatSubscribe:
		message = "Invalid Subscribe Format"
	case ErrorFormatAuthorize:
		message = "Invalid Authorize Format"
	case ErrorFormatSubmit:
		message = "Invalid Submit Format"
	case ErrorFormatShare:
		message = "Invalid Share Format"
	case ErrorUnsupported:
		message = "Unsupported Method"
	case ErrorMultipleAuth:
		message = "Multiple Authorization"
	default:
		message = "Unknown Error"
	}
	return &ErrorReply{
		Code:    int(t),
		Message: message,
	}
}
