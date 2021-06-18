package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func LoadPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	err = HandleResponseStatus(resp.StatusCode)
	if err != nil {
		return nil, err
	}
	
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Body.Close()

	return bytes, nil
}

func HandleResponseStatus(statusCode int) error {
	switch {
	case isSuccess(statusCode):
		return nil
		
	case isRedirection(statusCode) :
		fmt.Println("redirect")
		return nil

	case isClientError(statusCode):
		return fmt.Errorf("Response status %d ", statusCode)

	case isServerError(statusCode):
		return fmt.Errorf("Response status %d ", statusCode)

	default:
		return nil
	}
}

func isSuccess(code int) bool {
	return code < 300
}

func isRedirection(code int) bool {
	return code >= 300 && code < 400
}

func isClientError(code int) bool {
	return code >= 400 && code < 500
}

func isServerError(code int) bool {
	return code >= 500
}
