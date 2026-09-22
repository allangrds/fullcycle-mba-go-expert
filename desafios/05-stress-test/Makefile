.PHONY: up stop destroy test

up: ## Builda (se necessário) e roda o teste de carga via Docker Compose, lendo URL/REQUESTS/CONCURRENCY do .env
	docker compose up --build stress-test

stop: ## Para os containers do projeto, sem removê-los
	docker compose stop

destroy: ## Remove containers, rede, imagens locais e volumes do projeto
	docker compose down --rmi local --volumes --remove-orphans

test: ## Builda (se necessário) e roda a suíte de testes via Docker Compose
	docker compose run --build --rm test
