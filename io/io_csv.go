package io

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadCSV(filename string) []map[string]string {
	csvFile, _ := os.Open(filename)
	records, err := csv.NewReader(csvFile).Read()
	if err != nil {
		log.Fatal(err)
	}

	columns := records[0]
	rows := records[1:]

	var output []map[string]string

	record := make(map[string]string)

	for _, val := range rows {
		for idx, col := range columns {
			record[col] = val[idx]
		}
		output = append(output, record)
	}

	return output
}



