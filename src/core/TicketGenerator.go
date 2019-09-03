package main

import (
	"fmt"
	log "logger"
	"math/rand"
	"net/http"
	"sort"
)

type ticket struct {
	ID     int     `json:"ID"`
	Lines  [][]int `json:"Lines"`
	Status bool    `json:"StatusChecked"`
}

type lineResult struct {
	Line       []int `json:"Line"`
	LineResult int   `json:"LineResult"`
}

func (l *lineResult) calculateLineResult(line []int) {
	l.Line = line
	l.LineResult = 0
	if line[0]+line[1]+line[2] == 2 {
		l.LineResult = 10
	} else if line[0] == line[1] && line[1] == line[2] {
		l.LineResult = 5
	} else if line[0] != line[1] && line[0] != line[2] {
		l.LineResult = 1
	}
}

func sendTicketResult(w http.ResponseWriter, id int) error {
	t, err := fetchTicketFromDB(id)
	if err != nil {
		log.Error.Println("Failed to fetch ticket ", id, "from DB while fetching ticket result")
		return err
	}
	log.Info.Println("Generating ticket result for id: ", id)
	var arrLineResult []lineResult
	var result lineResult
	for _, line := range t.Lines {
		result.calculateLineResult(line)
		arrLineResult = append(arrLineResult, result)
	}
	err = updateTicketStatusinDB(id)
	if err != nil {
		log.Error.Println("Error setting status in database, err: ", err)
		return err
	}

	sort.SliceStable(arrLineResult, func(i, j int) bool {
		return arrLineResult[i].LineResult < arrLineResult[j].LineResult
	})

	fmt.Fprintf(w, "ID:%d %+v", t.ID, arrLineResult)
	return nil
}

func generateNewTicket(w http.ResponseWriter, count int) error {
	var t1 ticket
	var err error
	for i := 0; i < count; i++ {
		t1.Lines = append(t1.Lines, generateRandomLine())
	}
	t1.Status = false
	t1.ID, err = insertNewTicketIntoDB(t1)
	if err != nil {
		log.Error.Println("Failed to create ticket into DB")
		return err
	}
	fmt.Fprintf(w, "%+v", t1)
	return nil
}

func generateRandomLine() []int {
	var line []int
	for i := 0; i < configvalues.NumbersInALine; i++ {
		line = append(line, rand.Intn(3))
	}
	return line
}

func addLinesToTicket(w http.ResponseWriter, ID, count int) error {
	t, err := fetchTicketFromDB(ID)
	if err != nil {
		return err
	}
	if t.Status == true {
		fmt.Fprintf(w, "Ticket result has already been generated. Cannot append more lines")
		return nil
	}
	length := len(t.Lines)
	for i := 0; i < count; i++ {
		t.Lines = append(t.Lines, generateRandomLine())
	}
	err = addLinesToDB(ID, t.Lines[length:])
	if err != nil {
		log.Error.Println("error adding additional lines to DB ", err)
		return err
	}
	fmt.Fprintf(w, "%+v", t)
	return nil
}
