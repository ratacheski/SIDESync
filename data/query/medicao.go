package query

import "sidegosynchronizer/core/storage"

const CamposMedicao = `id,
data_medicao,
potencia_ativa,
potencia_reativa,
medidor_id`

var BuscaUltimaLeituraMedidor = `SELECT ` + storage.AdicionaAlias(CamposMedicao, "m") +
	` FROM sideufg_medicao as m WHERE m.medidor_id = $1 ORDER BY m.data_medicao desc LIMIT 1`

//InsertMedicoesPrefix Prefixo da Query de Insert das medicoes
var InsertMedicoesPrefix = `INSERT INTO public.sideufg_medicao (` + CamposMedicao + `) VALUES `
