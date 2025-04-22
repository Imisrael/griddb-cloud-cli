package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type Columns struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	TimePrecision string `json:"timePrecision,omitempty"`
}

type ContainerInfoColumns struct {
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	TimePrecision string   `json:"timePrecision,omitempty"`
	Index         []string `json:"index"`
}

type CloudResults struct {
	Offset  int       `json:"offset"`
	Limit   int       `json:"limit"`
	Total   int       `json:"total"`
	Rows    [][]any   `json:"rows"`
	Columns []Columns `json:"columns"`
}

type TQLResults struct {
	Offset           int       `json:"offset"`
	Limit            int       `json:"limit"`
	Total            int       `json:"total"`
	ResponseSizeByte int32     `json:"responseSizeByte,omitempty"`
	Results          [][]any   `json:"results"`
	Columns          []Columns `json:"columns"`
}

type SqlResults struct {
	ResponseSizeByte float32   `json:"responseSizeByte"`
	Results          [][]any   `json:"results"`
	Columns          []Columns `json:"columns"`
}

type QueryData struct {
	Name  string
	Type  string
	Value any
}

type ContainersList struct {
	Names  []string `json:"names"`
	Offset int      `json:"offset"`
	Limit  int      `json:"limit"`
	Total  int      `json:"total"`
}

type ContainerInfo struct {
	ContainerName string                 `json:"container_name"`
	ContainerType string                 `json:"container_type"`
	RowKey        bool                   `json:"rowkey"`
	Columns       []ContainerInfoColumns `json:"columns"`
}

type ErrorMsg struct {
	Version      string `json:"version"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func AddBasicAuth(req *http.Request) {
	user := viper.Get("cloud_username").(string)
	pass := viper.Get("cloud_pass").(string)
	req.SetBasicAuth(user, pass)
}

func MakeNewRequest(method, endpoint string, body io.Reader) (req *http.Request, e error) {

	url := viper.Get("url").(string)
	req, err := http.NewRequest(method, url+endpoint, body)
	if err != nil {
		fmt.Println("error with request:", err)
		return req, err
	}

	AddBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func CheckForErrors(resp *http.Response) {

	if resp.StatusCode == 403 {
		log.Fatal("(403) IP Connection Error. Is this IP Address Whitelisted?")
	}

	if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error with reading body! ", err)
		}

		var errorMsg ErrorMsg

		if err := json.Unmarshal(body, &errorMsg); err != nil {
			panic(err)
		}
		if resp.StatusCode == 400 {
			log.Fatal("400 Error: " + errorMsg.ErrorMessage)
		} else if resp.StatusCode == 401 {
			fmt.Println(errorMsg.ErrorMessage)
			log.Fatal("Authentication Error. Please check your username and password in your config file ")
		} else if resp.StatusCode == 404 {
			log.Fatal("404 (not found) - Does this container exist?")
		} else if resp.StatusCode == 500 {
			log.Fatal("500 error! " + errorMsg.ErrorMessage)
		} else {
			log.Fatal("Unknown Error. Please try again + ")
		}
	}

}
