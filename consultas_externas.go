package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ResultadoConsultaExterna struct {
	Tipo     string
	Sucesso  bool
	Mensagem string
	Dados    map[string]string
	DataHora string
}

type viacepResposta struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	Erro        bool   `json:"erro"`
}

type brasilAPICNPJResposta struct {
	CNPJ              string `json:"cnpj"`
	RazaoSocial       string `json:"razao_social"`
	NomeFantasia      string `json:"nome_fantasia"`
	DescricaoSituacao string `json:"descricao_situacao_cadastral"`
	DataInicio        string `json:"data_inicio_atividade"`
	CNAEFiscal        int    `json:"cnae_fiscal"`
	CNAEDescricao     string `json:"cnae_fiscal_descricao"`
	Logradouro        string `json:"logradouro"`
	Numero            string `json:"numero"`
	Complemento       string `json:"complemento"`
	Bairro            string `json:"bairro"`
	Municipio         string `json:"municipio"`
	UF                string `json:"uf"`
	CEP               string `json:"cep"`
}

func limparSomenteNumeros(valor string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(valor, "")
}

func (app *App) telaConsultasExternas(w http.ResponseWriter, r *http.Request) {
	resultado := ResultadoConsultaExterna{}

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")

		switch acao {
		case "consultar_cep":
			resultado = consultarCEPViaCEP(r.FormValue("cep"))
		case "consultar_cnpj":
			resultado = consultarCNPJBrasilAPI(r.FormValue("cnpj"))
		default:
			resultado = ResultadoConsultaExterna{
				Tipo:     "Erro",
				Sucesso:  false,
				Mensagem: "Ação de consulta não reconhecida.",
				DataHora: time.Now().Format("02/01/2006 15:04"),
			}
		}
	}

	tpl := template.Must(template.New("consultas").Parse(htmlBase(consultasExternasHTML)))
	dados := map[string]any{
		"Titulo":    "Consultas externas",
		"Resultado": resultado,
	}
	tpl.Execute(w, dados)
}

