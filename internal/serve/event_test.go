package serve

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/csp"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/spf13/viper"
	"github.com/steinfletcher/apitest"
)

func blackHole(c chan *itslog.Event) {
	<-c
}

func checkEq(t *testing.T, expected *itslog.Event) func(c chan *itslog.Event) {
	return func(c chan *itslog.Event) {
		go func() {
			select {
			case v := <-c:
				if !(v.Event == expected.Event && v.Source == expected.Source) {
					t.Error()
				}
			case <-time.After(1 * time.Second):
				t.Error()
			}
		}()
	}
}

func setup(consumer func(chan *itslog.Event)) *gin.Engine {
	var ch_evt_out = make(chan *itslog.Event)
	router := gin.Default()
	apiV1 := router.Group("/v1")
	apiV1.PUT("event/:appID/:eventID", Event(ch_evt_out))
	// This drains the channel so we don't have to worry about
	// it as part of the testing.
	go consumer(ch_evt_out)
	return router
}

/*
 * This test checks whether or not the PUT is picked up by the framework.
 * There is no checking of the value. It just makes sure the phone was picked up.
 * blackHole makes sure the channel is drained from the API handler, or otherwise
 * it will hang forever, waiting for the channel communication to terminate.
 */
func TestPutMessage(t *testing.T) {
	router := setup(blackHole)
	apitest.New().
		Handler(router).
		Put("/v1/event/us.me.lewiston/forage-bagels").
		Expect(t).
		Status(http.StatusOK).
		End()
}

/*
 * This test is a bit more nuanced. checkEq drains the channel, and then
 * it makes sure that the event that was read from the channel is identical to the
 * event that is passed in for comparison.
 * If they are different, an error is thrown, and the test fails.
 */
func TestPutMessageEq(t *testing.T) {
	// Check that we read the expected event on the channel
	source := "us.me.lewiston"
	event := "forage-bagels"
	router := setup(checkEq(t, &itslog.Event{Source: source, Event: event}))
	apitest.New().
		Handler(router).
		Put(fmt.Sprintf("/v1/event/%s/%s", source, event)).
		Expect(t).
		Status(http.StatusOK).
		End()
}

/*
 * This final test checks that we are storing things to SQLite.
 * It involved adding two contrived functions to the SQLC code.
 * That functions/querys, TestEventPairExists and TestDictionaryPairExists,
 * makes sure that a given pair of values are present in both tables.
 */
func TestPutMessageToDb(t *testing.T) {
	storage := &sqlite.SqliteStorage{
		Path: ":memory:",
	}
	viper.Set("app.hash_seed", 42)
	err := storage.Init()

	if err != nil {
		panic(err)
	}

	// FIXME: add these constants to the configuration
	consumer := func(ch_evt chan *itslog.Event) {
		ch_eb := make(chan csp.EventBuffers)
		go csp.Enqueue(ch_evt, ch_eb, 1, 1)
		go csp.FlushBuffersOnce(ch_eb, storage)
	}

	// Check that we read the expected event on the channel
	source := "us.me.lewiston"
	event := "forage-bagels"
	router := setup(consumer)
	apitest.New().
		Handler(router).
		Put(fmt.Sprintf("/v1/event/%s/%s", source, event)).
		Expect(t).
		Status(http.StatusOK).
		End()

	// Wait for Enqueue to flush the buffer w/ a 1-second timeout
	time.Sleep(2 * time.Second)
	// Check if the value is in both the events and dictionary tables.
	result := storage.TestEventExists(source, event)
	if result != 1 {
		t.Errorf("PairExists: %d", result)
	}
}
