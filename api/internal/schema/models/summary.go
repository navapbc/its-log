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
// Date      int64
// KeyID     string
// Operation string
// Tags      sql.NullString
// Value     sql.NullString
// Count     float64

func (ils *ItslogSummary) ReturnHash() string {
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
	return hashed
}

func (ils *ItslogSummary) UpdateHash() {
	hashed := ils.ReturnHash()
	ils.Hash = sql.NullString{String: hashed, Valid: true}
}

func (ils ItslogSummary) IsHashSameAs(other ItslogSummary) bool {
	return ils.ReturnHash() == other.ReturnHash()
}

// func (ils ItslogSummary) IsComputedSameAsExisting() bool {
// 	thisHash := ils.HashItslogSummary()
// 	otherHash := ils.Hash
// 	return thisHash == otherHash.String
// }
