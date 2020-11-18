package infrastructure

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// DoRequest is a generic method to do requests and return the data as []bytes
func DoRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	return body, err
}
