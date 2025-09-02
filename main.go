package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"runtime"

	"di_importer/routes"

	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed web/views/*
var content embed.FS

func main() {
	r := gin.Default()

	htmlTemplates := template.Must(template.ParseFS(content, "web/views/*"))
	r.SetHTMLTemplate(htmlTemplates)

	r.StaticFS("/static", http.FS(content))

	routes.SetupRoutes(r)

	port := "8080"

	openBrowser("http://localhost:" + port)

	if err := r.Run(":" + port); err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		os.Exit(1)
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		fmt.Println("Não é possível abrir o navegador automaticamente")
		return
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		fmt.Println("Erro ao abrir navegador:", err)
	}
}
