package main

func (app *App) criarTabelaDadosExternos() error {
	_, err := app.DB.Exec(`
		CREATE TABLE IF NOT EXISTS dados_externos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			reuniao_id INTEGER NOT NULL UNIQUE,

			cpf_cnpj TEXT DEFAULT '',
			cep TEXT DEFAULT '',
			razao_social TEXT DEFAULT '',
			nome_fantasia TEXT DEFAULT '',
			situacao_cadastral TEXT DEFAULT '',
			cnae TEXT DEFAULT '',
			descricao_cnae TEXT DEFAULT '',

			endereco TEXT DEFAULT '',
			bairro TEXT DEFAULT '',
			municipio TEXT DEFAULT '',
			uf TEXT DEFAULT '',
			codigo_ibge TEXT DEFAULT '',

			nome_imovel TEXT DEFAULT '',
			car TEXT DEFAULT '',
			situacao_car TEXT DEFAULT '',
			codigo_incra TEXT DEFAULT '',
			ccir TEXT DEFAULT '',
			nirf_cafir TEXT DEFAULT '',
			matricula TEXT DEFAULT '',
			area_total TEXT DEFAULT '',

			embargo_ibama TEXT DEFAULT 'nao_consultado',
			ibama_detalhes TEXT DEFAULT '',
			ceis_status TEXT DEFAULT 'nao_consultado',
			ceis_detalhes TEXT DEFAULT '',
			cnep_status TEXT DEFAULT 'nao_consultado',
			cnep_detalhes TEXT DEFAULT '',
			sncr_status TEXT DEFAULT 'nao_consultado',
			sncr_detalhes TEXT DEFAULT '',
			sigef_geo TEXT DEFAULT 'nao_consultado',

			fonte_consulta TEXT DEFAULT '',
			link_fonte TEXT DEFAULT '',
			observacoes TEXT DEFAULT '',
			ultima_investigacao_online TEXT DEFAULT '',
			atualizado_em TEXT DEFAULT '',

			FOREIGN KEY(reuniao_id) REFERENCES reunioes(id)
		)
	`)

	if err != nil {
		return err
	}

	return app.criarColunasDadosExternos()
}

func (app *App) criarColunasDadosExternos() error {
	colunas := []string{
		"ALTER TABLE dados_externos ADD COLUMN nome_imovel TEXT DEFAULT ''",
		"ALTER TABLE dados_externos ADD COLUMN situacao_car TEXT DEFAULT ''",
		"ALTER TABLE dados_externos ADD COLUMN ibama_detalhes TEXT DEFAULT ''",
		"ALTER TABLE dados_externos ADD COLUMN ceis_status TEXT DEFAULT 'nao_consultado'",
		"ALTER TABLE dados_externos ADD COLUMN ceis_detalhes TEXT DEFAULT ''",
		"ALTER TABLE dados_externos ADD COLUMN cnep_status TEXT DEFAULT 'nao_consultado'",
		"ALTER TABLE dados_externos ADD COLUMN cnep_detalhes TEXT DEFAULT ''",
		"ALTER TABLE dados_externos ADD COLUMN sncr_status TEXT DEFAULT 'nao_consultado'",
		"ALTER TABLE dados_externos ADD COLUMN sncr_detalhes TEXT DEFAULT ''",
		"ALTER TABLE dados_externos ADD COLUMN link_fonte TEXT DEFAULT ''",
		"ALTER TABLE dados_externos ADD COLUMN ultima_investigacao_online TEXT DEFAULT ''",
	}

	for _, comando := range colunas {
		_, _ = app.DB.Exec(comando)
	}

	return nil
}
