package slogsyslog

import (
	"bytes"
	"log/slog"
	"testing"
	"time"
)

func TestLevelToPriority(t *testing.T) {
	testCases := [...]struct {
		name  string
		level slog.Level
		want  int64
	}{
		{
			name:  "Debug",
			level: slog.LevelDebug,
			want:  7,
		},
		{
			name:  "Info",
			level: slog.LevelInfo,
			want:  6,
		},
		{
			name:  "Warning",
			level: slog.LevelWarn,
			want:  4,
		},
		{
			name:  "Error",
			level: slog.LevelError,
			want:  3,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if lvl := levelToPriority(tc.level); lvl != tc.want {
				t.Errorf("levelToPriority(%s) = %d; want %d", tc.level, lvl, tc.want)
			}
		})
	}
}

func TestAppendAttr(t *testing.T) {
	testCases := [...]struct {
		name string
		attr slog.Attr
		want []byte
	}{
		{
			name: "Time",
			attr: slog.Time("time", time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)),
			want: []byte("time=\"2000-01-02T03:04:05.000Z\""),
		},
		{
			name: "Slice",
			attr: slog.Any("slice", []byte("Test")),
			want: []byte("slice=\"Test\""),
		},
		{
			name: "Group",
			attr: slog.Group("group", slog.String("foo", "bar")),
			want: []byte("group.foo=\"bar\""),
		},
		{
			name: "Default",
			attr: slog.Int("int", 1),
			want: []byte("int=\"1\""),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := make([]byte, 0, 1024)
			buf = appendAttr(buf, nil, tc.attr)
			if !bytes.Equal(buf, tc.want) {
				t.Errorf("appendAttr(buf, <nil>, %v) = %s; want %s", tc.attr, buf, tc.want)
			}
		})
	}
}

func TestAppendKey(t *testing.T) {
	testCases := [...]struct {
		name   string
		prefix []byte
		key    string
		want   []byte
	}{
		{
			name:   "Plain",
			prefix: nil,
			key:    "foo",
			want:   []byte("foo=\""),
		},
		{
			name:   "Prefix",
			prefix: []byte("foo."),
			key:    "bar",
			want:   []byte("foo.bar=\""),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := make([]byte, 0, 1024)
			buf = appendKey(buf, tc.prefix, tc.key)
			if !bytes.Equal(buf, tc.want) {
				t.Errorf("appendKey(buf, %s, %q) = %s; want %s", tc.prefix, tc.key, buf, tc.want)
			}
		})
	}
}

func TestAppendByteSlice(t *testing.T) {
	testCases := [...]struct {
		name  string
		value []byte
		want  []byte
	}{
		{
			name:  "Plain",
			value: []byte("foo"),
			want:  []byte("foo"),
		},
		{
			name:  "Escaped",
			value: []byte("\\f\"o]o"),
			want:  []byte("\\\\f\\\"o\\]o"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := make([]byte, 0, 1024)
			buf = appendByteSlice(buf, tc.value)
			if !bytes.Equal(buf, tc.want) {
				t.Errorf("appendByteSlice(buf, %s) = %s; want %s", tc.value, buf, tc.want)
			}
		})
	}
}
