package main

type ReuniaoComResumo struct {
	Reuniao
	Resumo ResumoChecklist
}

func (app *App) resumoDaReuniao(reuniao Reuniao) ResumoChecklist {
	var resumo ResumoChecklist

	err := app.DB.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'Pendente' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'Recebido' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'Não se aplica' THEN 1 ELSE 0 END), 0)
		FROM checklist_itens
		WHERE reuniao_id = ?
	`, reuniao.ID).Scan(
		&resumo.Total,
		&resumo.Pendentes,
		&resumo.Recebidos,
		&resumo.NaoSeAplica,
	)

	if err != nil {
		return ResumoChecklist{}
	}

	if resumo.Total > 0 {
		resumo.PercentualConcluido = ((resumo.Recebidos + resumo.NaoSeAplica) * 100) / resumo.Total
	}

	return resumo
}
