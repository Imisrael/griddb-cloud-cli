package getContainers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(getContainersCmd)
}

func getContainers() {

	client := &http.Client{}
	req, err := cmd.MakeNewRequest("GET", "/containers?limit=100", nil)
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

	var listOfContainers cmd.ContainersList

	if err := json.Unmarshal(body, &listOfContainers); err != nil {
		panic(err)
	}

	for i, name := range listOfContainers.Names {
		fmt.Println(strconv.Itoa(i) + ": " + name)
	}
}

var getContainersCmd = &cobra.Command{
	Use:   "getContainers",
	Short: "Get a list of all of the containers",
	Long:  "The limit is set to 100 and is not configurable",
	Run: func(cmd *cobra.Command, args []string) {
		getContainers()
	},
}
