package e2e

import (
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/endpoints"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

const DETERMINISTIC_ITERATIONS = 3
const SEQUENCE_NAME = "full" // or "default"
const EXTRA_SLEEP_TIME = 5   // seconds

func ExpectedResults(iterations int) (int, map[string]int, map[string]int) {

	var countTotal int = 27 * iterations
	var countByTags = map[string]int{
		"claim.v1":   3 * iterations,
		"claim.v2":   3 * iterations,
		"claim.v3":   3 * iterations,
		"eob.v1":     3 * iterations,
		"eob.v2":     3 * iterations,
		"eob.v3":     3 * iterations,
		"patient.v1": 3 * iterations,
		"patient.v2": 3 * iterations,
		"patient.v3": 3 * iterations,
	}

	var countCombinations = map[string]int{
		"v2":      9 * iterations,
		"eob":     9 * iterations,
		"claim":   9 * iterations,
		"v1":      9 * iterations,
		"v3":      9 * iterations,
		"patient": 9 * iterations,
	}

	return countTotal, countByTags, countCombinations
}

// func truncateToDay(t time.Time) time.Time {
// 	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
// }
// func offsetDate(offset int) int64 {
// 	d := time.Now().AddDate(0, 0, -1*offset)
// 	return truncateToDay(d).Unix()
// }

func sleepy() {
	sleepSec := viper.GetInt("buffer.flushwaitsec") + EXTRA_SLEEP_TIME
	time.Sleep(time.Duration(sleepSec) * time.Second)
}

func DeterministicTest(t *testing.T, iterations int, dateOffset int) {
	date := types.NewILTimeToday()
	date.SubtractDays(dateOffset)

	log.Printf("DeterministicTest: %s\n", date.AsYYYYMMDD())
	gleCount := GenerateLogEvents(iterations, 10, date)
	glevCount := GenerateLogEventsWithValues(iterations, 10, date)
	gclCount := GenerateClusteredLogs(t, iterations, 10, date)

	total := gleCount + glevCount + gclCount
	// Make sure we flush the logs.
	sleepy()
	RunSequence(t, dateOffset, SEQUENCE_NAME)

	countTotal, countByTags, countCombinations := ExpectedResults(iterations)
	if !CheckSummaryValue("count.total", "%", countTotal, date) {
		t.Fail()
	}
	if !CheckSummaryValue("count.total", "%", total, date) {
		t.Fail()
	}

	for tags, value := range countByTags {
		if !CheckSummaryValue("count.by_tags", tags, value, date) {
			t.Fail()
		}
	}
	for tags, value := range countCombinations {
		if !CheckSummaryValue("count.combinations", tags, value, date) {
			t.Fail()
		}
	}
}

func StressTest(t *testing.T, iterations int, jitter int, dateOffset int) {
	date := types.NewILTimeToday()
	log.Println("running for date: " + date.AsYYYYMMDD())

	log.Printf("== Running stress test: %d ==\n", iterations)
	gleCount := GenerateLogEvents(iterations, jitter, date)
	glevCount := GenerateLogEventsWithValues(iterations, jitter, date)
	gclCount := GenerateClusteredLogs(t, iterations, jitter, date)
	total := gleCount + glevCount + gclCount
	sleepy()
	RunSequence(t, dateOffset, SEQUENCE_NAME)
	log.Printf("total events: %d\n", total)
}

func Setup(t *testing.T, dateOffset int) *types.Storage {
	date := types.NewILTimeToday()
	date.SubtractDays(dateOffset)

	s := types.NewStorage("pupper")
	s.SetDateILT(date)
	s.Init()
	// DEBUG LOG
	// t.Log("setting up: " + s.Filename)
	return s
}

func Cleanup(s *types.Storage) {
	s.Delete()
}

func RunTests(t *testing.T) {
	gin.DefaultWriter = io.Discard
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	// Run the server
	go endpoints.Serve(types.ServeParams{
		Mode: "debug",
	})

	// This is a deterministic test.
	// We should be able to get the same results every time it runs.
	// This can run blocking.
	t.Log("== Running deterministic tests ==")
	Cleanup(Setup(t, 0))
	DeterministicTest(t, DETERMINISTIC_ITERATIONS, 0)
	Cleanup(Setup(t, 0))

	t.Log("== Running a week of deterministic tests ==")
	// Offsets should be 1-5
	offsets := makeRange(1, 5)
	for offset := range offsets {
		Cleanup(Setup(t, offset))
		Setup(t, offset)
		DeterministicTest(t, DETERMINISTIC_ITERATIONS, offset)
		Cleanup(Setup(t, offset))
	}

	// This stresses the concurrent/parallel nature of the server, with
	// multiple loggers at once, as well as running the ETL while under heavy
	// logging load.
	t.Log("== Running parallel stress tests ==")
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func() {
			StressTest(t, DETERMINISTIC_ITERATIONS*i+1, 10, 0)
			wg.Done()
		}()
	}

	wg.Wait()
	Cleanup(Setup(t, 0))

}
