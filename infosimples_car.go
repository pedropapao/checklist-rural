package main

import (
	"net/http"
	"strconv"
	"strings"
)

type TelaInfoSimplesCAR struct {
	Titulo         string
	Reuniao        Reuniao
	Dados          DadosExternosReuniao
	Config         ConfigInfoSimples
	ConfigCaminho  string
	TokenMascarado string
	CAR            string
	Erro           string
	Sucesso        string
	Resposta       InfoSimplesResposta
	JSONBonito     string
	Consultou      bool
}

func consultarInfoSimplesCARDemonstrativo(car string) (InfoSimplesResposta, error) {
	config := lerConfigInfoSimples()

	car = normalizarCAR(car)
	if err := validarCARParaConsulta(car); err != nil {
		return InfoSimplesResposta{}, err
	}

	return consultarInfoSimplesPOSTForm(config.CARDemonstrativoURL, config.Token, map[string]string{
		"car": car,
	})
}

func consultarInfoSimplesCARDemonstrativoPDF(car string) (InfoSimplesResposta, error) {
	config := lerConfigInfoSimples()

	car = normalizarCAR(car)
	if err := validarCARParaConsulta(car); err != nil {
		return InfoSimplesResposta{}, err
	}

	return consultarInfoSimplesPOSTForm(config.CARDemonstrativoPDFURL, config.Token, map[string]string{
		"car": car,
	})
}

func (app *App) telaConfigInfoSimples(w http.ResponseWriter, r *http.Request) {
	mensagem := ""
	erro := ""

	config := lerConfigInfoSimples()

	if r.Method == http.MethodPost {
		config.Token = strings.TrimSpace(r.FormValue("token"))
		config.CARDemonstrativoURL = strings.TrimSpace(r.FormValue("car_demonstrativo_url"))
		config.CARDemonstrativoPDFURL = strings.TrimSpace(r.FormValue("car_demonstrativo_pdf_url"))
		config.CARDownloadShapefileURL = strings.TrimSpace(r.FormValue("car_download_shapefile_url"))

		if config.CARDemonstrativoURL == "" {
			config.CARDemonstrativoURL = urlPadraoInfoSimplesCARDemonstrativo
		}

		if config.CARDemonstrativoPDFURL == "" {
			config.CARDemonstrativoPDFURL = urlPadraoInfoSimplesCARDemonstrativoPDF
		}

		if config.CARDownloadShapefileURL == "" {
			config.CARDownloadShapefileURL = urlPadraoInfoSimplesCARDownloadShapefile
		}

		err := salvarConfigInfoSimples(config)
		if err != nil {
			erro = "Erro ao salvar configuração: " + err.Error()
		} else {
			mensagem = "Configuração da InfoSimples salva com sucesso."
		}
	}

	caminho, _ := caminhoConfigInfoSimples()

	dados := map[string]any{
		"Titulo":         "Configuração InfoSimples",
		"Config":         config,
		"TokenMascarado": tokenMascarado(config.Token),
		"CaminhoConfig":  caminho,
		"Mensagem":       mensagem,
		"Erro":           erro,
	}

	tpl := templateMust("infosimples_config", infosimplesConfigHTML)
	tpl.Execute(w, dados)
}

func (app *App) telaInfoSimplesCAR(w http.ResponseWriter, r *http.Request) {
	var tela TelaInfoSimplesCAR

	tela.Titulo = "InfoSimples - CAR Demonstrativo"
	tela.Config = lerConfigInfoSimples()
	tela.TokenMascarado = tokenMascarado(tela.Config.Token)
	tela.ConfigCaminho, _ = caminhoConfigInfoSimples()

	reuniaoID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if reuniaoID > 0 {
		reuniao, err := app.buscarReuniaoPorID(reuniaoID)
		if err == nil {
			tela.Reuniao = reuniao
			tela.Dados = app.buscarDadosExternosPorReuniao(reuniaoID)
			tela.CAR = tela.Dados.CAR
		}
	}

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		tela.CAR = strings.TrimSpace(r.FormValue("car"))

		resposta, err := consultarInfoSimplesCARDemonstrativo(tela.CAR)
		tela.Consultou = true
		tela.Resposta = resposta
		tela.JSONBonito = jsonBonitoInfoSimples(resposta.Raw)

		if err != nil {
			tela.Erro = err.Error()
		} else {
			tela.Sucesso = "Consulta realizada. Confira o retorno antes de salvar qualquer informação na reunião."

			if acao == "consultar_salvar" {
				if tela.Reuniao.ID <= 0 {
					tela.Erro = "Para salvar, abra esta tela a partir de uma reunião."
				} else {
					item, err := primeiroCARDemonstrativoInfoSimples(resposta)
					if err != nil {
						tela.Erro = err.Error()
					} else {
						dadosAtualizados := aplicarCARDemonstrativoInfoSimplesNosDados(tela.Dados, item, resposta)
						err = app.salvarDadosExternos(dadosAtualizados)
						if err != nil {
							tela.Erro = "Erro ao salvar dados na reunião: " + err.Error()
						} else {
							tela.Dados = dadosAtualizados
							tela.Sucesso = "Consulta realizada e dados do CAR salvos na investigação da reunião."
						}
					}
				}
			}
		}
	}

	tpl := templateMust("infosimples_car", infosimplesCARHTML)
	tpl.Execute(w, tela)
}
