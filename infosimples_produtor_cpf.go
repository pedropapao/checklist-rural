package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const urlInfoSimplesReceitaFederalCPF = "https://api.infosimples.com/api/v2/consultas/receita-federal/cpf"

type TelaInfoSimplesCPF struct {
	Titulo          string
	TokenMascarado  string
	CPF             string
	Birthdate       string
	Consultou       bool
	Sucesso         string
	Erro            string
	Resposta        InfoSimplesResposta
	JSONBonito      string
	NomeProdutor    string
	SituacaoCPF     string
	LinkNovaReuniao string
}

func somenteNumerosCPF(valor string) string {
	var b strings.Builder

	for _, r := range valor {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func validarCPFCompleto(cpf string) error {
	cpf = somenteNumerosCPF(cpf)

	if cpf == "" {
		return fmt.Errorf("informe o CPF")
	}

	if len(cpf) != 11 {
		return fmt.Errorf("CPF precisa ter 11 números")
	}

	todosIguais := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			todosIguais = false
			break
		}
	}

	if todosIguais {
		return fmt.Errorf("CPF inválido")
	}

	digito := func(pos int) byte {
		soma := 0
		peso := pos + 1

		for i := 0; i < pos; i++ {
			soma += int(cpf[i]-'0') * peso
			peso--
		}

		resto := (soma * 10) % 11
		if resto == 10 {
			resto = 0
		}

		return byte(resto) + '0'
	}

	if digito(9) != cpf[9] || digito(10) != cpf[10] {
		return fmt.Errorf("CPF inválido pelo dígito verificador")
	}

	return nil
}

func normalizarNascimentoCPF(valor string) (string, error) {
	valor = strings.TrimSpace(valor)

	if valor == "" {
		return "", fmt.Errorf("informe a data de nascimento")
	}

	formatos := []string{
		"2006-01-02",
		"02/01/2006",
		"2/1/2006",
		"02-01-2006",
		"02012006",
		"20060102",
	}

	for _, formato := range formatos {
		data, err := time.Parse(formato, valor)
		if err == nil {
			if data.After(time.Now()) {
				return "", fmt.Errorf("data de nascimento não pode ser futura")
			}

			return data.Format("2006-01-02"), nil
		}
	}

	limpo := strings.ReplaceAll(valor, "-", "")
	limpo = strings.ReplaceAll(limpo, "/", "")
	limpo = strings.ReplaceAll(limpo, ".", "")
	limpo = strings.ReplaceAll(limpo, " ", "")

	if len(limpo) != 8 {
		return "", fmt.Errorf("data de nascimento inválida. Use ISO 8601: AAAA-MM-DD. Exemplo: 1985-04-22")
	}

	var layout string

	if strings.HasPrefix(limpo, "19") || strings.HasPrefix(limpo, "20") {
		layout = "20060102"
	} else {
		layout = "02012006"
	}

	data, err := time.Parse(layout, limpo)
	if err != nil {
		return "", fmt.Errorf("data de nascimento inválida. Use ISO 8601: AAAA-MM-DD. Exemplo: 1985-04-22")
	}

	if data.After(time.Now()) {
		return "", fmt.Errorf("data de nascimento não pode ser futura")
	}

	return data.Format("2006-01-02"), nil
}

func consultarInfoSimplesReceitaFederalCPF(cpf string, birthdate string) (InfoSimplesResposta, error) {
	config := lerConfigInfoSimples()

	cpf = somenteNumerosCPF(cpf)
	if err := validarCPFCompleto(cpf); err != nil {
		return InfoSimplesResposta{}, err
	}

	birthdate, err := normalizarNascimentoCPF(birthdate)
	if err != nil {
		return InfoSimplesResposta{}, err
	}

	return consultarInfoSimplesPOSTForm(urlInfoSimplesReceitaFederalCPF, config.Token, map[string]string{
		"cpf":       cpf,
		"birthdate": birthdate,
	})
}

func primeiroCPFInfoSimples(resp InfoSimplesResposta) (map[string]any, error) {
	var itens []map[string]any

	if len(resp.Data) == 0 || string(resp.Data) == "null" {
		return nil, fmt.Errorf("resposta sem dados do CPF")
	}

	if err := json.Unmarshal(resp.Data, &itens); err != nil {
		return nil, fmt.Errorf("não consegui interpretar o retorno do CPF: %w", err)
	}

	if len(itens) == 0 {
		return nil, fmt.Errorf("nenhum dado retornado para o CPF")
	}

	return itens[0], nil
}

func valorInfoSimplesCPF(m map[string]any, chaves ...string) string {
	for _, chave := range chaves {
		if v, ok := m[chave]; ok && v != nil {
			txt := strings.TrimSpace(fmt.Sprint(v))
			if txt != "" && txt != "<nil>" {
				return txt
			}
		}
	}

	return ""
}

func montarLinkNovaReuniaoCPF(nome string, cpf string, situacao string) string {
	q := url.Values{}
	q.Set("produtor", nome)

	obs := "Consulta InfoSimples Receita Federal / CPF\n"
	obs += "CPF: " + somenteNumerosCPF(cpf) + "\n"

	if situacao != "" {
		obs += "Situação CPF: " + situacao + "\n"
	}

	q.Set("observacoes", obs)

	return "/nova-reuniao?" + q.Encode()
}

