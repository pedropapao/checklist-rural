package main

import (
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

func abrirNavegador(url string) {
	var comando *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		comando = exec.Command("xdg-open", url)
	case "windows":
		comando = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		comando = exec.Command("open", url)
	default:
		return
	}

	_ = comando.Start()
}

func checkValor(r *http.Request, nome string) string {
	if r.FormValue(nome) == "sim" {
		return "sim"
	}

	return "nao"
}

func simNao(valor string) string {
	if valor == "sim" {
		return "Sim"
	}

	return "Não"
}

func limparNomeArquivo(texto string) string {
	texto = strings.TrimSpace(texto)
	texto = strings.ToLower(texto)

	trocas := map[string]string{
		"á": "a",
		"à": "a",
		"ã": "a",
		"â": "a",
		"é": "e",
		"ê": "e",
		"í": "i",
		"ó": "o",
		"ô": "o",
		"õ": "o",
		"ú": "u",
		"ç": "c",
	}

	for antigo, novo := range trocas {
		texto = strings.ReplaceAll(texto, antigo, novo)
	}

	invalidos := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", ".", ",", ";"}

	for _, invalido := range invalidos {
		texto = strings.ReplaceAll(texto, invalido, "")
	}

	texto = strings.ReplaceAll(texto, " ", "_")

	for strings.Contains(texto, "__") {
		texto = strings.ReplaceAll(texto, "__", "_")
	}

	if texto == "" {
		texto = "arquivo"
	}

	return texto
}
