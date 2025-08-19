package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"di_importer/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)

	port := "8080"

	// Abre navegador automaticamente
	openBrowser("http://localhost:" + port)

	// Inicia servidor
	if err := r.Run(":" + port); err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		os.Exit(1)
	}
}

// openBrowser tenta abrir a URL no navegador padrão
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
