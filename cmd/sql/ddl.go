package sql

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Imisrael/griddb-cloud-cli/cmd"
	"github.com/spf13/cobra"
)

func init() {
	sqlCmd.AddCommand(ddlCmd)
	ddlCmd.Flags().StringVarP(&userSqlString, "string", "s", "", "SQL STRING")
	ddlCmd.MarkFlagRequired("string")
}

func runCreate() {
	client := &http.Client{}

	var s = SqlString(userSqlString)
	s.wrapInDblQuoteAndStmt()
	fmt.Println(s)

	convert := []byte(s)
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
	cmd.CheckForErrors(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error with reading body! ", err)
	}
	fmt.Println(string(body))
}

var ddlCmd = &cobra.Command{
	Use:     "ddl",
	Aliases: []string{"create"},
	Short:   "Run DDL statements (CREATE TABLE, DROP TABLE, ALTER TABLE)",
	Example: `griddb-cloud-cli sql ddl -s "DROP TABLE device1"`,
	Run: func(cmd *cobra.Command, args []string) {
		runCreate()
	},
}
