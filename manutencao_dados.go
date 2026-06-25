package main

func (app *App) corrigirClassificacaoReunioesAntigas() error {
	_, err := app.DB.Exec(`
		UPDATE reunioes
		SET classificacao_produtor = CASE
			WHEN COALESCE(renda_anual, 0) <= 0 THEN 'Não informado'
			WHEN renda_anual <= 500000 THEN 'Pequeno produtor'
			WHEN renda_anual <= 3500000 THEN 'Médio produtor'
			ELSE 'Grande produtor'
		END
		WHERE classificacao_produtor IS NULL
		   OR classificacao_produtor = ''
		   OR classificacao_produtor = 'Não disponível'
	`)

	return err
}
