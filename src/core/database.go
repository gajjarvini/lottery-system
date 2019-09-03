package main

import (
	"database/sql"
	"errors"
	"fmt"
	log "logger"
	"net/http"

	"github.com/mattn/go-sqlite3"
)

/*createTableInDB adds the newly created ticket in database*/
func createTableInDB() (err error) {
	var sqlite3Driver string
	sql.Register(sqlite3Driver, &sqlite3.SQLiteDriver{})
	database, err = sql.Open(sqlite3Driver, "./lotteryTickets.db")
	if err != nil {
		log.Error.Println("Error opening DB, err", err)
		return err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS ticket (id INTEGER NOT NULL PRIMARY KEY, status BOOL)")
	if err != nil {
		log.Error.Println("Error creating table", err)
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		log.Error.Println("Error creating table", err)
		return err
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS ORDERS ( ID2 int,   value int,   FOREIGN KEY (ID2) REFERENCES ticket (id));")
	if err != nil {
		log.Error.Println("Error creating table", err)
	}
	statement.Exec()
	if err != nil {
		log.Error.Println("Error creating table", err)
	}
	return nil
}

/*fetchTicketFromDB searches and fetches a ticket from the data base for the provided id*/
func fetchTicketFromDB(id int) (t1 ticket, err error) {
	rows, err := database.Query("select C.id, C.status, O.value from ticket C inner join ORDERS O on C.id = O.ID2 where C.id = $1", id)
	if err != nil {
		log.Error.Println("error querrying", err)
		return t1, err
	}
	var arr []int
	var val int
	for rows.Next() {
		err = rows.Scan(&(t1.ID), &(t1.Status), &val)
		if err != nil {
			log.Error.Println(err)
			return t1, err
		}
		arr = append(arr, val)
	}
	if t1.ID == 0 {
		return t1, errors.New("Ticket not found")
	}
	count := configvalues.NumbersInALine
	for i := 0; i <= (len(arr) - count); i = i + count {
		t1.Lines = append(t1.Lines, arr[i:i+count])
	}
	return t1, nil
}

/*getTicketlistFromDB  responds the client with all the tickets in the database*/
func getTicketlistFromDB(w http.ResponseWriter) error {
	var ticketIds []int
	rows, err := database.Query("select C.id from ticket C")
	if err != nil {
		log.Error.Println("error querrying", err)
		return err
	}
	var val int
	for rows.Next() {
		err = rows.Scan(&val)
		if err != nil {
			log.Error.Println(err)
			return err
		}
		ticketIds = append(ticketIds, val)
	}
	//The following code sends complete tickets
	/*var tickets []ticket
	var t ticket
	for _, id := range ticketIds {
		t, err = fetchTicketFromDB(id)
		if err != nil {
			log.Error.Println("Erro getting list of tickets, err: ", err)
			return err
		}
		tickets = append(tickets, t)
	}
	*/
	fmt.Fprintf(w, "%+v", ticketIds)
	return nil
}

/*insertNewTicketIntoDB pushes the newly created ticket into database*/
func insertNewTicketIntoDB(t ticket) (id int, err error) {
	statement, err := database.Prepare("INSERT INTO ticket (status) VALUES (?)")
	if err != nil {
		log.Error.Println(err)
		return id, err
	}
	result, err := statement.Exec(t.Status)
	if err != nil {
		log.Error.Println(err)
		return id, err
	}
	id64, err := result.LastInsertId()
	if err != nil {
		log.Error.Println(err)
	}
	id = int(id64)
	err = addLinesToDB(id, t.Lines)
	if err != nil {
		log.Error.Println(err)
		return id, err
	}
	return id, nil
}

/*addLinesToDB appends n number of lines to the ticket with provided id in database*/
func addLinesToDB(id int, lines [][]int) error {
	statement, _ := database.Prepare("INSERT INTO ORDERS (ID2, value) VALUES (?, ?)")
	for _, line := range lines {
		for _, value := range line {
			_, err := statement.Exec(id, value)
			if err != nil {
				log.Error.Println(err)
				return err
			}
		}
	}
	return nil
}

/*updateTicketStatusinDB updates result status of ticket addLinesToDB in database*/
func updateTicketStatusinDB(id int) error {
	statement, _ := database.Prepare("update ticket SET status = 1 where ticket.id = ?")
	_, err := statement.Exec(id)
	if err != nil {
		log.Error.Println(err)
		return err
	}
	return nil
}
