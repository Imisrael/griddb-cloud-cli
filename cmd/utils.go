package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

func addBasicAuth(req *http.Request) {
	user := viper.Get("cloud_username").(string)
	pass := viper.Get("cloud_pass").(string)
	req.SetBasicAuth(user, pass)
}

func makeNewRequest(method, endpoint string, body io.Reader) (req *http.Request, e error) {

	url := viper.Get("url").(string)
	req, err := http.NewRequest(method, url+endpoint, body)
	if err != nil {
		fmt.Println("error with request:", err)
		return req, err
	}

	addBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}
