package query

const CamposMedidor = `id,
denominacao,
tipo,
latitude,
longitude,
data_primeira_leitura,
data_ultima_leitura,
ultima_data_sincronizada,
medidor_enel_id`

//UpsertMedidoresPrefix Prefixo da Query de Upsert dos medidores
var UpsertMedidores = `INSERT INTO public.sideufg_medidor 
(id, denominacao, data_primeira_leitura, data_ultima_leitura) VALUES ($1,$2,$3,$4) 
ON CONFLICT ON CONSTRAINT sideufg_medidor_pkey 
		DO UPDATE SET data_primeira_leitura = $3, 
		data_ultima_leitura = $4;`



//InsertMedicoesPrefix Prefixo da Query de Insert das medicoes
var UpdateUltimaData = `UPDATE public.sideufg_medidor SET ultima_data_sincronizada = $1 
		WHERE id = $2;`

var BuscaMedidor = `SELECT `+CamposMedidor+` FROM
				public.sideufg_medidor WHERE id = $1;`