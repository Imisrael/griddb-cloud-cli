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
	sqlCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&sqlString, "string", "s", "", "SQL STRING")
	createCmd.MarkFlagRequired("string")
}

func runCreate() {
	sqlString = wrapInDoubleQuotes(sqlString)
	client := &http.Client{}

	stmt := "[{\"stmt\": " + sqlString + " }]"
	fmt.Println(stmt)

	convert := []byte(stmt)
	buf := bytes.NewBuffer(convert)

	url := "/sql/ddl"
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

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Run a sql command",
	Long:    "Run a DDL Sql Command to Create or Alter or Update your SQL Tables",
	Example: "\"CREATE TABLE IF NOT EXISTS pyIntPart2 (date TIMESTAMP NOT NULL PRIMARY KEY, value STRING) WITH (expiration_type='PARTITION',expiration_time=10,expiration_time_unit='DAY') PARTITION BY RANGE (date) EVERY (5, DAY);\"",
	Run: func(cmd *cobra.Command, args []string) {
		runCreate()
	},
}
