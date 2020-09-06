package query

//UpsertMedidoresPrefix Prefixo da Query de Upsert dos medidores
var UpsertMedidoresPrefix = `INSERT INTO public.sideufg_medidor 
(id, denominacao, data_primeira_leitura, data_ultima_leitura) VALUES `

//UpsertMedidoresSufix Sufixo da Query de Upsert dos medidores
var UpsertMedidoresSufix = ` ON CONFLICT (id) DO UPDATE 
		SET data_primeira_leitura = EXCLUDED.data_primeira_leitura, 
		data_ultima_leitura = EXCLUDED.data_ultima_leitura`
