package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func extrairSiteReceipts(resp InfoSimplesResposta) []string {
	var links []string

	if len(resp.SiteReceipts) > 0 && string(resp.SiteReceipts) != "null" {
		_ = json.Unmarshal(resp.SiteReceipts, &links)
	}

	if len(links) > 0 {
		return links
	}

	var itens []map[string]any
	if len(resp.Data) > 0 {
		_ = json.Unmarshal(resp.Data, &itens)
	}

	for _, item := range itens {
		if link, ok := item["site_receipt"].(string); ok && strings.TrimSpace(link) != "" {
			links = append(links, link)
		}
	}

	return links
}

func primeiroLinkComExtensao(links []string, extensao string) string {
	extensao = strings.ToLower(extensao)

	for _, link := range links {
		if strings.Contains(strings.ToLower(link), extensao) {
			return link
		}
	}

	return ""
}

func pastaDocumentosInfoSimplesReuniao(pastaDados string, reuniaoID int) (string, error) {
	pasta := filepath.Join(
		pastaDados,
		"infosimples",
		fmt.Sprintf("reuniao_%d", reuniaoID),
	)

	err := os.MkdirAll(pasta, 0755)
	if err != nil {
		return "", err
	}

	return pasta, nil
}

func baixarArquivoURL(urlArquivo string, caminhoDestino string) error {
	urlArquivo = strings.TrimSpace(urlArquivo)
	if urlArquivo == "" {
		return fmt.Errorf("URL do arquivo vazia")
	}

	resp, err := http.Get(urlArquivo)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("erro ao baixar arquivo: HTTP %d", resp.StatusCode)
	}

	tmp := caminhoDestino + ".tmp"

	arquivo, err := os.Create(tmp)
	if err != nil {
		return err
	}

	_, err = io.Copy(arquivo, resp.Body)
	fecharErr := arquivo.Close()

	if err != nil {
		_ = os.Remove(tmp)
		return err
	}

	if fecharErr != nil {
		_ = os.Remove(tmp)
		return fecharErr
	}

	return os.Rename(tmp, caminhoDestino)
}

func nomeArquivoCARDemonstrativoPDF(car string) string {
	car = normalizarCAR(car)
	if car == "" {
		car = "car"
	}

	agora := time.Now().Format("20060102_150405")
	return "demonstrativo_car_" + car + "_" + agora + ".pdf"
}
