package main

import (
	"fmt"
	"log"
	"logisticToolsMetrics/config"
	"logisticToolsMetrics/processor"
	"logisticToolsMetrics/reader"
	"logisticToolsMetrics/writer"
	"os"
	"time"

	"github.com/h2non/filetype"
	"github.com/xuri/excelize/v2"
)

func main() {
	files, err := os.ReadDir(config.Path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir(){
			continue
		}
		iniTime := time.Now()
        buf, err := os.ReadFile(fmt.Sprintf("%s/%s",config.Path,file.Name()))
		if err != nil {
			log.Fatal(err)
		}
  		kind, err := filetype.Match(buf)
		  if err != nil {
			log.Fatal(err)
		}
		if kind.Extension == "xlsx" {
			f, err := excelize.OpenFile(fmt.Sprintf("%s/%s",config.Path,file.Name()))
			if err != nil {		
				log.Fatal(err)
			}
			sheetName := f.GetSheetName(0)
			dataRows, err := reader.ReadData(f)			
			if err != nil {
				log.Fatal(err)
			}
			writer.WriteRowResumeHeader(f,sheetName)
			metrics, rowMetrics := processor.ProcessData(dataRows)		
			
			for _, rowMetric := range rowMetrics.List{				
				writer.WriteRowResume(f, sheetName, rowMetric.RowIndex, rowMetric)
			}

			writer.WriteMetrics(metrics, f)
			err = f.Save()
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
		}
		endTime := time.Now()
		fmt.Println("Execution time (s) for file: ", file.Name()," - ", endTime.Sub(iniTime).Seconds())
	}
}