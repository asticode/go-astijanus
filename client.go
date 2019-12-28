package astijanus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/asticode/go-astikit"
	"github.com/pkg/errors"
)

// Client represents the client
type Client struct {
	addr string
	s    *astikit.HTTPSender
}

// New creates a new client
func New(c Configuration) *Client {
	return &Client{
		addr: c.Addr,
		s:    astikit.NewHTTPSender(c.Sender),
	}
}

func (c *Client) send(ctx context.Context, method, url string, reqPayload interface{}) (m Message, err error) {
	// Create body
	var body io.Reader
	if reqPayload != nil {
		// Marshal
		buf := &bytes.Buffer{}
		if err = json.NewEncoder(buf).Encode(reqPayload); err != nil {
			err = errors.Wrapf(err, "astijanus: marshaling payload of %s request to %s failed", method, url)
			return
		}

		// Set body
		body = buf
	}

	// Create request
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, method, c.addr+url, body); err != nil {
		err = errors.Wrapf(err, "astijanus: creating %s request to %s failed", method, url)
		return
	}

	// Send
	var resp *http.Response
	if resp, err = c.s.Send(req); err != nil {
		err = errors.Wrapf(err, "astijanus: sending %s request to %s failed", req.Method, req.URL.Path)
		return
	}
	defer resp.Body.Close()

	// Process status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("astijanus: invalid status code %d", resp.StatusCode)
		return
	}

	// Unmarshal
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		err = errors.Wrap(err, "astijanus: unmarshaling response payload failed")
		return
	}

	// Check error
	if m.Error != nil {
		err = fmt.Errorf("astijanus: error %+v in response payload", *m.Error)
		return
	}
	return
}
