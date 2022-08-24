package poolmanager

import (
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	Pool_size  int
	Tasks_Chan chan string
	Quit       chan int
	Mutex      *sync.Mutex
	IsActive   bool
	Tasks      []string
}

/**
 * Create a new routine pool with number of active goroutines equal to size
 * @param size - size of routine pool
 */
func Createpool(size int) *Pool {
	fmt.Println("Creating New Routine Pool")

	var new_pool Pool
	new_pool.Pool_size = size

	new_pool.Mutex = &sync.Mutex{}
	new_pool.IsActive = false
	new_pool.Tasks_Chan = make(chan string)
	new_pool.Quit = make(chan int)

	return &new_pool
}

/**
 * Enqueue all items in tasks list into the channel, wait to be processed
 */
func (p *Pool) pool_enqueue() {
	for _, item := range p.Tasks {
		if p.IsActive == false {
			break
		}
		p.Tasks_Chan <- item
	}
}

/**
 * Start processing the to be processed channel
 * targ_funct - the function to be executed by the threadpool
 */
func (p *Pool) Pool_start(targ_func func(*Pool)) {
	fmt.Println("Starting Go Routines")
	p.IsActive = true
	for i := 0; i < p.Pool_size; i++ {
		go targ_func(p)
	}
}

/**
 * Stop execution and kill all goroutines in the pool
 */
func (p *Pool) Quit_pool() {
	p.IsActive = false
	for i := 0; i < p.Pool_size; i++ {
		p.Quit <- 1
	}
}

/**
 * Begin pool scheduler to periodicaly enqueue tasks into the to be executed channel
 * @param time_period - time period of the scheduler
 * @param tasks - list of tasks to be executed
 */
func (p *Pool) Pool_Scheduler_Start(time_period time.Duration, tasks []string) {
	fmt.Println("Starting Pool Scheduler")

	p.Tasks = append(p.Tasks, tasks...)

	go func() {
		for p.IsActive {
			fmt.Println("Enqueing Tasks")
			p.pool_enqueue()

			time.Sleep(time_period)
		}
	}()
}
