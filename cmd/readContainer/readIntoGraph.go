package readContainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

func init() {
	readContainerCmd.AddCommand(readIntoGraph)
	readIntoGraph.Flags().IntVar(&offset, "offset", 0, "How many rows you'd like to offset in your query")
	readIntoGraph.Flags().IntVar(&limit, "limit", 100, "How many rows you'd like to limit")
	readIntoGraph.Flags().BoolVarP(&pretty, "pretty", "p", false, "Print the JSON with Indent rules")
	readIntoGraph.Flags().BoolVar(&raw, "raw", false, "When enabled, will simply output direct results from GridDB Cloud")
}

func getGraphData(containerName string) [][]cmd.QueryData {
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
	var results cmd.CloudResults

	if err := json.Unmarshal(body, &results); err != nil {
		panic(err)
	}
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

	return data
}

func graphIt(data [][]cmd.QueryData) {

	var rows [][]float64 = make([][]float64, len(data))

	for i := range data {
		rows[i] = make([]float64, 1)
		for j := range data[i] {
			if data[i][j].Type == "FLOAT" || data[i][j].Type == "INTEGER" || data[i][j].Type == "DOUBLE" {
				fmt.Println(i, j)
				fmt.Println(data[i][j])
				rows[i] = append(rows[i], data[i][j].Value.(float64))
			}
		}
	}
	fmt.Println(rows)

	graph := asciigraph.PlotMany(rows, asciigraph.Height(20))
	fmt.Println(graph)
}

var readIntoGraph = &cobra.Command{
	Use:   "graph",
	Short: "Read container",
	Long:  "Read container and print out table",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatal("you may only read from one container at a time")
		} else if len(args) == 1 {
			data := getGraphData(args[0])
			graphIt(data)
		} else {
			log.Fatal("Please include the container name you'd like to read from!")
		}

	},
}
