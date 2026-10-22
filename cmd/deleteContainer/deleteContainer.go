package deleteContainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Imisrael/griddb-cloud-cli/cmd"
	"github.com/Imisrael/griddb-cloud-cli/cmd/containerInfo"
	"github.com/cqroot/prompt"
	"github.com/spf13/cobra"
)

var (
	multipleContainers []string
	force              bool
)

func init() {
	cmd.RootCmd.AddCommand(deleteContainerCmd)
	deleteContainerCmd.Flags().StringSliceVarP(&multipleContainers, "multi", "m", []string{}, "All containers to be deleted 'col1,col2,col3'")
	deleteContainerCmd.Flags().BoolVarP(&force, "force", "f", false, "Force delete (no prompt)")
}

type ContainerToBeDeleted string

func (c *ContainerToBeDeleted) wrapInDblQuotesAndBracket() {
	*c = "\"" + *c + "\""
	*c = "[" + *c + "]"
}

func (c *ContainerToBeDeleted) wrapInDblQuotes() {
	*c = "\"" + *c + "\""
}

func deleteRequest(payload []byte) {
	client := &http.Client{}

	buf := bytes.NewBuffer(payload)

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

func deleteAllContainers() {
	payload, err := json.Marshal(multipleContainers)
	if err != nil {
		log.Fatal(err)
	}

	if !force {
		ans, err := prompt.New().Ask("Delete these containers?\n" + strings.Join(multipleContainers, "\n")).
			Choose([]string{"NO", "YES"})
		cmd.CheckErr(err)
		if ans == "NO" {
			log.Fatal("Aboring")
		}
	}

	deleteRequest(payload)
}

func deleteContainer(containerName string) {

	info := containerInfo.GetContainerInfo(containerName)

	var containerToDelete = ContainerToBeDeleted(containerName)
	containerToDelete.wrapInDblQuotesAndBracket()
	payload := []byte(containerToDelete)

	if force {
		deleteRequest(payload)
		return
	}

	make, err := prompt.New().Ask("Delete Container? \n" + string(info)).
		Choose([]string{"NO", "YES"})
	cmd.CheckErr(err)

	if make == "NO" {
		log.Fatal("Aborting")
	} else {

		deleteRequest(payload)

	}
}

var deleteContainerCmd = &cobra.Command{
	Use:   "delete [container]",
	Short: "Delete one or more containers",
	Long:  "Delete a single container by name, or several at once with -m. Prompts for confirmation unless -f is set.",
	Example: `griddb-cloud-cli delete device1
griddb-cloud-cli delete -m "device1,device2,device3"
griddb-cloud-cli delete -f -m "device1,device2"`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(multipleContainers) > 0 {
			deleteAllContainers()
			return
		}
		if len(args) == 0 {
			log.Fatalln("Give me a container name, or use -m for several")
		}
		if strings.Contains(args[0], ",") {
			log.Fatalln("Did you forget the '-m'/'multi' flag?")
		}
		deleteContainer(args[0])
	},
}
