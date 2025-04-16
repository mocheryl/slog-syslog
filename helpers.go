package slogsyslog

import (
	"bytes"
	"encoding"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"time"
)

// levelToPriority turns slog level into syslog priority.
func levelToPriority(l slog.Level) int64 {
	var lvl int64 = 3
	if l <= slog.LevelDebug {
		lvl = 7
	} else if l <= slog.LevelInfo {
		lvl = 6
	} else if l <= slog.LevelWarn {
		lvl = 4
	}

	return lvl
}

// appendAttr formats slog's attributes into syslog's structured data.
func appendAttr(buf, prefix []byte, a slog.Attr) []byte {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return buf
	}

	// Most formatting principals for attributes' keys and values have been taken
	// from the standard library's text handler because essentially they look
	// very similar.
	switch a.Value.Kind() {
	case slog.KindTime:
		buf = appendKey(buf, prefix, a.Key)

		n := len(buf)
		t := a.Value.Time().Truncate(time.Millisecond).Add(time.Millisecond / 10)
		buf = t.AppendFormat(buf, time.RFC3339Nano)
		buf = append(buf[:n+attrTimePrefixLen], buf[n+attrTimePrefixLen+1:]...)
		buf = append(buf, '"')
	case slog.KindAny:
		buf = appendKey(buf, prefix, a.Key)

		val := a.Value.Any()
		if src, ok := val.(*slog.Source); ok {
			buf = append(buf, structuredEscape.Replace(src.File)...)
			buf = append(buf, ':')
			buf = strconv.AppendInt(buf, int64(src.Line), 10)
			buf = append(buf, '"')
			break
		}

		if tm, ok := val.(encoding.TextMarshaler); ok {
			data, err := tm.MarshalText()
			if err != nil {
				buf = append(buf, '!', 'E', 'R', 'R', 'O', 'R', ':')
				buf = append(buf, structuredEscape.Replace(fmt.Sprintf("%v", err))...)
				buf = append(buf, '"')
				break
			}

			buf = appendByteSlice(buf, data)
			buf = append(buf, '"')
			break
		}

		if bs, ok := val.([]byte); ok {
			buf = appendByteSlice(buf, bs)
			buf = append(buf, '"')
			break
		}

		t := reflect.TypeOf(val)
		if t != nil && t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
			buf = appendByteSlice(buf, reflect.ValueOf(val).Bytes())
			buf = append(buf, '"')
			break
		}

		buf = append(buf, structuredEscape.Replace(fmt.Sprintf("%+v", val))...)
		buf = append(buf, '"')
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return buf
		}

		var groupPrefix []byte
		if a.Key != "" {
			groupPrefix = make([]byte, 0, len(a.Key)+len(prefix)+1)
			groupPrefix = append(groupPrefix, a.Key...)
			groupPrefix = append(groupPrefix, '.')
			groupPrefix = append(groupPrefix, prefix...)
		} else {
			groupPrefix = prefix
		}

		for _, ga := range attrs {
			buf = appendAttr(buf, groupPrefix, ga)
		}
	default:
		buf = appendKey(buf, prefix, a.Key)
		buf = append(buf, structuredEscape.Replace(a.Value.String())...)
		buf = append(buf, '"')
	}

	return buf
}

// appendKey adds attribute key to the syslog's structured data.
func appendKey(buf, prefix []byte, key string) []byte {
	buf = append(buf, prefix...)
	buf = append(buf, key...)
	buf = append(buf, '=', '"')

	return buf
}

// appendByteSlice inserts safely escaped b as a structured data value into the
// provided buffer.
func appendByteSlice(buf, b []byte) []byte {
	return append(buf, bytes.ReplaceAll(
		bytes.ReplaceAll(
			bytes.ReplaceAll(b, []byte{'\\'}, []byte{'\\', '\\'}), []byte{'"'}, []byte{'\\', '"'}), []byte{']'}, []byte{'\\', ']'})...)
}
