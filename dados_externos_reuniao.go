package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type DadosExternosReuniao struct {
	ID        int
	ReuniaoID int

	CNPJ              string
	CEP               string
	RazaoSocial       string
	NomeFantasia      string
	SituacaoCadastral string
	CNAE              string
	DescricaoCNAE     string

	Endereco   string
	Bairro     string
	Municipio  string
	UF         string
	CodigoIBGE string

	NomeImovel  string
	CAR         string
	SituacaoCAR string
	CodigoINCRA string
	CCIR        string
	NIRFCAFIR   string
	Matricula   string
	AreaTotal   string

	EmbargoIBAMA  string
	IbamaDetalhes string
	CEISStatus    string
	CEISDetalhes  string
	CNEPStatus    string
	CNEPDetalhes  string
	SNCRStatus    string
	SNCRDetalhes  string
	SIGEFGEO      string

	FonteConsulta            string
	LinkFonte                string
	Observacoes              string
	UltimaInvestigacaoOnline string
	AtualizadoEm             string
}

func (app *App) buscarDadosExternosPorReuniao(reuniaoID int) DadosExternosReuniao {
	var d DadosExternosReuniao

	err := app.DB.QueryRow(`
		SELECT
			COALESCE(id, 0),
			COALESCE(reuniao_id, 0),

			COALESCE(cpf_cnpj, ''),
			COALESCE(cep, ''),
			COALESCE(razao_social, ''),
			COALESCE(nome_fantasia, ''),
			COALESCE(situacao_cadastral, ''),
			COALESCE(cnae, ''),
			COALESCE(descricao_cnae, ''),

			COALESCE(endereco, ''),
			COALESCE(bairro, ''),
			COALESCE(municipio, ''),
			COALESCE(uf, ''),
			COALESCE(codigo_ibge, ''),

			COALESCE(nome_imovel, ''),
			COALESCE(car, ''),
			COALESCE(situacao_car, ''),
			COALESCE(codigo_incra, ''),
			COALESCE(ccir, ''),
			COALESCE(nirf_cafir, ''),
			COALESCE(matricula, ''),
			COALESCE(area_total, ''),

			COALESCE(embargo_ibama, 'nao_consultado'),
			COALESCE(ibama_detalhes, ''),
			COALESCE(ceis_status, 'nao_consultado'),
			COALESCE(ceis_detalhes, ''),
			COALESCE(cnep_status, 'nao_consultado'),
			COALESCE(cnep_detalhes, ''),
			COALESCE(sncr_status, 'nao_consultado'),
			COALESCE(sncr_detalhes, ''),
			COALESCE(sigef_geo, 'nao_consultado'),

			COALESCE(fonte_consulta, ''),
			COALESCE(link_fonte, ''),
			COALESCE(observacoes, ''),
			COALESCE(ultima_investigacao_online, ''),
			COALESCE(atualizado_em, '')
		FROM dados_externos
		WHERE reuniao_id = ?
	`, reuniaoID).Scan(
		&d.ID,
		&d.ReuniaoID,

		&d.CNPJ,
		&d.CEP,
		&d.RazaoSocial,
		&d.NomeFantasia,
		&d.SituacaoCadastral,
		&d.CNAE,
		&d.DescricaoCNAE,

		&d.Endereco,
		&d.Bairro,
		&d.Municipio,
		&d.UF,
		&d.CodigoIBGE,

		&d.NomeImovel,
		&d.CAR,
		&d.SituacaoCAR,
		&d.CodigoINCRA,
		&d.CCIR,
		&d.NIRFCAFIR,
		&d.Matricula,
		&d.AreaTotal,

		&d.EmbargoIBAMA,
		&d.IbamaDetalhes,
		&d.CEISStatus,
		&d.CEISDetalhes,
		&d.CNEPStatus,
		&d.CNEPDetalhes,
		&d.SNCRStatus,
		&d.SNCRDetalhes,
		&d.SIGEFGEO,

		&d.FonteConsulta,
		&d.LinkFonte,
		&d.Observacoes,
		&d.UltimaInvestigacaoOnline,
		&d.AtualizadoEm,
	)

	if err == sql.ErrNoRows {
		d.ReuniaoID = reuniaoID
		d.EmbargoIBAMA = "nao_consultado"
		d.CEISStatus = "nao_consultado"
		d.CNEPStatus = "nao_consultado"
		d.SNCRStatus = "nao_consultado"
		d.SIGEFGEO = "nao_consultado"
		return d
	}

	return d
}

