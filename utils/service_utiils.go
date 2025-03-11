package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// UnmarshalJSON unmarshals the request body into the given interface
// It returns an error if the request body is not valid JSON
// It also closes the request body
func UnmarshalJSON(r *http.Request, v interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// GetPathParam returns the value of the path parameter with the given key
// It returns an error if the path parameter is not found
func GetPathParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// MarshalJSON marshals the given interface into a JSON byte slice
// It returns an error if the interface is not valid JSON
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Response creates an http.Response from the given interface and error
// If there is an error, it returns a response with status code 500 and the error message
// Otherwise, it returns a response with status code 200 and the marshaled interface
func ServerResponse(v interface{}, err error) *http.Response {
	if err != nil {
		//treat errors with specific types to return specific status codes
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(err.Error())),
		}
	}

	data, err := MarshalJSON(v)
	if err != nil {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(err.Error())),
		}
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(data)),
	}
}

func ResponseFailed(resp *http.Response) bool {
	return resp.StatusCode < 200 || resp.StatusCode >= 300
}

// GetErrorFromResponse extracts error message from error response
func GetErrorFromResponse(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	errMsg, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return errors.New(string(errMsg))
}
