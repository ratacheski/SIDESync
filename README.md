# SIDESync

SIDESync é o sincronizador desenvolvido para trazer as medições dos dispositivos CCK para o banco de dados do SIDE da UFG.

## Instalação

Dentro do repositório do projeto rode o comando [go build](https://golang.org/cmd/go/) para compilar os pacotes e dependências do projeto

```bash
go build
```

## Utilização

O Sistema possui dois arquivos de configuração que podem ser personalizados. O primeiro é o **config.json**:


```json
{
    "host": "HOST_DO_BANCO",
    "password": "SENHA_DO_USUARIO_DO_BANCO",
    "user": "USUARIO_DO_BANCO",
    "dbName": "NOME_DO_BANCO_POSTGRESQL",
    "port": "PORTA_DO_BANCO_POSTGRESQL",
    "logConfig": "./logConfig.json",
    "urlConexaoCCK": "URL_BASE_DO_CCK_WEBSERVICE",
    "intervaloSincronizacao": "INTERVALO_EM_SEGUNDOS_ENTRE_AS_SINCRONIZACOES"
}
```

O segundo é o **logConfig.json** que define as configurações de log da aplicação:


```json
{
  "Diretorio": "DIRETORIO_DE_ARMAZENAMENTO_DOS_LOGS",
  "NomeArquivo": "NOME_INICIAL_DO_ARQUIVO",
  "TamanhoMaximo": "TAMANHO_EM_MBs",
  "RotacionarPorPeriodo": "BOOLEANO_QUE_INDICA_SE_DEVERÁ_ROTACIONAR",
  "PeriodoRotacionar": "VALOR_EM_HORAS_PARA_ROTACIONAR",
  "HorarioLocal": "LOGAR_COM_HORARIO_LOCAL_OU_GMT0",
  "NivelLog": "NIVEL_DE_LOG",
  "Formato": {
    "Formatar": true,
    "DataHora": "2006-01-02 15:04:05.00000",
    "DataHoraCompleto": true,
    "Ordenar": false
  }
}
```

## Contribuição
Pull requests são bem vindos. Para mudanças significativas, por gentileza abra uma issue primeiramente para que possamos discutir as mudanças.

## Licença
[MIT](https://choosealicense.com/licenses/mit/)
