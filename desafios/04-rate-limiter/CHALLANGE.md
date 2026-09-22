Objetivo

Desenvolver um Rate Limiter em Go que funcione como um middleware para controlar o fluxo de requisições de um serviço web.

O sistema deve ser capaz de limitar o tráfego com base no IP do solicitante ou em um Token de acesso, utilizando o Redis para persistência e orquestração.



Regras de Negócio (Critérios de Limitação)

O Rate Limiter deve operar com base em duas estratégias principais:

Limitação por IP: Restringe o número máximo de requisições por segundo recebidas de um único endereço IP.

Limitação por Token: Restringe o número máximo de requisições por segundo baseadas em um token de acesso único.

    Header: O token deve ser verificado no header API_KEY: <TOKEN>.

Precedência (Regra de Ouro): As configurações do Token devem se sobrepor às do IP.

    Exemplo: Se o limite global por IP for 10 req/s, mas um token específico tiver permissão para 100 req/s, o sistema deve respeitar o limite do token (100 req/s).



Comportamento em Caso de Bloqueio

Quando o limite de requisições for excedido (seja por IP ou Token):

    Resposta HTTP: O servidor deve retornar imediatamente o Status Code 429.

    Corpo da Resposta: Deve retornar exatamente a mensagem: "you have reached the maximum number of requests or actions allowed within a certain time frame"

    Tempo de Bloqueio: O IP ou Token infrator deve ficar bloqueado por um tempo configurável (ex: 5 minutos). Durante esse período, todas as novas requisições devem ser rejeitadas.



Requisitos Técnicos e Arquitetura

Middleware: A lógica do Rate Limiter deve ser implementada como um middleware que envolve o servidor web.

Persistência (Redis): As informações de contagem e controle de tempo devem ser armazenadas em um banco de dados Redis (usando imagem Docker).

Design Pattern (Strategy):

    Você deve implementar o padrão Strategy para a persistência.
    Embora o Redis seja obrigatório neste desafio, sua arquitetura deve permitir trocar facilmente o Redis por outro mecanismo de persistência no futuro, apenas mudando a implementação da estratégia.

Desacoplamento: A lógica de negócio do limiter deve estar separada da lógica do middleware.

Configuração: Todas as definições devem ser feitas via variáveis de ambiente (ou arquivo .env na raiz), incluindo:

    Limite máximo de requisições por segundo.
    Tempo de bloqueio (em caso de excesso).
    Configurações de conexão com Redis.



Entregável

Código Fonte: Repositório contendo a implementação completa.

Infraestrutura:

    Dockerfile para a aplicação.
    docker-compose.yaml que suba a aplicação (na porta 8080) e o banco Redis.

Documentação (README): Explique como configurar o limiter, como alterar as estratégias e como executar o projeto.

Testes Automatizados: Testes que demonstrem a eficácia do limiter e a lógica de precedência (Token > IP).



Regras de Entrega

  Repositório Único: O repositório deve conter apenas o projeto em questão.

  Branch Principal: Todo o código deve estar na branch main.

  Execução: O avaliador deve conseguir rodar o projeto e os testes utilizando apenas o Docker/Docker Compose.
