package astijanus

import (
	"errors"
	"fmt"
	"strconv"
)

// Mountpoint types
const (
	StreamingMountpointTypeLive     = "live"
	StreamingMountpointTypeOnDemand = "ondemand"
	StreamingMountpointTypeRTP      = "rtp"
	StreamingMountpointTypeRSTP     = "rstp"
)

// StreamingHandle represents a streaming handle
type StreamingHandle struct {
	*pluginHandle
}

func (s *Session) NewStreamingHandle() (h *StreamingHandle, err error) {
	// Create streaming handle
	h = &StreamingHandle{}

	// Create plugin handle
	if h.pluginHandle, err = s.newPluginHandle("streaming"); err != nil {
		err = fmt.Errorf("astijanus: creating plugin handle failed: %w", err)
		return
	}
	return
}

type Mountpoint struct {
	Audio            bool   `json:"audio,omitempty"`
	AudioPayloadType int    `json:"audiopt,omitempty"`
	AudioPort        int    `json:"audioport,omitempty"`
	AudioRTPMap      string `json:"audiortpmap,omitempty"`
	Description      string `json:"description,omitempty"`
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Permanent        bool   `json:"permanent,omitempty"`
	Type             string `json:"type,omitempty"`
	Video            bool   `json:"video,omitempty"`
	VideoFMTP        string `json:"videofmtp,omitempty"`
	VideoPayloadType int    `json:"videopt,omitempty"`
	VideoPort        int    `json:"videoport,omitempty"`
	VideoRTPMap      string `json:"videortpmap,omitempty"`
}

// Create creates a new mountpoint
func (h *StreamingHandle) Create(m Mountpoint) (err error) {
	// Send
	if _, err = h.send(Message{
		Body: MessageMountpoint{
			MessageBody: MessageBody{
				Request: "create",
			},
			Mountpoint: m,
		},
		Janus:       "message",
		Transaction: "create",
	}); err != nil {
		err = fmt.Errorf("astijanus: sending to streaming handle failed: %w", err)
		return
	}
	return
}

// List lists all mountpoint ids
func (h *StreamingHandle) List() (ms []int, err error) {
	// Send
	var m Message
	if m, err = h.send(Message{
		Body: MessageBody{
			Request: "list",
		},
		Janus:       "message",
		Transaction: "list",
	}); err != nil {
		err = fmt.Errorf("astijanus: sending to streaming handle failed: %w", err)
		return
	}

	// Check list
	if m.PluginData == nil || m.PluginData.Data == nil {
		err = errors.New("main: no list in plugin data")
		return
	}

	// Loop through list items
	for _, v := range m.PluginData.Data.List {
		ms = append(ms, v.ID)
	}
	return
}

// Destroy deletes a mountpoint
func (h *StreamingHandle) Destroy(id int) (err error) {
	// Send
	if _, err = h.send(Message{
		Body: MessageMountpoint{
			MessageBody: MessageBody{
				Request: "destroy",
			},
			Mountpoint: Mountpoint{ID: id},
		},
		Janus:       "message",
		Transaction: "destroy",
	}); err != nil {
		err = fmt.Errorf("astijanus: sending to streaming handle failed: %w", err)
		return
	}
	return
}

// Watch requests to watch the provided id
func (h *StreamingHandle) Watch(id int) (err error) {
	// Send
	if _, err = h.send(Message{
		Body: MessageWatch{
			MessageBody: MessageBody{
				Request: "watch",
			},
			ID:         id,
			OfferVideo: true,
		},
		Janus:       "message",
		Transaction: strconv.Itoa(id),
	}); err != nil {
		err = fmt.Errorf("astijanus: sending to streaming handle failed: %w", err)
		return
	}
	return
}

// Start requests to start with the provided jsep
func (h *StreamingHandle) Start(jsep *MessageJSEP) (err error) {
	// Send
	if _, err = h.send(Message{
		Body: MessageWatch{
			MessageBody: MessageBody{
				Request: "start",
			},
		},
		JSEP:        jsep,
		Janus:       "message",
		Transaction: "start",
	}); err != nil {
		err = fmt.Errorf("astijanus: sending to streaming handle failed: %w", err)
		return
	}
	return
}

// Trickle requests to trickle the provided candidate
func (h *StreamingHandle) Trickle(c *MessageCandidate) (err error) {
	// Send
	if _, err = h.send(Message{
		Candidate:   c,
		Janus:       "trickle",
		Transaction: "trickle",
	}); err != nil {
		err = fmt.Errorf("astijanus: sending to streaming handle failed: %w", err)
		return
	}
	return
}
