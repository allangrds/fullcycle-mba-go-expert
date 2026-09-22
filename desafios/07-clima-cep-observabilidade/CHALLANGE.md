Objetivo

Desenvolver um sistema distribuído em Go composto por dois microsserviços (Serviço A e Serviço B) que cooperam para consultar o clima de uma cidade baseada no CEP. O diferencial deste desafio é a implementação de Observabilidade utilizando OpenTelemetry (OTEL) e Zipkin para realizar o rastreamento distribuído (Distributed Tracing) das requisições.
Arquitetura do sistema

O sistema é composto por:

    Serviço A (Input): Recebe a requisição do usuário, valida o CEP e encaminha para o Serviço B.
    Serviço B (Orquestração): Recebe o CEP, identifica a cidade, consulta a temperatura e realiza as conversões.
    OTEL + Zipkin: Infraestrutura de coleta e visualização dos traços.

Requisitos técnicos: Serviço A (Input)

Este serviço é a porta de entrada. Ele deve ser exposto via HTTP e comunicar-se com o Serviço B.

    Endpoint: Deve aceitar requisições via POST.
    Payload de entrada: O corpo da requisição deve seguir o formato JSON:

    {
      "cep": "29902555"
    }

    Validação:
        O CEP deve ser recebido como String.
        O CEP deve conter exatamente 8 dígitos.
    Comportamento:
        Válido: Encaminha a requisição para o Serviço B via HTTP.
        Inválido: Se o CEP não tiver 8 dígitos ou não for string, retornar:
            Código HTTP: 422
            Mensagem: invalid zipcode

Requisitos técnicos: Serviço B (Orquestração)

Este serviço é responsável pela lógica de negócio.

    Entrada: Recebe um CEP válido de 8 dígitos (enviado pelo Serviço A).
    Localização: Consulta uma API externa (como ViaCEP) para obter o nome da cidade.
    Clima: Consulta uma API externa (como WeatherAPI) para obter a temperatura atual da cidade.
    Conversão: Retorna a temperatura formatada em Celsius, Fahrenheit e Kelvin.
    Respostas (Output):

Sucesso — 200 OK

Deve retornar a cidade e as temperaturas formatadas:

{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}

Erros
Status	Mensagem	Quando ocorre
422	invalid zipcode	CEP com formato inválido
404	can not find zipcode	CEP com formato correto, mas não encontrado
Requisitos de observabilidade (OTEL + Zipkin)

Você deve instrumentar ambos os serviços para garantir o rastreamento completo da requisição.

    Tracing distribuído: Implemente o tracing de forma que seja possível visualizar no Zipkin o fluxo completo:

    Request → Serviço A → Serviço B
    Spans específicos: Além do tracing automático das requisições web, você deve criar Spans manuais para medir o tempo de resposta de:
        Busca de CEP (API externa de localização).
        Busca de temperatura (API externa de clima).
    Infraestrutura: Utilize um OTEL Collector para receber os dados dos serviços e enviá-los ao Zipkin.

Dicas e fórmulas
APIs sugeridas

    ViaCEP — localização
    WeatherAPI — clima

Fórmulas de conversão

    Celsius para Fahrenheit: F = C × 1.8 + 32
    Celsius para Kelvin: K = C + 273

Infraestrutura e entrega
Requisitos de Docker

O projeto deve ser totalmente executável via Docker Compose. O arquivo docker-compose.yaml deve subir:

    Serviço A
    Serviço B
    OTEL Collector
    Zipkin

Entregável

    Código fonte: Repositório contendo a implementação dos serviços A e B.
    Docker Compose: Arquivo configurado para rodar todo o ecossistema.
    Documentação (README):
        Instruções de como realizar a requisição POST no Serviço A.
        Instruções de como acessar o Zipkin para visualizar os traços.

Regras de entrega

    Repositório exclusivo: O repositório deve conter apenas o projeto em questão.
    Branch principal: Todo o código deve estar na branch main.
