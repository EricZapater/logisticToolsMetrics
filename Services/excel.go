package services

import (
	"log"
	model "logisticToolsMetrics/Model"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/xuri/excelize/v2"
)

const layout = "02/Jan/06 15:04"

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatalf("Error al cargar la zona horaria: %v", err)
	}
}

func replaceMonths(dateStr string) string {
	// Map of Spanish month abbreviations to English
	spanishToEnglishMonths := map[string]string{
		"ene": "Jan",
		"feb": "Feb",
		"mar": "Mar",
		"abr": "Apr",
		"may": "May",
		"jun": "Jun",
		"jul": "Jul",
		"ago": "Aug",
		"sep": "Sep",
		"oct": "Oct",
		"nov": "Nov",
		"dic": "Dec",
	}

	for spanish, english := range spanishToEnglishMonths {
		if strings.Contains(dateStr, spanish) {
			dateStr = strings.ReplaceAll(dateStr, spanish, english)
			break // Assume only one month abbreviation is present in the date string
		}
	}
	return dateStr
}

func OpenFile(path string)(*excelize.File, error){
	f, err := excelize.OpenFile(path)
	if err != nil {		
		return nil, err
	}
	return f, nil
}


func ProcessData(f *excelize.File)([][]string, error){

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil,err
	}	
	f.SetCellValue(sheet,"Q1","ExecutionTime (h)")
	f.SetCellValue(sheet,"R1","Deviation (h)")
	f.SetCellValue(sheet,"S1","LeadTime (d)")
	f.SetCellValue(sheet,"T1","PlanningTime (d)")
	f.SetCellValue(sheet,"U1","CycleTime (d)")
	f.SetCellValue(sheet,"V1", "DevelopmentTime (d)")
	f.SetCellValue(sheet,"W1","VerifyingTime (d)")	
	metriques := new(model.Metriques)	
	for i, row := range rows {
		if len(row)<13 || i == 0 || row[1] != "Historia"{
			continue
		}
		estimation := "0"
		sumestimation := "0"	
		fmt.Println(i, " - ", len(row), " - ", row[3])	
		if len(row)>=14{
			estimation = row[13]
			sumestimation = row[14]
		}
		dataRow := &model.DataRow{
			Project: row[0],
			Type: row[1],
			Key: row[2],
			Resume: row[3],
			Status: row[4],
			Sprint: row[5],
			ToStart: replaceMonths(row[6]),
			InProgress: replaceMonths(row[7]),
			InReview: replaceMonths(row[8]),
			ToVerify: replaceMonths(row[9]),
			Ready: replaceMonths(row[10]),
			TimesInRework: row[11],
			StoryPoints: row[12],
			Estimation: estimation,
			SumEstimation: sumestimation,
		}
		//Check estimation and dev time
		if dataRow.Type == "Épica"{
			continue
		}
		m := evaluateRow(i, *dataRow, f)
		if m.Projecte != ""{
			exist := metriques.Exists(m.Projecte, m.Mida);
			if exist {
				metriques.Update(m)
			}else{
				metriques.AddMetrica(m)
			}
		}
		
		
	}
	AddMetrica(*metriques, f)
	return rows, nil
}

func evaluateRow(i int, dr model.DataRow, f *excelize.File)model.Metrica{
	style, _ := f.NewStyle(&excelize.Style{Fill: excelize.Fill{	Type:"pattern",
																	Color: []string{"#FF3333"}, 
																	Pattern: 1}})
	sheet := f.GetSheetName(0)	
	if(dr.ToStart == "" || dr.InProgress == "" || dr.InReview == "" || dr.ToVerify == "" || dr.Ready == ""){		
		if err := f.SetCellStyle(sheet, fmt.Sprintf("A%d", i+1), fmt.Sprintf("O%d", i+1), style); err != nil {
			fmt.Println("Error applying style to row:", err)
		}
		return model.Metrica{}
	}
	if dr.Estimation == ""{
		dr.Estimation = "0:0"
	}
	sestimation := strings.Split(dr.Estimation, ":")
	estimation, err := strconv.ParseFloat(sestimation[0],64)
	if err != nil {
		log.Fatal(err)
	}


	iniLeadTime, _ := time.ParseInLocation(layout, dr.ToStart, location)
	endPlanningTime, _ := time.ParseInLocation(layout, dr.InProgress, location)
	endDevelopmentTime, _ := time.ParseInLocation(layout, dr.InReview, location)
	endCycleTime, _ := time.ParseInLocation(layout, dr.ToVerify, location)
	endLeadTime, _ := time.ParseInLocation(layout, dr.Ready, location)	
	
	executionTime := endLeadTime.Sub(endPlanningTime).Hours()
	
	deviation := executionTime - estimation 
	leadTime := endLeadTime.Sub(iniLeadTime).Hours()/24
	planningTime := endPlanningTime.Sub(iniLeadTime).Hours()/24
	cycleTime := endCycleTime.Sub(endPlanningTime).Hours()/24
	developmentTime := endDevelopmentTime.Sub(endPlanningTime).Hours()/24
	verifyingTime := endLeadTime.Sub(endCycleTime).Hours()/24
	
	if err := f.SetCellValue(sheet, fmt.Sprintf("Q%d", i+1), executionTime); err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("R%d", i+1), deviation); err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("S%d", i+1), leadTime); err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("T%d", i+1), planningTime); err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("U%d", i+1), cycleTime); err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("V%d", i+1), developmentTime); err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellValue(sheet, fmt.Sprintf("W%d", i+1), verifyingTime); err != nil {
		log.Fatal(err)
	}
	
	metrica := model.Metrica{
		Projecte: dr.Project,      
		Mida: dr.StoryPoints,            
		NumHUS: 1,
		LeadTime: leadTime,      
		PlanningTime: planningTime,
		CycleTime: cycleTime,
		DevelopmentTime: developmentTime, 
		VerifyingTime: verifyingTime,
	}
	return metrica
}

func AddMetrica(m model.Metriques, f *excelize.File){
	sheetName := "Metriques"
	f.NewSheet(sheetName)
	sheetIdx,_ := f.GetSheetIndex(sheetName)
	f.SetActiveSheet(sheetIdx)
	titles := []string{"Projecte", "Mida", "NumHUS", "LeadTime", "PlanningTime","CycleTime","DevelopmentTime", "VerifyingTime"}
	for i, title := range titles {
		// Define column titles
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, title)
	}

	// Write data
	for i, m := range m.Projectes {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), m.Projecte)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), m.Mida)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), m.NumHUS)						
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), (m.LeadTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), (m.PlanningTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), (m.CycleTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), (m.DevelopmentTime/float64(m.NumHUS)))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), (m.VerifyingTime/float64(m.NumHUS)))
	}
}