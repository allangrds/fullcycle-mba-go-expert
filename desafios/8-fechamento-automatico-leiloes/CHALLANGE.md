Objetivo

Adicionar uma funcionalidade crítica ao sistema de leilões existente: o fechamento automático. Atualmente, o projeto permite criar leilões e dar lances, mas o leilão nunca expira. Sua missão é utilizar Goroutines para garantir que o leilão seja encerrado automaticamente após um tempo pré-definido.
Base do projeto (obrigatória)

O desafio deve ser implementado sobre o código-fonte disponibilizado pelo curso.

    Repositório base: https://github.com/devfullcycle/labs-auction-goexpert
    Contexto: Toda a rotina de criação de leilão (auction) e lances (bid) já está desenvolvida. A validação que impede lances em leilões fechados também já existe. Seu trabalho é implementar a rotina de fechamento.

Requisitos técnicos

Você deve modificar o processo de criação de leilão para incluir o agendamento do seu fechamento.
Configuração de tempo

    Crie uma função (ou ajuste as existentes) para determinar a duração do leilão baseada em parâmetros definidos em variáveis de ambiente (ex: AUCTION_DURATION).

Processamento em background (Goroutine)

    Implemente uma Goroutine que será iniciada assim que um leilão for criado.
    Essa rotina deve monitorar o tempo e, quando a duração do leilão expirar, deve realizar o update no banco de dados alterando o status do leilão para Closed (fechado).

Local da implementação

Concentre seus esforços no arquivo:

internal/infra/database/auction/create_auction.go

É aqui que a mágica deve acontecer.
Requisitos de testes

Implemente um teste automatizado que comprove que o fechamento está ocorrendo.
Cenário
Passo	Ação	Resultado esperado
1	Criar um leilão	Leilão criado com status aberto
2	Aguardar o tempo configurado em AUCTION_DURATION	—
3	Verificar o status do leilão	Status alterado automaticamente para Closed, sem intervenção manual
Dicas de implementação

    Lembre-se que estamos lidando com concorrência. Certifique-se de que sua solução não bloqueie a thread principal.
    Analise como o sistema atual verifica se um leilão é válido na rotina de criação de bid para entender a lógica de tempo e status.

Entregável

    Código fonte: O link para o seu repositório contendo a implementação completa.
    Documentação (README):
        Instruções para rodar o projeto.
        Instruções para configurar as variáveis de ambiente de tempo.
    Docker: O projeto deve estar configurado para rodar via Docker / Docker Compose.

Regras de entrega

    Repositório exclusivo: O repositório deve conter apenas o projeto em questão.
    Branch principal: Todo o código deve estar na branch main.
