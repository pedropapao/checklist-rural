package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ArquivoBackup struct {
	Nome       string
	Caminho    string
	Tamanho    int64
	Modificado string
}

func (app *App) telaBackups(w http.ResponseWriter, r *http.Request) {
	backups, err := app.listarBackups()
	if err != nil {
		http.Error(w, "Erro ao listar backups: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tpl := template.Must(template.New("backups").Parse(htmlBase(backupsHTML)))

	dados := map[string]any{
		"Titulo":  "Backups",
		"Backups": backups,
	}

	tpl.Execute(w, dados)
}

func (app *App) criarBackupBanco(w http.ResponseWriter, r *http.Request) {
	caminhoBanco := filepath.Join(app.PastaDados, "checklist_rural.db")
	pastaBackups := filepath.Join(app.PastaDados, "backups")

	err := os.MkdirAll(pastaBackups, 0755)
	if err != nil {
		http.Error(w, "Erro ao criar pasta de backup: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dataHora := time.Now().Format("2006-01-02_15-04-05")
	nomeBackup := fmt.Sprintf("checklist_rural_backup_%s.db", dataHora)
	caminhoBackup := filepath.Join(pastaBackups, nomeBackup)

	err = copiarArquivo(caminhoBanco, caminhoBackup)
	if err != nil {
		http.Error(w, "Erro ao criar backup: "+err.Error(), http.StatusInternalServerError)
		return
	}

	app.telaArquivoGerado(w, "Backup criado com sucesso", caminhoBackup)
}

func (app *App) telaConfirmarRestaurarBackup(w http.ResponseWriter, r *http.Request) {
	nomeBackup := r.URL.Query().Get("arquivo")
	if nomeBackup == "" {
		http.Error(w, "Arquivo de backup não informado", http.StatusBadRequest)
		return
	}

	caminhoBackup := filepath.Join(app.PastaDados, "backups", nomeBackup)

	_, err := os.Stat(caminhoBackup)
	if err != nil {
		http.Error(w, "Backup não encontrado", http.StatusBadRequest)
		return
	}

	tpl := template.Must(template.New("confirmarRestaurarBackup").Parse(htmlBase(confirmarRestaurarBackupHTML)))

	dados := map[string]any{
		"Titulo":  "Confirmar restauração",
		"Arquivo": nomeBackup,
		"Caminho": caminhoBackup,
	}

	tpl.Execute(w, dados)
}

func (app *App) restaurarBackupBanco(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/backups", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao ler formulário", http.StatusBadRequest)
		return
	}

	nomeBackup := r.FormValue("arquivo")
	if nomeBackup == "" {
		http.Error(w, "Arquivo de backup não informado", http.StatusBadRequest)
		return
	}

	caminhoBanco := filepath.Join(app.PastaDados, "checklist_rural.db")
	pastaBackups := filepath.Join(app.PastaDados, "backups")
	caminhoBackup := filepath.Join(pastaBackups, nomeBackup)

	_, err = os.Stat(caminhoBackup)
	if err != nil {
		http.Error(w, "Backup não encontrado", http.StatusBadRequest)
		return
	}

	dataHora := time.Now().Format("2006-01-02_15-04-05")
	nomeSeguranca := fmt.Sprintf("antes_de_restaurar_%s.db", dataHora)
	caminhoSeguranca := filepath.Join(pastaBackups, nomeSeguranca)

	err = copiarArquivo(caminhoBanco, caminhoSeguranca)
	if err != nil {
		http.Error(w, "Erro ao criar backup de segurança antes da restauração: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = copiarArquivo(caminhoBackup, caminhoBanco)
	if err != nil {
		http.Error(w, "Erro ao restaurar backup: "+err.Error(), http.StatusInternalServerError)
		return
	}

	app.telaArquivoGerado(w, "Backup restaurado com sucesso. Reinicie o sistema.", caminhoBackup)
}

func (app *App) listarBackups() ([]ArquivoBackup, error) {
	pastaBackups := filepath.Join(app.PastaDados, "backups")

	err := os.MkdirAll(pastaBackups, 0755)
	if err != nil {
		return nil, err
	}

	entradas, err := os.ReadDir(pastaBackups)
	if err != nil {
		return nil, err
	}

	var backups []ArquivoBackup

	for _, entrada := range entradas {
		if entrada.IsDir() {
			continue
		}

		info, err := entrada.Info()
		if err != nil {
			continue
		}

		caminho := filepath.Join(pastaBackups, entrada.Name())

		backups = append(backups, ArquivoBackup{
			Nome:       entrada.Name(),
			Caminho:    caminho,
			Tamanho:    info.Size(),
			Modificado: info.ModTime().Format("02/01/2006 15:04"),
		})
	}

	return backups, nil
}

func copiarArquivo(origem string, destino string) error {
	arquivoOrigem, err := os.Open(origem)
	if err != nil {
		return err
	}
	defer arquivoOrigem.Close()

	arquivoDestino, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer arquivoDestino.Close()

	_, err = io.Copy(arquivoDestino, arquivoOrigem)
	if err != nil {
		return err
	}

	return arquivoDestino.Sync()
}
