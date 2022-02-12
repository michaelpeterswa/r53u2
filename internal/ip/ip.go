package ip

import (
	"io/ioutil"
	"net/http"
)

func Get(s string) (string, error) {
	resp, err := http.Get(s)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
