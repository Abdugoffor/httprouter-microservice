package clickhouse

import (
	"fmt"
	"time"
)

type AuditEvent struct {
	UserID     uint64
	RoleID     uint64
	Permission string
	Method     string
	Path       string
	Allowed    bool
	StatusCode int
	IP         string
	UserAgent  string
}

func InitAuditTable() error {
	return Exec(`
		CREATE TABLE IF NOT EXISTS permission_audit_log (
			event_time   DateTime  DEFAULT now(),
			user_id      UInt64,
			role_id      UInt64,
			permission   String,
			method       String,
			path         String,
			allowed      UInt8,
			status_code  UInt16,
			ip           String,
			user_agent   String
		) ENGINE = MergeTree()
		ORDER BY (event_time, user_id)
	`)
}

func LogPermissionCheck(e AuditEvent) {
	go func() {
		allowed := uint8(0)
		if e.Allowed {
			allowed = 1
		}

		query := fmt.Sprintf(`
			INSERT INTO permission_audit_log
			(event_time, user_id, role_id, permission, method, path, allowed, status_code, ip, user_agent)
			VALUES ('%s', %d, %d, '%s', '%s', '%s', %d, %d, '%s', '%s')`,
			time.Now().UTC().Format("2006-01-02 15:04:05"),
			e.UserID, e.RoleID,
			sqlEscape(e.Permission),
			sqlEscape(e.Method),
			sqlEscape(e.Path),
			allowed,
			e.StatusCode,
			sqlEscape(e.IP),
			sqlEscape(e.UserAgent),
		)

		if err := Exec(query); err != nil {
			fmt.Printf("[audit] ClickHouse write error: %v\n", err)
		}
	}()
}

func sqlEscape(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' {
			out = append(out, '\\', '\'')
		} else {
			out = append(out, s[i])
		}
	}
	return string(out)
}
