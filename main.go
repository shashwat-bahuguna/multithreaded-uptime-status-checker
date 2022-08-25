package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"http-status-checker/requesthandler"
)

var URL string = "localhost:8080" // URL+Port of API

const N int = 10                      // Pool Size
const Check_Time_Period = time.Minute // Time Period of scheduler thread
const Indicate_Routines = true        // Verbose Routines

func main() {

	requesthandler.Time_Period = Check_Time_Period
	requesthandler.VERBOSE_ROUTINES = Indicate_Routines
	requesthandler.PoolSize = N

	fmt.Println("Starting Server at", URL)
	http.HandleFunc("/websites", requesthandler.HandleRequest)

	if err := http.ListenAndServe(URL, nil); err != nil {
		log.Printf("Error Encountered while enabling server: %v.\n", err)
	}
}

/*
 POST REQUEST: curl -X POST localhost:8080/websites -d '{"websites": ["google.com", "yahoo.com", "abcd.com"]}'
			   curl -X POST localhost:8080/websites -d '{"websites": ["google.com", "abcd.com"]}'

 GET REQUEST:  curl "localhost:8080/websites"
 		       curl 'localhost:8080/websites?name=google.com'
			   curl 'localhost:8080/websites?name=google.com&name=yahoo.com'
			   curl 'localhost:8080/websites?name=facebook.com'
			   curl 'localhost:8080/websites?name=abcd.com'

*/
