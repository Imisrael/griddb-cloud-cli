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

	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

var (
	offset     int
	limit      int
	pretty     bool
	raw        bool
	height     int
	colToGraph string
)

func init() {
	cmd.RootCmd.AddCommand(readContainerCmd)
	readContainerCmd.Flags().IntVar(&offset, "offset", 0, "How many rows you'd like to offset in your query")
	readContainerCmd.Flags().IntVar(&limit, "limit", 100, "How many rows you'd like to limit")
	readContainerCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Print the JSON with Indent rules")
	readContainerCmd.Flags().BoolVar(&raw, "raw", false, "When enabled, will simply output direct results from GridDB Cloud")
	readContainerCmd.Flags().StringVar(&colToGraph, "colNames", "", "Which columns would you like to see charted (separated by commas!)")
}

func wrapInTqlObj(containerName string) string {

	//EXAMPLE [{"name" : "device1", "stmt" : "select * limit 100", "columns" : ["co", "humidity"], "hasPartialExecution" : true}]
	cols := unfurlUserColChoice()
	s := "[ { \"name\": \"" + containerName + "\", \"stmt\": \"select * limit " + strconv.Itoa(limit) + "\", \"columns\": " + cols + ", \"hasPartialExecution\": true }]"
	fmt.Println(s)
	return s
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

func readTql(containerName string, graph bool) [][]cmd.QueryData {
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
	cmd.CheckForErrors(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error with reading body! ", err)
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

	if raw {
		fmt.Println(string(body))
	} else if !graph {
		parseBody(body, pretty)
	}

	return data

}

func parseBody(body []byte, pretty bool) {
	var results []cmd.TQLResults

	if err := json.Unmarshal(body, &results); err != nil {
		panic(err)
	}
	//fmt.Println(results)
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

	if pretty {
		jso, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println(string(jso))
	} else {
		jso, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println(string(jso))
	}

}

var readContainerCmd = &cobra.Command{
	Use:   "read",
	Short: "Query container with TQL",
	Long:  "Read container and print contents in json format with --pretty",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatal("you may only read from one container at a time")
		} else if len(args) == 1 {
			readTql(args[0], false)
		} else {
			log.Fatal("Please include the container name you'd like to read from!")
		}

	},
}
