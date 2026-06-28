package main

func (app *App) criarColunasTriagem() error {
	colunas := []string{
		"ALTER TABLE reunioes ADD COLUMN possui_caf TEXT DEFAULT 'nao'",
		"ALTER TABLE reunioes ADD COLUMN finalidade_credito TEXT DEFAULT ''",
		"ALTER TABLE reunioes ADD COLUMN valor_pretendido REAL DEFAULT 0",
		"ALTER TABLE reunioes ADD COLUMN tem_orcamento TEXT DEFAULT 'nao'",
	}

	for _, comando := range colunas {
		_, _ = app.DB.Exec(comando)
	}

	return nil
}
