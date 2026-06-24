package main

import (
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
)

func (app *App) telaPastas(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("pastas").Parse(htmlBase(pastasHTML)))

	dados := map[string]any{
		"Titulo":       "Pastas do sistema",
		"PastaDados":   app.PastaDados,
		"PastaExports": filepath.Join(app.PastaDados, "exports"),
		"PastaBackups": filepath.Join(app.PastaDados, "backups"),
	}

	tpl.Execute(w, dados)
}

func (app *App) abrirPastaDados(w http.ResponseWriter, r *http.Request) {
	abrirPasta(app.PastaDados)
	http.Redirect(w, r, "/pastas", http.StatusSeeOther)
}

func (app *App) abrirPastaExports(w http.ResponseWriter, r *http.Request) {
	abrirPasta(filepath.Join(app.PastaDados, "exports"))
	http.Redirect(w, r, "/pastas", http.StatusSeeOther)
}

func (app *App) abrirPastaBackups(w http.ResponseWriter, r *http.Request) {
	abrirPasta(filepath.Join(app.PastaDados, "backups"))
	http.Redirect(w, r, "/pastas", http.StatusSeeOther)
}

func abrirPasta(caminho string) {
	var comando *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		comando = exec.Command("xdg-open", caminho)
	case "windows":
		comando = exec.Command("explorer", caminho)
	case "darwin":
		comando = exec.Command("open", caminho)
	default:
		return
	}

	_ = comando.Start()
}
