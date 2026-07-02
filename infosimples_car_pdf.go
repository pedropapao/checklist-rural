package main

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (app *App) telaInfoSimplesCARPDF(w http.ResponseWriter, r *http.Request) {
	tela := TelaInfoSimplesCAR{
		Titulo:         "InfoSimples - CAR Demonstrativo PDF",
		Config:         lerConfigInfoSimples(),
		TokenMascarado: tokenMascarado(lerConfigInfoSimples().Token),
	}

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

		resposta, err := consultarInfoSimplesCARDemonstrativoPDF(tela.CAR)
		tela.Consultou = true
		tela.Resposta = resposta
		tela.JSONBonito = jsonBonitoInfoSimples(resposta.Raw)

		if err != nil {
			tela.Erro = err.Error()
		} else {
			tela.Sucesso = "Consulta do PDF realizada. Confira o retorno e os links em site_receipts."

			if acao == "consultar_baixar" {
				if tela.Reuniao.ID <= 0 {
					tela.Erro = "Para baixar e salvar o PDF, abra esta tela a partir de uma reunião."
				} else {
					links := extrairSiteReceipts(resposta)
					linkPDF := primeiroLinkComExtensao(links, ".pdf")

					if linkPDF == "" {
						tela.Erro = "A InfoSimples respondeu com sucesso, mas não encontrei link PDF em site_receipts."
					} else {
						pasta, err := pastaDocumentosInfoSimplesReuniao(app.PastaDados, tela.Reuniao.ID)
						if err != nil {
							tela.Erro = "Erro ao criar pasta da InfoSimples: " + err.Error()
						} else {
							nome := nomeArquivoCARDemonstrativoPDF(tela.CAR)
							destino := filepath.Join(pasta, nome)

							err = baixarArquivoURL(linkPDF, destino)
							if err != nil {
								tela.Erro = "Erro ao baixar PDF: " + err.Error()
							} else {
								dados := tela.Dados
								dados.FonteConsulta = "InfoSimples CAR / Demonstrativo PDF"
								dados.LinkFonte = linkPDF
								dados.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")

								resumo := "\n\n--- PDF InfoSimples CAR / Demonstrativo ---\n" +
									"Data: " + time.Now().Format("02/01/2006 15:04") + "\n" +
									"CAR: " + normalizarCAR(tela.CAR) + "\n" +
									"Arquivo salvo em: " + destino + "\n" +
									"Preço informado pela API: R$ " + resposta.Header.Price + "\n"

								dados.Observacoes = strings.TrimSpace(dados.Observacoes + resumo)

								_ = app.salvarDadosExternos(dados)

								tela.Dados = dados
								tela.Sucesso = "PDF baixado e salvo na pasta da reunião: " + destino
							}
						}
					}
				}
			}
		}
	}

	tpl := templateMust("infosimples_car_pdf", infosimplesCARPDFHTML)
	tpl.Execute(w, tela)
}
