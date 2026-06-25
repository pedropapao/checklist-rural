package main

import "database/sql"

type App struct {
	DB         *sql.DB
	PastaDados string
}

type Reuniao struct {
	ID        int
	Produtor  string
	Telefone  string
	Municipio string
	UF        string

	Banco       string
	TipoProjeto string
	Atividade   string

	RendaAnual            float64
	ClassificacaoProdutor string

	CadastroBanco      string
	FinanciamentoAtivo string
	RestricaoCadastral string
	ImovelProprio      string
	ImovelArrendado    string
	TemCAR             string
	UsaAgua            string
	TemPecuaria        string
	TemInvestimento    string
	TemObra            string
	TemSupressao       string
	PrecisaZARC        string

	Observacoes string
	CriadoEm    string
}
