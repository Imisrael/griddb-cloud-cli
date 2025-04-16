package sql

import (
	"log"

	"github.com/spf13/cobra"
	"griddb.net/griddb-cloud-cli/cmd"
)

var (
	sqlString string
	pretty    bool
	raw       bool
)

func init() {
	cmd.RootCmd.AddCommand(sqlCmd)
}

func wrapInDoubleQuotes(sqlString string) string {
	newString := "\"" + sqlString + "\""
	return newString
}

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Run a sql command",
	Long:  "Run SQL Against your GridDB Cloud DB. Must choose whether to run DDL, DML or DDL Update",
	Run: func(cmd *cobra.Command, args []string) {

		log.Fatal("Please use one of the following subcommands: query, update, create")
	},
}