func (app *App) telaInfoSimplesCPF(w http.ResponseWriter, r *http.Request) {
	config := lerConfigInfoSimples()

	tela := TelaInfoSimplesCPF{
		Titulo:         "Consulta CPF antes da triagem",
		TokenMascarado: tokenMascarado(config.Token),
	}

	if r.Method == http.MethodPost {
		tela.CPF = strings.TrimSpace(r.FormValue("cpf"))
		tela.Birthdate = strings.TrimSpace(r.FormValue("birthdate"))

		resposta, err := consultarInfoSimplesReceitaFederalCPF(tela.CPF, tela.Birthdate)
		tela.Consultou = true
		tela.Resposta = resposta
		tela.JSONBonito = jsonBonitoInfoSimples(resposta.Raw)

		if err != nil {
			tela.Erro = err.Error()
		} else if resposta.Code != 200 {
			tela.Erro = fmt.Sprintf("InfoSimples retornou código %d: %s", resposta.Code, resposta.CodeMessage)
			if len(resposta.Errors) > 0 {
				tela.Erro += " - " + strings.Join(resposta.Errors, "; ")
			}
		} else {
			item, err := primeiroCPFInfoSimples(resposta)
			if err != nil {
				tela.Erro = err.Error()
			} else {
				nomeProdutor := valorInfoSimplesCPF(item, "nome", "nome_completo", "name")
				situacaoCPF := valorInfoSimplesCPF(item, "situacao", "situacao_cadastral", "status")

				if nomeProdutor == "" {
					nomeProdutor = valorInfoSimplesCPF(item, "cpf_nome", "nome_pf")
				}

				destino := montarLinkNovaReuniaoCPF(nomeProdutor, tela.CPF, situacaoCPF)

				http.Redirect(w, r, destino, http.StatusSeeOther)
				return
			}
		}
	}

	tpl := templateMust("infosimples_cpf", infosimplesCPFHTML)
	tpl.Execute(w, tela)
}

const infosimplesCPFHTML = `
<section class="card">
	<h2>Consulta CPF antes da triagem</h2>

	<p class="pequeno">
		Use esta tela para consultar o produtor pessoa física antes de abrir a reunião.
		A consulta exige CPF e data de nascimento.
	</p>

	<div class="alerta perigo">
		<strong>Atenção:</strong> esta consulta pode consumir crédito da InfoSimples.
		Confira CPF e data de nascimento antes de consultar.
	</div>

	<div class="grade">
		<div>
			<strong>Token InfoSimples:</strong><br>
			{{.TokenMascarado}}
		</div>
		<div>
			<strong>Consulta:</strong><br>
			Receita Federal / CPF
		</div>
	</div>

	{{if .Erro}}
		<div class="alerta perigo">{{.Erro}}</div>
	{{end}}

	{{if .Sucesso}}
		<div class="alerta sucesso">{{.Sucesso}}</div>
	{{end}}

	<form method="post" class="formulario">
		<label>CPF do produtor</label>
		<input type="text" name="cpf" value="{{.CPF}}" placeholder="000.000.000-00">

		<label>Data de nascimento</label>
		<input type="text" name="birthdate" value="{{.Birthdate}}" placeholder="DD/MM/AAAA">

		<div class="barra-acoes">
			<button class="botao" type="submit">Consultar CPF</button>
			<a class="botao secundario" href="/nova-reuniao">Triagem manual</a>
			<a class="botao secundario" href="/">Início</a>
		</div>
	</form>
</section>

{{if .Consultou}}
<section class="card">
	<h3>Resumo da consulta</h3>

	<div class="grade">
		<div><strong>Código:</strong><br>{{.Resposta.Code}}</div>
		<div><strong>Mensagem:</strong><br>{{.Resposta.CodeMessage}}</div>
		<div><strong>Consumiu crédito:</strong><br>{{.Resposta.Header.Billable}}</div>
		<div><strong>Preço informado:</strong><br>R$ {{.Resposta.Header.Price}}</div>
		<div><strong>Serviço:</strong><br>{{.Resposta.Header.Service}}</div>
		<div><strong>Itens retornados:</strong><br>{{.Resposta.DataCount}}</div>
	</div>
</section>

{{if .NomeProdutor}}
<section class="card">
	<h3>Dados encontrados</h3>

	<div class="grade">
		<div>
			<strong>Nome do produtor:</strong><br>
			{{.NomeProdutor}}
		</div>

		<div>
			<strong>Situação CPF:</strong><br>
			{{if .SituacaoCPF}}{{.SituacaoCPF}}{{else}}Não informado no retorno{{end}}
		</div>

		<div>
			<strong>CPF consultado:</strong><br>
			{{.CPF}}
		</div>
	</div>

	<div class="barra-acoes">
		<a class="botao" href="{{.LinkNovaReuniao}}">Iniciar triagem com esse produtor</a>
	</div>
</section>
{{end}}

<section class="card">
	<h3>JSON bruto retornado</h3>
	<p class="pequeno">
		Este bloco serve para conferirmos exatamente quais campos a InfoSimples devolveu.
		Depois podemos melhorar o preenchimento automático com base nesse retorno real.
	</p>
	<pre style="white-space: pre-wrap; overflow:auto; max-height:600px;">{{.JSONBonito}}</pre>
</section>
{{end}}
`
