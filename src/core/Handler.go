package main

import (
	"fmt"
	log "logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*handleGetTicketByID serves the requested ticket*/
func handleGetTicketByID(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Serving a get ticket by id method")
	vars := mux.Vars(r)
	id := vars["id"]
	ID, err := strconv.Atoi(id)
	if err != nil {
		log.Error.Println("Error in Atoi operation, err :", err)
		return
	}
	log.Info.Println("Fethching ticket with id : ", ID)
	//getTicketById(strconv(strconv.Atoi(count)))
	t1, err := fetchTicketFromDB(ID)
	if err != nil {
		log.Error.Println("Error fetching ticket from DB, Ticket ID", ID, "err :", err)
		fmt.Fprintf(w, "%+v", err)
		return
	}
	fmt.Fprintf(w, "%+v", t1)
	log.Info.Println("Served ticket with id : ", ID)
}

/*handleGetTicketListRequest serves the list of all the ticket IDs*/
func handleGetTicketListRequest(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Serving Get Tickets  list request")

	err := getTicketlistFromDB(w)
	if err != nil {
		log.Error.Println("Error fetching ticket list from DB, err :", err)
		return
	}
	log.Info.Println("Served the  ticket list ")
}

/*handleCreateTicketRequest creates and serves a new ticket*/
func handleCreateTicketRequest(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Serving create new ticket request")
	vars := mux.Vars(r)
	count := vars["count"]
	Count, err := strconv.Atoi(count) // pending
	if err != nil {
		log.Error.Println("Error in Atoi operation, err :", err)
		return
	}
	err = generateNewTicket(w, Count)
	if err != nil {
		log.Error.Println("handlePostRequest is failed")
		return
	}
	log.Info.Println("Successfully created a new ticket")
}

/*handleAddNewLinesRequest appends new lines to existing ticket*/
func handleAddNewLinesRequest(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Serving add lines to ticket request")
	vars := mux.Vars(r)
	id := vars["id"]

	ID, err := strconv.Atoi(id)
	if err != nil {
		log.Error.Println("Error in Atoi operation, err :", err)
		return
	}

	log.Info.Println("Adding lines to ticket id :", ID)

	count := vars["count"]
	Count, err := strconv.Atoi(count)
	if err != nil {
		log.Error.Println("Error in Atoi operation, err :", err)
		return
	}

	err = addLinesToTicket(w, ID, Count)
	if err != nil {
		log.Error.Println("Invalid ID")
		fmt.Fprintf(w, "%+v", err)
		return
	}

	log.Info.Println("Successfully appended new lines to the ticket id :", ID)
}

/*handleGetStatusRequest generates and serves the ticket result*/
func handleGetStatusRequest(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Serving get ticket result request")
	vars := mux.Vars(r)
	id := vars["id"]

	ID, err := strconv.Atoi(id)
	if err != nil {
		log.Error.Println("Error in Atoi operation, err :", err)
		return
	}

	err = sendTicketResult(w, ID)
	if err != nil {
		log.Error.Println("Invalid ID")
		fmt.Fprintf(w, "%+v", err)
		return
	}
	log.Info.Println("Sucessfully served ticket result")
}

/*defaultHandler handles the invalid requests*/
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Invalid path, request and path are :", r.Method, r.URL.Path)
	fmt.Fprintf(w, "Invalid Path")
}

func writeErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusLocked)
	fmt.Fprintf(w, "Failed")
}
