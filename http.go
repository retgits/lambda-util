// Package util implements utility methods
package util

// The imports
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPResponse is the response type for the HTTP requests
type HTTPResponse struct {
	Body       map[string]interface{}
	StatusCode int
	Headers    http.Header
}

// HTTPPost executes a POST request to a URL and returns the response body as a JSON object
func HTTPPost(URL string, encoding string, postData string, header http.Header) (HTTPResponse, error) {
	return httpcall(URL, "POST", encoding, postData, header)
}

// HTTPPatch executes a PATCH request to a URL and returns the response body as a JSON object
func HTTPPatch(URL string, encoding string, postData string, header http.Header) (HTTPResponse, error) {
	return httpcall(URL, "PATCH", encoding, postData, header)
}

// HTTPGet executes a GET request to a URL and returns the response body as a JSON object
func HTTPGet(URL string, encoding string, header http.Header) (HTTPResponse, error) {
	return httpcall(URL, "GET", encoding, "", header)
}

// httpcall executes an HTTP request request to a URL and returns the response body as a JSON object
func httpcall(URL string, requestType string, encoding string, payload string, header http.Header) (HTTPResponse, error) {
	// Instantiate a response object
	httpresponse := HTTPResponse{}

	// Prepare placeholders for the request and the error object
	req := &http.Request{}
	var err error

	// Create a request
	if len(payload) > 0 {
		req, err = http.NewRequest(requestType, URL, strings.NewReader(payload))
		if err != nil {
			return httpresponse, fmt.Errorf("error while creating HTTP request: %s", err.Error())
		}
	} else {
		req, err = http.NewRequest(requestType, URL, nil)
		if err != nil {
			return httpresponse, fmt.Errorf("error while creating HTTP request: %s", err.Error())
		}
	}

	// Associate the headers with the request
	if header != nil {
		req.Header = header
	}

	// Execute the HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return httpresponse, fmt.Errorf("error while performing HTTP request: %s", err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return httpresponse, err
	}

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return httpresponse, fmt.Errorf("error while unmarshaling HTTP response to JSON: %s", err.Error())
	}

	httpresponse.Body = data
	httpresponse.Headers = res.Header
	httpresponse.StatusCode = res.StatusCode

	return httpresponse, nil
}
