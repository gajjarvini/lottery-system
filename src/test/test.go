package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Enter server ip")
	var serverIp string
	fmt.Scanf("%s",&serverIp)
	for {
		fmt.Println("\n\n")
		fmt.Println("1. Get Ticket by id")
		fmt.Println("2. Generate new tickets")
		fmt.Println("3. Get status of a ticket ")
		fmt.Println("4. Add lines to ticket")
		fmt.Println("5. Get list of all the tickets")
		var method int
		fmt.Scanf("%d", &method)
		switch method {
		case 1:
			var id string
			fmt.Println("Enter Ticket id")
			fmt.Scanf("%s", &id)
			resp, err := http.Get("http://" + serverIp + ":8090/ticket/" + id)
			if err != nil {
				fmt.Println("Error is ", err)
				return
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(respBody))
		case 2:
			fmt.Println("Enter count for lines to be added")
			var count string
			fmt.Scanf("%s", &count)
			resp, err := http.Post("http://" + serverIp + ":8090/ticket/"+count, "text/plain", bytes.NewBuffer(nil))
			if err != nil {
				fmt.Println("Error is ", err)
				return
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(respBody))
		case 3:
			var id string
			fmt.Println("Enter ticket id")
			fmt.Scanf("%s", &id)
			req, err := http.NewRequest("PUT", "http://" + serverIp + ":8090/status/"+id, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			req.Header.Set("Content-Type", "plain/text")

			// Do the request
			timeout := time.Duration(30 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			response, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer response.Body.Close()

			respBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(respBody))

		case 4:
			var id, count string
			fmt.Println("Enter ticket id and count of lines")
			fmt.Scanf("%s \n %s", &id, &count)
			req, err := http.NewRequest("PUT", "http://" + serverIp + ":8090/ticket/"+id+"/"+count, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			req.Header.Set("Content-Type", "plain/text")

			// Do the request
			timeout := time.Duration(30 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			response, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer response.Body.Close()

			respBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(respBody))
		case 5:
			resp, err := http.Get("http://" + serverIp + ":8090/ticket")
			if err != nil {
				fmt.Println("Error is ", err)
				return
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(respBody))

		}
	}
}
