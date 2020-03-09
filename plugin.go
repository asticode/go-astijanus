package astijanus

import (
	"context"
	"fmt"
	"net/http"
)

type pluginHandle struct {
	id int
	s  *Session
}

func (s *Session) newPluginHandle(plugin string) (h *pluginHandle, err error) {
	// Attach plugin
	var id int
	if id, err = s.attachPlugin("janus.plugin." + plugin); err != nil {
		err = fmt.Errorf("astijanus: attaching plugin failed: %w", err)
		return
	}

	// Create handle
	h = &pluginHandle{
		id: id,
		s:  s,
	}
	return
}

func (h *pluginHandle) send(reqPayload interface{}) (m Message, err error) {
	// Send
	if err = h.s.c.send(context.Background(), http.MethodPost, fmt.Sprintf("/%d/%d", h.s.id, h.id), reqPayload, &m); err != nil {
		err = fmt.Errorf("astijanus: sending failed: %w", err)
		return
	}
	return
}
