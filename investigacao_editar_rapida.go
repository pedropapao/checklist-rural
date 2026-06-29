package main

import (
	"html/template"
	"net/http"
	"time"
)

func (app *App) editarInvestigacaoRapida(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dados := app.buscarDadosExternosPorReuniao(reuniao.ID)
	mensagem := ""

	if r.Method == http.MethodPost {
		dados.ReuniaoID = reuniao.ID

		dados.CNPJ = r.FormValue("cnpj")
		dados.CEP = r.FormValue("cep")
		dados.RazaoSocial = r.FormValue("razao_social")
		dados.NomeFantasia = r.FormValue("nome_fantasia")
		dados.SituacaoCadastral = r.FormValue("situacao_cadastral")
		dados.Endereco = r.FormValue("endereco")
		dados.Bairro = r.FormValue("bairro")
		dados.Municipio = r.FormValue("municipio")
		dados.UF = r.FormValue("uf")

		dados.NomeImovel = r.FormValue("nome_imovel")
		dados.CAR = r.FormValue("car")
		dados.SituacaoCAR = r.FormValue("situacao_car")
		dados.CodigoINCRA = r.FormValue("codigo_incra")
		dados.CCIR = r.FormValue("ccir")
		dados.NIRFCAFIR = r.FormValue("nirf_cafir")
		dados.Matricula = r.FormValue("matricula")
		dados.AreaTotal = r.FormValue("area_total")

		dados.EmbargoIBAMA = r.FormValue("embargo_ibama")
		dados.CEISStatus = r.FormValue("ceis_status")
		dados.CNEPStatus = r.FormValue("cnep_status")
		dados.SNCRStatus = r.FormValue("sncr_status")
		dados.SIGEFGEO = r.FormValue("sigef_geo")

		dados.FonteConsulta = r.FormValue("fonte_consulta")
		dados.LinkFonte = r.FormValue("link_fonte")
		dados.Observacoes = r.FormValue("observacoes")
		dados.AtualizadoEm = time.Now().Format("02/01/2006 15:04")

		if dados.UltimaInvestigacaoOnline == "" {
			dados.UltimaInvestigacaoOnline = dados.AtualizadoEm
		}

		err := app.salvarDadosExternos(dados)
		if err != nil {
			http.Error(w, "Erro ao salvar investigação rápida: "+err.Error(), http.StatusInternalServerError)
			return
		}

		mensagem = "Investigação rápida salva com sucesso."
		dados = app.buscarDadosExternosPorReuniao(reuniao.ID)
	}

	tpl := template.Must(template.New("investigacao_editar_rapida").Parse(htmlBase(investigacaoEditarRapidaHTML)))

	tpl.Execute(w, map[string]any{
		"Titulo":   "Editar investigação rápida",
		"Reuniao":  reuniao,
		"Dados":    dados,
		"Mensagem": mensagem,
	})
}

