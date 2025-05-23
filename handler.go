package zeroslog

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Handler struct {
	out      io.Writer
	timeFmt  string
	minLevel slog.Level
	color    bool
	mu       sync.Mutex
	attrs    []slog.Attr
	groups   []string
}

type Option func(*Handler)

func WithTimeFormat(format string) Option {
	return func(h *Handler) {
		h.timeFmt = format
	}
}

func WithOutput(w io.Writer) Option {
	return func(h *Handler) {
		h.out = w
	}
}

func WithMinLevel(level slog.Level) Option {
	return func(h *Handler) {
		h.minLevel = level
	}
}

func WithColors() Option {
	return func(h *Handler) {
		h.color = true
	}
}

func New(opts ...Option) *Handler {
	h := &Handler{
		out:      os.Stdout,
		timeFmt:  time.RFC3339,
		minLevel: slog.LevelInfo,
		color:    false,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Handler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.minLevel
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	if !h.Enabled(nil, r.Level) {
		return nil
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	levelStr := chooseLevelStr(r.Level, h.color)
	buf.WriteString(levelStr)

	// ----- timestamp -----
	buf.WriteByte('[')
	tsPtr := tsPool.Get().(*[]byte)
	tsBuf := (*tsPtr)[:0]
	tsBuf = r.Time.AppendFormat(tsBuf, h.timeFmt)
	buf.Write(tsBuf)
	*tsPtr = tsBuf
	tsPool.Put(tsPtr)
	buf.Write([]byte{']', ' '})

	// ----- message -----
	buf.WriteString(r.Message)
	padRunes(buf, r.Message, 50)

	// ----- groups -----
	if len(h.groups) > 0 {
		buf.WriteString("[")
		for i, g := range h.groups {
			if i > 0 {
				buf.WriteString(".")
			}
			buf.WriteString(g)
		}
		buf.Write([]byte{']', ' '})
	}

	// ----- previous attrs -----
	for _, a := range h.attrs {
		if h.color {
			buf.WriteString(levelColorCode(r.Level))
			buf.WriteString(a.Key)
			buf.WriteString(cReset)
		} else {
			buf.WriteString(a.Key)
		}
		buf.WriteByte('=')
		appendVal(buf, a.Value.Any())
		buf.WriteByte(' ')
	}

	// ----- attrs -----
	r.Attrs(func(a slog.Attr) bool {
		if h.color {
			buf.WriteString(levelColorCode(r.Level))
			buf.WriteString(a.Key)
			buf.WriteString(cReset)
		} else {
			buf.WriteString(a.Key)
		}
		buf.WriteByte('=')
		appendVal(buf, a.Value.Any())
		buf.WriteByte(' ')
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	_, _ = h.out.Write(buf.Bytes())
	h.mu.Unlock()

	bufPool.Put(buf)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &Handler{
		out:      h.out,
		timeFmt:  h.timeFmt,
		minLevel: h.minLevel,
		color:    h.color,
		attrs:    newAttrs,
		groups:   h.groups,
	}
}

func (h *Handler) WithGroup(group string) slog.Handler {
	newGroups := append([]string{}, h.groups...)
	if group != "" {
		newGroups = append(newGroups, group)
	}
	return &Handler{
		out:      h.out,
		timeFmt:  h.timeFmt,
		minLevel: h.minLevel,
		color:    h.color,
		attrs:    h.attrs,
		groups:   newGroups,
	}
}
