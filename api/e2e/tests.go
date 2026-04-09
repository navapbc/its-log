package e2e

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/endpoints"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

const DETERMINISTIC_ITERATIONS = 10

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

func DeterministicTest(iterations int, dateOffset int) {
	date := time.Now().AddDate(0, 0, dateOffset).Format("2006-01-02")
	log.Println("running for date: " + date)
	gleCount := GenerateLogEvents(iterations, 10, date)
	glevCount := GenerateLogEventsWithValues(iterations, 10, date)
	gclCount := GenerateClusteredLogs(iterations, 10, date)

	total := gleCount + glevCount + gclCount
	// Make sure we flush the logs.
	time.Sleep(time.Duration(viper.GetInt("buffer.flushwaitsec")+1) * time.Second)
	RunDefaultSequence(dateOffset)

	countTotal, countByTags, countCombinations := ExpectedResults(DETERMINISTIC_ITERATIONS)
	CheckSummaryValue("count.total", "%", countTotal, date)
	CheckSummaryValue("count.total", "%", total, date)

	for tags, value := range countByTags {
		CheckSummaryValue("count.by_tags", tags, value, date)
	}
	for tags, value := range countCombinations {
		CheckSummaryValue("count.combinations", tags, value, date)
	}
}
func StressTest(iterations int, jitter int, dateOffset int) {
	date := time.Now().AddDate(0, 0, dateOffset).Format("2006-01-02")
	log.Println("running for date: " + date)

	log.Printf("== Running stress test: %d ==\n", iterations)
	gleCount := GenerateLogEvents(iterations, jitter, date)
	glevCount := GenerateLogEventsWithValues(iterations, jitter, date)
	gclCount := GenerateClusteredLogs(iterations, jitter, date)
	total := gleCount + glevCount + gclCount
	time.Sleep(time.Duration(viper.GetInt("buffer.flushwaitsec")+1) * time.Second)
	RunDefaultSequence(dateOffset)
	log.Printf("total events: %d\n", total)
}

func Setup(dateOffset int) *types.Storage {
	date := time.Now().AddDate(0, 0, dateOffset).Format("2006-01-02")
	s := types.NewStorage("pupper")
	s.SetDate(date)
	s.Init()
	return s
}

func Cleanup(s *types.Storage) {
	s.Delete()
}

func RunTests() {
	gin.DefaultWriter = io.Discard
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	// Run the server
	go endpoints.Serve(types.ServeParams{
		Mode: "debug",
	})

	// This is a deterministic test.
	// We should be able to get the same results every time it runs.
	// This can run blocking.
	log.Println("== Running deterministic tests ==")
	Cleanup(Setup(0))
	DeterministicTest(DETERMINISTIC_ITERATIONS, 0)
	Cleanup(Setup(0))

	log.Println("== Running a week of deterministic tests ==")
	for offset := range 5 {
		Cleanup(Setup(offset))
		Setup(offset)
		DeterministicTest(DETERMINISTIC_ITERATIONS, -1*offset)
	}

	// This stresses the concurrent/parallel nature of the server, with
	// multiple loggers at once, as well as running the ETL while under heavy
	// logging load.
	log.Println("== Running parallel stress tests ==")
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func() {
			StressTest(30*i+1, 10, 0)
			wg.Done()
		}()
	}

	wg.Wait()
	Cleanup(Setup(0))

}