func (app *App) salvarDadosExternos(d DadosExternosReuniao) error {
	d.AtualizadoEm = time.Now().Format("02/01/2006 15:04")

	_, err := app.DB.Exec(`
		INSERT INTO dados_externos (
			reuniao_id,

			cpf_cnpj,
			cep,
			razao_social,
			nome_fantasia,
			situacao_cadastral,
			cnae,
			descricao_cnae,

			endereco,
			bairro,
			municipio,
			uf,
			codigo_ibge,

			nome_imovel,
			car,
			situacao_car,
			codigo_incra,
			ccir,
			nirf_cafir,
			matricula,
			area_total,

			embargo_ibama,
			ibama_detalhes,
			ceis_status,
			ceis_detalhes,
			cnep_status,
			cnep_detalhes,
			sncr_status,
			sncr_detalhes,
			sigef_geo,

			fonte_consulta,
			link_fonte,
			observacoes,
			ultima_investigacao_online,
			atualizado_em
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(reuniao_id) DO UPDATE SET
			cpf_cnpj = excluded.cpf_cnpj,
			cep = excluded.cep,
			razao_social = excluded.razao_social,
			nome_fantasia = excluded.nome_fantasia,
			situacao_cadastral = excluded.situacao_cadastral,
			cnae = excluded.cnae,
			descricao_cnae = excluded.descricao_cnae,

			endereco = excluded.endereco,
			bairro = excluded.bairro,
			municipio = excluded.municipio,
			uf = excluded.uf,
			codigo_ibge = excluded.codigo_ibge,

			nome_imovel = excluded.nome_imovel,
			car = excluded.car,
			situacao_car = excluded.situacao_car,
			codigo_incra = excluded.codigo_incra,
			ccir = excluded.ccir,
			nirf_cafir = excluded.nirf_cafir,
			matricula = excluded.matricula,
			area_total = excluded.area_total,

			embargo_ibama = excluded.embargo_ibama,
			ibama_detalhes = excluded.ibama_detalhes,
			ceis_status = excluded.ceis_status,
			ceis_detalhes = excluded.ceis_detalhes,
			cnep_status = excluded.cnep_status,
			cnep_detalhes = excluded.cnep_detalhes,
			sncr_status = excluded.sncr_status,
			sncr_detalhes = excluded.sncr_detalhes,
			sigef_geo = excluded.sigef_geo,

			fonte_consulta = excluded.fonte_consulta,
			link_fonte = excluded.link_fonte,
			observacoes = excluded.observacoes,
			ultima_investigacao_online = excluded.ultima_investigacao_online,
			atualizado_em = excluded.atualizado_em
	`,
		d.ReuniaoID,

		d.CNPJ,
		d.CEP,
		d.RazaoSocial,
		d.NomeFantasia,
		d.SituacaoCadastral,
		d.CNAE,
		d.DescricaoCNAE,

		d.Endereco,
		d.Bairro,
		d.Municipio,
		d.UF,
		d.CodigoIBGE,

		d.NomeImovel,
		d.CAR,
		d.SituacaoCAR,
		d.CodigoINCRA,
		d.CCIR,
		d.NIRFCAFIR,
		d.Matricula,
		d.AreaTotal,

		d.EmbargoIBAMA,
		d.IbamaDetalhes,
		d.CEISStatus,
		d.CEISDetalhes,
		d.CNEPStatus,
		d.CNEPDetalhes,
		d.SNCRStatus,
		d.SNCRDetalhes,
		d.SIGEFGEO,

		d.FonteConsulta,
		d.LinkFonte,
		d.Observacoes,
		d.UltimaInvestigacaoOnline,
		d.AtualizadoEm,
	)

	return err
}

