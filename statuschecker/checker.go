package statuschecker

import (
	"context"
	"log"
	"net/http"
	"time"
)

type StatusChecker interface {
	Check(ctx context.Context, hostname string) (status bool, err error)
}

type HttpChecker struct {
}

/**
 * Reciecer for http based website status check.
 * @param {string} ctx - context to be passed to http request
 * @param {string} hostnane - hostname of the website to be pinged
 */
func (h HttpChecker) Check(ctx context.Context, hostname string) (status bool,
	err_ret error) {
	url := "http://" + hostname

	var client = http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		log.Println("Error1: ", err)
		err_ret = err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error2: ", err)
		err_ret = err
		return
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		status = true
	}

	return
}
