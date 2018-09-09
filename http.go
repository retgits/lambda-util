// Package util implements utility methods
package util

// The imports
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HTTPResponse is the response type for the HTTP requests
type HTTPResponse struct {
	Body    map[string]interface{}
	Headers http.Header
}

// HTTPGet executes a GET request to a URL and returns the response body as a JSON object
func HTTPGet(URL string, header http.Header) (HTTPResponse, error) {
	httpresponse := HTTPResponse{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return httpresponse, fmt.Errorf("error while creating HTTP request: %s", err.Error())
	}

	if header != nil {
		req.Header = header
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return httpresponse, fmt.Errorf("error while performing HTTP request: %s", err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return httpresponse, fmt.Errorf("the HTTP request returned a non-OK response: %v", res.StatusCode)
	}

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

	return httpresponse, nil
}

// HTTPPost executes a POST request to a URL and returns the response body as a JSON object
func HTTPPost(URL string, encoding string, postData url.Values) (HTTPResponse, error) {
	httpresponse := HTTPResponse{}

	res, err := http.Post(URL, encoding, strings.NewReader(postData.Encode()))
	if err != nil {
		return httpresponse, fmt.Errorf("error while performing HTTP request: %s", err.Error())
	}

	if res.StatusCode != 200 {
		return httpresponse, fmt.Errorf("the HTTP request returned a non-OK response: %v", res.StatusCode)
	}

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

	return httpresponse, nil
}
