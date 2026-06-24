package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	pastaDados, err := criarPastaDados()
	if err != nil {
		log.Fatal(err)
	}

	caminhoBanco := filepath.Join(pastaDados, "checklist_rural.db")

	db, err := sql.Open("sqlite", caminhoBanco)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &App{
		DB:         db,
		PastaDados: pastaDados,
	}

	err = app.criarTabelas()
	if err != nil {
		log.Fatal(err)
	}

	err = app.criarTabelaItensChecklist()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", app.telaInicial)
	http.HandleFunc("/nova-reuniao", app.telaNovaReuniao)
	http.HandleFunc("/salvar-reuniao", app.salvarReuniao)

	http.HandleFunc("/editar-reuniao", app.telaEditarReuniao)
	http.HandleFunc("/atualizar-reuniao", app.atualizarReuniao)

	http.HandleFunc("/confirmar-excluir", app.telaConfirmarExcluir)
	http.HandleFunc("/excluir-reuniao", app.excluirReuniao)

	http.HandleFunc("/reunioes", app.listarReunioes)
	http.HandleFunc("/backups", app.telaBackups)
	http.HandleFunc("/pastas", app.telaPastas)
	http.HandleFunc("/abrir-pasta-dados", app.abrirPastaDados)
	http.HandleFunc("/abrir-pasta-exports", app.abrirPastaExports)
	http.HandleFunc("/abrir-pasta-backups", app.abrirPastaBackups)
	http.HandleFunc("/criar-backup", app.criarBackupBanco)
	http.HandleFunc("/confirmar-restaurar-backup", app.telaConfirmarRestaurarBackup)
	http.HandleFunc("/restaurar-backup", app.restaurarBackupBanco)
	http.HandleFunc("/detalhes", app.telaDetalhes)

	http.HandleFunc("/checklist", app.telaChecklist)
	http.HandleFunc("/gerar-itens-checklist", app.gerarItensChecklist)
	http.HandleFunc("/checklist-controle", app.telaChecklistControle)
	http.HandleFunc("/salvar-itens-checklist", app.salvarItensChecklist)

	http.HandleFunc("/whatsapp", app.telaWhatsApp)
	http.HandleFunc("/whatsapp-pendencias", app.telaWhatsAppPendencias)
	http.HandleFunc("/abrir-whatsapp-pendencias", app.abrirWhatsAppPendencias)

	http.HandleFunc("/exportar-checklist-txt", app.exportarChecklistTXT)
	http.HandleFunc("/exportar-checklist-controle-txt", app.exportarChecklistControleTXT)
	http.HandleFunc("/exportar-whatsapp-txt", app.exportarWhatsAppTXT)
	http.HandleFunc("/exportar-whatsapp-pendencias-txt", app.exportarWhatsAppPendenciasTXT)

	endereco := "http://localhost:8080"

	fmt.Println("Checklist Rural iniciado.")
	fmt.Println("Banco de dados:", caminhoBanco)
	fmt.Println("Acesse:", endereco)

	abrirNavegador(endereco)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
