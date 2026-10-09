# Clima por CEP com Observabilidade (OpenTelemetry + Zipkin)

Sistema distribuído em Go composto por dois microsserviços que consultam o clima de uma cidade a partir de um CEP, com rastreamento distribuído via OpenTelemetry Collector e Zipkin.

## Arquitetura

```
Cliente ──POST──▶ Serviço A (:8080) ──GET──▶ Serviço B (:8081) ──▶ ViaCEP
                       │                          │            └──▶ WeatherAPI
                       └────── OTLP/gRPC ─────────┴──▶ OTEL Collector ──▶ Zipkin (:9411)
```

- **Serviço A (input)**: recebe `POST /` com `{"cep": "..."}`, valida que o CEP é uma string com 8 dígitos e encaminha para o Serviço B.
- **Serviço B (orquestração)**: busca a cidade no ViaCEP, a temperatura na WeatherAPI e devolve Celsius, Fahrenheit e Kelvin.
- **OTEL Collector**: recebe os traces dos serviços (OTLP gRPC) e exporta para o Zipkin.
- **Zipkin**: visualização dos traces.

### Estrutura do projeto (Clean Architecture)

```
cmd/
  service-a/                  composição e inicialização do Serviço A
  service-b/                  composição e inicialização do Serviço B
internal/
  platform/                   código compartilhado de infraestrutura (telemetria, servidor HTTP, env)
  servicea/
    domain/                   entidades, value objects e erros de domínio
    usecase/                  caso de uso e porta (interface) do gateway
    infra/gateway/            cliente HTTP do Serviço B
    infra/web/                handler e rotas HTTP
  serviceb/
    domain/                   CEP, temperatura (conversões) e erros de domínio
    usecase/                  caso de uso e portas (LocationFinder, TemperatureFinder)
    infra/gateway/            clientes ViaCEP e WeatherAPI
    infra/tracing/            decorators que criam os spans manuais
    infra/web/                handler e rotas HTTP
```

As dependências apontam sempre para dentro (`infra → usecase → domain`). Os casos de uso dependem de interfaces, e os spans manuais são adicionados por decorators sem alterar os gateways nem as regras de negócio.

## Pré-requisitos

- Docker e Docker Compose
- Uma chave gratuita da [WeatherAPI](https://www.weatherapi.com/)

## Como executar

1. Crie o arquivo `.env` com a sua chave da WeatherAPI:

   ```bash
   cp .env.example .env
   ```

   Edite `.env` e substitua o valor de `WEATHER_API_KEY`.

2. Suba todo o ecossistema:

   ```bash
   docker compose up -d --build
   ```

   Serão iniciados: `service-a` (8080), `service-b` (8081), `otel-collector` (4317/4318) e `zipkin` (9411).

## Como fazer a requisição ao Serviço A

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'
```

### Respostas

| Cenário                                      | Status | Corpo                                                                 |
|----------------------------------------------|--------|-----------------------------------------------------------------------|
| CEP válido e encontrado                      | 200    | `{"city":"Linhares","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`      |
| CEP não é string ou não tem 8 dígitos        | 422    | `invalid zipcode`                                                     |
| CEP com formato válido, mas não encontrado   | 404    | `can not find zipcode`                                                |

Exemplos:

```bash
curl -X POST http://localhost:8080 -d '{"cep": "123"}'
curl -X POST http://localhost:8080 -d '{"cep": 29902555}'
curl -X POST http://localhost:8080 -d '{"cep": "99999999"}'
```

Conversões utilizadas: `F = C × 1.8 + 32` e `K = C + 273`.

## Como visualizar os traces no Zipkin

1. Faça uma ou mais requisições ao Serviço A.
2. Acesse **http://localhost:9411**.
3. Clique em **Run Query** (opcionalmente filtre por `serviceName = service-a`).
4. Abra um trace para ver o fluxo completo:

```
service-a: post /
└── service-a: http get                         (chamada ao Serviço B)
    └── service-b: get /weather/{cep}
        ├── service-b: busca-cep-viacep             (span manual)
        │   └── service-b: http get                 (ViaCEP)
        └── service-b: busca-temperatura-weatherapi (span manual)
            └── service-b: http get                 (WeatherAPI)
```

Os spans `busca-cep-viacep` e `busca-temperatura-weatherapi` medem o tempo de resposta das APIs externas de localização e de clima.

## Testes

```bash
go test ./...
```

## Encerrando

```bash
docker compose down
```
