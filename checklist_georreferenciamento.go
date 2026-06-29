package main

import "strings"

func (app *App) adicionarItensGeorreferenciamentoAoChecklist(itens []ItemChecklist, r Reuniao) []ItemChecklist {
	resumo := app.resumoArquivosGeorreferenciamento(r.ID)

	if resumo.Total == 0 {
		itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
			Grupo:      "Georreferenciamento e localização",
			Item:       "Mapa, croqui, KML, GeoJSON ou arquivo de localização da área do projeto",
			Status:     "Pendente",
			Explicacao: "Importar ou solicitar arquivo de localização da área do projeto. Pode ser KML, KMZ, GeoJSON, croqui, planta, shapefile, PDF ou outro documento que ajude a identificar a área financiada.",
		})

		return itens
	}

	itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
		Grupo:      "Georreferenciamento e localização",
		Item:       "Conferência dos arquivos georreferenciados importados",
		Status:     "Pendente",
		Explicacao: "Verificar se os arquivos importados correspondem à área correta do projeto, imóvel rural, CAR, talhão, gleba ou benfeitoria financiada.",
	})

	if possuiArquivoVetorialGeoref(resumo) {
		itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
			Grupo:      "Georreferenciamento e localização",
			Item:       "Conferir coordenadas, polígono ou talhão do arquivo vetorial",
			Status:     "Pendente",
			Explicacao: "Quando houver KML, KMZ, GeoJSON, GPX, CSV ou shapefile, conferir se as coordenadas representam corretamente a área do projeto e se batem com CAR, matrícula, croqui ou informação do produtor.",
		})
	}

	return itens
}

func possuiArquivoVetorialGeoref(resumo ResumoArquivosGeoref) bool {
	for _, arquivo := range resumo.Arquivos {
		tipo := strings.ToLower(arquivo.Tipo)
		nome := strings.ToLower(arquivo.NomeOriginal)

		if strings.Contains(tipo, "kml") ||
			strings.Contains(tipo, "kmz") ||
			strings.Contains(tipo, "geojson") ||
			strings.Contains(tipo, "gpx") ||
			strings.Contains(tipo, "csv") ||
			strings.Contains(tipo, "shape") ||
			strings.Contains(nome, ".kml") ||
			strings.Contains(nome, ".kmz") ||
			strings.Contains(nome, ".geojson") ||
			strings.Contains(nome, ".gpx") ||
			strings.Contains(nome, ".csv") ||
			strings.Contains(nome, ".shp") {
			return true
		}
	}

	return false
}
