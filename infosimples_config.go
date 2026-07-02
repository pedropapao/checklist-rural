package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const urlPadraoInfoSimplesCARDemonstrativo = "https://api.infosimples.com/api/v2/consultas/car/demonstrativo"
const urlPadraoInfoSimplesCARDemonstrativoPDF = "https://api.infosimples.com/api/v2/consultas/car/demonstrativo-pdf"
const urlPadraoInfoSimplesCARDownloadShapefile = "https://api.infosimples.com/api/v2/consultas/car/download-shapefile"

type ConfigInfoSimples struct {
	Token                   string
	CARDemonstrativoURL     string
	CARDemonstrativoPDFURL  string
	CARDownloadShapefileURL string
}

func caminhoConfigInfoSimples() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pasta := filepath.Join(home, "Documentos", "ChecklistRural")
	if err := os.MkdirAll(pasta, 0755); err != nil {
		return "", err
	}

	return filepath.Join(pasta, "infosimples.env"), nil
}

func lerConfigInfoSimples() ConfigInfoSimples {
	c := ConfigInfoSimples{
		CARDemonstrativoURL:     urlPadraoInfoSimplesCARDemonstrativo,
		CARDemonstrativoPDFURL:  urlPadraoInfoSimplesCARDemonstrativoPDF,
		CARDownloadShapefileURL: urlPadraoInfoSimplesCARDownloadShapefile,
	}

	caminho, err := caminhoConfigInfoSimples()
	if err != nil {
		return c
	}

	arquivo, err := os.Open(caminho)
	if err != nil {
		return c
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		linha := strings.TrimSpace(scanner.Text())
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}

		partes := strings.SplitN(linha, "=", 2)
		if len(partes) != 2 {
			continue
		}

		chave := strings.TrimSpace(partes[0])
		valor := strings.Trim(strings.TrimSpace(partes[1]), `"`)

		switch chave {
		case "INFOSIMPLES_TOKEN":
			c.Token = valor
		case "INFOSIMPLES_CAR_DEMONSTRATIVO_URL":
			c.CARDemonstrativoURL = valor
		case "INFOSIMPLES_CAR_DEMONSTRATIVO_PDF_URL":
			c.CARDemonstrativoPDFURL = valor
		case "INFOSIMPLES_CAR_DOWNLOAD_SHAPEFILE_URL":
			c.CARDownloadShapefileURL = valor
		}
	}

	if strings.TrimSpace(c.CARDemonstrativoURL) == "" {
		c.CARDemonstrativoURL = urlPadraoInfoSimplesCARDemonstrativo
	}

	if strings.TrimSpace(c.CARDemonstrativoPDFURL) == "" {
		c.CARDemonstrativoPDFURL = urlPadraoInfoSimplesCARDemonstrativoPDF
	}

	if strings.TrimSpace(c.CARDownloadShapefileURL) == "" {
		c.CARDownloadShapefileURL = urlPadraoInfoSimplesCARDownloadShapefile
	}

	if strings.TrimSpace(c.CARDemonstrativoPDFURL) == "" {
		c.CARDemonstrativoPDFURL = urlPadraoInfoSimplesCARDemonstrativoPDF
	}

	if strings.TrimSpace(c.CARDownloadShapefileURL) == "" {
		c.CARDownloadShapefileURL = urlPadraoInfoSimplesCARDownloadShapefile
	}

	return c
}

func salvarConfigInfoSimples(c ConfigInfoSimples) error {
	caminho, err := caminhoConfigInfoSimples()
	if err != nil {
		return err
	}

	if strings.TrimSpace(c.CARDemonstrativoURL) == "" {
		c.CARDemonstrativoURL = urlPadraoInfoSimplesCARDemonstrativo
	}

	linhas := []string{
		"# Configurações privadas da InfoSimples",
		"# Não envie este arquivo para o GitHub.",
		"",
		"INFOSIMPLES_TOKEN=" + strings.TrimSpace(c.Token),
		"INFOSIMPLES_CAR_DEMONSTRATIVO_URL=" + strings.TrimSpace(c.CARDemonstrativoURL),
		"INFOSIMPLES_CAR_DEMONSTRATIVO_PDF_URL=" + strings.TrimSpace(c.CARDemonstrativoPDFURL),
		"INFOSIMPLES_CAR_DOWNLOAD_SHAPEFILE_URL=" + strings.TrimSpace(c.CARDownloadShapefileURL),
		"",
	}

	return os.WriteFile(caminho, []byte(strings.Join(linhas, "\n")), 0600)
}
