// Package slogsyslog implements [log/slog.Handler] for syslog server.
package slogsyslog

import (
	"context"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Options sets the syslog logging options.
type Options struct {
	// Level is the level at which we log at.
	Level slog.Leveler

	// Network protocol to use when connecting to a syslog server.
	Network string

	// Address of the syslog server.
	Address string

	// DialTimeout is duration after which connecting to a syslog server
	// timeouts.
	DialTimeout time.Duration

	// WriteTimeout is duration after which writing to a syslog server timeouts.
	WriteTimeout time.Duration

	// Facility with which we are logging.
	Facility Facility

	// Tag with which we are logging.
	Tag string
}

// SyslogHandler is a structured log [log/slog.Handler] implementation that
// writes messages to a syslog server.
type SyslogHandler struct {
	// mu protects the connection.
	mu *sync.Mutex

	// opts are options for this log.
	opts Options

	// formatter used for writing messages.
	formatter messageFormatter

	// hostname is the host's name we send when connected to a remote syslog
	// server.
	hostname string

	// conn is the syslog connection.
	conn net.Conn

	// prefix value keys with group(s).
	prefix []byte

	// preformat is a pre-generated value of attributes.
	preformat []byte
}

// New creates a new syslog slog [log/slog.Handler]. By default it will log at
// [log/slog.LevelInfo] level to a UNIX datagram socket located at /dev/log.
func New(opts *Options) (*SyslogHandler, error) {
	h := &SyslogHandler{mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	if h.opts.Network == "" {
		h.opts.Network = "unixgram"
	}
	if h.opts.Address == "" {
		h.opts.Address = filepath.Join(string(filepath.Separator), "dev", "log")
	}
	if h.opts.Facility <= 0 {
		h.opts.Facility = Kern
	}
	if h.opts.Tag == "" {
		h.opts.Tag = os.Args[0]
	}

	// TODO: Support other syslog message formats.
	var local bool
	if h.opts.Network == "unixgram" || h.opts.Network == "unix" {
		h.formatter = localFormat
		local = true
	} else {
		h.formatter = goFormat
		h.hostname, _ = os.Hostname()
		local = false
	}

	// TODO: Add TLS support.
	conn, err := net.DialTimeout(h.opts.Network, h.opts.Address, h.opts.DialTimeout)
	if err != nil {
		return nil, err
	}
	h.conn = conn

	if h.hostname == "" && local {
		h.hostname = conn.LocalAddr().String()
	}

	return h, nil
}

func (s *SyslogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= s.opts.Level.Level()
}

func (s *SyslogHandler) Handle(ctx context.Context, r slog.Record) error {
	bufp := allocBuf()
	buf := *bufp

	buf = s.formatter(ctx, buf, r, formatOptions{
		Hostname:  s.hostname,
		Facility:  s.opts.Facility,
		Tag:       s.opts.Tag,
		Prefix:    s.prefix,
		Preformat: s.preformat,
	})

	s.mu.Lock()
	if s.opts.WriteTimeout > 0 {
		s.conn.SetWriteDeadline(time.Now().Add(s.opts.WriteTimeout))
	}
	// FIXME: Don't write to a closed connection.
	_, err := s.conn.Write(buf)
	s.mu.Unlock()
	*bufp = buf
	freeBuf(bufp)
	return err
}

func (s *SyslogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return s
	}

	prefix := make([]byte, 0, len([]byte(name))+len(s.prefix)+1)
	prefix = append(prefix, name...)
	prefix = append(prefix, '.')
	prefix = append(prefix, s.prefix...)

	return &SyslogHandler{
		mu:        s.mu,
		opts:      s.opts,
		formatter: s.formatter,
		hostname:  s.hostname,
		conn:      s.conn,
		prefix:    prefix,
		preformat: s.preformat,
	}
}

func (s *SyslogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return s
	}

	preformat := s.preformat
	for _, a := range attrs {
		preformat = appendAttr(preformat, s.prefix, a)
		preformat = append(preformat, ' ')
	}

	return &SyslogHandler{
		mu:        s.mu,
		opts:      s.opts,
		formatter: s.formatter,
		hostname:  s.hostname,
		conn:      s.conn,
		prefix:    s.prefix,
		preformat: preformat,
	}
}

func (s *SyslogHandler) Close() error {
	s.mu.Lock()
	err := s.conn.Close()
	s.mu.Unlock()

	return err
}
