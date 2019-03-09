package astijanus

type Message struct {
	Body        interface{}        `json:"body,omitempty"`
	Candidate   *MessageCandidate  `json:"candidate,omitempty"`
	Data        *MessageData       `json:"data,omitempty"`
	Error       *MessageError      `json:"error,omitempty"`
	Janus       string             `json:"janus,omitempty"`
	JSEP        *MessageJSEP       `json:"jsep,omitempty"`
	Plugin      string             `json:"plugin,omitempty"`
	PluginData  *MessagePluginData `json:"plugindata,omitempty"`
	Transaction string             `json:"transaction,omitempty"`
}

type MessageBody struct {
	Request string `json:"request"`
}

type MessageMountpoint struct {
	MessageBody
	Mountpoint
}

type MessageWatch struct {
	MessageBody
	ID         int  `json:"id"`
	OfferVideo bool `json:"offer_video"`
}

type MessageCandidate struct {
	Candidate     string `json:"candidate"`
	Completed     bool   `json:"completed,omitempty"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

type MessageData struct {
	ID int `json:"id"`
}

type MessageError struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

type MessageJSEP struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type MessagePluginData struct {
	Data   *PluginData `json:"data"`
	Plugin string      `json:"plugin"`
}

type PluginData struct {
	Error     string        `json:"error,omitempty"`
	ErrorCode int           `json:"error_code,omitempty"`
	Result    *PluginResult `json:"result,omitempty"`
	Streaming string        `json:"streaming,omitempty"`
}

const (
	PluginResultStatusPreparing = "preparing"
)

type PluginResult struct {
	Status string `json:"status,omitempty"`
}
