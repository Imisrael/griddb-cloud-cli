package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type Columns struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	TimePrecision string `json:"timePrecision,omitempty"`
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
