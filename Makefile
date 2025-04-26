.PHONY: run-local-build
run-local-build:
	docker compose -f docker-compose-mock.yml up --build
	docker compose -f docker-compose-local.yml up --build

.PHONY: run-local
run-local:
	docker compose -f docker-compose-mock.yml up
	docker compose -f docker-compose-local.yml up

.PHONY: run-build
run-build:
	docker compose up --build
