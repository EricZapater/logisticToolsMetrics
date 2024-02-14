package processor

import (
	"logisticToolsMetrics/config"
	Model "logisticToolsMetrics/model"
	"strconv"
	"strings"
	"time"
)

// Función para procesar las filas de datos y calcular las métricas
func ProcessData(dataRows []Model.DataRow) (Model.Metrics, Model.Metrics) {
    metrics := Model.Metrics{}
	rowMetrics := Model.Metrics{}
    for _ , dr := range dataRows {
        metric := evaluateRow(dr)				
		rowMetrics.AddMetric(metric)
        if metrics.Exists(metric.Project, metric.Size) {
            metrics.UpdateMetric(metric)
        } else {
            metrics.AddMetric(metric)
        }
    }
    return metrics, rowMetrics
}


func evaluateRow(dr Model.DataRow) (Model.Metric) {
    
    deviation := calculateDeviation(dr.Estimation, dr.InProgress, dr.Ready)
    leadTime := calculateLeadTime(dr.ToStart, dr.Ready)
    planningTime := calculatePlanningTime(dr.ToStart, dr.InProgress)
    cycleTime := calculateCycleTime(dr.InProgress, dr.ToVerify)
    developmentTime := calculateDevelopmentTime(dr.InProgress, dr.InReview)
    verifyingTime := calculateVerifyingTime(dr.ToVerify, dr.Ready)
	executionTime := calculateExecutionTime(dr.InProgress, dr.Ready)

    return Model.Metric{
		RowIndex: dr.RowIndex,
        Project:         dr.Project,
		Size: dr.StoryPoints,
		NumHUS: 1,
        LeadTime:        leadTime,
        PlanningTime:    planningTime,
        CycleTime:       cycleTime,
        DevelopmentTime: developmentTime,
        VerifyingTime:   verifyingTime,
		Deviation: deviation,
		ExecutionTime: executionTime,
    }
}

func replaceMonths(dateStr string) string {	
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
			break
		}
	}
	return dateStr
}

func calculateDeviation(estimation, inProgress, ready string) float64 {
	if estimation == ""{
		estimation = "0:0"
	}
	sestimation := strings.Split(estimation, ":")
	//drestimation, _ := strconv.ParseFloat(sestimation[0],64)
	
    estimationHours, _ := strconv.ParseFloat(sestimation[0],64)
	inProgressTime, _ := time.ParseInLocation(config.Layout, replaceMonths(inProgress), config.Location)
	readyTime, _ := time.ParseInLocation(config.Layout, replaceMonths(ready), config.Location)
	executionTime := readyTime.Sub(inProgressTime).Hours()    

    return executionTime - estimationHours
}

func calculateExecutionTime(inProgress, ready string) float64 {
	inProgressTime, _ := time.ParseInLocation(config.Layout, replaceMonths(inProgress), config.Location)
	readyTime, _ := time.ParseInLocation(config.Layout, replaceMonths(ready), config.Location)
	return readyTime.Sub(inProgressTime).Hours() 	
}

func calculateLeadTime(toStart, ready string) float64 {
    toStartTime, _ := time.ParseInLocation(config.Layout, replaceMonths(toStart), config.Location)
    readyTime, _ := time.ParseInLocation(config.Layout, replaceMonths(ready), config.Location)

    return readyTime.Sub(toStartTime).Hours() /24
}

func calculatePlanningTime(toStart, inProgress string) float64 {
    toStartTime, _ := time.ParseInLocation(config.Layout, replaceMonths(toStart), config.Location)
    inProgressTime, _ := time.ParseInLocation(config.Layout, replaceMonths(inProgress), config.Location)

    return inProgressTime.Sub(toStartTime).Hours() /24
}

func calculateCycleTime(inProgress, toVerify string) float64 {
    inProgressTime, _ := time.ParseInLocation(config.Layout, replaceMonths(inProgress), config.Location)
    toVerifyTime, _ := time.ParseInLocation(config.Layout, replaceMonths(toVerify), config.Location)

    return toVerifyTime.Sub(inProgressTime).Hours() /24
}

func calculateDevelopmentTime(inProgress, inReview string) float64 {
    inProgressTime, _ := time.ParseInLocation(config.Layout, replaceMonths(inProgress), config.Location)
    inReviewTime, _ := time.ParseInLocation(config.Layout, replaceMonths(inReview), config.Location)
    return inReviewTime.Sub(inProgressTime).Hours() /24
}

func calculateVerifyingTime(toVerify, ready string)float64{	
	inToVerify, _ := time.ParseInLocation(config.Layout, replaceMonths(toVerify), config.Location)
	inReadyTime, _ := time.ParseInLocation(config.Layout, replaceMonths(ready), config.Location)
	return inReadyTime.Sub(inToVerify).Hours()/24
}