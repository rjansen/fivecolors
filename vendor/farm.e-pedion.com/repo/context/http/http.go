package http

import (
	"context"
	haki "farm.e-pedion.com/repo/context"
	"farm.e-pedion.com/repo/context/media/json"
	"farm.e-pedion.com/repo/logger"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

//ResponseWriter is a wrapper function to store status and body length of the request
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 200 if the response has
	// not been written (as this is the default response code in net/http)
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
}

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		ResponseWriter: w,
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(s int) {
	w.status = s
	w.ResponseWriter.WriteHeader(s)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.Written() {
		// The status will be 200 if WriteHeader has not been called yet
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.status != 0
}

func (w *responseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		if !w.Written() {
			// The status will be 200 if WriteHeader has not been called yet
			w.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

//SimpleHTTPHandler is a contract for fast http handlers
type SimpleHTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

//HTTPHandlerFunc is a function to handle fasthttp requrests
type HTTPHandlerFunc func(http.ResponseWriter, *http.Request) error

//HandleRequest is the contract with HTTPHandler interface
func (h HTTPHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

//HTTPHandler is a contract for fast http handlers
type HTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

//Handler wraps a library handler func nto a http handler func
func Handler(handler HTTPHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, h *http.Request) {
		handler(w, h)
	}
}

func errorHandle(handler HTTPHandlerFunc, w http.ResponseWriter, r *http.Request) error {
	if err := handler(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

//Error wraps the provided HTTPHandlerFunc with exception control
func Error(handler HTTPHandlerFunc) HTTPHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return errorHandle(handler, w, r)
	}
}

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func logHandle(handler HTTPHandlerFunc, w http.ResponseWriter, r *http.Request) error {
	tid := uuid.NewV4().String()
	r = r.WithContext(context.WithValue(r.Context(), "tid", tid))
	start := time.Now()
	logger.Info("contex.Request",
		logger.String("tid", tid),
		logger.String("method", r.Method),
		logger.String("path", r.URL.Path),
	)
	logger.Debug("context.Context",
		logger.Bool("ctxIsNil", r.Context() == nil),
	)
	r = r.WithContext(context.WithValue(r.Context(), "log", logger.Get()))
	rw := NewResponseWriter(w)
	var err error
	if err = handler(rw, r); err != nil {
		logger.Error("contex.LogHandler.Error",
			logger.String("tid", tid),
			logger.String("method", r.Method),
			logger.String("path", r.URL.Path),
			logger.Err(err),
		)
	}
	response := rw.(ResponseWriter)
	logger.Info("context.Response",
		logger.String("tid", tid),
		logger.String("method", r.Method),
		logger.String("path", r.URL.Path),
		logger.String("status", http.StatusText(response.Status())),
		logger.Int("size", response.Size()),
		logger.Duration("requestTime", time.Since(start)),
	)
	return err
}

//Log wraps the provided HTTPHandlerFunc with access logging control
func Log(handler HTTPHandlerFunc) HTTPHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return logHandle(handler, w, r)
	}
}

type LogHandler func(http.ResponseWriter, *http.Request) error

func (h LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logHandle(HTTPHandlerFunc(h), w, r)
}

//ReadByContentType reads data from context using the Content-Type header to define the media type
func ReadByContentType(r *http.Request, data interface{}) error {
	contentType := r.Header.Get(haki.ContentTypeHeader)
	switch {
	case strings.Contains(contentType, json.ContentType):
		return ReadJSON(r, data)
	// case strings.Contains(contentType, proto.ContentType):
	// 	return ReadProtoBuff(r, data)
	default:
		return haki.ErrInvalidContentType
	}
}

//WriteByAccept writes data to context using the Accept header to define the media type
func WriteByAccept(w http.ResponseWriter, r *http.Request, status int, result interface{}) error {
	contentType := r.Header.Get(haki.AcceptHeader)
	switch {
	case strings.Contains(contentType, json.ContentType):
		return JSON(w, status, result)
	// case bytes.Contains(contentType, []byte(proto.ContentType)):
	// 	return ProtoBuff(ctx, status, result)
	default:
		return haki.ErrInvalidAccept
	}
}

//ReadJSON unmarshals from provided context a json media into data
func ReadJSON(r *http.Request, data interface{}) error {
	if err := json.Unmarshal(r.Body, data); err != nil {
		return err
	}
	return nil
}

func Bytes(w http.ResponseWriter, status int, result []byte) error {
	w.Header().Set(haki.ContentTypeHeader, "application/octet-stream")
	w.WriteHeader(status)
	_, err := w.Write(result)
	if err != nil {
		return err
	}
	return nil
}

func JSON(w http.ResponseWriter, status int, result interface{}) error {
	w.Header().Set(haki.ContentTypeHeader, json.ContentType)
	w.WriteHeader(status)
	if err := json.Marshal(w, result); err != nil {
		return err
	}
	return nil
}

func Status(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}

func Err(w http.ResponseWriter, err error) error {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return err
}

type BaseHandler struct {
}

func (h BaseHandler) JSON(w http.ResponseWriter, status int, result interface{}) error {
	return JSON(w, status, result)
}

func (h BaseHandler) Status(w http.ResponseWriter, status int) error {
	return Status(w, status)
}

func (h BaseHandler) Err(w http.ResponseWriter, err error) error {
	return Err(w, err)
}
