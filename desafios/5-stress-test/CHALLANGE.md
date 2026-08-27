Criar um sistema de CLI (Command Line Interface) em Go para realizar testes de carga em serviços web. O usuário deverá fornecer a URL do serviço, o número total de requisições e a quantidade de chamadas simultâneas. Ao final, o sistema deve gerar um relatório detalhado da execução.



Entrada de Parâmetros

A aplicação deve aceitar os seguintes parâmetros via linha de comando:

    --url: URL do serviço a ser testado.

    --requests: Número total de requisições a serem realizadas.

    --concurrency: Número de chamadas simultâneas.



Requisitos Técnicos



    Execução do Teste:

    Realizar requisições HTTP para a URL especificada.
    Distribuir as requisições de acordo com o nível de concorrência definido.
    Garantir que o número total de requisições (--requests) seja cumprido exato.

​2. Geração de Relatório: Ao final da execução, o sistema deve apresentar no console as seguintes informações:

    Tempo total gasto na execução.
    Quantidade total de requests realizados.
    Quantidade de requests com status HTTP 200.
    Distribuição de outros códigos de status HTTP (ex: quantidade de 404, 500, etc...).



Execução via Docker (Obrigatório)

A aplicação deve ser empacotada em uma imagem Docker para facilitar a execução. O comando de teste deve seguir o formato abaixo:

docker run <sua-imagem-docker> --url=http://google.com --requests=1000 --concurrency=10



Entregável

    Código Fonte: Repositório contendo a implementação do CLI.

    Dockerfile: Arquivo de configuração para construção da imagem.

    README: Instruções detalhadas de como buildar a imagem e executar o comando de teste.



Regras de Entrega

    Repositório Exclusivo: O repositório deve conter apenas o projeto em questão.

    Branch Principal: Todo o código deve estar na branch main.
