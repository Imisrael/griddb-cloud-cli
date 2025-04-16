package sql

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

var (
	sqlString string
)

func init() {
	cmd.RootCmd.AddCommand(sqlCmd)

	sqlCmd.Flags().StringVarP(&sqlString, "string", "s", "", "SQL STRING")
	sqlCmd.MarkFlagRequired("string")
}

func wrapInDoubleQuotes(sqlString string) string {
	newString := "\"" + sqlString + "\""
	return newString
}

func runSql() {
	sqlString = wrapInDoubleQuotes(sqlString)
	client := &http.Client{}

	stmt := "[{\"stmt\": " + sqlString + " }]"
	fmt.Println(stmt)

	convert := []byte(stmt)
	buf := bytes.NewBuffer(convert)

	req, err := cmd.MakeNewRequest("POST", "/sql/dml/query", buf)
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

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Run a sql command",
	Long:  "Run SQL against your DB",
	Run: func(cmd *cobra.Command, args []string) {
		runSql()
	},
}
