package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func consultarInfoSimplesCARDownloadShapefile(car string) (InfoSimplesResposta, error) {
	config := lerConfigInfoSimples()

	car = normalizarCAR(car)
	if err := validarCARParaConsulta(car); err != nil {
		return InfoSimplesResposta{}, err
	}

	return consultarInfoSimplesPOSTForm(config.CARDownloadShapefileURL, config.Token, map[string]string{
		"car": car,
	})
}

func nomeArquivoCARShapefileZIP(car string) string {
	car = normalizarCAR(car)
	if car == "" {
		car = "car"
	}

	agora := time.Now().Format("20060102_150405")
	return "shapefile_car_" + car + "_" + agora + ".zip"
}

func (app *App) registrarZIPInfoSimplesNoGeorreferenciamento(reuniaoID int, caminhoZIP string, car string, resp InfoSimplesResposta) error {
	nomeOriginal := filepath.Base(caminhoZIP)
	nomeSalvo := filepath.Base(caminhoZIP)

	observacao := fmt.Sprintf(
		"Arquivo baixado automaticamente pela InfoSimples CAR / Download Shapefile.\nCAR: %s\nData: %s\nPreço informado pela API: R$ %s",
		normalizarCAR(car),
		time.Now().Format("02/01/2006 15:04"),
		resp.Header.Price,
	)

	_, err := app.DB.Exec(`
		INSERT INTO arquivos_georef
			(reuniao_id, nome_original, nome_salvo, tipo, caminho, observacao, criado_em)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
	`, reuniaoID, nomeOriginal, nomeSalvo, "zip", caminhoZIP, observacao, time.Now().Format("02/01/2006 15:04"))

	return err
}

func (app *App) telaInfoSimplesCARShapefile(w http.ResponseWriter, r *http.Request) {
	tela := TelaInfoSimplesCAR{
		Titulo:         "InfoSimples - CAR Download Shapefile",
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

		resposta, err := consultarInfoSimplesCARDownloadShapefile(tela.CAR)
		tela.Consultou = true
		tela.Resposta = resposta
		tela.JSONBonito = jsonBonitoInfoSimples(resposta.Raw)

		if err != nil {
			tela.Erro = err.Error()
		} else {
			tela.Sucesso = "Consulta do shapefile realizada. Confira o retorno antes de salvar."

			if acao == "consultar_baixar" {
				if tela.Reuniao.ID <= 0 {
					tela.Erro = "Para baixar e registrar o shapefile, abra esta tela a partir de uma reunião."
				} else {
					links := extrairSiteReceipts(resposta)
					linkZIP := primeiroLinkComExtensao(links, ".zip")

					if linkZIP == "" {
						tela.Erro = "A InfoSimples respondeu com sucesso, mas não encontrei link .zip em site_receipts."
					} else {
						pasta, err := pastaGeorreferenciamentoReuniao(tela.Reuniao.ID)
						if err != nil {
							tela.Erro = "Erro ao criar pasta de georreferenciamento: " + err.Error()
						} else {
							nome := nomeArquivoCARShapefileZIP(tela.CAR)
							destino := filepath.Join(pasta, nome)

							err = baixarArquivoURL(linkZIP, destino)
							if err != nil {
								tela.Erro = "Erro ao baixar ZIP do shapefile: " + err.Error()
							} else {
								_ = os.Chmod(destino, 0600)

								err = app.registrarZIPInfoSimplesNoGeorreferenciamento(tela.Reuniao.ID, destino, tela.CAR, resposta)
								if err != nil {
									tela.Erro = "ZIP baixado, mas erro ao registrar no georreferenciamento: " + err.Error()
								} else {
									dados := tela.Dados
									dados.CAR = normalizarCAR(tela.CAR)
									dados.FonteConsulta = "InfoSimples CAR / Download Shapefile"
									dados.LinkFonte = linkZIP
									dados.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")

									resumo := "\n\n--- Shapefile InfoSimples CAR ---\n" +
										"Data: " + time.Now().Format("02/01/2006 15:04") + "\n" +
										"CAR: " + normalizarCAR(tela.CAR) + "\n" +
										"Arquivo ZIP salvo em: " + destino + "\n" +
										"Preço informado pela API: R$ " + resposta.Header.Price + "\n"

									dados.Observacoes = strings.TrimSpace(dados.Observacoes + resumo)
									_ = app.salvarDadosExternos(dados)

									tela.Dados = dados
									tela.Sucesso = "ZIP do shapefile baixado e registrado no georreferenciamento da reunião."
								}
							}
						}
					}
				}
			}
		}
	}

	tpl := templateMust("infosimples_car_shapefile", infosimplesCARShapefileHTML)
	tpl.Execute(w, tela)
}
