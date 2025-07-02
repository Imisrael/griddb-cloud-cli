package migrate

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Imisrael/griddb-cloud-cli/cmd"
	"github.com/Imisrael/griddb-cloud-cli/cmd/createContainer"
	"github.com/Imisrael/griddb-cloud-cli/cmd/putRow"
	"github.com/spf13/cobra"
)

var (
	force bool
)

func init() {
	cmd.RootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVarP(&force, "force", "f", false, "Force create (no prompt)")
}

func readCSVFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseCSV(data []byte) (*csv.Reader, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1 //gets rid of checking for certain column length in csv
	return reader, nil
}

func readAllRows(csvFileName string) ([][]string, error) {
	data, err := readCSVFile(csvFileName)
	if err != nil {
		log.Fatal("Error reading file:", err)
		return nil, err
	}
	reader, err := parseCSV(data)
	if err != nil {
		log.Fatal("Error creating CSV reader:", err)
		return nil, err
	}

	allRows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Chop off first four rows, they are meta data for import tool
	allRows = allRows[4:]
	return allRows, nil
}

func mapping(cols []cmd.ContainerInfoColumns) []string {

	var types []string
	for _, col := range cols {
		types = append(types, col.Type)
	}
	return types
}

func processCSV(allRows [][]string, typeMapping []string, containerName, fileName string) {

	var stringOfValues string = "["
	for idx, row := range allRows {
		if idx != 0 {
			stringOfValues = stringOfValues + ", ["
		}
		for i, record := range row {

			// Convert Type method here converts the values to what the http request expects
			// for example, timestmap is converted to the format it likes
			// and strings are encapusulated into double quotes, etc
			if i == 0 {
				stringOfValues = stringOfValues + putRow.ConvertType(typeMapping[i], record)
			} else {
				stringOfValues = stringOfValues + ",  " + putRow.ConvertType(typeMapping[i], record)
			}
		}
		stringOfValues = stringOfValues + "]"
	}

	stringOfValues = "[" + stringOfValues + "]"
	// We can enter entire files at a time because they seem to be relatively short
	fmt.Println("inserting into (" + containerName + "). csv: " + fileName)
	putMultiRows(stringOfValues, containerName)

}

func putMultiRows(arrayString, containerName string) {

	url := "/containers/" + containerName + "/rows"
	convert := []byte(arrayString)
	buf := bytes.NewBuffer(convert)

	client := &http.Client{}
	req, err := cmd.MakeNewRequest("PUT", url, buf)
	if err != nil {
		log.Fatal("Error making new request", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("error with client DO: ", err)
	}

	cmd.CheckForErrors(resp)

	fmt.Println(resp.Status)
}

func migrate(dirName string) {

	c, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	var propFile string
	for _, entry := range c {
		name := entry.Name()
		if strings.Contains(name, ".json") && name != "gs_export.json" {
			propFile = dirName + "/" + name
		}

	}

	conInfo, csvFiles := createContainer.ParseJson(propFile)

	createContainer.Create(conInfo, force)
	containerName := conInfo.ContainerName
	types := mapping(conInfo.Columns)

	//var mapRows = make(map[string][][]string)

	for _, file := range csvFiles {
		fileName := dirName + "/" + file
		allRows, err := readAllRows(fileName)
		if err != nil {
			log.Fatal(err)
		}
		// mapRows[file] = allRows
		processCSV(allRows, types, containerName, fileName)
	}

}

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Migrate from GridDB CE Export Files to Cloud",
	Long:    "Use the export tool on your GridDB CE Instance to create the dir output of csv files and a properties file and then migrate that table to GridDB Cloud",
	Example: "griddb-cloud-cli migrate <directory>",
	Run: func(cmd *cobra.Command, args []string) {
		migrate(args[0])
	},
}
