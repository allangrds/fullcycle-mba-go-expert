.PHONY: up stop destroy test logs logs-a logs-b logs-otel ps

up: ## Builda (se necessário) e sobe todos os serviços em primeiro plano (Ctrl+C para parar)
	docker compose up --build

stop: ## Para os containers do projeto, sem removê-los
	docker compose stop

destroy: ## Derruba containers, remove imagens locais, volumes e redes do projeto
	docker compose down --rmi local --volumes --remove-orphans

logs: ## Acompanha os logs de todos os serviços
	docker compose logs -f

logs-a: ## Acompanha os logs do service-a
	docker compose logs -f service-a

logs-b: ## Acompanha os logs do service-b
	docker compose logs -f service-b

logs-otel: ## Acompanha os logs do otel-collector (útil pra confirmar que os spans estão chegando)
	docker compose logs -f otel-collector

ps: ## Lista os containers do projeto
	docker compose ps

test: ## Builda (se necessário) e roda toda a suite de testes
	docker compose run --build --rm test
