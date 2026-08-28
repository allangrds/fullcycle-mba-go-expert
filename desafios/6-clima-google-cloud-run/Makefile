.PHONY: up stop destroy test logs ps

up: ## Builda (se necessário) e sobe a API em primeiro plano (Ctrl+C para parar)
	docker compose up --build

stop: ## Para os containers do projeto, sem removê-los
	docker compose stop

destroy: ## Derruba containers, remove imagens locais (app/test), volumes e redes do projeto
	docker compose down --rmi local --volumes --remove-orphans

logs: ## Acompanha os logs do app
	docker compose logs -f app

ps: ## Lista os containers do projeto
	docker compose ps

test: ## Builda (se necessário) e roda toda a suite de testes
	docker compose run --build --rm test
