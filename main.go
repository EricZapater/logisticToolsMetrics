package main

import (
	"fmt"
	"log"
	services "logisticToolsMetrics/Services"
	"os"
	"time"

	"github.com/h2non/filetype"
)


const path = "./Files"

func main(){
	
	//fmt.Println(location)
	files, err := os.ReadDir(path)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
		iniTime := time.Now()
        buf, _ := os.ReadFile(fmt.Sprintf("%s/%s",path,file.Name()))
  		kind, _ := filetype.Match(buf)
		if kind.Extension == "xlsx" {
			f, err := services.OpenFile(fmt.Sprintf("%s/%s",path,file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			_, err = services.ProcessData(f)
			if err != nil{
				log.Fatal(err)
			}
			if err := f.SaveAs(fmt.Sprintf("%s/%s",path,file.Name())); err != nil {
				log.Fatal(err)
			}
		
			f.Close()
		}
		endTime := time.Now()
		fmt.Println("Execution time (s) for file: ", file.Name()," - ", endTime.Sub(iniTime).Seconds())
    }
}