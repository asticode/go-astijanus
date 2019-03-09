package astijanus

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
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

// Watch requests to watch the provided id
func (h *StreamingHandle) Watch(id int) (err error) {
	// Send
	if _, err = h.s.c.send(context.Background(), http.MethodPost, fmt.Sprintf("/%d/%d", h.s.id, h.id), Message{
		Body: &MessageBody{
			ID:         id,
			OfferVideo: true,
			Request:    "watch",
		},
		Janus:       "message",
		Transaction: strconv.Itoa(id),
	}); err != nil {
		err = errors.Wrap(err, "astijanus: sending failed")
		return
	}
	return
}

// Start requests to start
func (h *StreamingHandle) Start(jsep *MessageJSEP) (err error) {
	// Send
	if _, err = h.s.c.send(context.Background(), http.MethodPost, fmt.Sprintf("/%d/%d", h.s.id, h.id), Message{
		Body: &MessageBody{
			Request: "start",
		},
		JSEP:        jsep,
		Janus:       "message",
		Transaction: "start",
	}); err != nil {
		err = errors.Wrap(err, "astijanus: sending failed")
		return
	}
	return
}

// Trickle requests to trickle
func (h *StreamingHandle) Trickle(c *MessageCandidate) (err error) {
	// Send
	if _, err = h.s.c.send(context.Background(), http.MethodPost, fmt.Sprintf("/%d/%d", h.s.id, h.id), Message{
		Candidate:   c,
		Janus:       "trickle",
		Transaction: "trickle",
	}); err != nil {
		err = errors.Wrap(err, "astijanus: sending failed")
		return
	}
	return
}
