package main

import (
	"os"
	"path/filepath"
)

func criarPastaDados() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pasta := filepath.Join(home, "Documentos", "ChecklistRural")

	err = os.MkdirAll(pasta, 0755)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(pasta, "exports"), 0755)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(pasta, "backups"), 0755)
	if err != nil {
		return "", err
	}

	return pasta, nil
}

func (app *App) criarTabelas() error {
	sqlTabelas := `
	CREATE TABLE IF NOT EXISTS reunioes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		produtor TEXT NOT NULL,
		telefone TEXT,
		municipio TEXT,
		uf TEXT,
		banco TEXT,
		tipo_projeto TEXT,
		atividade TEXT,
		cadastro_banco TEXT DEFAULT 'nao',
		financiamento_ativo TEXT DEFAULT 'nao',
		restricao_cadastral TEXT DEFAULT 'nao',
		imovel_proprio TEXT DEFAULT 'nao',
		imovel_arrendado TEXT DEFAULT 'nao',
		tem_car TEXT DEFAULT 'nao',
		usa_agua TEXT DEFAULT 'nao',
		tem_pecuaria TEXT DEFAULT 'nao',
		tem_investimento TEXT DEFAULT 'nao',
		tem_obra TEXT DEFAULT 'nao',
		tem_supressao TEXT DEFAULT 'nao',
		precisa_zarc TEXT DEFAULT 'nao',
		observacoes TEXT,
		criado_em TEXT NOT NULL
	);
	`

	_, err := app.DB.Exec(sqlTabelas)
	if err != nil {
		return err
	}

	colunas := []string{
		"cadastro_banco TEXT DEFAULT 'nao'",
		"financiamento_ativo TEXT DEFAULT 'nao'",
		"restricao_cadastral TEXT DEFAULT 'nao'",
		"imovel_proprio TEXT DEFAULT 'nao'",
		"imovel_arrendado TEXT DEFAULT 'nao'",
		"tem_car TEXT DEFAULT 'nao'",
		"usa_agua TEXT DEFAULT 'nao'",
		"tem_pecuaria TEXT DEFAULT 'nao'",
		"tem_investimento TEXT DEFAULT 'nao'",
		"tem_obra TEXT DEFAULT 'nao'",
		"tem_supressao TEXT DEFAULT 'nao'",
		"precisa_zarc TEXT DEFAULT 'nao'",
	}

	for _, coluna := range colunas {
		_, _ = app.DB.Exec("ALTER TABLE reunioes ADD COLUMN " + coluna)
	}

	return nil
}
