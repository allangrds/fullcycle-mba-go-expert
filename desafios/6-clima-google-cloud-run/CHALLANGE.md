Objetivo

Desenvolver um sistema em Go que receba um CEP, identifique a cidade correspondente e retorne o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin). O requisito final é que este sistema esteja publicado e acessível no Google Cloud Run.
Requisitos funcionais

    Entrada: O sistema deve receber um CEP válido de 8 dígitos.
    Identificação de localização: O sistema deve realizar a busca do CEP para encontrar o nome da localização (cidade).
    Consulta de clima: A partir da localização, o sistema deve consultar a temperatura atual.
    Conversão: O sistema deve retornar as temperaturas formatadas em Celsius, Fahrenheit e Kelvin.

Especificações da API (contrato)
Cenário 1: Sucesso — 200 OK

{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}

Cenários de falha
Cenário	Condição	Status	Mensagem
Formato inválido	CEP sem 8 dígitos ou com caracteres inválidos	422	invalid zipcode
CEP não encontrado	CEP com formato correto, mas inexistente na base de dados	404	can not find zipcode
Fórmulas de conversão

    Celsius para Fahrenheit: F = C × 1.8 + 32
    Celsius para Kelvin: K = C + 273

Dicas de APIs externas

Você pode utilizar as seguintes APIs (ou similares) para obter os dados:

    Localização: ViaCEP
    Temperatura: WeatherAPI

Requisitos de infraestrutura e deploy

    Docker: O projeto deve possuir um Dockerfile para containerização.
    Cloud Run: A aplicação deve ser implantada no Google Cloud Run (pode utilizar o free tier).
    Testes: Devem ser implementados testes automatizados que comprovem o funcionamento das conversões e das requisições.

Entregável

    Código fonte: Link do repositório no GitHub.
    URL de acesso: O endereço ativo da aplicação no Google Cloud Run (deve constar no README).
    Testes: O projeto deve conter testes automatizados.

Regras de entrega

    Repositório exclusivo: O repositório deve conter apenas o projeto em questão.
    Branch principal: Todo o código deve estar na branch main.
    README: O arquivo deve conter:
        A URL do sistema rodando no Cloud Run.
        Instruções de como rodar os testes e a aplicação localmente via Docker.