func montarDadosExternosDoFormulario(r *http.Request, reuniaoID int) DadosExternosReuniao {
	return DadosExternosReuniao{
		ReuniaoID: reuniaoID,

		CNPJ:              r.FormValue("cnpj"),
		CEP:               r.FormValue("cep"),
		RazaoSocial:       r.FormValue("razao_social"),
		NomeFantasia:      r.FormValue("nome_fantasia"),
		SituacaoCadastral: r.FormValue("situacao_cadastral"),
		CNAE:              r.FormValue("cnae"),
		DescricaoCNAE:     r.FormValue("descricao_cnae"),

		Endereco:   r.FormValue("endereco"),
		Bairro:     r.FormValue("bairro"),
		Municipio:  r.FormValue("municipio"),
		UF:         r.FormValue("uf"),
		CodigoIBGE: r.FormValue("codigo_ibge"),

		NomeImovel:  r.FormValue("nome_imovel"),
		CAR:         r.FormValue("car"),
		SituacaoCAR: r.FormValue("situacao_car"),
		CodigoINCRA: r.FormValue("codigo_incra"),
		CCIR:        r.FormValue("ccir"),
		NIRFCAFIR:   r.FormValue("nirf_cafir"),
		Matricula:   r.FormValue("matricula"),
		AreaTotal:   r.FormValue("area_total"),

		EmbargoIBAMA:  r.FormValue("embargo_ibama"),
		IbamaDetalhes: r.FormValue("ibama_detalhes"),
		CEISStatus:    r.FormValue("ceis_status"),
		CEISDetalhes:  r.FormValue("ceis_detalhes"),
		CNEPStatus:    r.FormValue("cnep_status"),
		CNEPDetalhes:  r.FormValue("cnep_detalhes"),
		SNCRStatus:    r.FormValue("sncr_status"),
		SNCRDetalhes:  r.FormValue("sncr_detalhes"),
		SIGEFGEO:      r.FormValue("sigef_geo"),

		FonteConsulta:            r.FormValue("fonte_consulta"),
		LinkFonte:                r.FormValue("link_fonte"),
		Observacoes:              r.FormValue("observacoes"),
		UltimaInvestigacaoOnline: r.FormValue("ultima_investigacao_online"),
	}
}

