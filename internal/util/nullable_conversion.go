package util

import (
	"database/sql"

	"github.com/google/uuid"
)

func NullUUIDToUUID(nu uuid.NullUUID) *uuid.UUID {
	if nu.Valid {
		return &nu.UUID
	}
	return nil
}

func NullStringToString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
