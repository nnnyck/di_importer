package services

import (
	"di_importer/models"
	"fmt"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

func ProcessExcel(file multipart.File) ([]models.DadosLinha, error) {
	xl, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	sheetName := xl.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("planilha sem aba")
	}

	rows, err := xl.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var linhas []models.DadosLinha

	for i, row := range rows {
		if i == 0 {
			continue
		}

		hawb := ""
		if len(row) > 2 {
			hawb = row[2]
		}

		somenteHAWB := true
		for j := 0; j < len(row); j++ {
			if j != 2 && row[j] != "" {
				somenteHAWB = false
				break
			}
		}

		if somenteHAWB && hawb != "" {
			linhas = append(linhas, models.DadosLinha{
				HAWB:  hawb,
				Index: i + 1,
			})
		}
	}

	fmt.Println("Linhas com identificadores HAWB apenas:", linhas)
	return linhas, nil
}

func UpdateExcelWithData(xl *excelize.File, sheetName string, dados []models.DadosExtraidos) error {
	for _, d := range dados {
		if d.DataChegada != "" {
			xl.SetCellValue(sheetName, fmt.Sprintf("A%d", d.LinhaIndex), d.DataChegada)
		}
		if d.MAWB != "" {
			xl.SetCellValue(sheetName, fmt.Sprintf("B%d", d.LinhaIndex), d.MAWB)
		}

		if d.VMLE != "" {
			xl.SetCellValue(sheetName, fmt.Sprintf("D%d", d.LinhaIndex), d.VMLE)
		}
		if d.Volumes != "" {
			xl.SetCellValue(sheetName, fmt.Sprintf("E%d", d.LinhaIndex), d.Volumes)
		}
		if d.PesoBruto != "" {
			xl.SetCellValue(sheetName, fmt.Sprintf("F%d", d.LinhaIndex), d.PesoBruto)
		}

		if d.Aeroporto != "" {
			xl.SetCellValue(sheetName, fmt.Sprintf("I%d", d.LinhaIndex), d.Aeroporto)
		}
	}
	return nil
}
