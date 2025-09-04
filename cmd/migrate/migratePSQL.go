package migrate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Imisrael/griddb-cloud-cli/cmd"
	"github.com/Imisrael/griddb-cloud-cli/cmd/createContainer"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migratePsqlCmd)
}

type PsqlSchema struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
}

func typeSwitcher(s string) string {
	switch s {
	case "bool":
		return "BOOL"
	case "char", "varchar", "text":
		return "STRING"
	case "int", "int2", "int4":
		return "INTEGER"
	case "int8":
		return "LONG"
	case "decimal", "real", "numeric":
		return "FLOAT"
	case "float", "float8":
		return "DOUBLE"
	case "timetz", "timestamptz":
		return "TIMESTAMP"
	default:
		return strings.ToUpper(s)

	}

}

func transformToConInfoCols(colSet []PsqlSchema) []cmd.ContainerInfoColumns {
	n := len(colSet)
	var conInfoCols = make([]cmd.ContainerInfoColumns, n)

	for idx, val := range colSet {
		conInfoCols[idx].Name = val.ColumnName
		conInfoCols[idx].Type = typeSwitcher(val.DataType)
	}
	return conInfoCols
}

func parseJson(jsonName string) map[string]cmd.ContainerInfo {
	schema, err := os.ReadFile(jsonName)
	if err != nil {
		log.Fatal(err)
	}
	var psqlSchema map[string][]PsqlSchema

	err = json.Unmarshal(schema, &psqlSchema)
	if err != nil {
		log.Fatal(err)
	}

	var conInfo = make(map[string]cmd.ContainerInfo)

	for tableName, schema := range psqlSchema {
		var info cmd.ContainerInfo
		info.ContainerName = tableName
		info.ContainerType = "COLLECTION"
		info.RowKey = true
		cols := transformToConInfoCols(schema)
		info.Columns = cols
		conInfo[tableName] = info
	}

	return conInfo

}

func migratePsql(dirName string) {

	c, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	var schemaFile string

	for _, entry := range c {
		name := entry.Name()
		if strings.Contains(name, ".json") {
			schemaFile = dirName + "/" + name
		}

	}

	fmt.Println(schemaFile)

	conInfo := parseJson(schemaFile)

	for tableName, info := range conInfo {
		createContainer.Create(info, force)
		types := mapping(info.Columns)
		associatedCsvFile := dirName + "/" + tableName + ".csv"

		allRows, err := readAllRows(associatedCsvFile)
		//chop off col headers
		allRows = allRows[1:]
		if err != nil {
			log.Fatal(err)
		}
		processCSV(allRows, types, tableName, associatedCsvFile)

	}

}

var migratePsqlCmd = &cobra.Command{
	Use:     "psql",
	Short:   "Migrate from PostgreSQL csv/json data to GridDB Cloud",
	Long:    "Use the json query ",
	Example: "griddb-cloud-cli migrate psql <directory>",
	Run: func(cmd *cobra.Command, args []string) {
		migratePsql(args[0])
	},
}
