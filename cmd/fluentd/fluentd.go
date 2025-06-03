package fluentd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Imisrael/griddb-cloud-cli/cmd"
	"github.com/Imisrael/griddb-cloud-cli/cmd/containerInfo"
	"github.com/Imisrael/griddb-cloud-cli/cmd/putRow"
	"github.com/spf13/cobra"
)

var (
	keys          []string
	containerName string
)

func init() {
	cmd.RootCmd.AddCommand(fluentCmd)
	fluentCmd.Flags().StringSliceVar(&keys, "keys", []string{}, "Key names aka header names from your cols in GridDB, separated by comma like this: 'col1,col2,col3'")
	fluentCmd.Flags().StringVarP(&containerName, "name", "n", "", "containerName to create and save into")
	fluentCmd.MarkFlagRequired("keys")
	fluentCmd.MarkFlagRequired("name")
}

func isDateValue(stringDate string) bool {
	_, err := time.Parse("01/02/2006", stringDate)
	return err == nil
}

func isNumValue(v string) bool {
	_, err := strconv.ParseInt(v, 10, 64)
	return err == nil
}

func isFloatValue(v string) bool {
	_, err := strconv.ParseFloat(v, 64)
	return err == nil
}

func createCols(mapOfTypes map[string]string) []cmd.ContainerInfoColumns {

	log.Println("gonna try to creates cols")

	// We need to add +1 because we're adding the NOW timestamp for every container
	var columnInfo []cmd.ContainerInfoColumns = make([]cmd.ContainerInfoColumns, len(mapOfTypes)+1)

	columnInfo[0].Name = "ts"
	columnInfo[0].Type = "TIMESTAMP"

	log.Println(columnInfo)

	i := 1
	for key, value := range mapOfTypes {
		columnInfo[i].Name = key
		columnInfo[i].Type = value
		i++
		log.Println(columnInfo)
	}
	log.Println(columnInfo)
	return columnInfo
}

func createContainer(conInfo cmd.ContainerInfo) {

	log.Println("New container going up")

	jsonContainerInfo, err := json.Marshal(conInfo)
	if err != nil {
		log.Println("Error", err)
	}
	log.Println(string(jsonContainerInfo))
	convert := []byte(jsonContainerInfo)
	buf := bytes.NewBuffer(convert)

	client := &http.Client{}
	req, err := cmd.MakeNewRequest("POST", "/containers", buf)
	if err != nil {
		log.Println("Error making new request", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error with client DO: ", err)
	}

	log.Println("Checking for errors")
	cmd.CheckForErrors(resp)

	log.Println(resp.Status)
}

func readTSVFile(filename string) ([]byte, error) {
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

func parseTSV(data []byte) (*csv.Reader, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = '\t'
	return reader, nil
}

func parseValues(str, key string) string {
	//--keys host,ident,pid,message
	// 2025/05/28 19:05:06 [ aptupdate CRON 6559 pam_unix(cron:session): session closed for user root]
	// 2025/05/28 19:05:06 Val:
	// 2025/05/28 19:05:06 Val: aptupdate
	// 2025/05/28 19:05:06 Val: CRON
	// 2025/05/28 19:05:06 Val: 6559
	// 2025/05/28 19:05:06 Val: pam_unix(cron:session): session closed for user root

	// for each key, check the value of the corresponding value
	// then once determined, create table with these col names and types
	// and then finally push these values with the timestamps

	log.Printf("key: %s, val: %s", key, str)

	switch {
	case isDateValue(str):
		return "TIMESTAMP"
	case isNumValue(str):
		return "LONG"
	case isFloatValue(str):
		return "DOUBLE"
	default:
		return "STRING"

	}

}

func putMultiRows(arrayString, containerName string) {

	url := "/containers/" + containerName + "/rows"
	convert := []byte(arrayString)
	buf := bytes.NewBuffer(convert)

	client := &http.Client{}
	req, err := cmd.MakeNewRequest("PUT", url, buf)
	if err != nil {
		fmt.Println("Error making new request", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error with client DO: ", err)
	}

	cmd.CheckForErrors(resp)

	fmt.Println(resp.Status)
}

func ingest(fileName string) {

	fmt.Println("Reading " + fileName)

	data, err := readTSVFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	reader, err := parseTSV(data)
	if err != nil {
		fmt.Println("Error creating CSV reader:", err)
		return
	}
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	log.Println("Starting to iterate over lines")

	if containerInfo.ContainerExists(containerName) {
		for _, line := range lines {
			log.Println("Pushing Row (Container exists!)")
			jsonStr := putRow.BuildPutRowStrNonInteractive(containerName, line)
			log.Println(jsonStr)
			putMultiRows(jsonStr, containerName)
		}

	} else {
		var mapOfTypes map[string]string = make(map[string]string)
		line := lines[0]
		log.Println(line)
		log.Println(keys)

		var vals = []string{}
		for _, val := range line {
			if val != "" {
				log.Println("Val: " + val)
				vals = append(vals, val)
			}
		}
		log.Println(len(keys))
		log.Println(len(vals))
		for idx, v := range vals {
			gtype := parseValues(v, keys[idx])
			mapOfTypes[keys[idx]] = gtype
		}
		log.Println(mapOfTypes)

		var conInfo cmd.ContainerInfo
		colInfos := createCols(mapOfTypes)

		conInfo.ContainerName = containerName
		conInfo.ContainerType = "TIME_SERIES"
		conInfo.RowKey = true
		conInfo.Columns = colInfos
		createContainer(conInfo)

		for _, line := range lines {
			log.Println("Pushing Row")
			jsonStr := putRow.BuildPutRowStrNonInteractive(containerName, line)
			log.Println(jsonStr)
			putMultiRows(jsonStr, containerName)
		}

	}

}

var fluentCmd = &cobra.Command{
	Use:     "fluentd",
	Short:   "Ingest a `csv` file to a new or existing container",
	Long:    "Ingesting a csv file. You decide the type mapping an placement based on the interactive CLI options",
	Example: "griddb-cloud-cli fluentd --keys col1,col2,col3",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.OpenFile("/tmp/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
		ingest(args[0])
	},
}
