package deleteContainer

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/cqroot/prompt"
	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
	"griddb.net/griddb-cloud-cli/cmd/containerInfo"
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

	info := containerInfo.GetContainerInfo(containerName)

	make, err := prompt.New().Ask("Delete Container? \n" + string(info)).
		Choose([]string{"NO", "YES"})
	cmd.CheckErr(err)

	if make == "NO" {
		log.Fatal("Aborting")
	} else {

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

		fmt.Println(resp.StatusCode, "Successfully Deleted")
	}
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
