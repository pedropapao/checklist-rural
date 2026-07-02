package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ResultadoFluxoInfoSimples struct {
	Titulo     string
	Reuniao    Reuniao
	Dados      DadosExternosReuniao
	CAR        string
	Executou   bool
	Sucesso    bool
	Erro       string
	Mensagens  []string
	LinkMapa   string
	LinkGeoref string
}

func (app *App) telaInfoSimplesFluxoAutomatico(w http.ResponseWriter, r *http.Request) {
	reuniaoID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || reuniaoID <= 0 {
		http.Error(w, "ID da reunião inválido", http.StatusBadRequest)
		return
	}

	reuniao, err := app.buscarReuniaoPorID(reuniaoID)
	if err != nil {
		http.Error(w, "Reunião não encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	dados := app.buscarDadosExternosPorReuniao(reuniaoID)

	tela := ResultadoFluxoInfoSimples{
		Titulo:     "Atualização automática InfoSimples",
		Reuniao:    reuniao,
		Dados:      dados,
		CAR:        dados.CAR,
		LinkMapa:   "/mapa-georef?id=" + strconv.Itoa(reuniaoID),
		LinkGeoref: "/georreferenciamento?id=" + strconv.Itoa(reuniaoID),
	}

	if r.Method == http.MethodPost {
		tela.Executou = true
		car := strings.TrimSpace(r.FormValue("car"))
		tela.CAR = car

		mensagens, err := app.executarFluxoAutomaticoInfoSimples(reuniaoID, car)
		tela.Mensagens = mensagens

		if err != nil {
			tela.Erro = err.Error()
			tela.Sucesso = false
		} else {
			tela.Sucesso = true
			tela.Dados = app.buscarDadosExternosPorReuniao(reuniaoID)
		}
	}

	tpl := templateMust("infosimples_fluxo_automatico", infosimplesFluxoAutomaticoHTML)
	tpl.Execute(w, tela)
}

func (app *App) executarFluxoAutomaticoInfoSimples(reuniaoID int, car string) ([]string, error) {
	var mensagens []string

	car = normalizarCAR(car)
	if err := validarCARParaConsulta(car); err != nil {
		return mensagens, err
	}

	dados := app.buscarDadosExternosPorReuniao(reuniaoID)

	mensagens = append(mensagens, "Iniciando atualização automática pela InfoSimples.")
	mensagens = append(mensagens, "CAR usado: "+car)

	// 1. Demonstrativo
	respDemonstrativo, err := consultarInfoSimplesCARDemonstrativo(car)
	if err != nil {
		return mensagens, fmt.Errorf("erro no CAR / Demonstrativo: %w", err)
	}
	if respDemonstrativo.Code != 200 {
		return mensagens, fmt.Errorf("CAR / Demonstrativo retornou código %d: %s", respDemonstrativo.Code, respDemonstrativo.CodeMessage)
	}

	itemDemonstrativo, err := primeiroCARDemonstrativoInfoSimples(respDemonstrativo)
	if err != nil {
		return mensagens, err
	}

	dados = aplicarCARDemonstrativoInfoSimplesNosDados(dados, itemDemonstrativo, respDemonstrativo)
	if err := app.salvarDadosExternos(dados); err != nil {
		return mensagens, fmt.Errorf("erro ao salvar demonstrativo na reunião: %w", err)
	}

	mensagens = append(mensagens, "CAR / Demonstrativo consultado e salvo. Preço informado: R$ "+respDemonstrativo.Header.Price)

	// 2. CAR / Imóvel
	respImovel, err := consultarInfoSimplesCARImovel(car)
	if err != nil {
		return mensagens, fmt.Errorf("erro no CAR / Imóvel: %w", err)
	}
	if respImovel.Code != 200 {
		return mensagens, fmt.Errorf("CAR / Imóvel retornou código %d: %s", respImovel.Code, respImovel.CodeMessage)
	}

	itemImovel, err := primeiroCARImovelInfoSimples(respImovel)
	if err != nil {
		return mensagens, err
	}

	dados = aplicarCARImovelInfoSimplesNosDados(dados, itemImovel, respImovel)
	if err := app.salvarDadosExternos(dados); err != nil {
		return mensagens, fmt.Errorf("erro ao salvar CAR / Imóvel na reunião: %w", err)
	}

	mensagens = append(mensagens, "CAR / Imóvel consultado e salvo. Preço informado: R$ "+respImovel.Header.Price)

	// 3. Demonstrativo PDF
	respPDF, err := consultarInfoSimplesCARDemonstrativoPDF(car)
	if err != nil {
		return mensagens, fmt.Errorf("erro no Demonstrativo PDF: %w", err)
	}
	if respPDF.Code != 200 {
		return mensagens, fmt.Errorf("Demonstrativo PDF retornou código %d: %s", respPDF.Code, respPDF.CodeMessage)
	}

	linksPDF := extrairSiteReceipts(respPDF)
	linkPDF := primeiroLinkComExtensao(linksPDF, ".pdf")

	if linkPDF == "" {
		mensagens = append(mensagens, "Demonstrativo PDF consultado, mas não encontrei link .pdf no retorno.")
	} else {
		pastaPDF, err := pastaDocumentosInfoSimplesReuniao(app.PastaDados, reuniaoID)
		if err != nil {
			return mensagens, fmt.Errorf("erro ao criar pasta do PDF: %w", err)
		}

		nomePDF := nomeArquivoCARDemonstrativoPDF(car)
		destinoPDF := filepath.Join(pastaPDF, nomePDF)

		if err := baixarArquivoURL(linkPDF, destinoPDF); err != nil {
			return mensagens, fmt.Errorf("erro ao baixar PDF: %w", err)
		}

		dados = app.buscarDadosExternosPorReuniao(reuniaoID)
		dados.FonteConsulta = "InfoSimples - fluxo automático"
		dados.LinkFonte = linkPDF
		dados.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")
		dados.Observacoes = strings.TrimSpace(dados.Observacoes + "\n\n--- PDF InfoSimples baixado automaticamente ---\nArquivo: " + destinoPDF + "\nPreço informado: R$ " + respPDF.Header.Price + "\n")
		_ = app.salvarDadosExternos(dados)

		mensagens = append(mensagens, "Demonstrativo PDF baixado: "+destinoPDF)
		mensagens = append(mensagens, "Preço PDF informado: R$ "+respPDF.Header.Price)
	}

	// 4. Download Shapefile
	respShape, err := consultarInfoSimplesCARDownloadShapefile(car)
	if err != nil {
		return mensagens, fmt.Errorf("erro no Download Shapefile: %w", err)
	}
	if respShape.Code != 200 {
		return mensagens, fmt.Errorf("Download Shapefile retornou código %d: %s", respShape.Code, respShape.CodeMessage)
	}

	linksShape := extrairSiteReceipts(respShape)
	linkZIP := primeiroLinkComExtensao(linksShape, ".zip")

	if linkZIP == "" {
		return mensagens, fmt.Errorf("Download Shapefile consultado, mas não encontrei link .zip no retorno")
	}

	pastaGeo, err := pastaGeorreferenciamentoReuniao(reuniaoID)
	if err != nil {
		return mensagens, fmt.Errorf("erro ao criar pasta de georreferenciamento: %w", err)
	}

	nomeZIP := nomeArquivoCARShapefileZIP(car)
	destinoZIP := filepath.Join(pastaGeo, nomeZIP)

	if err := baixarArquivoURL(linkZIP, destinoZIP); err != nil {
		return mensagens, fmt.Errorf("erro ao baixar ZIP do shapefile: %w", err)
	}

	if err := app.registrarZIPInfoSimplesNoGeorreferenciamento(reuniaoID, destinoZIP, car, respShape); err != nil {
		return mensagens, fmt.Errorf("ZIP baixado, mas erro ao registrar no georreferenciamento: %w", err)
	}

	dados = app.buscarDadosExternosPorReuniao(reuniaoID)
	dados.CAR = car
	dados.FonteConsulta = "InfoSimples - fluxo automático"
	dados.LinkFonte = linkZIP
	dados.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")
	dados.Observacoes = strings.TrimSpace(dados.Observacoes + "\n\n--- Shapefile InfoSimples baixado automaticamente ---\nArquivo ZIP: " + destinoZIP + "\nPreço informado: R$ " + respShape.Header.Price + "\n")
	_ = app.salvarDadosExternos(dados)

	mensagens = append(mensagens, "ZIP do shapefile baixado e registrado no georreferenciamento: "+destinoZIP)
	mensagens = append(mensagens, "Preço Shapefile informado: R$ "+respShape.Header.Price)

	mensagens = append(mensagens, "Fluxo automático finalizado com sucesso.")

	return mensagens, nil
}
