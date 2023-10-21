//go:build !binary_log
// +build !binary_log

package logger_test

import (
	"errors"
	"flag"
	"time"

	"github.com/rs/zerolog"

	"myapp/util/logger"
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

	return logger.New(true)
}

// Simple logging example using the Print function in the log package
// Note that both Print and Printf are at the debug log level by default
func ExampleLogger_Print() {
	l := setup()
	l.Print("hello world")

	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

// Simple logging example using the Printf function in the log package
func ExampleLogger_Printf() {
	l := setup()
	l.Printf("hello %s", "world")

	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

// Example of a log with no particular "level"
func ExampleLogger_Log() {
	l := setup()
	l.Log().Msg("hello world")

	// Output: {"time":1199811905,"message":"hello world"}
}

// Example of a conditional level based on the presence of an error.
func ExampleLogger_Err() {
	l := setup()
	err := errors.New("some error")
	l.Err(err).Msg("hello world")
	l.Err(nil).Msg("hello world")

	// Output: {"level":"error","error":"some error","time":1199811905,"message":"hello world"}
	// {"level":"info","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "trace")
func ExampleLogger_Trace() {
	l := setup()
	l.Trace().Msg("hello world")

	// Output: {"level":"trace","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "debug")
func ExampleLogger_Debug() {
	l := setup()
	l.Debug().Msg("hello world")

	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "info")
func ExampleLogger_Info() {
	l := setup()
	l.Info().Msg("hello world")

	// Output: {"level":"info","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "warn")
func ExampleLogger_Warn() {
	l := setup()
	l.Warn().Msg("hello world")

	// Output: {"level":"warn","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "error")
func ExampleLogger_Error() {
	l := setup()
	l.Error().Msg("hello world")

	// Output: {"level":"error","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "fatal")
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

// This example uses command-line flags to demonstrate various outputs
// depending on the chosen log level.
func Example() {
	l := setup()
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	l.Debug().Msg("This message appears only when log level set to Debug")
	l.Info().Msg("This message appears when log level set to Debug or Info")

	if e := l.Debug(); e.Enabled() {
		// Compute log output only if enabled.
		value := "bar"
		e.Str("foo", value).Msg("some debug message")
	}

	// Output: {"level":"info","time":1199811905,"message":"This message appears when log level set to Debug or Info"}
}
