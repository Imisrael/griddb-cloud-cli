package main

import (
	"griddb.net/griddb-cloud-cli/cmd"
	_ "griddb.net/griddb-cloud-cli/cmd/checkConnection"
	_ "griddb.net/griddb-cloud-cli/cmd/containerInfo"
	_ "griddb.net/griddb-cloud-cli/cmd/createContainer"
	_ "griddb.net/griddb-cloud-cli/cmd/listContainers"
	_ "griddb.net/griddb-cloud-cli/cmd/readContainer"
	_ "griddb.net/griddb-cloud-cli/cmd/sql"
)

func main() {
	cmd.Execute()
}
