package slogsyslog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"
)

var testTime = time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)

type name struct {
	First, Last string
}

func (n name) String() string { return n.Last + ", " + n.First }

type text struct {
	s string
}

func (t text) String() string { return t.s } // should be ignored

func (t text) MarshalText() ([]byte, error) {
	if t.s == "" {
		return nil, errors.New("text: empty string")
	}
	return []byte(fmt.Sprintf("text{%q}", t.s)), nil
}

func TestGoFormat(t *testing.T) {
	pid := strconv.Itoa(os.Getpid())
	testCases := [...]struct {
		name string
		attr slog.Attr
		want []byte
	}{
		{
			name: "Unquoted",
			attr: slog.Int("a", 1),
			want: []byte("<6>2000-01-02T03:04:05Z localhost test[" + pid + "]: [a=\"1\"] a message\n"),
		},
		{
			name: "Quoted",
			attr: slog.String("x = y", `qu"o`),
			want: []byte("<6>2000-01-02T03:04:05Z localhost test[" + pid + "]: [x = y=\"qu\\\"o\"] a message\n"),
		},
		{
			name: "String",
			attr: slog.Any("name", name{"Ren", "Hoek"}),
			want: []byte("<6>2000-01-02T03:04:05Z localhost test[" + pid + "]: [name=\"Hoek, Ren\"] a message\n"),
		},
		{
			name: "Struct",
			attr: slog.Any("x", &struct{ A, b int }{A: 1, b: 2}),
			want: []byte("<6>2000-01-02T03:04:05Z localhost test[" + pid + "]: [x=\"&{A:1 b:2}\"] a message\n"),
		},
		{
			name: "Marshaler",
			attr: slog.Any("t", text{"abc"}),
			want: []byte("<6>2000-01-02T03:04:05Z localhost test[" + pid + "]: [t=\"text{\\\"abc\\\"}\"] a message\n"),
		},
		{
			name: "Error",
			attr: slog.Any("t", text{""}),
			want: []byte("<6>2000-01-02T03:04:05Z localhost test[" + pid + "]: [t=\"!ERROR:text: empty string\"] a message\n"),
		},
		{
			name: "Nil",
			attr: slog.Any("a", nil),
			want: []byte("<6>2000-01-02T03:04:05Z localhost test[" + pid + "]: [a=\"<nil>\"] a message\n"),
		},
	}

	opts := formatOptions{
		Hostname: "localhost",
		Tag:      "test",
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := slog.NewRecord(testTime, slog.LevelInfo, "a message", 0)
			r.AddAttrs(tc.attr)

			buf := make([]byte, 0, 1024)
			buf = goFormat(context.Background(), buf, r, opts)
			if !bytes.Equal(buf, tc.want) {
				t.Errorf("defaultFormat(ctx, buf, %v, %v) = %s; want %s", r, opts, buf, tc.want)
			}
		})
	}
}

func TestLocalFormat(t *testing.T) {
	pid := strconv.Itoa(os.Getpid())
	testCases := [...]struct {
		name string
		attr slog.Attr
		want []byte
	}{
		{
			name: "Unquoted",
			attr: slog.Int("a", 1),
			want: []byte("<6>Jan  2 03:04:05 test[" + pid + "]: [a=\"1\"] a message\n"),
		},
		{
			name: "Quoted",
			attr: slog.String("x = y", `qu"o`),
			want: []byte("<6>Jan  2 03:04:05 test[" + pid + "]: [x = y=\"qu\\\"o\"] a message\n"),
		},
		{
			name: "String",
			attr: slog.Any("name", name{"Ren", "Hoek"}),
			want: []byte("<6>Jan  2 03:04:05 test[" + pid + "]: [name=\"Hoek, Ren\"] a message\n"),
		},
		{
			name: "Struct",
			attr: slog.Any("x", &struct{ A, b int }{A: 1, b: 2}),
			want: []byte("<6>Jan  2 03:04:05 test[" + pid + "]: [x=\"&{A:1 b:2}\"] a message\n"),
		},
		{
			name: "Marshaler",
			attr: slog.Any("t", text{"abc"}),
			want: []byte("<6>Jan  2 03:04:05 test[" + pid + "]: [t=\"text{\\\"abc\\\"}\"] a message\n"),
		},
		{
			name: "Error",
			attr: slog.Any("t", text{""}),
			want: []byte("<6>Jan  2 03:04:05 test[" + pid + "]: [t=\"!ERROR:text: empty string\"] a message\n"),
		},
		{
			name: "Nil",
			attr: slog.Any("a", nil),
			want: []byte("<6>Jan  2 03:04:05 test[" + pid + "]: [a=\"<nil>\"] a message\n"),
		},
	}

	opts := formatOptions{
		Tag: "test",
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := slog.NewRecord(testTime, slog.LevelInfo, "a message", 0)
			r.AddAttrs(tc.attr)

			buf := make([]byte, 0, 1024)
			buf = localFormat(context.Background(), buf, r, opts)
			if !bytes.Equal(buf, tc.want) {
				t.Errorf("localFormat(ctx, buf, %v, %v) = %s; want %s", r, opts, buf, tc.want)
			}
		})
	}
}
