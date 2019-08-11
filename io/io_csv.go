package io

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func ReadCSV(filename string) {
	csvFile, _ := os.Open(filename)
	r := csv.NewReader(csvFile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)
	}
}
