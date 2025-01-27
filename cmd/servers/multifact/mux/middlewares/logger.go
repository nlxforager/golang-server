package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type MiddleWare func(http.HandlerFunc) http.HandlerFunc

// Wrap
// .Wrap(f1,f2,f3) => f1 => f2 => f3
func (mw MiddleWare) Wrap(nexts ...MiddleWare) MiddleWare {
	for _, next := range nexts {
		mw = mw.wrap(next)
	}
	return mw
}

func (mw MiddleWare) wrap(next MiddleWare) MiddleWare {
	if mw == nil {
		return next
	}
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return next(mw(handlerFunc))
	}
}

var LogMiddleware MiddleWare = func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now() // as close as possible to receive time

		var sb strings.Builder
		sb.WriteString(currentTime.Format("2006-01-02 15:04:05.000") + "................................REQUEST................................\n")

		sb.WriteString(fmt.Sprintf("%s %s %s\n", r.Method, r.URL.String(), r.Proto))

		for name, values := range r.Header {
			for _, value := range values {
				sb.WriteString(fmt.Sprintf("%s: %s\n", name, value))
			}
		}

		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				sb.WriteString(fmt.Sprintf("Error reading body: %v\n", err))
				http.Error(w, "Unable to read request body", http.StatusInternalServerError)
				return
			}

			sb.WriteString(fmt.Sprintf("%s\n", string(bodyBytes)))

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status
			headers:        make(http.Header),
			body:           bytes.NewBuffer(nil),
		}

		handlerFunc(lrw, r)
		sb.WriteString(time.Now().Format("2006-01-02 15:04:05.000") + "................................RESPONSE................................\n")

		sb.WriteString(fmt.Sprintf("HTTP/?.? %d %s\n", lrw.statusCode, http.StatusText(lrw.statusCode)))

		for name, values := range lrw.Header() {
			for _, value := range values {
				sb.WriteString(fmt.Sprintf("%s: %s\n", name, value))
			}
		}

		sb.WriteString(fmt.Sprintf("%s\n", lrw.body.String()))

		endTime := time.Now() // as close as possible to send out time
		d := endTime.Sub(currentTime)

		sb.WriteString(endTime.Format("2006-01-02 15:04:05.000") + "................................END................................" + d.String())
		fmt.Println(sb.String())
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	headers    http.Header
	body       *bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	// Capture the response body
	lrw.body.Write(data)
	return lrw.ResponseWriter.Write(data)
}
