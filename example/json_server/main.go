package main

import (
	"fmt"
	"github.com/streamingfast/dhttp/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/streamingfast/derr"
	"github.com/streamingfast/dhttp"
	"github.com/streamingfast/logging"
	"github.com/streamingfast/validator"
	"go.uber.org/zap"
)

var zlog, _ = logging.ApplicationLogger("example", "github.com/streamingfast/dhttp/example/json_server", logging.WithLogLevelSwitcherServerAutoStart())

func main() {
	router := mux.NewRouter()

	// Test with 'curl http://localhost:8080/healthz'
	router.Path("/healthz").Handler(dhttp.JSONHandler(getHealth))

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(middleware.NewCORSMiddleware(".*"))
	apiRouter.Use(middleware.NewTracingLoggingMiddleware(zlog))
	apiRouter.Use(middleware.NewLogRequestMiddleware(zlog))

	// Test with 'curl http://localhost:8080/api/v1/todos?user=john'
	apiRouter.Methods("GET").Path("/todos").Handler(dhttp.JSONHandler(getTodos))

	// Test with "curl -X PUT -d '{"id": "abc"}' http://localhost:8080/api/v1/todo"
	apiRouter.Methods("PUT").Path("/todos").Handler(dhttp.JSONHandler(putTodo))

	errorLogger, err := zap.NewStdLogAt(zlog, zap.ErrorLevel)
	if err != nil {
		panic(fmt.Errorf("unable to create error logger: %w", err))
	}

	serv := &http.Server{
		Addr:     "0.0.0.0:8080",
		Handler:  router,
		ErrorLog: errorLogger,
	}

	go func() {
		zlog.Info("serving HTTP", zap.String("listen_addr", "0.0.0.0:8080"))
		zlog.Info("endpoints")
		zlog.Info(" curl http://localhost:8080/healthz")
		zlog.Info(" curl http://localhost:8080/api/v1/todos?user=john")
		zlog.Info(` curl -X PUT -d '{"id": "abc"}' http://localhost:8080/api/v1/todos`)

		// FIXME: Drain connection when app is terminating as a finalizer step
		serv.ListenAndServe()
	}()

	// Wait until Ctrc-C is hit, in your own application, it should be tied to lifecycle
	// like a shutter.Shutter.
	<-derr.SetupSignalHandler(500 * time.Millisecond)
}

// Would normally go in `get_health.go` file

type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

func getHealth(r *http.Request) (out interface{}, err error) {
	return HealthResponse{Healthy: true}, nil
}

// Would normally go in `get_todos.go` file

type GetTodosParams struct {
	User string `schema:"user"`
}

type GetTodosResponse struct {
	IDs []string `json:"ids"`
}

func getTodos(r *http.Request) (out interface{}, err error) {
	ctx := r.Context()
	request := GetTodosParams{}
	err = dhttp.ExtractRequest(ctx, r, &request, dhttp.NewRequestValidator(validator.Rules{
		"user": []string{"required"},
	}))
	if err != nil {
		return nil, err
	}

	logging.Logger(ctx, zlog).Info("getting todo from request", zap.String("user", request.User))

	return GetTodosResponse{IDs: []string{request.User}}, nil
}

// Would normally go in `put_todo.go` file

type PutTodoRequest struct {
	ID string `json:"id"`
}

type TodosResponse struct {
	IDs []string `json:"ids"`
}

func putTodo(r *http.Request) (out interface{}, err error) {
	ctx := r.Context()
	logger := logging.Logger(ctx, zlog)

	request := PutTodoRequest{}
	err = dhttp.ExtractJSONRequest(ctx, r, &request, dhttp.NewJSONRequestValidator(validator.Rules{
		"id": []string{"required"},
	}))
	if err != nil {
		return nil, err
	}

	// You got here a specific logger for this request, contains a TraceID as well as any fields
	// you have added.
	logger.Debug("adding new todo in backend")

	return &TodosResponse{IDs: []string{request.ID}}, nil
}
