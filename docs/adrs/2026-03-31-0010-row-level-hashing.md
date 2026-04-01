# Row-level hashing

Date: 2026-03-33

## status

Accepted

## context

Depending on how the data is consolidated and/or used, we may want to have a "fingerprint" on rows. This would be a hash.

Possible reasons for this include, but are not limited to...

### data validation

Knowing that a summary from a given date in the past was not modified can be asserted by hashing values and comparing the hashes. 

### uniqueness

Hashing a row is an easy way to assert that a value being inserted into the summary table is "the same" as a previously inserted value. This comes up when consolidating rows. Being able to compare row hashes makes it easy to assert that a given row is already present in a consolidated table.

### uniqueness, upstream

We may want to export data into log systems like Datadog. These timestream-based systems are not necessarily well-suited to our summarizations; they would rather work with raw event data. But, these systems do have deduplication filters, and the hash value would be a way to "re-insert" summary data on a periodic (daily, weekly) basis, but not have to worry about duplicating historic/prior summary data.

## hashing

The hash on a summary row is computed as follows

```
func (ils ItslogSummary) HashItslogSummary() string {
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
```

The last run timestamp is converted to an RFC3339 value, the components are searated by `-` characters, and a SHA1 is computed and converted to hex. For now, this captures the essence of a summary row, which when being copied forward from the past to the present, *should not change*. We do not hash the "last run" timestamp, because it could change... but the values themselves do not. That the values were computed at different times does not matter if the operation/tags/value/count are the same.

We only keep three decimal places of the `Count` field. This will typically be an integer, but in case it is not, we do go out to three places for hashing purposes. 

## decision

Rows in the summary table will be hashed to support uniqueness calculations.