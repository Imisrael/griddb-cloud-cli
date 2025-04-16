package sql

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

func init() {
	sqlCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&sqlString, "string", "s", "", "SQL STRING")
	queryCmd.MarkFlagRequired("string")
}

func runQuery() {
	client := &http.Client{}

	sqlString = wrapInDoubleQuotes(sqlString)
	stmt := "[{\"stmt\": " + sqlString + " }]"
	fmt.Println(stmt)

	convert := []byte(stmt)
	buf := bytes.NewBuffer(convert)

	url := "/sql/dml/query"
	req, err := cmd.MakeNewRequest("POST", url, buf)
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
	fmt.Println(string(body))
}

var queryCmd = &cobra.Command{
	Use:     "query",
	Short:   "Run a sql command",
	Long:    "Run a DML Sql Command to Query your GridDB Container",
	Example: "sql query -s \"SELECT * FROM pyIntPart2\"",
	Run: func(cmd *cobra.Command, args []string) {

		runQuery()
	},
}
