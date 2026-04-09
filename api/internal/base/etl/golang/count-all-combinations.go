package etl

import (
	"context"
	"log"
	"slices"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/navapbc/its-log/internal/types"

	// "gonum.org/v1/gonum/stat/combin"
	combi "github.com/mxschmitt/golang-combinations"
)

func unique[A comparable](input []A) []A {
	seen := make(map[A]bool)
	var result []A

	for _, v := range input {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func CountAllCombinations(etlP *types.RunEtlParams) error {
	allTags, err := etlP.Storage.Queries.GetDistinctTags(context.Background())
	if err != nil {
		panic(":panicohnoes: " + err.Error())
	}
	// Tags take the form
	// a.b.c
	// We want to find all individual tags.
	list_of_tags := make([]string, 0)
	for _, tag := range allTags {
		for _, part := range strings.Split(tag, ".") {
			list_of_tags = append(list_of_tags, part)
		}
	}
	// Uniquify
	list_of_tags = unique(list_of_tags)
	// Now, we want all the combinations, exluding the original tags.
	// E.g.
	// a
	// b
	// c
	// a.b
	// a.c
	// b.c
	counts := make(map[string]int, 0)
	// The Combinations call yields an array of arrays.
	combinations := combi.Combinations(list_of_tags, 0)
	// We now have
	// [
	//   [a]
	//   [a, b]
	//    ...
	// ]
	// for all of the tags.

	// Squirrel is an SQL builder. For this, we could use an ORM, but we don't have one.
	// So, a small library that takes structures and generates SQL strings works.
	// We'll use this to AND together each individual combination as a LIKE query,
	// so that we can get a count of the unique combinations in the events table.
	for _, comb := range combinations {
		// Start with an AND array.
		and_qs := make(sq.And, 0)

		// Comb is a combination, meaning an array.
		for _, tag := range comb {
			// Append a LIKE query into the AND array.
			and_qs = append(and_qs, sq.Like{"tags": "%" + tag + "%"})
		}

		// Start with a base query of counting everything on the table.
		// From there, we'll append a Where() of our And() objects.
		the_query := sq.Select("count(*)").From("itslog_events").Where(and_qs)
		// The ToSql call pulls out the SQL string with prepared `?` elements
		// and an array of arguments in the correct order.
		query_string, args, err := the_query.ToSql()
		// Run the query. Get back one value.
		the_count_row := etlP.Storage.GetDB().QueryRowContext(context.Background(), query_string, args...)
		var the_count int
		err = the_count_row.Scan(&the_count)
		if err != nil {
			// FIXME: Handle this error better.
			// If that query didn't work, panic.
			log.Println("combinations SQL event count err: " + err.Error())
			return err
		}
		// Build a map of the counts against the original name.
		slices.Sort(comb)
		joined := strings.Join(comb, ".")
		counts[joined] = the_count
	}

	for tag, count := range counts {
		// Expect all ETLs to return 0 or 1.
		if count > 0 && !slices.Contains(allTags, tag) {
			// FIXME: sqlc can do this insert.
			_, err := etlP.Storage.GetDB().ExecContext(context.Background(),
				`INSERT OR REPLACE INTO itslog_summary 
				(key_id, operation, tags, value, count)
				VALUES
				(?, 'count.combinations', ?, '', ?)`,
				etlP.KeyId, tag, count)
			if err != nil {
				log.Println("combinations SQL insert err: " + err.Error())
				return err
			}
		}
	}

	return nil
}
