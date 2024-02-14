package reader

import (
	Model "logisticToolsMetrics/model"

	"github.com/xuri/excelize/v2"
)

// Función para leer las filas de datos del archivo Excel
func ReadData(f *excelize.File) ([]Model.DataRow, error) {
    sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil,err
	}	
	dataRows := []Model.DataRow{}
	for i, row := range rows {		
		if len(row)<13 || i == 0 || row[1] != "Historia" || row[6] == "" || row[7] == "" || row[8] == "" || row[9] == "" || row[10] == ""{
			continue
		}
		estimation := "0"
		sumestimation := "0"	
		
		if len(row)>=14{
			estimation = row[13]
			sumestimation = row[14]
		}
		dataRow := Model.DataRow{
			RowIndex: i,
			Project: row[0],
			Type: row[1],
			Key: row[2],
			Resume: row[3],
			Status: row[4],
			Sprint: row[5],
			ToStart: row[6],
			InProgress: row[7],
			InReview: row[8],
			ToVerify: row[9],
			Ready: row[10],
			TimesInRework: row[11],
			StoryPoints: row[12],
			Estimation: estimation,
			SumEstimation: sumestimation,
		}
		dataRows = append(dataRows, dataRow)        
    }

    return dataRows, nil
}
