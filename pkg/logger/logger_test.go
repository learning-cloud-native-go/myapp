//go:build !binary_log
// +build !binary_log

package logger_test

import (
	"errors"
	"flag"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"myapp/pkg/logger"
)

// setup would normally be an init() function, however, there seems
// to be something awry with the testing framework when we set the
// global Logger from an init()
func setup() *logger.Logger {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = ""
	// In order to always output a static time to stdout for these
	// examples to pass, we need to override zerolog.TimestampFunc
	// and log.Logger globals -- you would not normally need to do this
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 5, 0, time.UTC)
	}

	return logger.NewConsole(true)
}

func ExampleLogger_Print() {
	l := setup()
	l.Print("hello world")

	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

func ExampleLogger_Printf() {
	l := setup()
	l.Printf("hello %s", "world")

	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

func ExampleLogger_Log() {
	l := setup()
	l.Log().Msg("hello world")

	// Output: {"time":1199811905,"message":"hello world"}
}

func ExampleLogger_Debug() {
	l := setup()
	l.Debug().Msg("hello world")

	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

func ExampleLogger_Info() {
	l := setup()
	l.Info().Msg("hello world")

	// Output: {"level":"info","time":1199811905,"message":"hello world"}
}

func ExampleLogger_Warn() {
	l := setup()
	l.Warn().Msg("hello world")

	// Output: {"level":"warn","time":1199811905,"message":"hello world"}
}

func ExampleLogger_Error() {
	l := setup()
	l.Error().Msg("hello world")

	// Output: {"level":"error","time":1199811905,"message":"hello world"}
}

func ExampleLogger_Fatal() {
	l := setup()
	err := errors.New("A repo man spends his life getting into tense situations")
	service := "myservice"

	l.Fatal().
		Err(err).
		Str("service", service).
		Msgf("Cannot start %s", service)

	// Outputs: {"level":"fatal","time":1199811905,"error":"A repo man spends his life getting into tense situations","service":"myservice","message":"Cannot start myservice"}
}

func Example() {
	l := setup()
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	l.Debug().Msg("This message appears only when log level set to Debug")
	l.Info().Msg("This message appears when log level set to Debug or Info")

	if e := l.Debug(); e.Enabled() {
		value := "bar"
		e.Str("foo", value).Msg("some debug message")
	}

	// Output: {"level":"info","time":1199811905,"message":"This message appears when log level set to Debug or Info"}
}

func TestReqLogFromRequest(t *testing.T) {
	t.Parallel()

	l := setup()

	tests := []struct {
		name string
		req  *http.Request
		want string
	}{
		{
			"success",
			httptest.NewRequest(http.MethodGet, "https://api.geocodify.com/v2/reverse?api_key=xxxxx&lat=34.052&lng=-118.243", nil),
			"GET /v2/reverse?api_key=xxxxx&lat=34.052&lng=-118.243 HTTP/1.1\r\nHost: api.geocodify.com\r\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := l.ReqLogFromRequest(tt.req); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Logger.ReqLogFromRequest() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
