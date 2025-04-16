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
	sqlType   string
)

func init() {
	cmd.RootCmd.AddCommand(sqlCmd)
	sqlCmd.Flags().StringVarP(&sqlType, "type", "t", "query", "choose what kind of SQL command you want to conduct: query, update, or create")
	sqlCmd.Flags().StringVarP(&sqlString, "string", "s", "", "SQL STRING")
	sqlCmd.MarkFlagRequired("string")
}

func wrapInDoubleQuotes(sqlString string) string {
	newString := "\"" + sqlString + "\""
	return newString
}

func getURLSuffix(sqlType string) string {
	switch sqlType {
	case "query":
		return "/dml/query"
	case "update":
		return "/dml/update"
	case "create":
		return "/ddl"
	default:
		return "/dml/query"
	}
}

func runSql() {
	sqlString = wrapInDoubleQuotes(sqlString)
	client := &http.Client{}

	stmt := "[{\"stmt\": " + sqlString + " }]"
	fmt.Println(stmt)

	convert := []byte(stmt)
	buf := bytes.NewBuffer(convert)

	urlSuffix := getURLSuffix(sqlType)
	url := "/sql/" + urlSuffix

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

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Run a sql command",
	Long:  "Run SQL Against your GridDB Cloud DB. Must choose whether to run DDL, DML or DDL Update",
	Run: func(cmd *cobra.Command, args []string) {

		runSql()
	},
}
