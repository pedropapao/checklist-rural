package main

import (
	"strconv"
	"strings"
)

func normalizarValorMonetario(texto string) (float64, error) {
	texto = strings.TrimSpace(texto)

	if texto == "" {
		return 0, nil
	}

	texto = strings.ReplaceAll(texto, "R$", "")
	texto = strings.ReplaceAll(texto, " ", "")
	texto = strings.ReplaceAll(texto, ".", "")
	texto = strings.ReplaceAll(texto, ",", ".")

	return strconv.ParseFloat(texto, 64)
}

func classificarProdutorPorRBA(rba float64) string {
	if rba <= 0 {
		return "Não informado"
	}

	if rba <= 500000 {
		return "Pequeno produtor"
	}

	if rba <= 3500000 {
		return "Médio produtor"
	}

	return "Grande produtor"
}

func (app *App) criarColunasClassificacaoProdutor() error {
	colunas := []string{
		"ALTER TABLE reunioes ADD COLUMN renda_anual REAL DEFAULT 0",
		"ALTER TABLE reunioes ADD COLUMN classificacao_produtor TEXT DEFAULT ''",
	}

	for _, comando := range colunas {
		_, _ = app.DB.Exec(comando)
	}

	return nil
}
