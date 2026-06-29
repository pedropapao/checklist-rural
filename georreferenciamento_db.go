package main

func (app *App) criarTabelaArquivosGeorreferenciamento() error {
	_, err := app.DB.Exec(`
		CREATE TABLE IF NOT EXISTS arquivos_georef (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			reuniao_id INTEGER NOT NULL,
			nome_original TEXT DEFAULT '',
			nome_salvo TEXT DEFAULT '',
			tipo TEXT DEFAULT '',
			caminho TEXT DEFAULT '',
			observacao TEXT DEFAULT '',
			criado_em TEXT DEFAULT '',

			FOREIGN KEY(reuniao_id) REFERENCES reunioes(id)
		)
	`)
	return err
}
