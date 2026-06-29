package main

import (
	"net/url"
	"strings"
)

func linkConsultaPublicaCAR(car string) string {
	car = strings.TrimSpace(car)

	base := "https://consulta.car.gov.br/"

	if car == "" {
		return base
	}

	// A consulta pública normalmente usa formulário na página.
	// Mesmo que o site não aceite query automática, deixamos o CAR no link
	// para facilitar copiar/colar e abrir direto a fonte oficial.
	v := url.Values{}
	v.Set("car", car)

	return base + "?" + v.Encode()
}

func textoOrientacaoConsultaPublicaCAR(car string) string {
	car = strings.TrimSpace(car)

	if car == "" {
		return "Abra a Consulta Pública do CAR e pesquise por UF, município ou número de registro no CAR."
	}

	return "Abra a Consulta Pública do CAR e pesquise pelo número: " + car
}
