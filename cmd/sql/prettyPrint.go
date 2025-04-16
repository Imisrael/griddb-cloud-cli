package sql

import (
	"encoding/json"
	"fmt"

	"griddb.net/griddb-cloud-cli/cmd"
)

func prettyPrint(body []byte, pretty bool) []byte {
	var results []cmd.SqlResults

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
		return jso
	} else {
		jso, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error", err)
		}
		return jso
	}

}
