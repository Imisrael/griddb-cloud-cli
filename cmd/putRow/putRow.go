package putRow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cqroot/prompt"
	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
	"griddb.net/griddb-cloud-cli/cmd/containerInfo"
)

func init() {
	cmd.RootCmd.AddCommand(putRowCmd)
}

func placeHolderVal(colType string) string {
	switch colType {
	case "TIMESTAMP":
		return "2016-01-16T10:25:00.253Z"
	case "BOOL":
		return "true"
	case "STRING":
		return "meter_1"
	case "INTEGER", "BYTE", "SHORT", "LONG":
		return "10"
	case "FLOAT", "DOUBLE":
		return "32.05"
	default:
		return "meter_1"
	}
}

func ConvertType(colType, val string) string {
	switch colType {

	case "TIMESTAMP":
		return val

	case "BOOL":
		return val

	case "STRING":
		return "\"" + val + "\""

	case "INTEGER", "BYTE", "SHORT", "LONG":
		return val

	case "FLOAT", "DOUBLE":
		return val

	default:
		return val
	}
}

func BuildPutRowContents(containerName string) string {

	info := containerInfo.GetContainerInfo(containerName)

	var containerInfo cmd.ContainerInfo
	var cols []cmd.ContainerInfoColumns
	if err := json.Unmarshal(info, &containerInfo); err != nil {
		panic(err)
	}
	cols = containerInfo.Columns

	var stringOfValues string = "[["

	fmt.Println("Container Name: " + containerName)

	for i, cont := range cols {
		defaultValue := placeHolderVal(cont.Type)
		val, err := prompt.New().Ask("Column " + strconv.Itoa(i+1) + " of " + strconv.Itoa(len(cols)) + "\n Column Name: " + cont.Name + "\n Column Type: " + cont.Type).Input(defaultValue)
		cmd.CheckErr(err)
		if i == 0 {
			stringOfValues = stringOfValues + ConvertType(cont.Type, val)
		} else {
			stringOfValues = stringOfValues + ",  " + ConvertType(cont.Type, val)
		}

	}

	stringOfValues = stringOfValues + "]]"
	return stringOfValues

}

func put(containerName string) {

	conInfo := BuildPutRowContents(containerName)

	fmt.Println(conInfo)

	make, err := prompt.New().Ask("Add the Following to container " + containerName + "?").
		Choose([]string{"YES", "NO"})
	cmd.CheckErr(err)

	if make == "NO" {
		log.Fatal("Aborting")
	} else {

		convert := []byte(conInfo)
		buf := bytes.NewBuffer(convert)

		client := &http.Client{}
		req, err := cmd.MakeNewRequest("PUT", "/containers/"+containerName+"/rows", buf)
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

}

var putRowCmd = &cobra.Command{
	Use:   "put",
	Short: "Interactive walkthrough to push a row",
	Long:  "A series of CLI prompts to create your griddb container",
	Run: func(cmd *cobra.Command, args []string) {
		put(args[0])
	},
}
