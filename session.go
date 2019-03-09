package astijanus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Session represents a session
type Session struct {
	c      *Client
	cancel context.CancelFunc
	ctx    context.Context
	id     int
}

func newSession(ctx context.Context, id int, c *Client) (s *Session) {
	// Create session
	s = &Session{
		c:  c,
		id: id,
	}

	// Create context
	s.ctx, s.cancel = context.WithCancel(ctx)
	return
}

// NewSession creates a new session
func (c *Client) NewSession(ctx context.Context) (s *Session, err error) {
	// Send
	var m Message
	if m, err = c.send(context.Background(), http.MethodPost, "", Message{
		Janus:       "create",
		Transaction: "create-session",
	}); err != nil {
		err = errors.Wrap(err, "astijanus: sending failed")
		return
	}

	// No data
	if m.Data == nil {
		err = errors.New("astijanus: no data in message")
		return
	}

	// Create session
	s = newSession(ctx, m.Data.ID, c)
	return
}

// Close closes the session properly
func (s *Session) Close() error {
	s.cancel()
	return nil
}

// LongPollCallback represents long poll callbacks indexed by event name
type LongPollCallbacks struct {
	StreamingPreparing func(transaction string, jsep *MessageJSEP) error
	Unknown            func(m Message) error
}

// LongPoll starts long polling
func (s *Session) LongPoll(cbs LongPollCallbacks) (err error) {
	for {
		if err = s.longPoll(cbs); err != nil {
			return
		}
	}
}

func (s *Session) longPoll(cbs LongPollCallbacks) (err error) {
	// Send
	var m Message
	if m, err = s.c.send(s.ctx, http.MethodGet, fmt.Sprintf("/%d", s.id), nil); err != nil {
		err = errors.Wrap(err, "astijanus: sending failed")
		return
	}

	// No plugin data
	if m.PluginData == nil || m.PluginData.Data == nil {
		if cbs.Unknown != nil {
			if err = cbs.Unknown(m); err != nil {
				err = errors.Wrap(err, "astijanus: executing callback on unknown long poll event failed")
				return
			}
		}
		return
	}

	// Check plugin error
	if m.PluginData.Data.Error != "" {
		err = errors.Wrapf(err, "astijanus: long poll plugin %s error %d with message %s", m.PluginData.Plugin, m.PluginData.Data.ErrorCode, m.PluginData.Data.Error)
		return
	}

	// Switch on status
	switch m.PluginData.Data.Result.Status {
	case "preparing":
		// No jsep
		if m.JSEP == nil {
			err = errors.New("astijanus: no jsep")
			return
		}

		// Execute callback
		if cbs.StreamingPreparing != nil {
			if err = cbs.StreamingPreparing(m.Transaction, m.JSEP); err != nil {
				err = errors.Wrap(err, "astijanus: executing callback on streaming preparing long poll event failed")
				return
			}
		}
	}
	return
}

func (s *Session) attachPlugin(plugin string) (id int, err error) {
	// Send
	var m Message
	if m, err = s.c.send(context.Background(), http.MethodPost, fmt.Sprintf("/%d", s.id), Message{
		Janus:       "attach",
		Plugin:      plugin,
		Transaction: "attach-plugin",
	}); err != nil {
		err = errors.Wrap(err, "astijanus: sending failed")
		return
	}

	// No data
	if m.Data == nil {
		err = errors.New("astijanus: no data in message")
		return
	}

	// Set id
	id = m.Data.ID
	return
}
