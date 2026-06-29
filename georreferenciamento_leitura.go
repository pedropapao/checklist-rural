package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ResumoGeoref struct {
	Sucesso          bool
	Mensagem         string
	TipoArquivo      string
	QuantidadePontos int
	TemPonto         bool
	TemLinha         bool
	TemPoligono      bool
	NomesEncontrados []string
}

type kmlDocumento struct {
	XMLName    xml.Name       `xml:"kml"`
	Placemarks []kmlPlacemark `xml:"Document>Placemark"`
	Pastas     []kmlFolder    `xml:"Document>Folder"`
}

type kmlFolder struct {
	Placemarks []kmlPlacemark `xml:"Placemark"`
}

type kmlPlacemark struct {
	Name         string `xml:"name"`
	Point        string `xml:"Point>coordinates"`
	LineString   string `xml:"LineString>coordinates"`
	Polygon      string `xml:"Polygon>outerBoundaryIs>LinearRing>coordinates"`
	MultiPolygon string `xml:"MultiGeometry>Polygon>outerBoundaryIs>LinearRing>coordinates"`
}

func analisarArquivoGeoref(caminho string) ResumoGeoref {
	ext := strings.ToLower(filepath.Ext(caminho))

	switch ext {
	case ".kml":
		return analisarKML(caminho)
	case ".geojson", ".json":
		return analisarGeoJSON(caminho)
	default:
		return ResumoGeoref{
			Sucesso:     false,
			Mensagem:    "Leitura automática ainda não disponível para este tipo. O arquivo foi guardado normalmente.",
			TipoArquivo: tipoArquivoGeoref(caminho),
		}
	}
}

func analisarKML(caminho string) ResumoGeoref {
	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		return ResumoGeoref{
			Sucesso:  false,
			Mensagem: "Erro ao ler KML: " + err.Error(),
		}
	}

	var doc kmlDocumento
	err = xml.Unmarshal(conteudo, &doc)
	if err != nil {
		return ResumoGeoref{
			Sucesso:  false,
			Mensagem: "Erro ao interpretar KML: " + err.Error(),
		}
	}

	placemarks := []kmlPlacemark{}
	placemarks = append(placemarks, doc.Placemarks...)

	for _, pasta := range doc.Pastas {
		placemarks = append(placemarks, pasta.Placemarks...)
	}

	resumo := ResumoGeoref{
		Sucesso:     true,
		TipoArquivo: "KML",
	}

	for _, p := range placemarks {
		if strings.TrimSpace(p.Name) != "" {
			resumo.NomesEncontrados = append(resumo.NomesEncontrados, strings.TrimSpace(p.Name))
		}

		if strings.TrimSpace(p.Point) != "" {
			resumo.TemPonto = true
			resumo.QuantidadePontos += contarCoordenadasKML(p.Point)
		}

		if strings.TrimSpace(p.LineString) != "" {
			resumo.TemLinha = true
			resumo.QuantidadePontos += contarCoordenadasKML(p.LineString)
		}

		if strings.TrimSpace(p.Polygon) != "" {
			resumo.TemPoligono = true
			resumo.QuantidadePontos += contarCoordenadasKML(p.Polygon)
		}

		if strings.TrimSpace(p.MultiPolygon) != "" {
			resumo.TemPoligono = true
			resumo.QuantidadePontos += contarCoordenadasKML(p.MultiPolygon)
		}
	}

	if resumo.QuantidadePontos == 0 {
		resumo.Mensagem = "KML lido, mas não encontrei coordenadas no formato esperado."
	} else {
		resumo.Mensagem = fmt.Sprintf("KML lido com sucesso. Foram identificadas aproximadamente %d coordenadas.", resumo.QuantidadePontos)
	}

	resumo.NomesEncontrados = limitarNomesGeoref(resumo.NomesEncontrados)

	return resumo
}

func contarCoordenadasKML(texto string) int {
	texto = strings.TrimSpace(texto)
	if texto == "" {
		return 0
	}

	partes := strings.Fields(texto)
	total := 0

	for _, parte := range partes {
		if strings.Contains(parte, ",") {
			total++
		}
	}

	return total
}

