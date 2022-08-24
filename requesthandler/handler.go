package requesthandler

import (
	"context"
	"encoding/json"
	"fmt"
	"http-status-checker/poolmanager"   // Manages the routine pool of the declared size
	"http-status-checker/statuschecker" // Checks the status of the website using an interface
	"log"
	"net/http"
	"time"
)

// To unmarshal incoming data in post request
type postdata struct {
	Websites []string `json:"websites"`
}

var PoolSize int = 10             // Thread Pool Size (These many threads will always be active)
var VERBOSE_ROUTINES bool = false // Whether to indicate beginning/ending of a goroutine
var Time_Period = time.Minute     // Time Period of website checker cycle

var mp map[string]string = nil          // Map storing current status of all websites
var global_pool *poolmanager.Pool = nil // Pointer to latest routine Pool

/**
 * Parse map into json string byte, returns errors if incurred
 * @param data -  map from string to string
 */
func parseMaptoJSON(data map[string]string) ([]byte, error) {
	if data == nil {
		return []byte{}, fmt.Errorf("Website list not declated yet! Please declare website list with a post request first.")
	}
	return json.Marshal(data)
}

/**
 * Processing the list of websites in a get request, returns the json byte array and error if incurred
 * @param websites - list of requested websites
 */
func process_get_request(websites []string) ([]byte, error) {

	data := map[string]string{}
	var ok bool = false
	for _, website := range websites {
		data[website], ok = mp[website]
		if ok != true {
			return []byte{}, fmt.Errorf("Unsupported Website Provided as Request Argument, %v", website)
		}
	}
	return parseMaptoJSON(data)
}

/**
 * To handle incoming requests at /websites endpoint. Handles both get and post requests.
 */
func HandleRequest(rw http.ResponseWriter, r *http.Request) {
	// s REQUEST: curl -X POST localhost:8080/websites -d '{"websites": ["google.com", "yahoo.com", "abcd.com"]}

	// GET REQUEST:  curl "localhost:8080/websites"
	//				 curl 'localhost:8080/websites?name="google.com"'

	if r.Method == http.MethodPost {
		fmt.Println("Received POST Request")
		decoder := json.NewDecoder(r.Body)
		var data postdata
		decoder.Decode(&data)

		mp = make(map[string]string)

		if global_pool != nil {
			// If a pool currently active, quit it before starting new pool
			global_pool.Quit_pool()
		}

		global_pool = poolmanager.Createpool(PoolSize)
		mp = map[string]string{}

		for _, website := range data.Websites {
			mp[website] = "NOT CHECKED YET"
		}

		global_pool.Pool_start(routinefunc)
		global_pool.Pool_Scheduler_Start(Time_Period, data.Websites)

	} else if r.Method == http.MethodGet {
		fmt.Println("Received GET Request")

		params := r.URL.Query()
		websites, ok := params["name"]

		fmt.Println("Params:", params)

		var err error
		var outp []byte

		if ok {
			outp, err = process_get_request(websites)
		} else {
			outp, err = parseMaptoJSON(mp)
		}

		if mp == nil {
			err = fmt.Errorf("Website list not declated yet! Please declare website list with a post request first.")
		}

		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte(err.Error()))
		} else {
			log.Println("Processed Output: ", string(outp))
			rw.Write(outp)
		}
	}
}

/**
 * To update the status of a given website in the stored map.
 * @param hostname hostname of target website
 * @param ctx context to be passed to get request
 * @param checker object supporting statuschecker interface
 */
func updateStatus(hostname string, routinepool *poolmanager.Pool, ctx context.Context, checker statuschecker.StatusChecker) (err error) {

	var status bool
	status, err = checker.Check(ctx, hostname)

	routinepool.Mutex.Lock()
	if routinepool.IsActive == true {
		if status == true {
			log.Println("Is Up")
			mp[hostname] = "UP"
		} else {
			log.Println("Is Down")
			mp[hostname] = "DOWN"
		}
	}
	routinepool.Mutex.Unlock()

	return
}

/**
 * Function to be executed by each go routine. Keeps the thread active until killed by the main thread.
 * @param {string} routinepool - Pool of goroutine threads
 */
func routinefunc(routinepool *poolmanager.Pool) {
	if VERBOSE_ROUTINES == true {
		log.Println("Routine Started")
		defer log.Println("Routine Ended")
	}

	for routinepool.IsActive {
		select {
		case <-routinepool.Quit:
			break

		case hostname, _ := <-routinepool.Tasks_Chan:
			log.Println("Starting Check for", hostname)
			err := updateStatus(hostname, routinepool, context.Background(), statuschecker.HttpChecker{})
			if err != nil {
				log.Println("Error while checking website status: ", err)
			}
		}
	}
}
