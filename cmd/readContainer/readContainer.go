package readContainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

var (
	offset int
	limit  int
)

func init() {
	cmd.RootCmd.AddCommand(readContainerCmd)
	readContainerCmd.Flags().IntVar(&offset, "offset", 0, "How many rows you'd like to offset in your query")
	readContainerCmd.Flags().IntVar(&limit, "limit", 100, "How many rows you'd like to limit")
}

func readContainer(containerName string) {
	client := &http.Client{}
	convert := []byte(
		"{   \"offset\" : " + strconv.Itoa(offset) + ",   \"limit\": " + strconv.Itoa(limit) + "}",
	)
	buf := bytes.NewBuffer(convert)

	req, err := cmd.MakeNewRequest("POST", "/containers/"+containerName+"/rows", buf)
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
	parseBody(body)
}

func parseBody(body []byte) {
	var results cmd.CloudResults

	if err := json.Unmarshal(body, &results); err != nil {
		panic(err)
	}
	//fmt.Println(results)
	var cols []cmd.Columns = results.Columns
	var rows [][]any = results.Rows
	var rowsLength int

	if len(rows) > 0 {
		rowsLength = len(rows)
	}

	var data [][]cmd.QueryData = make([][]cmd.QueryData, rowsLength)

	for i := range rows {
		data[i] = make([]cmd.QueryData, len(rows[i]))
		for j := range rows[i] {
			data[i][j].Name = cols[j].Name
			data[i][j].Type = cols[j].Type
			data[i][j].Value = rows[i][j]
		}
	}

	jso, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println(string(jso))

}

var readContainerCmd = &cobra.Command{
	Use:   "read",
	Short: "Read container",
	Long:  "Read container and print out table",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatal("you may only read from one container at a time")
		} else if len(args) == 1 {
			readContainer(args[0])
		} else {
			log.Fatal("Please include the container name you'd like to read from!")
		}

	},
}
