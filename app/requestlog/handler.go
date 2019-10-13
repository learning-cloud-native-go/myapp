// Copyright 2019 The Go Cloud Development Kit Authors
// Modified by Dumindu Madunuwan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Source: https://github.com/google/go-cloud/blob/master/server/requestlog/requestlog.go
package requestlog

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"myapp/util/logger"
)

type Handler struct {
	handler http.Handler
	logger  *logger.Logger
}

func NewHandler(h http.HandlerFunc, l *logger.Logger) *Handler {
	return &Handler{
		handler: h,
		logger:  l,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	le := &logEntry{
		ReceivedTime:      start,
		RequestMethod:     r.Method,
		RequestURL:        r.URL.String(),
		RequestHeaderSize: headerSize(r.Header),
		UserAgent:         r.UserAgent(),
		Referer:           r.Referer(),
		Proto:             r.Proto,
		RemoteIP:          ipFromHostPort(r.RemoteAddr),
	}

	if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		le.ServerIP = ipFromHostPort(addr.String())
	}
	r2 := new(http.Request)
	*r2 = *r
	rcc := &readCounterCloser{r: r.Body}
	r2.Body = rcc
	w2 := &responseStats{w: w}

	h.handler.ServeHTTP(w2, r2)

	le.Latency = time.Since(start)
	if rcc.err == nil && rcc.r != nil {
		// If the handler hasn't encountered an error in the Body (like EOF),
		// then consume the rest of the Body to provide an accurate rcc.n.
		io.Copy(ioutil.Discard, rcc)
	}
	le.RequestBodySize = rcc.n
	le.Status = w2.code
	if le.Status == 0 {
		le.Status = http.StatusOK
	}
	le.ResponseHeaderSize, le.ResponseBodySize = w2.size()
	h.logger.Info().
		Time("received_time", le.ReceivedTime).
		Str("method", le.RequestMethod).
		Str("url", le.RequestURL).
		Int64("header_size", le.RequestHeaderSize).
		Int64("body_size", le.RequestBodySize).
		Str("agent", le.UserAgent).
		Str("referer", le.Referer).
		Str("proto", le.Proto).
		Str("remote_ip", le.RemoteIP).
		Str("server_ip", le.ServerIP).
		Int("status", le.Status).
		Int64("resp_header_size", le.ResponseHeaderSize).
		Int64("resp_body_size", le.ResponseBodySize).
		Dur("latency", le.Latency).
		Msg("")
}
