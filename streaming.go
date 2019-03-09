package astijanus

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
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
	id int
	s  *Session
}

func newStreamingHandle(id int, s *Session) *StreamingHandle {
	return &StreamingHandle{
		id: id,
		s:  s,
	}
}

// NewStreamingHandle creates a new plugin handle
func (s *Session) NewStreamingHandle() (h *StreamingHandle, err error) {
	// Attach plugin
	var id int
	if id, err = s.attachPlugin("janus.plugin.streaming"); err != nil {
		err = errors.Wrap(err, "astijanus: attaching streaming plugin failed")
		return
	}

	// Create handle
	h = newStreamingHandle(id, s)
	return
}

func (h *StreamingHandle) send(reqPayload interface{}) (m Message, err error) {
	// Send
	if m, err = h.s.c.send(context.Background(), http.MethodPost, fmt.Sprintf("/%d/%d", h.s.id, h.id), reqPayload); err != nil {
		err = errors.Wrap(err, "astijanus: sending failed")
		return
	}
	return
}

type Mountpoint struct {
	Audio            bool   `json:"audio,omitempty"`
	AudioPayloadType int    `json:"audiopt,omitempty"`
	AudioPort        int    `json:"audioport,omitempty"`
	AudioRTPMap      string `json:"audiortpmap,omitempty"`
	Description      string `json:"description"`
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Permanent        bool   `json:"permanent,omitempty"`
	Type             string `json:"type"`
	Video            bool   `json:"video,omitempty"`
	VideoPayloadType int    `json:"videopt,omitempty"`
	VideoPort        int    `json:"videoport,omitempty"`
	VideoRTPMap      string `json:"videortpmap,omitempty"`
}

// Create creates a new mountpoint
func (h *StreamingHandle) CreateMountpoint(m Mountpoint) (err error) {
	// Send
	if _, err = h.send(Message{
		Body: MessageMountpoint{
			MessageBody: MessageBody{
				Request: "create",
			},
			Mountpoint: m,
		},
		Janus:       "message",
		Transaction: "create-mountpoint",
	}); err != nil {
		err = errors.Wrap(err, "astijanus: sending to streaming handle failed")
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
		err = errors.Wrap(err, "astijanus: sending to streaming handle failed")
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
		err = errors.Wrap(err, "astijanus: sending to streaming handle failed")
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
		err = errors.Wrap(err, "astijanus: sending to streaming handle failed")
		return
	}
	return
}
