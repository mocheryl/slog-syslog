package slogsyslog

import (
	"strconv"
	"strings"
)

// Facility is the log facility.
type Facility int

// Log facilities.
const (
	Kern Facility = iota << 3

	User

	Mail

	Daemon

	Auth

	Syslog

	LPR

	News

	UUCP

	Cron

	AuthPriv

	FTP

	_ // unused

	_ // unused

	_ // unused

	_ // unused

	Local0

	Local1

	Local2

	Local3

	Local4

	Local5

	Local6

	Local7
)

func (f Facility) String() string {
	switch f {
	case Kern:
		return "Kern"
	case User:
		return "User"
	case Mail:
		return "Mail"
	case Daemon:
		return "Daemon"
	case Auth:
		return "Auth"
	case Syslog:
		return "Syslog"
	case LPR:
		return "LPR"
	case News:
		return "News"
	case UUCP:
		return "UUCP"
	case Cron:
		return "Cron"
	case AuthPriv:
		return "AuthPriv"
	case FTP:
		return "FTP"
	case Local0:
		return "Local0"
	case Local1:
		return "Local1"
	case Local2:
		return "Local2"
	case Local3:
		return "Local3"
	case Local4:
		return "Local4"
	case Local5:
		return "Local5"
	case Local6:
		return "Local6"
	case Local7:
		return "Local7"
	default:
		return "Facility(" + strconv.FormatInt(int64(f), 10) + ")"
	}
}

// structuredEscape escapes all control characters in structured values.
var structuredEscape = strings.NewReplacer(`"`, `\"`, `\`, `\\`, `]`, `\]`)

const (
	// maxBufferSize is the maximum capacity of a byte slice we may return to
	// the buffer pool.
	maxBufferSize = 16 << 10

	// attrTimePrefixLen is the cutoff when adding time value to structured
	// data.
	attrTimePrefixLen = len("2006-01-02T15:04:05.000")
)
