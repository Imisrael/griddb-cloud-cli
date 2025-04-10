package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var containerName string

func init() {
	rootCmd.AddCommand(readContainerCmd)
	readContainerCmd.Flags().StringVarP(&containerName, "containerName", "n", "", "Container Name you'd like to read from")
	readContainerCmd.MarkFlagRequired("containerName")
}

func readContainer() {
	client := &http.Client{}
	convert := []byte("{   \"offset\" : 0,   \"limit\": 100 }")
	buf := bytes.NewBuffer(convert)

	req, err := makeNewRequest("POST", "/containers/"+containerName+"/rows", buf)
	if err != nil {
		fmt.Println("Error making new request", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error with client DO: ", err)
	}

	fmt.Println(resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error with reading body! ", err)
	}
	fmt.Println(string(body))
}

var readContainerCmd = &cobra.Command{
	Use:   "readContainer",
	Short: "Read container",
	Long:  "Read container and print out table",
	Run: func(cmd *cobra.Command, args []string) {
		readContainer()
	},
}
