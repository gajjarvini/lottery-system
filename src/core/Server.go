package main

import (
	"database/sql"
	"errors"
	"fmt"
	log "logger"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type configuration struct {
	PortNumber string `json:"PortNumber"`

	//count of numbers in one line
	NumbersInALine int `json:"NumbersInALine"`
}

//configvalues holds configuration values parsed from conf.json
var configvalues configuration
var database *sql.DB

func startHTTPSServer() error {
	r := mux.NewRouter()

	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/status/{id:[0-9]+}", handleGetStatusRequest).Methods("PUT")

	s := r.PathPrefix("/ticket").Subrouter()
	s.HandleFunc("", handleGetTicketListRequest).Methods("GET")
	s.HandleFunc("/", handleGetTicketListRequest).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}", handleGetTicketByID).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}/{count:[0-9]+}", handleAddNewLinesRequest).Methods("PUT")
	s.HandleFunc("/{count:[0-9]+}", handleCreateTicketRequest).Methods("POST")

	log.Info.Println("Starting server on port number :", configvalues.PortNumber)
	fmt.Println("Starting server on port number :", configvalues.PortNumber)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + configvalues.PortNumber,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Error.Println("Unable to Listen and Serve. Error: ", err)
		return err
	}
	return errors.New("http server is stopped")
}

func main() {

	err := OpenLoggerFile()
	if err != nil {
		os.Exit(1)
	}

	fileHandle := GetFilehandle()
	if fileHandle != nil {
		defer fileHandle.Close()
	} else {
		log.Error.Println("fileHandle should not be nil")
		os.Exit(1)
	}

	log.Info.Println("log file is crteated ")

	err = GetConfigFileValues()
	if err != nil {
		log.Error.Println("Parsing of config file is failed err:", err)
		os.Exit(1)
	}

	log.Info.Println("parsing config values is done ")
	err = createTableInDB()
	if err != nil {
		log.Error.Println("Error creating table, err :", err)
		os.Exit(1)
	}
	err = startHTTPSServer()
	if err != nil {
		log.Error.Println("error in http server")
	}
}