func (app *App) telaDadosExternosReuniao(w http.ResponseWriter, r *http.Request) {
	idTexto := r.URL.Query().Get("id")
	reuniaoID, err := strconv.Atoi(idTexto)
	if err != nil || reuniaoID <= 0 {
		http.Error(w, "ID da reunião inválido", http.StatusBadRequest)
		return
	}

	reuniao, err := app.buscarReuniaoPorID(reuniaoID)
	if err != nil {
		http.Error(w, "Reunião não encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	resultado := ResultadoConsultaExterna{}

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		dados := montarDadosExternosDoFormulario(r, reuniaoID)

		switch acao {
		case "salvar":
			if dados.FonteConsulta == "" {
				dados.FonteConsulta = "Informado manualmente"
			}

			err = app.salvarDadosExternos(dados)
			if err != nil {
				http.Error(w, "Erro ao salvar dados externos: "+err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dados-externos-reuniao?id="+strconv.Itoa(reuniaoID), http.StatusSeeOther)
			return

		case "consultar_cep":
			resultado = consultarCEPViaCEP(dados.CEP)
			if resultado.Sucesso {
				dados.CEP = resultado.Dados["CEP"]
				dados.Endereco = resultado.Dados["Logradouro"]
				dados.Bairro = resultado.Dados["Bairro"]
				dados.Municipio = resultado.Dados["Município"]
				dados.UF = resultado.Dados["UF"]
				dados.CodigoIBGE = resultado.Dados["Código IBGE"]
				dados.FonteConsulta = "ViaCEP"
				dados.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")
				_ = app.salvarDadosExternos(dados)
			}

		case "consultar_cnpj":
			resultado = consultarCNPJBrasilAPI(dados.CNPJ)
			if resultado.Sucesso {
				dados.CNPJ = resultado.Dados["CNPJ"]
				dados.RazaoSocial = resultado.Dados["Razão social"]
				dados.NomeFantasia = resultado.Dados["Nome fantasia"]
				dados.SituacaoCadastral = resultado.Dados["Situação cadastral"]
				dados.CNAE = resultado.Dados["CNAE principal"]
				dados.DescricaoCNAE = resultado.Dados["Descrição CNAE"]
				dados.Endereco = resultado.Dados["Endereço"]
				dados.Bairro = resultado.Dados["Bairro"]
				dados.Municipio = resultado.Dados["Município"]
				dados.UF = resultado.Dados["UF"]
				dados.CEP = resultado.Dados["CEP"]
				dados.FonteConsulta = "BrasilAPI CNPJ"
				dados.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")
				_ = app.salvarDadosExternos(dados)
			}

		case "consultar_ceis_cnep":
			ceis := consultarCEISPortalTransparencia(dados.CNPJ)
			cnep := consultarCNEPPortalTransparencia(dados.CNPJ)

			if ceis.Encontrou {
				dados.CEISStatus = "sim"
			} else {
				dados.CEISStatus = "nao"
			}
			dados.CEISDetalhes = ceis.Resumo

			if cnep.Encontrou {
				dados.CNEPStatus = "sim"
			} else {
				dados.CNEPStatus = "nao"
			}
			dados.CNEPDetalhes = cnep.Resumo

			dados.FonteConsulta = "Portal da Transparência - CEIS/CNEP"
			dados.LinkFonte = "https://portaldatransparencia.gov.br/sancoes/consulta"
			dados.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")

			_ = app.salvarDadosExternos(dados)

			resultado = ResultadoConsultaExterna{
				Tipo:     "CEIS/CNEP",
				Sucesso:  true,
				Mensagem: "Consulta CEIS/CNEP finalizada. Confira os campos de detalhes.",
				DataHora: time.Now().Format("02/01/2006 15:04"),
			}
		}
	}

	dadosExternos := app.buscarDadosExternosPorReuniao(reuniaoID)

	tpl := template.Must(template.New("dados_externos").Parse(htmlBase(dadosExternosReuniaoHTML)))
	dadosTela := map[string]any{
		"Titulo":    "Investigação online",
		"Reuniao":   reuniao,
		"Dados":     dadosExternos,
		"Resultado": resultado,
	}

	tpl.Execute(w, dadosTela)
}

const dadosExternosReuniaoHTML = `
<h2>Investigação online do produtor/imóvel</h2>

<p class="pequeno">
	Reunião: <strong>{{.Reuniao.Produtor}}</strong> — {{.Reuniao.Municipio}}/{{.Reuniao.UF}}
</p>

<p>
	<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Voltar para detalhes</a>
</p>

{{if .Resultado.Tipo}}
<div class="card destaque">
	<h3>Resultado da consulta - {{.Resultado.Tipo}}</h3>
	<p><strong>Status:</strong> {{if .Resultado.Sucesso}}Sucesso{{else}}Atenção{{end}}</p>
	<p><strong>Mensagem:</strong> {{.Resultado.Mensagem}}</p>
	<p><strong>Data/hora:</strong> {{.Resultado.DataHora}}</p>
</div>
{{end}}

<form method="POST" action="/dados-externos-reuniao?id={{.Reuniao.ID}}">

<div class="card destaque">
	<h3>1. Consulta online por CNPJ</h3>

	<p class="pequeno">
		Use quando o produtor for pessoa jurídica, empresa rural, associação, cooperativa, agroindústria ou fornecedor.
	</p>

	<label>CNPJ</label>
	<input type="text" name="cnpj" value="{{.Dados.CNPJ}}" placeholder="Ex: 00000000000191">

	<div class="grid">
		<button type="submit" name="acao" value="consultar_cnpj">
			Consultar CNPJ
		</button>

		<button type="submit" name="acao" value="salvar">
			Salvar investigação
		</button>
	</div>

	<label>Razão social / nome encontrado</label>
	<input type="text" name="razao_social" value="{{.Dados.RazaoSocial}}">

	<label>Nome fantasia</label>
	<input type="text" name="nome_fantasia" value="{{.Dados.NomeFantasia}}">

	<label>Situação cadastral</label>
	<input type="text" name="situacao_cadastral" value="{{.Dados.SituacaoCadastral}}">

	<label>CNAE principal</label>
	<input type="text" name="cnae" value="{{.Dados.CNAE}}">

	<label>Descrição CNAE</label>
	<input type="text" name="descricao_cnae" value="{{.Dados.DescricaoCNAE}}">
</div>

<div class="card">
	<h3>2. Endereço / localização</h3>

	<label>CEP</label>
	<input type="text" name="cep" value="{{.Dados.CEP}}" placeholder="Ex: 01001000">

	<button type="submit" name="acao" value="consultar_cep">
		Consultar CEP
	</button>

	<label>Endereço</label>
	<input type="text" name="endereco" value="{{.Dados.Endereco}}">

	<label>Bairro</label>
	<input type="text" name="bairro" value="{{.Dados.Bairro}}">

	<div class="grid">
		<div>
			<label>Município</label>
			<input type="text" name="municipio" value="{{.Dados.Municipio}}">
		</div>

		<div>
			<label>UF</label>
			<input type="text" name="uf" value="{{.Dados.UF}}">
		</div>
	</div>

	<label>Código IBGE</label>
	<input type="text" name="codigo_ibge" value="{{.Dados.CodigoIBGE}}">
</div>

<div class="card">
	<h3>3. Imóvel rural</h3>

	<label>Nome do imóvel rural</label>
	<input type="text" name="nome_imovel" value="{{.Dados.NomeImovel}}" placeholder="Ex: Fazenda Santa Maria">

	<label>CAR</label>
	<input type="text" name="car" value="{{.Dados.CAR}}" placeholder="Número do CAR">

	<label>Situação do CAR</label>
	<input type="text" name="situacao_car" value="{{.Dados.SituacaoCAR}}" placeholder="Ex: ativo, pendente, aguardando análise">

	<label>Código do imóvel rural / INCRA</label>
	<input type="text" name="codigo_incra" value="{{.Dados.CodigoINCRA}}">

	<label>CCIR</label>
	<input type="text" name="ccir" value="{{.Dados.CCIR}}">

	<label>NIRF / CAFIR</label>
	<input type="text" name="nirf_cafir" value="{{.Dados.NIRFCAFIR}}">

	<label>Matrícula</label>
	<input type="text" name="matricula" value="{{.Dados.Matricula}}">

	<label>Área total</label>
	<input type="text" name="area_total" value="{{.Dados.AreaTotal}}" placeholder="Ex: 125,30 ha">
</div>

<div class="card">
	<h3>4. Alertas e bases públicas</h3>

	<label>Embargo Ibama</label>
	<select name="embargo_ibama">
		<option value="nao_consultado" {{if eq .Dados.EmbargoIBAMA "nao_consultado"}}selected{{end}}>Não consultado</option>
		<option value="nao" {{if eq .Dados.EmbargoIBAMA "nao"}}selected{{end}}>Não identificado</option>
		<option value="sim" {{if eq .Dados.EmbargoIBAMA "sim"}}selected{{end}}>Sim, possui alerta/embargo</option>
	</select>

	<label>Detalhes Ibama</label>
	<textarea name="ibama_detalhes" rows="3">{{.Dados.IbamaDetalhes}}</textarea>

	<button type="submit" name="acao" value="consultar_ceis_cnep">
		Consultar CEIS/CNEP no Portal da Transparência
	</button>

	<p class="pequeno">
		Exige token da API do Portal da Transparência configurado no sistema.
	</p>

	<label>CEIS - Cadastro de empresas inidôneas/suspensas</label>
	<select name="ceis_status">
		<option value="nao_consultado" {{if eq .Dados.CEISStatus "nao_consultado"}}selected{{end}}>Não consultado</option>
		<option value="nao" {{if eq .Dados.CEISStatus "nao"}}selected{{end}}>Não identificado</option>
		<option value="sim" {{if eq .Dados.CEISStatus "sim"}}selected{{end}}>Possui registro/alerta</option>
	</select>

	<label>Detalhes CEIS</label>
	<textarea name="ceis_detalhes" rows="3">{{.Dados.CEISDetalhes}}</textarea>

	<label>CNEP - Cadastro de empresas punidas</label>
	<select name="cnep_status">
		<option value="nao_consultado" {{if eq .Dados.CNEPStatus "nao_consultado"}}selected{{end}}>Não consultado</option>
		<option value="nao" {{if eq .Dados.CNEPStatus "nao"}}selected{{end}}>Não identificado</option>
		<option value="sim" {{if eq .Dados.CNEPStatus "sim"}}selected{{end}}>Possui registro/alerta</option>
	</select>

	<label>Detalhes CNEP</label>
	<textarea name="cnep_detalhes" rows="3">{{.Dados.CNEPDetalhes}}</textarea>

	<label>SNCR/CNIR</label>
	<select name="sncr_status">
		<option value="nao_consultado" {{if eq .Dados.SNCRStatus "nao_consultado"}}selected{{end}}>Não consultado</option>
		<option value="nao" {{if eq .Dados.SNCRStatus "nao"}}selected{{end}}>Não identificado</option>
		<option value="sim" {{if eq .Dados.SNCRStatus "sim"}}selected{{end}}>Encontrado/conferido</option>
	</select>

	<label>Detalhes SNCR/CNIR</label>
	<textarea name="sncr_detalhes" rows="3">{{.Dados.SNCRDetalhes}}</textarea>

	<label>SIGEF/GEO</label>
	<select name="sigef_geo">
		<option value="nao_consultado" {{if eq .Dados.SIGEFGEO "nao_consultado"}}selected{{end}}>Não consultado</option>
		<option value="nao" {{if eq .Dados.SIGEFGEO "nao"}}selected{{end}}>Não identificado</option>
		<option value="sim" {{if eq .Dados.SIGEFGEO "sim"}}selected{{end}}>Sim, possui GEO/SIGEF</option>
	</select>
</div>

<div class="card">
	<h3>5. Fonte da investigação</h3>

	<label>Fonte da consulta</label>
	<input type="text" name="fonte_consulta" value="{{.Dados.FonteConsulta}}" placeholder="Ex: BrasilAPI, ViaCEP, Ibama, SNCR, documento do produtor">

	<label>Link da fonte oficial ou consulta</label>
	<input type="text" name="link_fonte" value="{{.Dados.LinkFonte}}" placeholder="Cole aqui o link da consulta ou documento">

	<label>Última investigação online</label>
	<input type="text" name="ultima_investigacao_online" value="{{.Dados.UltimaInvestigacaoOnline}}" readonly>

	<label>Observações gerais da investigação</label>
	<textarea name="observacoes" rows="5">{{.Dados.Observacoes}}</textarea>

	<p class="pequeno">
		Última atualização salva: {{.Dados.AtualizadoEm}}
	</p>

	<button type="submit" name="acao" value="salvar">
		Salvar investigação
	</button>
</div>

</form>
`
