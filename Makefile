.PHONY: up down build logs health test

up:
	docker compose up --build

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
		dietician.local/services/profile-service/... \
		dietician.local/services/progress-service/... \
		dietician.local/services/recommendation-service/... \
		dietician.local/services/medical-service/... \

test-v:
	@echo "=== Running all tests (verbose) ==="
	go test -v dietician.local/packages/... \
		dietician.local/services/gateway/... \
		dietician.local/services/account-service/... \
		dietician.local/services/profile-service/... \
		dietician.local/services/progress-service/... \
		dietician.local/services/recommendation-service/... \
		dietician.local/services/medical-service/... \
