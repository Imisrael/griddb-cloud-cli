package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func addBasicAuth(req *http.Request) {
	user := viper.Get("cloud_username").(string)
	pass := viper.Get("cloud_pass").(string)
	req.SetBasicAuth(user, pass)
}

func makeNewRequest(method, endpoint string) (req *http.Request, e error) {

	url := viper.Get("url").(string)
	req, err := http.NewRequest(method, url+endpoint, nil)
	if err != nil {
		fmt.Println("error with request:", err)
		return req, err
	}

	addBasicAuth(req)
	return req, nil
}
