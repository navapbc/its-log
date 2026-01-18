package serve

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/internal/csp"
	"github.com/jadudm/its-log/internal/fsdb"
	"github.com/jadudm/its-log/internal/itslog"
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

func setup(consumer func(chan *itslog.Event)) (*gin.Engine, string) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	var ch_evt_out = make(chan *itslog.Event)
	// This drains the channel so we don't have to worry about
	// it as part of the testing.
	go consumer(ch_evt_out)

	router := gin.Default()
	apiV1 := router.Group("/v1")
	permissions := []itslog.PermissionType{itslog.Log, itslog.Test}
	apiV1.Use(AuthMiddleWare(permissions))
	apiV1.PUT("se/:source/:event", Event(itslog.SE, ch_evt_out))
	// Mock the env setup.
	// This implies we have a JSON structure in an APIKEY variable.
	key := "12345678901234561234567890123456"
	m := gin.H{
		"app_id":     "test",
		"key_id":     "test",
		"permission": "log",
		"key":        key,
	}
	keystr, _ := json.Marshal(m)
	os.Setenv("ITSLOG_APIKEY_TEST", string(keystr))
	// Read the API keys in, so they can be found by the middleware.
	itslog.GetApiKeys()
	log.Printf("found %d keys", len(itslog.LiveKeys))
	// Return the router and the API key
	return router, key
}

/*
 * This test checks whether or not the PUT is picked up by the framework.
 * There is no checking of the value. It just makes sure the phone was picked up.
 * blackHole makes sure the channel is drained from the API handler, or otherwise
 * it will hang forever, waiting for the channel communication to terminate.
 */
func TestPutMessage(t *testing.T) {
	router, key := setup(blackHole)

	apitest.New().
		Handler(router).
		Put("/v1/se/us.me.lewiston/forage-bagels").
		Headers(map[string]string{"x-api-key": key}).
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
	router, key := setup(checkEq(t, &itslog.Event{Source: source, Event: event}))
	apitest.New().
		Handler(router).
		Put(fmt.Sprintf("/v1/se/%s/%s", source, event)).
		Headers(map[string]string{"x-api-key": key}).
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
// FIXME: Using RAM won't work, because the memory DB gets
// recreated further down. So, the DB created here, and the DB
// that data is stored to are different. I think. Need to rethink this.
func TestPutMessageToDb(t *testing.T) {
	viper.Set("hash.seed", 42)
	s := &fsdb.SqliteStorage{
		Kind:  fsdb.InMemory,
		AppId: "test",
	}
	err := s.Init()
	if err != nil {
		t.Fatalf("error: %s\n", err.Error())
	}

	// FIXME: add these constants to the configuration
	consumer := func(ch_evt chan *itslog.Event) {
		ch_eb := make(chan csp.EventBuffers)
		go csp.Enqueue(ch_evt, ch_eb, 1, 1)
		go csp.FlushBuffersOnce(s, ch_eb)
	}

	// Check that we read the expected event on the channel
	source := "us.me.lewiston"
	event := "forage-bagels"
	router, key := setup(consumer)
	apitest.New().
		Handler(router).
		Put(fmt.Sprintf("/v1/se/%s/%s", source, event)).
		Headers(map[string]string{"x-api-key": key}).
		Expect(t).
		Status(http.StatusOK).
		End()

	// Wait for Enqueue to flush the buffer w/ a 1-second timeout
	time.Sleep(2 * time.Second)
	// Flushing closes the DB.
	// With in-memory testing, this erases the DB.

	// Check if the value is in both the events and dictionary tables.
	result := s.TestEventExists(source, event)
	if result != 1 {
		t.Errorf("PairExists: %d", result)
	}
}
