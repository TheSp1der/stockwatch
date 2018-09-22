/******************************************************************************
	httpclient.go
	http client functions for interacting with remote http web servers
******************************************************************************/
package main

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

func httpGet(url string, headers httpHeader) ([]byte, error) {
	var (
		err    error          // error handler
		client http.Client    // http client
		req    *http.Request  // http request
		res    *http.Response // http response
		output []byte         // output
	)

	// set timeouts
	client = http.Client{
		Timeout: time.Duration(time.Second * 2),
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Duration(time.Second * 2),
			}).Dial,
			TLSHandshakeTimeout: time.Duration(time.Second * 2),
		},
	}

	// setup request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return output, err
	}

	// setup headers
	if len(headers) > 0 {
		for _, header := range headers {
			req.Header.Set(header.Name, header.Value)
		}
	}

	// perform the request
	if res, err = client.Do(req); err != nil {
		return output, err
	}

	// close the connection upon function closure
	defer res.Body.Close()

	// extract response body
	if output, err = ioutil.ReadAll(res.Body); err != nil {
		return output, err
	}

	// check status
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return output, errors.New("non-successful status code received [" + strconv.Itoa(res.StatusCode) + "]")
	}

	return output, nil
}
