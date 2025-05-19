package createContainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Imisrael/griddb-cloud-cli/cmd"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/input"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(createContainerCmd)
}

func ColDeclaration(numberOfCols int, timeseries bool) []cmd.ContainerInfoColumns {
	if numberOfCols < 1 {
		log.Fatal("Please pick a number greater than 0 for the number of cols")
	}
	var columnInfo []cmd.ContainerInfoColumns = make([]cmd.ContainerInfoColumns, numberOfCols)
	for i := range numberOfCols {

		colName, err := prompt.New().Ask("Col name For col #" + strconv.Itoa(i+1)).Input("temperature")
		cmd.CheckErr(err)
		var colType string

		// if it's timeseries and first col, it's always set to ROWKEY=true and timestamp type
		if timeseries && i == 0 {
			colType, err = prompt.New().Ask("Col #" + strconv.Itoa(i+1) + "(TIMESTAMP CONTAINERS ARE LOCKED TO TIMESTAMP FOR THEIR ROWKEY)").
				Choose([]string{"TIMESTAMP"})
			cmd.CheckErr(err)
		} else {
			// User inputs col type for every other scenario
			colType, err = prompt.New().Ask("Column Type for col #" + strconv.Itoa(i+1)).
				Choose([]string{"BOOL", "STRING", "INTEGER", "LONG", "FLOAT", "DOUBLE", "TIMESTAMP", "GEOMETRY"})
			cmd.CheckErr(err)
		}

		// if collection type and first col, you must set an index type.
		if !timeseries && i == 0 {
			colIndex, err := prompt.New().Ask("Column Index Type" + strconv.Itoa(i+1)).
				Choose([]string{"none (default)", "TREE", "SPATIAL"})
			cmd.CheckErr(err)
			var indexArr []string = make([]string, 1)
			if colIndex == "none (default)" {
				columnInfo[i].Index = nil
			} else {
				indexArr[0] = colIndex
				columnInfo[i].Index = indexArr
			}

		}

		columnInfo[i].Name = colName
		columnInfo[i].Type = colType

	}
	return columnInfo
}

func InteractiveContainerInfo(ingest bool, header []string) cmd.ContainerInfo {

	var retVal cmd.ContainerInfo
	var rowkey bool = true
	var timeseries bool = false

	containerName, err := prompt.New().Ask("Container Name:").Input("device2")
	cmd.CheckErr(err)

	containerType, err := prompt.New().Ask("Choose:").
		Choose([]string{"COLLECTION", "TIME_SERIES"})
	cmd.CheckErr(err)

	if containerType == "COLLECTION" {
		rk, err := prompt.New().Ask("Row Key?").
			Choose([]string{"true", "false"})
		cmd.CheckErr(err)
		val, err := strconv.ParseBool(rk)
		if err != nil {
			fmt.Println(err)
		}
		rowkey = val
	} else {
		timeseries = true
	}

	var colInfos []cmd.ContainerInfoColumns

	if ingest {
		colInfos = ColDeclaration(len(header), timeseries)
	} else {
		numberOfCols, err := prompt.New().Ask("How Many Columns for this Container?").
			Input("", input.WithInputMode(input.InputInteger))
		cmd.CheckErr(err)

		numOfCols, err := strconv.Atoi(numberOfCols)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		colInfos = ColDeclaration(numOfCols, timeseries)
	}

	retVal.ContainerName = containerName
	retVal.ContainerType = containerType
	retVal.RowKey = rowkey
	retVal.Columns = colInfos

	return retVal

}

func Create(conInfo cmd.ContainerInfo) {

	jsoPrettyPrint, err := json.MarshalIndent(conInfo, "", "    ")
	if err != nil {
		fmt.Println("Error", err)
	}

	make, err := prompt.New().Ask("Make Container? \n" + string(jsoPrettyPrint)).
		Choose([]string{"YES", "NO"})
	cmd.CheckErr(err)

	if make == "NO" {
		log.Fatal("Aborting")
	} else {
		jsonContainerInfo, err := json.Marshal(conInfo)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println(string(jsonContainerInfo))
		convert := []byte(jsonContainerInfo)
		buf := bytes.NewBuffer(convert)

		client := &http.Client{}
		req, err := cmd.MakeNewRequest("POST", "/containers", buf)
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

var createContainerCmd = &cobra.Command{
	Use:   "create",
	Short: "Interactive walkthrough to create a container",
	Long:  "A series of CLI prompts to create your griddb container",
	Run: func(cmd *cobra.Command, args []string) {

		var ingest bool = false
		conInfo := InteractiveContainerInfo(ingest, nil)
		Create(conInfo)
	},
}
