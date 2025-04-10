package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkConnectionCmd)
}

func checkConnection() {

	client := &http.Client{}
	req, err := makeNewRequest("GET", "/checkConnection")
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

var checkConnectionCmd = &cobra.Command{
	Use:   "checkConnection",
	Short: "Testing the test command",
	Long:  "like ya",
	Run: func(cmd *cobra.Command, args []string) {
		checkConnection()
	},
}
