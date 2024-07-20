package slogsyslog

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// formatOptions are options passed to a formatter.
type formatOptions struct {
	// Hostname is the host's name we send when connected to a remote syslog
	// server.
	Hostname string

	// Facility with which we are logging.
	Facility Facility

	// Tag with which we are logging.
	Tag string

	// Prefix value keys with group(s).
	Prefix []byte

	// preformat is a pre-generated value of attributes.
	Preformat []byte
}

// messageFormatter outputs a log message based on the input options.
type messageFormatter func(ctx context.Context, buf []byte, r slog.Record, opts formatOptions) []byte

// goFormat outputs a message in a format as used by the syslog package from the
// standard library.
func goFormat(_ context.Context, buf []byte, r slog.Record, opts formatOptions) []byte {
	buf = append(buf, '<')
	buf = strconv.AppendInt(buf, int64(opts.Facility)|levelToPriority(r.Level), 10)
	buf = append(buf, '>')

	var timestamp time.Time
	if !r.Time.IsZero() {
		timestamp = r.Time
	} else {
		timestamp = time.Now()
	}

	buf = append(buf, timestamp.Format(time.RFC3339)...)
	buf = append(buf, ' ')
	buf = append(buf, opts.Hostname...)
	buf = append(buf, ' ')
	buf = append(buf, opts.Tag...)
	buf = append(buf, '[')
	buf = strconv.AppendInt(buf, int64(os.Getpid()), 10)
	buf = append(buf, ']', ':', ' ')

	if r.NumAttrs() > 0 || len(opts.Preformat) > 0 {
		buf = append(buf, '[')
		buf = append(buf, opts.Preformat...)

		r.Attrs(func(a slog.Attr) bool {
			buf = appendAttr(buf, opts.Prefix, a)
			buf = append(buf, ' ')
			return true
		})
		buf = bytes.TrimSuffix(buf, []byte{' '})

		buf = append(buf, ']', ' ')
	}

	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()

		buf = append(buf, f.File...)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(f.Line), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, r.Message...)
	if !strings.HasSuffix(r.Message, "\n") {
		buf = append(buf, '\n')
	}

	return buf
}

// localFormat outputs a message formatted for a syslog server listening on the
// localhost.
func localFormat(_ context.Context, buf []byte, r slog.Record, opts formatOptions) []byte {
	buf = append(buf, '<')
	buf = strconv.AppendInt(buf, int64(opts.Facility)|levelToPriority(r.Level), 10)
	buf = append(buf, '>')

	var timestamp time.Time
	if !r.Time.IsZero() {
		timestamp = r.Time
	} else {
		timestamp = time.Now()
	}

	buf = append(buf, timestamp.Format(time.Stamp)...)
	buf = append(buf, ' ')
	buf = append(buf, opts.Tag...)
	buf = append(buf, '[')
	buf = strconv.AppendInt(buf, int64(os.Getpid()), 10)
	buf = append(buf, ']', ':', ' ')

	if r.NumAttrs() > 0 || len(opts.Preformat) > 0 {
		buf = append(buf, '[')
		buf = append(buf, opts.Preformat...)

		r.Attrs(func(a slog.Attr) bool {
			buf = appendAttr(buf, opts.Prefix, a)
			buf = append(buf, ' ')
			return true
		})
		buf = bytes.TrimSuffix(buf, []byte{' '})

		buf = append(buf, ']', ' ')
	}

	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()

		buf = append(buf, f.File...)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(f.Line), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, r.Message...)
	if !strings.HasSuffix(r.Message, "\n") {
		buf = append(buf, '\n')
	}

	return buf
}
