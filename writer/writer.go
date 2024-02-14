package writer

import (
	"fmt"
	Model "logisticToolsMetrics/model"

	"github.com/xuri/excelize/v2"
)

func WriteMetrics(metrics Model.Metrics, f *excelize.File){
	sheetName := "Results"
	f.NewSheet(sheetName)
	writeHeader(f, sheetName)
	writeData(f, sheetName, metrics)
}

func writeHeader(f *excelize.File, sheetName string){
	titles := []string{"Project", "Size", "NumHUS", "LeadTime", "PlanningTime","CycleTime","DevelopmentTime", "VerifyingTime"}
	for i, title := range titles {		
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, title)
	}
}

func writeData(f *excelize.File, sheetName string, metrics Model.Metrics){
	for i, m := range metrics.List {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), m.Project)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), m.Size)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), m.NumHUS)						
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), (m.LeadTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), (m.PlanningTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), (m.CycleTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), (m.DevelopmentTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), (m.VerifyingTime/float64(m.NumHUS)))
	}
}

func WriteRowResume(f *excelize.File, sheetname string, rowIndex int, metric Model.Metric){
	f.SetCellValue(sheetname, fmt.Sprintf("Q%d", rowIndex+1), metric.ExecutionTime)	
	f.SetCellValue(sheetname, fmt.Sprintf("R%d", rowIndex+1), metric.Deviation)
	f.SetCellValue(sheetname, fmt.Sprintf("S%d", rowIndex+1), metric.LeadTime)
	f.SetCellValue(sheetname, fmt.Sprintf("T%d", rowIndex+1), metric.PlanningTime)	
	f.SetCellValue(sheetname, fmt.Sprintf("U%d", rowIndex+1), metric.CycleTime)
	f.SetCellValue(sheetname, fmt.Sprintf("V%d", rowIndex+1), metric.DevelopmentTime)
	f.SetCellValue(sheetname, fmt.Sprintf("W%d", rowIndex+1), metric.VerifyingTime)
}
func WriteRowResumeHeader(f *excelize.File, sheetname string){
	f.SetCellValue(sheetname,"Q1","ExecutionTime (h)")
	f.SetCellValue(sheetname,"R1","Deviation (h)")
	f.SetCellValue(sheetname,"S1","LeadTime (d)")
	f.SetCellValue(sheetname,"T1","PlanningTime (d)")
	f.SetCellValue(sheetname,"U1","CycleTime (d)")
	f.SetCellValue(sheetname,"V1", "DevelopmentTime (d)")
	f.SetCellValue(sheetname,"W1","VerifyingTime (d)")	
}