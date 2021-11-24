package dhttp_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/streamingfast/derr"
	"github.com/streamingfast/dhttp"
	"github.com/streamingfast/logging"
	"github.com/streamingfast/validator"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()
var _ = logging.ApplicationLogger("dhttp", "github.com/streamingfast/dhttp_test", &zlog)

func Example_JSONServer() {
	router := mux.NewRouter()

	// Test with 'curl http://localhost:8080/healthz'
	router.Path("/healthz").Handler(dhttp.JSONHandler(getHealth))

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(dhttp.NewCORSMiddleware(".*"))
	apiRouter.Use(dhttp.AddOpenCensusMiddleware)
	apiRouter.Use(dhttp.AddLoggerMiddleware)
	apiRouter.Use(dhttp.LogMiddleware)
	apiRouter.Use(dhttp.AddTraceMiddleware)

	// Test with 'curl http://localhost:8080/api/v1/todos?user=john'
	apiRouter.Methods("GET").Path("/todos").Handler(dhttp.JSONHandler(getTodos))

	// Test with "curl -X PUT -d '{"id": "abc"}' http://localhost:8080/api/v1/todo"
	apiRouter.Methods("PUT").Path("/todos").Handler(dhttp.JSONHandler(putTodo))

	errorLogger, err := zap.NewStdLogAt(zlog, zap.ErrorLevel)
	if err != nil {
		panic(fmt.Errorf("unable to create error logger: %w", err))
	}

	server := &http.Server{
		Addr:     "0.0.0.0:8080",
		Handler:  router,
		ErrorLog: errorLogger,
	}

	go func() {
		zlog.Info("serving HTTP", zap.String("listen_addr", "0.0.0.0:8080"))

		// FIXME: Drain connection when app is terminating as a finalizer step
		server.ListenAndServe()
	}()

	// Wait until Ctrc-C is hit, in your own application, it should be tied to lifecycle
	// like a shutter.Shutter.
	<-derr.SetupSignalHandler(500 * time.Millisecond)

	// Output: Completed
	fmt.Println("Completed")
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
	IDs []string
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

	request := PutTodoRequest{}
	err = dhttp.ExtractJSONRequest(ctx, r, &request, dhttp.NewJSONRequestValidator(validator.Rules{
		"id": []string{"required"},
	}))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
