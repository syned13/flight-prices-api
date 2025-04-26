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

.PHONY: run-with-mocks
run-with-mocks:
	docker compose -f docker-compose-mock.yml -f docker-compose-local.yml up 

.PHONY: run-with-mocks-build
run-with-mocks-build:
	docker compose -f docker-compose-mock.yml -f docker-compose-local.yml up --build
