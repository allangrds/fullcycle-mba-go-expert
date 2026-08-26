# ☸️ Deploy no Kubernetes com Go

> Esta aula mostra como pegar uma aplicação Go simples (um servidor HTTP) e colocá-la para rodar de forma resiliente dentro de um cluster Kubernetes. O foco não é a aplicação em si — ela é propositalmente mínima — e sim entender **Deployment**, **Probes** (startup/readiness/liveness) e **Service**, que são os blocos mais fundamentais para qualquer deploy em produção no k8s.

Se você nunca mexeu com Kubernetes antes, não se preocupe: cada conceito abaixo começa com uma analogia do mundo real antes de qualquer YAML.

---

## 📑 Sumário

- [🤔 O que é Kubernetes?](#-o-que-é-kubernetes)
- [🐳 Docker vs Kubernetes](#-docker-vs-kubernetes)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [📦 Dockerfile de Dev vs Dockerfile.prod](#-dockerfile-de-dev-vs-dockerfileprod)
- [📚 Conceitos Fundamentais do Kubernetes](#-conceitos-fundamentais-do-kubernetes)
  - [Pod](#pod)
  - [Deployment](#deployment)
  - [ReplicaSet](#replicaset-o-que-o-deployment-gerencia-por-baixo)
  - [Resources: requests vs limits](#resources-requests-vs-limits)
  - [Probes: startupProbe, readinessProbe, livenessProbe](#probes-startupprobe-readinessprobe-livenessprobe)
  - [Service e tipo LoadBalancer](#service-e-tipo-loadbalancer)
- [⚙️ Como as Peças se Conectam](#️-como-as-peças-se-conectam)
- [✅ Boas Práticas Presentes no Projeto](#-boas-práticas-presentes-no-projeto)
- [🚧 O que Poderia Ser Melhorado / Inconsistências Reais](#-o-que-poderia-ser-melhorado--inconsistências-reais)
- [⚠️ Principais Problemas ao Trabalhar com Kubernetes](#️-principais-problemas-ao-trabalhar-com-kubernetes)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs: Kubernetes vs Alternativas](#️-trade-offs-kubernetes-vs-alternativas)
- [🎯 Quando Usar Kubernetes](#-quando-usar-kubernetes)
- [💼 Perguntas de Entrevista Respondidas](#-perguntas-de-entrevista-respondidas)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é Kubernetes?

Imagine que você é o gerente de um restaurante bem movimentado. Você não fica pessoalmente contratando e demitindo garçons a cada minuto — você define uma **regra**: "eu sempre quero 3 garçons no salão". Se um garçom passa mal e sai, você automaticamente chama outro para substituí-lo. Se o movimento aumenta muito, você contrata mais gente. Se um garçom novo ainda está sendo treinado, você não deixa ele atender mesas sozinho até estar pronto.

O Kubernetes (também chamado de **k8s**, porque tem 8 letras entre o "k" e o "s") é esse "gerente" para os seus containers. Você não fica rodando `docker run` manualmente em cada servidor — você **declara** o estado desejado ("eu quero 3 cópias do meu servidor rodando, sempre") e o Kubernetes garante que isso aconteça, 24 horas por dia, sem você precisar olhar.

Isso é chamado de **infraestrutura declarativa**: você não escreve "faça isso, depois isso, depois aquilo" (imperativo); você escreve "eu quero que o mundo esteja assim" e uma ferramenta cuida de fazer (e manter) esse estado acontecer.

```
Você (imperativo, sem k8s)           Kubernetes (declarativo)
─────────────────────────           ─────────────────────────
docker run app                       "quero 3 réplicas de app"
# um container caiu, ninguém viu     ↓
docker run app  (de novo, manual)    k8s detecta a queda
# precisa de mais 1 réplica?         e sobe outro Pod sozinho
docker run app  (de novo)            automaticamente
```

---

## 🐳 Docker vs Kubernetes

Uma dúvida muito comum de quem está começando: "se eu já tenho Docker, para que preciso de Kubernetes?".

| | Docker | Kubernetes |
|---|---|---|
| **O que é** | Ferramenta para **empacotar** e **rodar** um container | Ferramenta para **orquestrar** múltiplos containers em múltiplas máquinas |
| **Escopo** | Uma aplicação, um container, geralmente uma máquina | Um cluster inteiro (várias máquinas/nós) |
| **Responsabilidade** | Build da imagem, isolamento de processo | Onde rodar, quantas réplicas, reiniciar se cair, expor rede, rolling updates |
| **Analogia** | A caixa de mudança (empacota o que você precisa) | A empresa de logística que decide onde cada caixa vai, substitui a caixa se ela se perder, garante que sempre tenham caixas suficientes no destino |

Eles não competem — **trabalham juntos**. O Docker (ou outro *container runtime*) constrói e roda a imagem; o Kubernetes decide *onde*, *quantas cópias* e *o que fazer quando algo dá errado*. Nesta aula, o `Dockerfile.prod` cuida da primeira parte, e os manifests em `k8s/` cuidam da segunda.

---

## 🗂️ Estrutura do Projeto

```
aulas/19-deploy-k8s/
├── main.go              # servidor HTTP mínimo (a aplicação que será deployada)
├── go.mod                # módulo Go (nome do módulo inconsistente com a pasta — ver nota abaixo)
├── Dockerfile             # imagem de DESENVOLVIMENTO, usada com docker-compose
├── Dockerfile.prod       # imagem de PRODUÇÃO (multi-stage, imagem final `scratch`)
├── docker-compose.yaml   # ambiente de dev local, monta o código como volume
└── k8s/
    ├── deployment.yaml   # declara: 3 réplicas do servidor + probes de saúde
    └── service.yaml      # expõe os Pods na rede via LoadBalancer
```

> **Nota curiosa:** o `go.mod` declara `module github.com/devfullcycle/21-deploy-k8s` e a imagem publicada no Docker Hub é `wesleywillians/21-deploy-k8s:latest`, enquanto a pasta da aula se chama `19-deploy-k8s`. É só um detalhe de nomenclatura da numeração do curso — não afeta o funcionamento, mas é um bom lembrete de como inconsistências de nome se acumulam em projetos reais e podem confundir quem lê o código depois.

`main.go` é propositalmente simples — só para termos algo real para deployar:

```go
package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	http.ListenAndServe(":8080", nil)
}
```

Ele sobe um servidor HTTP na porta `8080` que responde `"Hello World"` em qualquer rota. Esse é o mesmo endpoint (`/`, porta `8080`) que os probes do Kubernetes vão usar para saber se a aplicação está saudável — fica mais fácil entender os manifests sabendo disso.

---

## 📦 Dockerfile de Dev vs Dockerfile.prod

Repare que a pasta tem **dois** Dockerfiles com propósitos bem diferentes.

### `Dockerfile` (desenvolvimento)

```dockerfile
FROM golang:latest

WORKDIR /app

CMD ["tail", "-f", "/dev/null"]
```

Ele **não roda a aplicação**. `tail -f /dev/null` é um truque comum para manter o container vivo indefinidamente sem fazer nada — o container fica "ligado, esperando", e o código real é montado de fora como **volume** pelo `docker-compose.yaml`:

```yaml
services:
  goapp:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - .:/app
```

Isso permite abrir um terminal dentro do container (`docker exec`) e rodar `go run main.go` manualmente, com hot-reload de arquivos, sem precisar reconstruir a imagem a cada alteração — ótimo para o dia a dia de desenvolvimento, péssimo para produção (a imagem nem contém o binário compilado).

### `Dockerfile.prod` (produção)

```dockerfile
FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server .

FROM scratch
COPY --from=builder /app/server .
CMD ["./server"]
```

Esse é um **multi-stage build**: duas etapas (`FROM ... as builder` e um segundo `FROM`) dentro do mesmo Dockerfile.

1. **Estágio `builder`**: usa a imagem completa `golang:latest` (que tem compilador, ferramentas, tudo) só para **compilar** o binário `server`.
   - `CGO_ENABLED=0`: desliga o CGo, forçando um binário **estaticamente linkado** (não depende de bibliotecas C do sistema operacional, como a `glibc`). Sem isso, o binário não rodaria numa imagem `scratch`, que não tem absolutamente nada nela.
   - `-ldflags="-w -s"`: remove informações de debug e a tabela de símbolos do binário, deixando-o menor.
2. **Estágio final (`FROM scratch`)**: `scratch` é uma imagem **vazia** — zero sistema operacional, zero shell, zero bibliotecas. Só copiamos o binário `server` já compilado do estágio anterior.

| | Imagem completa (`golang:latest`) | Imagem `scratch` |
|---|---|---|
| **Tamanho** | Centenas de MB | Poucos MB (só o binário) |
| **Superfície de ataque** | Grande (shell, pacotes, libs que podem ter CVEs) | Mínima (nada para explorar além do próprio binário) |
| **Debug dentro do container** | Fácil (`docker exec` com shell disponível) | Praticamente impossível (não tem `sh`, `ls`, nada) |
| **Uso recomendado** | Build/desenvolvimento | Runtime de produção |

Esse padrão (compilar numa imagem grande, rodar numa imagem mínima) é **muito comum em Go** justamente porque Go compila para um binário único e estático — outras linguagens (Node, Python, Ruby) normalmente não conseguem usar `scratch` porque precisam do runtime da linguagem presente na imagem final.

---

## 📚 Conceitos Fundamentais do Kubernetes

### Pod

Um **Pod** é a menor unidade que o Kubernetes gerencia. Pense nele como "uma cápsula" que contém um (ou mais) containers que sempre rodam juntos, compartilhando rede e armazenamento. Na prática, na maioria dos casos (como o desta aula) um Pod = um container.

Você raramente cria Pods diretamente à mão — eles são criados **através** de outros objetos, como o `Deployment` que veremos a seguir. Isso importa porque um Pod sozinho, se morrer, **não volta** — ninguém está de olho nele. É o Deployment que garante que sempre existam Pods suficientes vivos.

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: server
  template:
    metadata:
      labels:
        app: server
    spec:
      containers:
      - name: server
        image: wesleywillians/21-deploy-k8s:latest
        ...
```

Um `Deployment` é a declaração "eu quero N cópias (réplicas) deste container rodando, sempre". Aqui `replicas: 3` diz ao Kubernetes: mantenha **3 Pods** rodando essa mesma imagem, o tempo todo.

Por baixo dos panos, o Deployment não gerencia os Pods diretamente — ele cria e gerencia um `ReplicaSet` (ver próxima seção), que é quem de fato garante a contagem de réplicas. O Deployment adiciona por cima disso a capacidade de fazer **rolling updates**: quando você troca a versão da imagem, o Kubernetes sobe Pods novos gradualmente e derruba os antigos aos poucos, **sem downtime**, em vez de matar tudo e recriar de uma vez.

O campo `template` descreve como cada Pod deve ser — é basicamente "a receita" que o Deployment usa para criar cada réplica. Repare que `template.metadata.labels.app: server` precisa bater com `selector.matchLabels.app: server` — é assim que o Deployment sabe *quais* Pods são "dele".

### ReplicaSet (o que o Deployment gerencia por baixo)

Você não vê um `ReplicaSet` neste projeto porque ele **não é escrito manualmente** — o `Deployment` o cria automaticamente. Pense nele como o "funcionário de confiança" do Deployment: a única responsabilidade do ReplicaSet é contar quantos Pods com determinado label existem e criar/destruir Pods até bater o número desejado (`replicas: 3`). O Deployment usa o ReplicaSet como ferramenta para implementar rolling updates (cria um ReplicaSet novo para a versão nova, reduz o antigo gradualmente).

```
Deployment (o que você declara)
   └── gerencia → ReplicaSet (garante a contagem)
                     └── gerencia → Pods (as instâncias rodando de fato)
```

### Resources: requests vs limits

```yaml
resources:
  limits:
    memory: "32Mi"
    cpu: "100m"
```

Esse bloco diz ao Kubernetes quanto de CPU/memória o container pode consumir **no máximo** (`limits`). `100m` significa "100 millicores", ou seja, 0.1 de um núcleo de CPU; `32Mi` são 32 mebibytes de memória.

Só que existe um segundo campo, `requests`, que **não aparece** no `deployment.yaml` desta aula — e essa ausência é intencional para discutirmos o problema:

| Campo | Para que serve |
|---|---|
| `requests` | O mínimo que o Pod **precisa** para ser agendado num nó. O Kubernetes usa isso para decidir *em qual nó* colocar o Pod (soma dos requests de tudo que já roda ali não pode passar da capacidade do nó). |
| `limits` | O máximo que o Pod **pode** consumir. Se ultrapassar o limite de memória, o container é morto (`OOMKilled`); se ultrapassar CPU, ele é apenas "estrangulado" (throttled), não morto. |

Sem `requests`, o Kubernetes assume que o Pod não precisa de nada garantido, o que pode fazer o *scheduler* colocar vários Pods "famintos" no mesmo nó e todos competirem por recursos ao mesmo tempo — voltamos a esse ponto na seção [O que Poderia Ser Melhorado](#-o-que-poderia-ser-melhorado--inconsistências-reais).

### Probes: startupProbe, readinessProbe, livenessProbe

Essa é a parte mais importante da aula. Probes são "verificações de saúde" — formas do Kubernetes perguntar "você está bem?" para cada Pod, periodicamente. O `deployment.yaml` usa os três tipos que existem:

```yaml
# startup probe
startupProbe:
  httpGet:
    path: /
    port: 8080
  periodSeconds: 10
  failureThreshold: 10

readinessProbe:
  httpGet:
    path: /
    port: 8080
  periodSeconds: 10
  failureThreshold: 2
  timeoutSeconds: 5

livenessProbe:
  httpGet:
    path: /
    port: 8080
  periodSeconds: 10
  failureThreshold: 3
  timeoutSeconds: 5
  successThreshold: 1
```

Pense numa analogia de um restaurante contratando um garçom novo:

- **startupProbe** — "Será que ele terminou o treinamento inicial e já consegue atender uma mesa?" Enquanto essa pergunta não tiver resposta "sim", **nenhuma outra probe roda** — o Kubernetes dá um tempo generoso para a aplicação inicializar (aqui, até `10 tentativas × 10s = 100s` antes de desistir) sem correr o risco de um app lento na subida ser confundido com um app quebrado.
- **readinessProbe** — "Ele está pronto para receber um cliente **agora**?" Um garçom pode estar vivo e bem, mas ocupado limpando uma mesa — nesse momento você não manda cliente novo para ele. Se essa probe falhar, o Pod **continua vivo**, mas é **removido temporariamente do Service** (não recebe mais tráfego) até voltar a responder com sucesso.
- **livenessProbe** — "Ele desmaiou?" Se essa probe falhar repetidamente, o Kubernetes entende que o processo travou de um jeito que só reiniciar resolve, e **mata e recria o container**.

| Probe | Pergunta que responde | Se falhar, o que acontece |
|---|---|---|
| `startupProbe` | "Já terminou de inicializar?" | Bloqueia as outras probes até passar; se nunca passar, o container é reiniciado |
| `readinessProbe` | "Está pronto para tráfego agora?" | Pod é tirado do Service (para de receber requisições), mas continua rodando |
| `livenessProbe` | "Ainda está vivo/funcional?" | Container é **reiniciado** |

Comparando os parâmetros usados no YAML:

| Parâmetro | startupProbe | readinessProbe | livenessProbe |
|---|---|---|---|
| `periodSeconds` (intervalo entre checagens) | 10s | 10s | 10s |
| `failureThreshold` (falhas seguidas até agir) | 10 | 2 | 3 |
| `timeoutSeconds` (tempo até considerar falha) | — (padrão 1s) | 5s | 5s |
| `successThreshold` (sucessos seguidos para "curar") | — (padrão 1) | — (padrão 1) | 1 (explícito) |

Note que o `startupProbe` tem uma tolerância bem maior (`failureThreshold: 10`, ou seja, até 100 segundos) — faz sentido, pois inicialização costuma ser mais lenta e imprevisível que o funcionamento normal. Já o `readinessProbe` é mais rígido (`failureThreshold: 2`, só 20 segundos) porque queremos tirar rapidamente da fila de tráfego um Pod que não está respondendo bem.

**Sem probes configuradas**, o Kubernetes só sabe se o **processo** morreu (o container crashou de vez) — ele não tem como saber se a aplicação travou, ficou lenta demais, ou ainda está inicializando. É por isso que probes são consideradas essenciais em qualquer deploy real.

### Service e tipo LoadBalancer

```yaml
apiVersion: v1
kind: Service
metadata:
  name: serversvc
spec:
  type: LoadBalancer
  selector:
    app: server
  ports:
  - port: 8080
    targetPort: 8080
```

Pods são efêmeros: eles nascem, morrem, são recriados com **IPs diferentes** o tempo todo (imagine o Deployment reiniciando um Pod às 3h da manhã — o IP muda, mas ninguém deveria perceber). Se você tentasse acessar os Pods diretamente pelo IP, sua aplicação quebraria a cada reinício.

O `Service` resolve isso: é um **endereço estável** que sempre aponta para o conjunto de Pods certo, usando o `selector` (`app: server`) para descobrir dinamicamente quais Pods atender — os mesmos que o `Deployment` rotulou com `labels.app: server`. O Service atualiza sua lista de destinos automaticamente conforme Pods sobem e descem.

Existem três tipos principais de Service:

| Tipo | Acessível de onde | Uso típico |
|---|---|---|
| `ClusterIP` (padrão) | Só de dentro do cluster | Comunicação interna entre serviços (ex.: backend → banco de dados) |
| `NodePort` | De fora, via `IP-do-nó:porta-alta` | Testes locais, ambientes sem load balancer de nuvem disponível |
| `LoadBalancer` | De fora, via um IP público/load balancer dedicado | Expor um serviço diretamente para a internet, tipicamente em cloud (AWS, GCP, Azure provisionam um LB real) |

A aula usa `LoadBalancer` porque o objetivo é expor o servidor HTTP para fora do cluster, como uma aplicação real de internet faria. Em clusters locais (Minikube, Kind, Docker Desktop) esse tipo às vezes fica pendente (`<pending>`) esperando um provedor de load balancer, ou precisa de uma ferramenta adicional (ex.: `minikube tunnel`) para funcionar de verdade — vale saber disso ao testar localmente.

---

## ⚙️ Como as Peças se Conectam

```
                         Internet / cliente
                                │
                                ▼
                    ┌───────────────────────┐
                    │   Service: serversvc   │   type: LoadBalancer
                    │   port 8080 → 8080     │
                    └───────────┬───────────┘
                                │ selector: app=server
              ┌─────────────────┼─────────────────┐
              ▼                 ▼                 ▼
        ┌───────────┐    ┌───────────┐     ┌───────────┐
        │   Pod 1    │    │   Pod 2    │    │   Pod 3    │   ← Deployment "server"
        │ label:     │    │ label:     │    │ label:     │      (replicas: 3)
        │ app=server │    │ app=server │    │ app=server │
        │            │    │            │    │            │
        │ container: │    │ container: │    │ container: │
        │ imagem     │    │ imagem     │    │ imagem     │
        │ 21-deploy- │    │ 21-deploy- │    │ 21-deploy- │
        │ k8s:latest │    │ k8s:latest │    │ k8s:latest │
        │ :8080      │    │ :8080      │    │ :8080      │
        └───────────┘    └───────────┘     └───────────┘
        cada Pod é monitorado por startupProbe → readinessProbe → livenessProbe
```

O `Deployment` garante que sempre existam 3 Pods rodando a imagem correta; cada Pod é continuamente checado pelas probes; o `Service` observa quais Pods têm o label `app: server` **e estão prontos** (passaram na `readinessProbe`) e distribui o tráfego recebido só entre esses.

---

## ✅ Boas Práticas Presentes no Projeto

1. **Multi-stage build** (`Dockerfile.prod`) — separa a etapa de compilação da etapa de execução, evitando carregar ferramentas de build (compilador, cache do Go) na imagem final.
2. **Imagem final `scratch`** — reduz drasticamente o tamanho da imagem e a superfície de ataque, já que não existe shell, gerenciador de pacotes nem bibliotecas extras para um invasor explorar.
3. **Binário estaticamente linkado** (`CGO_ENABLED=0`) — necessário para rodar em `scratch`, e como bônus elimina dependências dinâmicas de sistema operacional que poderiam divergir entre ambientes.
4. **Múltiplas réplicas** (`replicas: 3`) — se um Pod cair ou estiver sendo atualizado, os outros dois continuam atendendo, evitando indisponibilidade total.
5. **As três probes configuradas** — dá ao Kubernetes visão real do estado da aplicação (inicializando / pronto para tráfego / travado), em vez de depender só de "o processo ainda existe?".

---

## 🚧 O que Poderia Ser Melhorado / Inconsistências Reais

Pare e pense antes de ler cada item — são problemas reais presentes nos arquivos desta aula, ótimos para praticar revisão crítica:

- **Falta `resources.requests`** — só `limits` está definido. Sem `requests`, o *scheduler* do Kubernetes não tem garantia de quanto de CPU/memória reservar para o Pod, o que pode causar má distribuição entre os nós do cluster (ver seção [Resources](#resources-requests-vs-limits)).
- **Tag de imagem `:latest`** (`wesleywillians/21-deploy-k8s:latest`) — usar `:latest` em produção significa que você não sabe, com certeza, qual versão do código está rodando em cada Pod, e um rolling update pode não disparar corretamente porque o Kubernetes pode não perceber que a imagem "mudou" (o nome da tag é o mesmo).
- **Sem `Namespace` definido** — todos os recursos vão parar no namespace `default`. Em um cenário real, isolar por ambiente/equipe (`namespace: producao`, `namespace: staging`) evita colisão de nomes e facilita controle de acesso.
- **Nome do módulo Go inconsistente** — `go.mod` referencia `21-deploy-k8s`, mas a pasta é `19-deploy-k8s` (e a imagem publicada segue o nome do módulo). Não quebra nada tecnicamente, mas é o tipo de detalhe que confunde quem lê o repositório depois.
- **Configuração 100% hardcoded** — a porta `8080` e o comportamento do servidor estão fixos no código, sem uso de `ConfigMap` ou `Secret` para externalizar configuração (fora do escopo desta aula introdutória, mas necessário em cenários reais — ver [Próximos Passos](#-próximos-passos)).

---

## ⚠️ Principais Problemas ao Trabalhar com Kubernetes

Erros comuns de quem está começando com k8s, com o problema e a correção:

**1. Esquecer probes**

❌ Deployment sem nenhuma probe configurada — o Kubernetes só percebe problemas quando o processo morre de vez, então uma aplicação travada (mas com processo vivo) continua recebendo tráfego indefinidamente.

✅ Configurar ao menos `readinessProbe` e `livenessProbe` (como neste projeto), para que o Kubernetes saiba diferenciar "vivo mas ocupado", "vivo e pronto" e "morto/travado".

**2. Usar `:latest` como tag de imagem**

❌ `image: minhaapp:latest` — impossível saber qual versão está rodando; rollback fica difícil ou impossível de rastrear.

✅ `image: minhaapp:v1.2.3` (ou o hash do commit) — cada deploy é rastreável e reversível com precisão.

**3. Não definir `requests`**

❌ Só `limits`, sem `requests` — o scheduler não reserva recursos, podendo empilhar Pods famintos no mesmo nó e gerar contenção de CPU/memória sob carga.

✅ Definir `requests` próximo do consumo real esperado, e `limits` como uma margem de segurança acima disso.

**4. Labels/selector que não batem**

❌ `Deployment.template.metadata.labels` diferente de `Service.spec.selector` — o Service simplesmente não encontra nenhum Pod, e o tráfego não chega em lugar nenhum, sem nenhum erro óbvio no log.

✅ Manter os labels consistentes (como `app: server` neste projeto, usado tanto no Deployment quanto no Service) e, idealmente, testar com `kubectl get endpoints <nome-do-service>` para confirmar que o Service encontrou Pods.

---

## ▶️ Como Executar

### Ambiente de desenvolvimento (local, com hot-reload manual)

```bash
# sobe o container de dev (mantém vivo com tail -f /dev/null)
docker-compose up -d

# entra no container e roda a aplicação manualmente
docker-compose exec goapp go run main.go

# acesse em http://localhost:8080
```

### Build da imagem de produção

```bash
docker build -f Dockerfile.prod -t minha-conta/19-deploy-k8s:v1 .
docker push minha-conta/19-deploy-k8s:v1   # se for usar num cluster remoto
```

> Se você for usar sua própria imagem, lembre de atualizar o campo `image:` em `k8s/deployment.yaml` para apontar para ela.

### Aplicando no cluster

```bash
# aplica Deployment e Service
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
# ou, de uma vez, aplicando toda a pasta:
kubectl apply -f k8s/

# acompanhar os Pods subindo (e ver as probes em ação)
kubectl get pods -w

# ver detalhes e eventos de um Pod específico (útil para debugar probes)
kubectl describe pod <nome-do-pod>

# conferir se o Service encontrou os Pods certos
kubectl get svc serversvc
kubectl get endpoints serversvc

# ver logs de um Pod
kubectl logs <nome-do-pod>
```

Se estiver usando um cluster local (Minikube, Kind, Docker Desktop) sem um load balancer real de nuvem, o `EXTERNAL-IP` do Service pode ficar como `<pending>` — nesses casos, use `kubectl port-forward svc/serversvc 8080:8080` para acessar localmente, ou `minikube tunnel` se estiver no Minikube.

---

## ⚖️ Trade-offs: Kubernetes vs Alternativas

| Abordagem | Vantagens | Desvantagens |
|---|---|---|
| **Kubernetes** | Self-healing, escalonamento, rolling updates, portável entre nuvens, ecossistema enorme | Curva de aprendizado alta, complexidade operacional, overhead de manter o próprio cluster (ou custo de um gerenciado) |
| **Docker Compose sozinho** | Simples, rápido para configurar, ótimo para dev local | Sem self-healing real entre múltiplas máquinas, sem escalonamento automático, não pensado para produção distribuída |
| **PaaS (Heroku, Cloud Run, Fly.io)** | Deploy extremamente simples, menos operação para você | Menos controle fino, pode ficar caro em escala, vendor lock-in mais forte |
| **VMs geridas manualmente** | Controle total | Você reimplementa manualmente tudo que o k8s já resolve (health check, restart, load balancing, deploy sem downtime) |

Kubernetes não é "melhor" em todo cenário — é uma ferramenta poderosa que resolve problemas de **escala e resiliência distribuída**, ao custo de complexidade operacional.

---

## 🎯 Quando Usar Kubernetes

**Faz sentido quando:**
- Você tem múltiplos serviços/microsserviços para orquestrar.
- Precisa de alta disponibilidade real (múltiplas réplicas, múltiplos nós/zonas).
- A equipe já tem ou está disposta a construir conhecimento operacional em k8s.
- Escalonamento automático baseado em carga é um requisito.

**Pode ser overkill quando:**
- É um projeto pequeno, um MVP, ou uma aplicação única sem necessidade real de múltiplas réplicas.
- O time é pequeno e não tem tempo/interesse em operar um cluster.
- Um PaaS simples (Cloud Run, Heroku, Render) resolveria o mesmo problema com muito menos complexidade.

---

## 💼 Perguntas de Entrevista Respondidas

**1. O que é Kubernetes e por que ele existe?**
É uma plataforma de orquestração de containers: você declara o estado desejado (quantas réplicas, quais recursos, como verificar saúde) e o Kubernetes garante que a realidade do cluster converge para esse estado continuamente, se recuperando sozinho de falhas. Existe porque rodar containers manualmente em várias máquinas, com self-healing, escalonamento e deploys sem downtime, é inviável de fazer à mão em escala.

**2. Qual a diferença entre Pod, ReplicaSet e Deployment?**
Pod é a menor unidade executável (um ou mais containers rodando juntos). ReplicaSet garante que exista sempre um número N de Pods de um determinado template rodando. Deployment gerencia ReplicaSets e adiciona a capacidade de fazer rolling updates/rollbacks controlados — na prática, você quase sempre interage com Deployments, não diretamente com Pods ou ReplicaSets.

**3. Qual a diferença entre `livenessProbe`, `readinessProbe` e `startupProbe`?**
`startupProbe` verifica se a aplicação terminou de inicializar (bloqueia as outras probes até passar). `readinessProbe` verifica se o Pod está pronto para receber tráfego *agora*; se falhar, o Pod sai do Service mas continua rodando. `livenessProbe` verifica se o processo ainda está funcional; se falhar repetidamente, o container é reiniciado.

**4. O que acontece se a `livenessProbe` falhar continuamente?**
O Kubernetes considera o container não-saudável e o **reinicia** (mata o container e sobe um novo, respeitando a política de restart do Pod). Se isso acontecer repetidamente em sequência muito rápida, o Kubernetes aplica um *backoff* exponencial (`CrashLoopBackOff`) antes de tentar de novo.

**5. Qual a diferença entre `resources.requests` e `resources.limits`?**
`requests` é o que o Pod precisa garantidamente para ser agendado — o *scheduler* usa isso para escolher em qual nó colocar o Pod. `limits` é o teto máximo de consumo permitido; ultrapassar o limite de memória mata o container (`OOMKilled`), ultrapassar o de CPU apenas o estrangula (throttling), sem matar.

**6. Qual a diferença entre `ClusterIP`, `NodePort` e `LoadBalancer`?**
`ClusterIP` (padrão) só é acessível de dentro do cluster — ideal para comunicação interna entre serviços. `NodePort` expõe uma porta alta em cada nó do cluster, acessível de fora sem depender de um load balancer de nuvem. `LoadBalancer` provisiona (tipicamente via um provedor de nuvem) um balanceador de carga externo dedicado, com IP público próprio — é o mais comum para expor um serviço diretamente à internet.

**7. Por que evitar a tag `:latest` em imagens de produção?**
Porque ela não identifica de forma única qual versão do código está rodando, dificulta rollback preciso e pode até impedir que o Kubernetes perceba que a imagem mudou durante um rolling update (o nome da tag continua o mesmo). O ideal é usar tags imutáveis, como uma versão semântica (`v1.2.3`) ou o hash do commit.

**8. O que é um rolling update e como o Deployment garante (quase) zero downtime?**
É a estratégia padrão de atualização do Deployment: em vez de derrubar todos os Pods antigos de uma vez, ele sobe Pods novos gradualmente (com a nova versão) e só derruba os antigos conforme os novos ficam prontos (passam na `readinessProbe`), mantendo sempre um número mínimo de Pods saudáveis atendendo tráfego durante toda a transição.

**9. Por que usar multi-stage build e uma imagem `scratch` para uma aplicação Go?**
Porque Go compila para um binário único, estaticamente linkado (com `CGO_ENABLED=0`), que não depende de nenhuma biblioteca do sistema operacional. Isso permite copiar só o binário final para uma imagem `scratch` (vazia), resultando numa imagem final minúscula e com superfície de ataque quase nula — sem shell, sem pacotes, sem nada além do próprio programa.

**10. Como o Kubernetes decide em qual nó (node) rodar um Pod?**
O componente `scheduler` do Kubernetes avalia, entre outros fatores, os `resources.requests` declarados pelo Pod contra a capacidade disponível em cada nó, além de regras de afinidade/anti-afinidade, taints/tolerations e restrições de recursos, e escolhe o nó que melhor atende a todos esses critérios simultaneamente.

**11. O que é um `Service` e por que os Pods não são acessados diretamente pelo IP?**
Pods são efêmeros — são recriados com frequência e cada recriação pode gerar um IP novo. Um `Service` fornece um endereço estável (nome DNS interno e, dependendo do tipo, um IP externo) que sempre aponta para o conjunto correto de Pods vivos e prontos, via `selector` de labels, então o cliente nunca precisa saber o IP de nenhum Pod individual.

**12. O que significa `CrashLoopBackOff`?**
É o estado que o Kubernetes reporta quando um container está falhando (crashando ou reprovando na `livenessProbe`) repetidamente, logo após reiniciar. O Kubernetes aumenta progressivamente o intervalo entre tentativas de restart (backoff exponencial) para não sobrecarregar o sistema tentando reiniciar algo que está falhando de forma consistente.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Cluster** | Conjunto de máquinas (nós) gerenciadas em conjunto pelo Kubernetes |
| **Node** | Uma máquina (física ou virtual) que faz parte do cluster e roda Pods |
| **Pod** | Menor unidade executável do Kubernetes; um ou mais containers rodando juntos |
| **Deployment** | Objeto que declara quantas réplicas de um Pod devem existir e gerencia atualizações |
| **ReplicaSet** | Objeto (geralmente gerenciado pelo Deployment) que garante a contagem de réplicas de Pods |
| **Service** | Endereço estável de rede que direciona tráfego para o conjunto correto de Pods |
| **Probe** | Verificação periódica de saúde de um container (`startup`, `readiness`, `liveness`) |
| **Label/Selector** | Sistema de "etiquetas" usado para associar objetos (ex.: Service encontra Pods via `selector` batendo com `labels`) |
| **Manifest** | Arquivo YAML que descreve declarativamente um recurso do Kubernetes |
| **kubectl** | CLI oficial para interagir com um cluster Kubernetes |
| **Rolling update** | Estratégia de atualização gradual de Pods, evitando downtime |
| **CrashLoopBackOff** | Estado de um container que está reiniciando repetidamente por falha |
| **Multi-stage build** | Técnica de Dockerfile com múltiplos `FROM`, separando build e runtime |
| **scratch** | Imagem Docker base completamente vazia, sem sistema operacional |

---

## 🚀 Próximos Passos

Esta aula cobre o essencial (Deployment, probes, Service), mas Kubernetes tem muito mais superfície. Para continuar estudando:

- [ ] **ConfigMap e Secret** — externalizar configuração e dados sensíveis em vez de hardcoded no código/imagem.
- [ ] **Ingress** — expor múltiplos serviços por HTTP/HTTPS através de um único ponto de entrada, com roteamento por domínio/path.
- [ ] **Namespace** — isolar recursos por ambiente ou equipe dentro do mesmo cluster.
- [ ] **HorizontalPodAutoscaler (HPA)** — escalar automaticamente o número de réplicas com base em métricas de uso (CPU, memória, ou métricas customizadas).
- [ ] **Helm** — gerenciador de pacotes para Kubernetes, útil para templatizar e versionar conjuntos de manifests.
- [ ] **PersistentVolume / PersistentVolumeClaim** — armazenamento persistente para aplicações com estado (bancos de dados, por exemplo).

### Conceitos relacionados no curso

- `aulas/14-cobra-cli` — CLI em Go, útil para entender ferramentas como o próprio `kubectl`.
- `aulas/18-clean-architecture` — organização de código que normalmente é a aplicação que acaba sendo deployada num cluster como este.
