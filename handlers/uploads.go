package handlers

import (
	"fmt"
	"net/http"

	"di_importer/models"
	"di_importer/services"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// HomeHandler exibe formulário simples
func HomeHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func UploadHandler(c *gin.Context) {
	sheetFile, err := c.FormFile("spreadsheet")
	if err != nil {
		c.String(http.StatusBadRequest, "Erro ao receber a planilha")
		return
	}

	file, err := sheetFile.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao abrir a planilha")
		return
	}
	defer file.Close()

	xl, err := excelize.OpenReader(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao abrir a planilha")
		return
	}

	sheetName := xl.GetSheetName(0)

	sheet, err := sheetFile.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao abrir a planilha")
		return
	}
	defer sheet.Close()

	linhas, err := services.ProcessExcel(sheet)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao processar a planilha")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Erro ao processar o formulário")
		return
	}

	pdfs := form.File["pdfs"]
	var dadosExtraidos []models.DadosExtraidos
	for _, pdfFile := range pdfs {
		dados, err := services.ProcessPDFWithLines(pdfFile, linhas)
		if err != nil {
			fmt.Println("Erro ao processar PDF:", err)
			continue
		}
		dadosExtraidos = append(dadosExtraidos, dados...)
	}

	if err := services.UpdateExcelWithData(xl, sheetName, dadosExtraidos); err != nil {
		c.String(http.StatusInternalServerError, "Erro ao atualizar a planilha")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="saida.xlsx"`)
	c.Header("Content-Transfer-Encoding", "binary")

	if err := xl.Write(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Erro ao gerar arquivo final")
		return
	}
}
