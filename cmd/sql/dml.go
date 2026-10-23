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
	sqlCmd.AddCommand(dmlCmd)
	dmlCmd.Flags().StringVarP(&userSqlString, "string", "s", "", "SQL STRING")
	dmlCmd.MarkFlagRequired("string")
}

func runUpdate() {
	client := &http.Client{}

	var s = SqlString(userSqlString)
	s.wrapInDblQuoteAndStmt()
	fmt.Println(s)
	convert := []byte(s)

	buf := bytes.NewBuffer(convert)

	url := "/sql/dml/update"
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

var dmlCmd = &cobra.Command{
	Use:     "dml",
	Aliases: []string{"update"},
	Short:   "Run DML statements (INSERT, UPDATE, DELETE)",
	Example: `griddb-cloud-cli sql dml -s "DELETE FROM device1 WHERE temp > 100"`,
	Run: func(cmd *cobra.Command, args []string) {
		runUpdate()
	},
}