func analisarGeoJSON(caminho string) ResumoGeoref {
	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		return ResumoGeoref{
			Sucesso:  false,
			Mensagem: "Erro ao ler GeoJSON: " + err.Error(),
		}
	}

	var raiz any
	err = json.Unmarshal(conteudo, &raiz)
	if err != nil {
		return ResumoGeoref{
			Sucesso:  false,
			Mensagem: "Erro ao interpretar GeoJSON/JSON: " + err.Error(),
		}
	}

	resumo := ResumoGeoref{
		Sucesso:     true,
		TipoArquivo: "GeoJSON/JSON",
	}

	analisarValorGeoJSON(raiz, &resumo)

	if resumo.QuantidadePontos == 0 {
		resumo.Mensagem = "GeoJSON/JSON lido, mas não encontrei coordenadas no formato esperado."
	} else {
		resumo.Mensagem = fmt.Sprintf("GeoJSON/JSON lido com sucesso. Foram identificadas aproximadamente %d coordenadas.", resumo.QuantidadePontos)
	}

	resumo.NomesEncontrados = limitarNomesGeoref(resumo.NomesEncontrados)

	return resumo
}

func analisarValorGeoJSON(valor any, resumo *ResumoGeoref) {
	switch v := valor.(type) {
	case map[string]any:
		if nome, ok := v["name"].(string); ok && strings.TrimSpace(nome) != "" {
			resumo.NomesEncontrados = append(resumo.NomesEncontrados, strings.TrimSpace(nome))
		}

		if props, ok := v["properties"].(map[string]any); ok {
			for _, chave := range []string{"name", "nome", "Name", "NOME"} {
				if nome, ok := props[chave].(string); ok && strings.TrimSpace(nome) != "" {
					resumo.NomesEncontrados = append(resumo.NomesEncontrados, strings.TrimSpace(nome))
				}
			}
		}

		if tipo, ok := v["type"].(string); ok {
			switch strings.ToLower(tipo) {
			case "point", "multipoint":
				resumo.TemPonto = true
			case "linestring", "multilinestring":
				resumo.TemLinha = true
			case "polygon", "multipolygon":
				resumo.TemPoligono = true
			}
		}

		if coords, ok := v["coordinates"]; ok {
			resumo.QuantidadePontos += contarCoordenadasGeoJSON(coords)
		}

		for _, filho := range v {
			analisarValorGeoJSON(filho, resumo)
		}

	case []any:
		for _, item := range v {
			analisarValorGeoJSON(item, resumo)
		}
	}
}

func contarCoordenadasGeoJSON(valor any) int {
	switch v := valor.(type) {
	case []any:
		if len(v) >= 2 {
			_, ok1 := v[0].(float64)
			_, ok2 := v[1].(float64)

			if ok1 && ok2 {
				return 1
			}
		}

		total := 0
		for _, item := range v {
			total += contarCoordenadasGeoJSON(item)
		}
		return total

	default:
		return 0
	}
}

func limitarNomesGeoref(nomes []string) []string {
	limpos := []string{}
	vistos := map[string]bool{}

	for _, nome := range nomes {
		nome = strings.TrimSpace(nome)
		if nome == "" {
			continue
		}

		chave := strings.ToLower(nome)
		if vistos[chave] {
			continue
		}

		vistos[chave] = true
		limpos = append(limpos, nome)

		if len(limpos) >= 10 {
			break
		}
	}

	return limpos
}

func textoResumoGeoref(resumo ResumoGeoref) string {
	partes := []string{}

	if resumo.Mensagem != "" {
		partes = append(partes, resumo.Mensagem)
	}

	partes = append(partes, "Tipo: "+resumo.TipoArquivo)

	if resumo.TemPonto {
		partes = append(partes, "Contém ponto(s)")
	}

	if resumo.TemLinha {
		partes = append(partes, "Contém linha(s)")
	}

	if resumo.TemPoligono {
		partes = append(partes, "Contém polígono(s)")
	}

	if resumo.QuantidadePontos > 0 {
		partes = append(partes, fmt.Sprintf("Coordenadas identificadas: %d", resumo.QuantidadePontos))
	}

	if len(resumo.NomesEncontrados) > 0 {
		partes = append(partes, "Nomes encontrados: "+strings.Join(resumo.NomesEncontrados, ", "))
	}

	return strings.Join(partes, "\n")
}
