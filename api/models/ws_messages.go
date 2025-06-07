package models

// wsErrMsg represents a WebSocket error message.
type wsErrMsg struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// NewWsErrMsg creates a new WebSocket error message with the specified message.
// The Type field is set to "error" to indicate that this is an error message.
// This function is used to send error messages over WebSocket connections.
func NewWsErrMsg(msg string) wsErrMsg {
	return wsErrMsg{
		Type: "error",
		Msg:  msg,
	}
}
