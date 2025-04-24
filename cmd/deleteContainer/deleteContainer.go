package deleteContainer

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(deleteContainerCmd)
}

type ContainerToBeDeleted string

func (c *ContainerToBeDeleted) wrapInDblQuotesAndBracket() {
	*c = "\"" + *c + "\""
	*c = "[" + *c + "]"
}

func deleteContainer(containerName string) {

	var containerToDelete = ContainerToBeDeleted(containerName)
	containerToDelete.wrapInDblQuotesAndBracket()

	client := &http.Client{}

	sliceOfBytes := []byte(containerToDelete)
	buf := bytes.NewBuffer(sliceOfBytes)

	req, err := cmd.MakeNewRequest("DELETE", "/containers", buf)
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

var deleteContainerCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Test your Connection with GridDB Cloud",
	Long:    "A response of 200 is ideal, 401 is an auth error",
	Example: "griddb-cloud-cli checkConnection",
	Run: func(cmd *cobra.Command, args []string) {
		deleteContainer(args[0])
	},
}
