package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mazzus/graphical-elevator-simulator/backend/src/elevator"
	"github.com/op/go-logging"
)

type webError struct {
	StatusCode int
	Message    string
	Cause      error
}

const module = "HTTPServer"

var log = logging.MustGetLogger("HTTP")

func HTTPServer(port int, safeElevator *elevator.SafeElevator) {
	router := mux.NewRouter()

	router.Path("/api/total").Methods("GET").HandlerFunc(requestLogger(Total(safeElevator)))
	router.Path("/api/button").Methods("POST").HandlerFunc(requestLogger(SetButton(safeElevator)))
	router.PathPrefix("/").Methods("GET").Handler(http.StripPrefix("/", http.FileServer(assetFS())))

	hostString := fmt.Sprintf(":%d", port)
	server := http.Server{
		Addr:           hostString,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Infof("HTTP Server started on port %v", port)
	log.Infof("The UI can be found by opening the following address in a browser: http://localhost:%v/index.html", port)
	if err := server.ListenAndServe(); err != nil {
		log.Error("http server returned an error!", err)
	}
}

func (err *webError) Error() string {
	if err.Cause != nil {
		return fmt.Sprintf("StatusCode: %d. msg: %s \nCaused by: %s", err.StatusCode, err.Message, err.Cause.Error())
	}

	return fmt.Sprintf("StatusCode: %d. msg: %s \nCaused by: nil", err.StatusCode, err.Message)
}

func errorHandler(f func(http.ResponseWriter, *http.Request) *webError) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			log.Warningf("Path: %s. Error: %s", r.RequestURI, err)
		}
	}
}

func requestLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Handling: %s", r.RequestURI)
		h(w, r)
		log.Debugf("Handled: %s", r.RequestURI)
	}
}

func allowCORS(methods []string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		if r.Method == "OPTIONS" {

			w.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ", "))

			return
		}
		h(w, r)
	}
}

func Total(safeElev *elevator.SafeElevator) func(http.ResponseWriter, *http.Request) {

	handlerFunction := func(w http.ResponseWriter, r *http.Request) *webError {

		encoder := json.NewEncoder(w)
		safeElev.Lock()
		err := encoder.Encode(safeElev.Elevator)
		safeElev.Unlock()
		if err != nil {
			return &webError{500, "Could not encode the elevator", err}
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		return nil
	}
	return errorHandler(handlerFunction)
}

func SetButton(safeElevator *elevator.SafeElevator) func(http.ResponseWriter, *http.Request) {
	handlerFunction := func(w http.ResponseWriter, r *http.Request) *webError {
		decoder := json.NewDecoder(r.Body)
		var body struct {
			Type  string
			Floor int
			Value bool
		}

		err := decoder.Decode(&body)
		if err != nil {
			return &webError{400, "Could not decode the request body", err}
		}
		switch body.Type {
		case "up":
			safeElevator.Lock()
			err = safeElevator.SetUpButton(body.Floor, body.Value)
			safeElevator.Unlock()
		case "down":
			safeElevator.Lock()
			err = safeElevator.SetDownButton(body.Floor, body.Value)
			safeElevator.Unlock()
		case "cabin":
			safeElevator.Lock()
			err = safeElevator.SetCabinButton(body.Floor, body.Value)
			safeElevator.Unlock()
		case "stop":
			safeElevator.Lock()
			safeElevator.SetStopButton(body.Value)
			safeElevator.Unlock()
		case "obstruction":
			safeElevator.Lock()
			safeElevator.SetObstructionButton(body.Value)
			safeElevator.Unlock()
		}
		if err != nil {
			return &webError{400, "Invalid button configuration", err}
		}
		return nil
	}
	return errorHandler(handlerFunction)
}
