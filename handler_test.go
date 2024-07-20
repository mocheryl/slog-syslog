package slogsyslog

import (
	"log/slog"
	"testing"
)

func TestNew(t *testing.T) {
	f, err := New(nil)
	if f == nil {
		t.Fatal("New(<nil>) cannot return nil")
	}
	if err != nil {
		t.Errorf("New(<nil>) = %v; want nil", err)
	}
}

func TestSyslogHandler_WithGroup(t *testing.T) {
	s, _ := New(nil)
	if s == nil {
		t.Fatal("New(<nil>) cannot return nil")
	}

	s = s.WithGroup("foo").(*SyslogHandler)
	if string(s.prefix) != "foo." {
		t.Fatalf("*SyslogHandler.prefix = %s, want %s", s.prefix, "foo.")
	}

	s = s.WithGroup("bar").(*SyslogHandler)
	if string(s.prefix) != "bar.foo." {
		t.Fatalf("*SyslogHandler.prefix = %s, want %s", s.prefix, "bar.foo.")
	}
}

func TestSyslogHandler_WithAttrs(t *testing.T) {
	s, _ := New(nil)
	if s == nil {
		t.Fatal("New(<nil>) cannot return nil")
	}

	s = s.WithAttrs([]slog.Attr{slog.String("foo", "bar")}).(*SyslogHandler)
	if string(s.preformat) != "foo=\"bar\" " {
		t.Fatalf("*SyslogHandler.preformat = %s, want %s", s.preformat, "foo=\"bar\" ")
	}

	s = s.WithAttrs([]slog.Attr{slog.String("bar", "foo")}).(*SyslogHandler)
	if string(s.preformat) != "foo=\"bar\" bar=\"foo\" " {
		t.Fatalf("*SyslogHandler.preformat = %s, want %s", s.preformat, "foo=\"bar\" bar=\"foo\" ")
	}
}
