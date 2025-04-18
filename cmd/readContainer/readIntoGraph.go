package readContainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

var (
	colToGraph string
)

func init() {
	readContainerCmd.AddCommand(readIntoGraph)
	readIntoGraph.Flags().IntVar(&offset, "offset", 0, "How many rows you'd like to offset in your query")
	readIntoGraph.Flags().IntVar(&limit, "limit", 100, "How many rows you'd like to limit")
	readIntoGraph.Flags().StringVar(&colToGraph, "colNames", "", "Which columns would you like to see charted (separated by commas!)")
}

func unfurlUserColChoice() string {
	if len(colToGraph) > 1 {
		removeSpace := strings.ReplaceAll(colToGraph, " ", "")
		turnToSlice := strings.Split(removeSpace, ",")
		s := "["
		for _, val := range turnToSlice {
			dblQuote := "\"" + val + "\","
			s = s + dblQuote
		}
		s = strings.Trim(s, ",")
		s = s + "]"
		//	fmt.Println(s)
		return s
	} else {
		return "null"
	}
}

func wrapInTqlObj(containerName string) string {

	//EXAMPLE [{"name" : "device1", "stmt" : "select * limit 100", "columns" : ["co", "humidity"], "hasPartialExecution" : true}]
	cols := unfurlUserColChoice()
	s := "[ { \"name\": \"" + containerName + "\", \"stmt\": \"select * limit " + strconv.Itoa(limit) + "\", \"columns\": " + cols + ", \"hasPartialExecution\": true }]"
	fmt.Println(s)
	return s
}

func readTql(containerName string) [][]cmd.QueryData {
	client := &http.Client{}

	stmt := wrapInTqlObj(containerName)
	stmtBytes := []byte(stmt)
	buf := bytes.NewBuffer(stmtBytes)

	req, err := cmd.MakeNewRequest("POST", "/tql/", buf)
	if err != nil {
		fmt.Println("Error making new request", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error with client DO: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error with reading body! ", err)
	}

	if resp.StatusCode == 400 {
		var errorMsg cmd.ErrorMsg
		if err := json.Unmarshal(body, &errorMsg); err != nil {
			panic(err)
		}
		log.Fatal("Reading Container ERROR: " + errorMsg.ErrorMessage)
	}

	var results []cmd.TQLResults

	if err := json.Unmarshal(body, &results); err != nil {
		panic(err)
	}
	var cols []cmd.Columns = results[0].Columns
	var rows [][]any = results[0].Results
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

	var m map[string][]float64 = make(map[string][]float64)

	for i := range data {
		for j := range data[i] {
			if data[i][j].Type == "FLOAT" || data[i][j].Type == "INTEGER" || data[i][j].Type == "DOUBLE" {
				m[data[i][j].Name] = append(m[data[i][j].Name], data[i][j].Value.(float64))
			}
		}
	}

	var rows [][]float64 = make([][]float64, len(m))
	var colNames []string = make([]string, len(m))

	var i int
	for rowName, rowValue := range m {
		rows[i] = make([]float64, len(rowValue))
		rows[i] = rowValue
		colNames[i] = rowName
		i++
	}

	graph := asciigraph.PlotMany(
		rows,
		asciigraph.Height(30),
		asciigraph.SeriesColors(asciigraph.Red, asciigraph.Green, asciigraph.Blue, asciigraph.Pink, asciigraph.Orange),
		asciigraph.SeriesLegends(colNames...),
		asciigraph.Caption("Series with legends"),
		asciigraph.Width(100),
	)
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
			data := readTql(args[0])
			graphIt(data)
		} else {
			log.Fatal("Please include the container name you'd like to read from!")
		}

	},
}
