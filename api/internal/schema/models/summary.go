package models

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"strconv"
	"strings"
)

// ID        int64
// LastRun   time.Time
// Date      string
// KeyID     string
// Operation string
// Tags      sql.NullString
// Value     sql.NullString
// Count     float64

func (ils *ItslogSummary) HashItslogSummary() string {
	fields := []string{
		ils.Date,
		ils.KeyID,
		ils.Operation,
		ils.Tags,
		ils.Value,
		strconv.FormatFloat(ils.Count, 'f', 3, 64),
	}
	joined := strings.Join(fields, "-")
	h := sha1.New()
	h.Write([]byte(joined))
	hashed := hex.EncodeToString(h.Sum(nil))
	ils.Hash = sql.NullString{String: hashed, Valid: true}
	return hashed
}

func (ils ItslogSummary) IsHashSameAs(other ItslogSummary) bool {
	thisHash := ils.HashItslogSummary()
	otherHash := other.HashItslogSummary()
	return thisHash == otherHash
}

// func (ils ItslogSummary) IsComputedSameAsExisting() bool {
// 	thisHash := ils.HashItslogSummary()
// 	otherHash := ils.Hash
// 	return thisHash == otherHash.String
// }
