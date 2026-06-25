.PHONY: up up-dev down build logs health test swag tidy

up:
	docker compose -f docker-compose.dev.yml up

# Allows passing arguments to specific commands, e.g., make up-dev --build service_name
ifneq (,$(filter $(firstword $(MAKECMDGOALS)),up-dev up-infra))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

up-dev:
	docker compose -f docker-compose.dev.yml up $(RUN_ARGS)

up-infra:
	docker compose -f docker-compose.infra.yml up $(RUN_ARGS)

up-d:
	docker compose up --build -d

down:
	docker compose down

down-v:
	docker compose down -v

build:
	docker compose build

logs:
	docker compose logs -f

health:
	@echo "=== Gateway Health ==="
	@curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || echo "Gateway unreachable"

health-all:
	@echo "=== Gateway ===" && curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || echo "unreachable"
	@echo "\n=== account ===" && curl -s http://localhost:8081/health | python3 -m json.tool 2>/dev/null || echo "unreachable"
	@echo "\n=== Profile ===" && curl -s http://localhost:8082/health | python3 -m json.tool 2>/dev/null || echo "unreachable"
	@echo "\n=== Progress ===" && curl -s http://localhost:8083/health | python3 -m json.tool 2>/dev/null || echo "unreachable"
	@echo "\n=== Recommendation ===" && curl -s http://localhost:8084/health | python3 -m json.tool 2>/dev/null || echo "unreachable"
	@echo "\n=== Medical ===" && curl -s http://localhost:8085/health | python3 -m json.tool 2>/dev/null || echo "unreachable"
	@echo "\n=== Tracking ===" && curl -s http://localhost:8086/health | python3 -m json.tool 2>/dev/null || echo "unreachable"

ps:
	docker compose ps

test:
	@echo "=== Running all tests ==="
	go test dietician.local/packages/... \
		dietician.local/services/gateway/... \
		dietician.local/services/account-service/... \
		dietician.local/services/progress-service/... \
		dietician.local/services/recommendation-service/... \
		dietician.local/services/medical-service/... \

test-v:
	@echo "=== Running all tests (verbose) ==="
	go test -v dietician.local/packages/... \
		dietician.local/services/gateway/... \
		dietician.local/services/account-service/... \
		dietician.local/services/progress-service/... \
		dietician.local/services/recommendation-service/... \
		dietician.local/services/medical-service/... \

swag:
	@for d in services/*; do \
		if [ -f "$$d/Makefile" ]; then \
			echo "=== Running make swag for $$d ==="; \
			$(MAKE) -C "$$d" swag || true; \
		fi \
	done

postman: swag
	@echo "=== Converting swagger to Postman collection ==="
	@node scripts/swag2postman.js

wire:
	@for d in services/*; do \
		if [ -f "$$d/Makefile" ]; then \
			echo "=== Running make wire for $$d ==="; \
			$(MAKE) -C "$$d" wire || true; \
		fi \
	done

tidy:
	@for d in services/*; do \
		if [ -f "$$d/go.mod" ]; then \
			echo "=== Running go mod tidy for $$d ==="; \
			(cd "$$d" && go mod tidy); \
		fi \
	done
