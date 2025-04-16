package main

import (
	"griddb.net/griddb-cloud-cli/cmd"
	_ "griddb.net/griddb-cloud-cli/cmd/Sql"
	_ "griddb.net/griddb-cloud-cli/cmd/checkConnection"
	_ "griddb.net/griddb-cloud-cli/cmd/readContainer"
)

func main() {
	cmd.Execute()
}