const investigacaoEditarRapidaHTML = `
<h2>Editar investigação rápida</h2>

<div class="barra-acoes">
	<a class="botao secundario" href="/investigacao?id={{.Reuniao.ID}}">Voltar para investigação</a>
	<a class="botao secundario" href="/dados-externos-reuniao?id={{.Reuniao.ID}}">Editar investigação completa</a>
</div>

{{if .Mensagem}}
<div class="card destaque">
	<p>{{.Mensagem}}</p>
</div>
{{end}}

<div class="card destaque">
	<h3>{{.Reuniao.Produtor}}</h3>
	<p class="pequeno">
		{{.Reuniao.Municipio}}/{{.Reuniao.UF}} — {{.Reuniao.Banco}} — {{.Reuniao.Atividade}}
	</p>
</div>

<form method="POST" action="/investigacao-editar?id={{.Reuniao.ID}}">
	<div class="card">
		<h3>1. CNPJ e endereço</h3>
		<p class="pequeno">Dados cadastrais básicos do produtor/empresa rural.</p>

		<div class="grid">
			<div>
				<label>CNPJ</label>
				<input type="text" name="cnpj" value="{{.Dados.CNPJ}}" placeholder="Somente CNPJ, sem CPF">
			</div>

			<div>
				<label>Situação cadastral</label>
				<input type="text" name="situacao_cadastral" value="{{.Dados.SituacaoCadastral}}" placeholder="Ex: ativa, inapta, baixada">
			</div>
		</div>

		<label>Razão social</label>
		<input type="text" name="razao_social" value="{{.Dados.RazaoSocial}}">

		<label>Nome fantasia</label>
		<input type="text" name="nome_fantasia" value="{{.Dados.NomeFantasia}}">

		<div class="grid">
			<div>
				<label>CEP</label>
				<input type="text" name="cep" value="{{.Dados.CEP}}">
			</div>

			<div>
				<label>Bairro</label>
				<input type="text" name="bairro" value="{{.Dados.Bairro}}">
			</div>
		</div>

		<label>Endereço</label>
		<input type="text" name="endereco" value="{{.Dados.Endereco}}">

		<div class="grid">
			<div>
				<label>Município</label>
				<input type="text" name="municipio" value="{{.Dados.Municipio}}">
			</div>

			<div>
				<label>UF</label>
				<input type="text" name="uf" value="{{.Dados.UF}}" maxlength="2">
			</div>
		</div>
	</div>

	<div class="card">
		<h3>2. Imóvel rural</h3>
		<p class="pequeno">Campos principais para CAR, matrícula, Incra, CCIR e área.</p>

		<label>Nome do imóvel</label>
		<input type="text" name="nome_imovel" value="{{.Dados.NomeImovel}}" placeholder="Ex: Fazenda Santa Maria">

		<label>CAR</label>
		<input type="text" name="car" value="{{.Dados.CAR}}" placeholder="Código/registro do CAR">

		<label>Situação do CAR</label>
		<input type="text" name="situacao_car" value="{{.Dados.SituacaoCAR}}" placeholder="Ex: ativo, pendente, informado manualmente">

		<div class="grid">
			<div>
				<label>Código INCRA</label>
				<input type="text" name="codigo_incra" value="{{.Dados.CodigoINCRA}}">
			</div>

			<div>
				<label>CCIR</label>
				<input type="text" name="ccir" value="{{.Dados.CCIR}}">
			</div>
		</div>

		<div class="grid">
			<div>
				<label>NIRF / CAFIR</label>
				<input type="text" name="nirf_cafir" value="{{.Dados.NIRFCAFIR}}">
			</div>

			<div>
				<label>Matrícula</label>
				<input type="text" name="matricula" value="{{.Dados.Matricula}}">
			</div>
		</div>

		<label>Área total</label>
		<input type="text" name="area_total" value="{{.Dados.AreaTotal}}" placeholder="Ex: 125,50 ha">
	</div>

	<div class="card">
		<h3>3. Alertas e bases públicas</h3>
		<p class="pequeno">Resumo manual ou resultado das consultas públicas.</p>

		<div class="grid">
			<div>
				<label>Embargo Ibama</label>
				<select name="embargo_ibama">
					<option value="" {{if eq .Dados.EmbargoIBAMA ""}}selected{{end}}>Não consultado</option>
					<option value="nao" {{if eq .Dados.EmbargoIBAMA "nao"}}selected{{end}}>Não encontrado</option>
					<option value="sim" {{if eq .Dados.EmbargoIBAMA "sim"}}selected{{end}}>Encontrado</option>
				</select>
			</div>

			<div>
				<label>CEIS</label>
				<input type="text" name="ceis_status" value="{{.Dados.CEISStatus}}" placeholder="não consultado / sem registro / com alerta">
			</div>

			<div>
				<label>CNEP</label>
				<input type="text" name="cnep_status" value="{{.Dados.CNEPStatus}}" placeholder="não consultado / sem registro / com alerta">
			</div>

			<div>
				<label>SNCR</label>
				<input type="text" name="sncr_status" value="{{.Dados.SNCRStatus}}">
			</div>
		</div>

		<label>SIGEF / GEO</label>
		<input type="text" name="sigef_geo" value="{{.Dados.SIGEFGEO}}" placeholder="Ex: informado, pendente, consulta manual, conferido">
	</div>

	<div class="card">
		<h3>4. Fonte e observações</h3>
		<p class="pequeno">Registre de onde veio a informação: produtor, documento, site, consulta pública ou API.</p>

		<label>Fonte da consulta</label>
		<input type="text" name="fonte_consulta" value="{{.Dados.FonteConsulta}}" placeholder="Ex: documento do produtor, recibo CAR, Ibama, consulta pública">

		<label>Link da fonte</label>
		<input type="text" name="link_fonte" value="{{.Dados.LinkFonte}}" placeholder="https://...">

		<label>Observações</label>
		<textarea name="observacoes" rows="8" placeholder="Anote pendências, fonte dos dados, conferências e alertas.">{{.Dados.Observacoes}}</textarea>
	</div>

	<div class="card">
		<button type="submit">Salvar investigação rápida</button>
		<a class="botao secundario" href="/investigacao?id={{.Reuniao.ID}}">Cancelar</a>
	</div>
</form>
`
