package e2e

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/endpoints"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

const DETERMINISTIC_ITERATIONS = 10

var count_total int = 27 * DETERMINISTIC_ITERATIONS
var count_by_tags = map[string]int{
	"claim.v1":   3 * DETERMINISTIC_ITERATIONS,
	"claim.v2":   3 * DETERMINISTIC_ITERATIONS,
	"claim.v3":   3 * DETERMINISTIC_ITERATIONS,
	"eob.v1":     3 * DETERMINISTIC_ITERATIONS,
	"eob.v2":     3 * DETERMINISTIC_ITERATIONS,
	"eob.v3":     3 * DETERMINISTIC_ITERATIONS,
	"patient.v1": 3 * DETERMINISTIC_ITERATIONS,
	"patient.v2": 3 * DETERMINISTIC_ITERATIONS,
	"patient.v3": 3 * DETERMINISTIC_ITERATIONS,
}

var count_combinations = map[string]int{
	"v2":      9 * DETERMINISTIC_ITERATIONS,
	"eob":     9 * DETERMINISTIC_ITERATIONS,
	"claim":   9 * DETERMINISTIC_ITERATIONS,
	"v1":      9 * DETERMINISTIC_ITERATIONS,
	"v3":      9 * DETERMINISTIC_ITERATIONS,
	"patient": 9 * DETERMINISTIC_ITERATIONS,
}

func DeterministicTest(iterations int) {
	gleCount := GenerateLogEvents(iterations, 10)
	glevCount := GenerateLogEventsWithValues(iterations, 10)
	gclCount := GenerateClusteredLogs(iterations, 10)
	total := gleCount + glevCount + gclCount
	// Make sure we flush the logs.
	time.Sleep(time.Duration(viper.GetInt("buffer.flushwaitsec")+1) * time.Second)
	RunDefaultSequence()

	CheckSummaryValue("count.total", "%", total)
	for tags, value := range count_by_tags {
		CheckSummaryValue("count.by_tags", tags, value)
	}
	for tags, value := range count_combinations {
		CheckSummaryValue("count.combinations", tags, value)
	}
}
func StressTest(iterations int, jitter int) {
	log.Printf("== Running stress test: %d ==\n", iterations)
	gleCount := GenerateLogEvents(iterations, jitter)
	glevCount := GenerateLogEventsWithValues(iterations, jitter)
	gclCount := GenerateClusteredLogs(iterations, jitter)
	total := gleCount + glevCount + gclCount
	time.Sleep(time.Duration(viper.GetInt("buffer.flushwaitsec")+1) * time.Second)
	RunDefaultSequence()
	log.Printf("total events: %d\n", total)
}

func Setup() *types.Storage {
	// Cleanup
	s := types.NewStorage("pupper")
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
	go endpoints.Serve()

	// Make sure nothing is hanging around
	s := Setup()
	Cleanup(s)

	// This is a deterministic test.
	// We should be able to get the same results every time it runs.
	// This can run blocking.
	log.Println("== Running deterministic tests ==")
	DeterministicTest(DETERMINISTIC_ITERATIONS)
	Cleanup(Setup())
	Setup()

	// This stresses the concurrent/parallel nature of the server, with
	// multiple loggers at once, as well as running the ETL while under heavy
	// logging load.
	log.Println("== Running parallel stress tests ==")
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			StressTest(50*i+1, 10*i+1)
		}()
	}
	wg.Wait()
	Cleanup(Setup())

	os.Exit(0)
}
