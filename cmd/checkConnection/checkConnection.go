package checkConnection

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/Imisrael/griddb-cloud-cli/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(checkConnectionCmd)
}

func checkConnection() {

	client := &http.Client{}
	req, err := cmd.MakeNewRequest("GET", "/checkConnection", nil)
	if err != nil {
		fmt.Println("Error making new request", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error with client DO: ", err)
	}

	cmd.CheckForErrors(resp)

	fmt.Println(resp.Status)
}

var checkConnectionCmd = &cobra.Command{
	Use:     "checkConnection",
	Short:   "Test your Connection with GridDB Cloud",
	Long:    "A response of 200 is ideal, 401 is an auth error",
	Example: "griddb-cloud-cli checkConnection",
	Run: func(cmd *cobra.Command, args []string) {
		checkConnection()
	},
}
