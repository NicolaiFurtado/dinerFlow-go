package utils

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

type AuditEntry struct {
	Count     int         `json:"count"`
	UserID    interface{} `json:"user_id"`
	Timestamp string      `json:"timestamp"`
	OldData   interface{} `json:"old_data"`
	NewData   interface{} `json:"new_data"`
}

// AppendAuditLog takes old audit JSON, adds a new entry, and returns the updated JSON
func AppendAuditLog(existingAudit sql.NullString, userId interface{}, oldData interface{}, newData interface{}) ([]byte, error) {
	var auditHistory []AuditEntry

	if existingAudit.Valid && len(existingAudit.String) > 0 {
		firstChar := strings.TrimSpace(existingAudit.String)[0]

		if firstChar == '{' {
			// Single object previously stored — convert it into an array
			var single AuditEntry
			if err := json.Unmarshal([]byte(existingAudit.String), &single); err != nil {
				return nil, err
			}
			auditHistory = append(auditHistory, single)
		} else if firstChar == '[' {
			// Already a valid array
			if err := json.Unmarshal([]byte(existingAudit.String), &auditHistory); err != nil {
				return nil, err
			}
		}
	}

	newAudit := AuditEntry{
		Count:     len(auditHistory) + 1, // 👈 Set count based on array size
		UserID:    userId,
		Timestamp: time.Now().Format(time.RFC3339),
		OldData:   oldData,
		NewData:   newData,
	}

	auditHistory = append(auditHistory, newAudit)
	return json.Marshal(auditHistory)
}
