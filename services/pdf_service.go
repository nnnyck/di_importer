package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"

	pdfpkg "rsc.io/pdf"

	"di_importer/models"
)

func ProcessPDFWithLines(pdf *multipart.FileHeader, linhas []models.DadosLinha) ([]models.DadosExtraidos, error) {
	fmt.Println("Processando PDF:", pdf.Filename)

	var resultados []models.DadosExtraidos

	// Abre o PDF
	file, err := pdf.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Lê todo o PDF em memória
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buf.Bytes())

	// Cria o leitor de PDF com rsc.io/pdf
	r, err := pdfpkg.NewReader(reader, int64(reader.Len()))
	if err != nil {
		return nil, err
	}

	// Extrai texto de todas as páginas
	var fullText strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		content := p.Content()
		for _, text := range content.Text {
			fullText.WriteString(text.S)
		}
	}
	texto := strings.Join(strings.Fields(fullText.String()), "") // normaliza

	// Regex para capturar campos
	reDataChegada := regexp.MustCompile(`Data\s*chegada\s*[-:>]*\s*(\d{2}/\d{2}/\d{4})`)
	reVMLE := regexp.MustCompile(`VMLE:\s*[A-Z\s]+([\d,.]+)`)
	reMAWB := regexp.MustCompile(`Ident\.\s*master\s*do\s*conhecimento\s*[-:>]*\s*(\d+)`)
	reHAWB := regexp.MustCompile(`Identificaçãodoconhecimento:?(\d+)`)
	reVolumes := regexp.MustCompile(`Quantidade:\s*(\d+)`)
	rePesoBruto := regexp.MustCompile(`PesoBruto:\s*([\d,.]+)`)
	rePesoLiquido := regexp.MustCompile(`PesoL[ií]quido:\s*([\d,.]+)`)

	formatPeso := func(p string) string {
		if p == "" {
			return ""
		}
		p = strings.ReplaceAll(p, ",", ".")
		f, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return p
		}
		s := fmt.Sprintf("%.3f", f)
		return strings.ReplaceAll(s, ".", ",")
	}

	// Processa cada linha da planilha
	for _, linha := range linhas {
		hawb := strings.ReplaceAll(linha.HAWB, " ", "")

		if hawb != "" && strings.Contains(texto, hawb) {
			fmt.Printf("✅ Encontrado HAWB %s no PDF %s (linha %d)\n",
				linha.HAWB, pdf.Filename, linha.Index)

			dados := models.DadosExtraidos{
				LinhaIndex: linha.Index,
				HAWB:       linha.HAWB,
			}

			// Extrai campos do PDF
			if match := reMAWB.FindStringSubmatch(texto); len(match) > 1 {
				dados.MAWB = match[1]
			}
			if match := reHAWB.FindStringSubmatch(texto); len(match) > 1 {
				dados.HAWB = match[1]
			}
			if match := reDataChegada.FindStringSubmatch(texto); len(match) > 1 {
				// converte dd/mm/yyyy -> dd/mm/yy
				dados.DataChegada = match[1][:6] + match[1][8:]
			}
			if match := reVMLE.FindStringSubmatch(texto); len(match) > 1 {
				dados.VMLE = match[1]
			}
			if match := reVolumes.FindStringSubmatch(texto); len(match) > 1 {
				if v, err := strconv.Atoi(match[1]); err == nil {
					dados.Volumes = fmt.Sprintf("%02d", v) // 2 dígitos
				} else {
					dados.Volumes = match[1]
				}
			}
			if match := rePesoBruto.FindStringSubmatch(texto); len(match) > 1 {
				dados.PesoBruto = formatPeso(match[1])
			}
			if match := rePesoLiquido.FindStringSubmatch(texto); len(match) > 1 {
				dados.PesoLiquido = formatPeso(match[1])
			}
			if strings.Contains(strings.ToLower(texto), "viracopos") {
				dados.Aeroporto = "VCP"
			}

			resultados = append(resultados, dados)
		}
	}

	return resultados, nil
}