func consultarCEPViaCEP(cepInformado string) ResultadoConsultaExterna {
	cep := limparSomenteNumeros(cepInformado)

	if len(cep) != 8 {
		return ResultadoConsultaExterna{
			Tipo:     "CEP",
			Sucesso:  false,
			Mensagem: "Informe um CEP com 8 números.",
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	cliente := http.Client{Timeout: 10 * time.Second}
	resp, err := cliente.Get(url)
	if err != nil {
		return ResultadoConsultaExterna{
			Tipo:     "CEP",
			Sucesso:  false,
			Mensagem: "Erro ao consultar ViaCEP: " + err.Error(),
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResultadoConsultaExterna{
			Tipo:     "CEP",
			Sucesso:  false,
			Mensagem: fmt.Sprintf("ViaCEP retornou status %d.", resp.StatusCode),
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	var dados viacepResposta
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		return ResultadoConsultaExterna{
			Tipo:     "CEP",
			Sucesso:  false,
			Mensagem: "Erro ao ler resposta do ViaCEP: " + err.Error(),
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	if dados.Erro {
		return ResultadoConsultaExterna{
			Tipo:     "CEP",
			Sucesso:  false,
			Mensagem: "CEP não encontrado no ViaCEP.",
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	return ResultadoConsultaExterna{
		Tipo:     "CEP",
		Sucesso:  true,
		Mensagem: "CEP consultado com sucesso.",
		DataHora: time.Now().Format("02/01/2006 15:04"),
		Dados: map[string]string{
			"CEP":         dados.CEP,
			"Logradouro":  dados.Logradouro,
			"Complemento": dados.Complemento,
			"Bairro":      dados.Bairro,
			"Município":   dados.Localidade,
			"UF":          dados.UF,
			"Código IBGE": dados.IBGE,
			"Fonte":       "ViaCEP",
		},
	}
}

func consultarCNPJBrasilAPI(cnpjInformado string) ResultadoConsultaExterna {
	cnpj := limparSomenteNumeros(cnpjInformado)

	if len(cnpj) != 14 {
		return ResultadoConsultaExterna{
			Tipo:     "CNPJ",
			Sucesso:  false,
			Mensagem: "Informe um CNPJ com 14 números.",
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	url := fmt.Sprintf("https://brasilapi.com.br/api/cnpj/v1/%s", cnpj)

	cliente := http.Client{Timeout: 15 * time.Second}
	resp, err := cliente.Get(url)
	if err != nil {
		return ResultadoConsultaExterna{
			Tipo:     "CNPJ",
			Sucesso:  false,
			Mensagem: "Erro ao consultar BrasilAPI: " + err.Error(),
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ResultadoConsultaExterna{
			Tipo:     "CNPJ",
			Sucesso:  false,
			Mensagem: "CNPJ não encontrado na BrasilAPI.",
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return ResultadoConsultaExterna{
			Tipo:     "CNPJ",
			Sucesso:  false,
			Mensagem: fmt.Sprintf("BrasilAPI retornou status %d.", resp.StatusCode),
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	var dados brasilAPICNPJResposta
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		return ResultadoConsultaExterna{
			Tipo:     "CNPJ",
			Sucesso:  false,
			Mensagem: "Erro ao ler resposta da BrasilAPI: " + err.Error(),
			DataHora: time.Now().Format("02/01/2006 15:04"),
		}
	}

	endereco := strings.TrimSpace(dados.Logradouro + ", " + dados.Numero)
	if dados.Complemento != "" {
		endereco += " - " + dados.Complemento
	}

	return ResultadoConsultaExterna{
		Tipo:     "CNPJ",
		Sucesso:  true,
		Mensagem: "CNPJ consultado com sucesso.",
		DataHora: time.Now().Format("02/01/2006 15:04"),
		Dados: map[string]string{
			"CNPJ":                dados.CNPJ,
			"Razão social":        dados.RazaoSocial,
			"Nome fantasia":       dados.NomeFantasia,
			"Situação cadastral":  dados.DescricaoSituacao,
			"Início da atividade": dados.DataInicio,
			"CNAE principal":      fmt.Sprintf("%d", dados.CNAEFiscal),
			"Descrição CNAE":      dados.CNAEDescricao,
			"Endereço":            endereco,
			"Bairro":              dados.Bairro,
			"Município":           dados.Municipio,
			"UF":                  dados.UF,
			"CEP":                 dados.CEP,
			"Fonte":               "BrasilAPI CNPJ",
		},
	}
}

const consultasExternasHTML = `
<h2>Consultas externas</h2>

<p class="pequeno">
	Use esta tela para buscar dados públicos que ajudam a preencher a pré-análise. 
	As consultas dependem de internet e das APIs externas estarem disponíveis.
</p>

<div class="grid">
	<div class="card">
		<h3>Consultar CEP</h3>

		<form method="POST" action="/consultas-externas">
			<input type="hidden" name="acao" value="consultar_cep">

			<label>CEP</label>
			<input type="text" name="cep" placeholder="Ex: 01001000">

			<button type="submit">Consultar CEP</button>
		</form>
	</div>

	<div class="card">
		<h3>Consultar CNPJ</h3>

		<form method="POST" action="/consultas-externas">
			<input type="hidden" name="acao" value="consultar_cnpj">

			<label>CNPJ</label>
			<input type="text" name="cnpj" placeholder="Ex: 00000000000191">

			<button type="submit">Consultar CNPJ</button>
		</form>
	</div>
</div>

{{if .Resultado.Tipo}}
<div class="card destaque">
	<h3>Resultado da consulta - {{.Resultado.Tipo}}</h3>

	<p><strong>Status:</strong> {{if .Resultado.Sucesso}}Sucesso{{else}}Atenção{{end}}</p>
	<p><strong>Mensagem:</strong> {{.Resultado.Mensagem}}</p>
	<p><strong>Data/hora:</strong> {{.Resultado.DataHora}}</p>

	{{if .Resultado.Dados}}
	<table>
		<thead>
			<tr>
				<th>Campo</th>
				<th>Valor</th>
			</tr>
		</thead>
		<tbody>
			{{range $chave, $valor := .Resultado.Dados}}
			<tr>
				<td>{{$chave}}</td>
				<td>{{$valor}}</td>
			</tr>
			{{end}}
		</tbody>
	</table>
	{{end}}
</div>
{{end}}

<div class="card">
	<h3>Próximas consultas que vamos adicionar</h3>
	<ul>
		<li>CAR / SICAR por número do CAR ou importação de arquivo.</li>
		<li>Embargos Ibama por CPF/CNPJ ou importação de base aberta.</li>
		<li>SNCR/CCIR por código do imóvel rural.</li>
		<li>Importação de planilha externa para preencher dados do imóvel.</li>
	</ul>
</div>
`
