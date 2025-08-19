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
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<h1>DI Importer</h1>
		<form action="/upload" method="post" enctype="multipart/form-data">
			<label>Selecione PDFs:</label><br>
			<input type="file" name="pdfs" multiple><br><br>
			<label>Selecione a planilha:</label><br>
			<input type="file" name="spreadsheet"><br><br>
			<input type="submit" value="Enviar">
		</form>
	`))
}

func UploadHandler(c *gin.Context) {
	// Recebe a planilha
	sheetFile, err := c.FormFile("spreadsheet")
	if err != nil {
		c.String(http.StatusBadRequest, "Erro ao receber a planilha")
		return
	}

	// Abrir planilha com excelize
	sheetPath := "./upload.xlsx"
	if err := c.SaveUploadedFile(sheetFile, sheetPath); err != nil {
		c.String(http.StatusInternalServerError, "Erro ao salvar a planilha")
		return
	}
	xl, err := excelize.OpenFile(sheetPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao abrir a planilha")
		return
	}
	sheetName := xl.GetSheetName(0)

	// Processa a planilha: cria slice de linhas com MAWB e/ou HAWB
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

	// Processa PDFs
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

	// Atualiza planilha
	if err := services.UpdateExcelWithData(xl, sheetName, dadosExtraidos); err != nil {
		c.String(http.StatusInternalServerError, "Erro ao atualizar a planilha")
		return
	}

	// Salvar e devolver arquivo atualizado
	outputPath := "./saida.xlsx"
	if err := xl.SaveAs(outputPath); err != nil {
		c.String(http.StatusInternalServerError, "Erro ao salvar a planilha final")
		return
	}

	c.File(outputPath)
}
