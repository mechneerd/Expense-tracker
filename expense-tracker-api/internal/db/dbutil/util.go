package dbutil

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func ParseUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	_ = u.Scan([]byte(s))
	return u
}

func TimestampToString(t pgtype.Timestamp) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}

func TextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

func DateToString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func FloatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan([]byte(fmt.Sprintf("%.2f", f)))
	return n
}

func StringToText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func ParseDate(s string) pgtype.Date {
	var d pgtype.Date
	if s == "" {
		return d
	}
	_ = d.Scan([]byte(s))
	return d
}

func ParseTimestamp(s string) pgtype.Timestamp {
	var t pgtype.Timestamp
	if s == "" {
		return t
	}
	_ = t.Scan([]byte(s))
	return t
}

func TextPtrToString(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func TimestampToTime(t pgtype.Timestamp) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
