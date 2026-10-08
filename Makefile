.PHONY: up down logs restart tidy test health demo step1 step2 step3 step4

up:
	docker compose -f compose.yaml up --build -d

down:
	docker compose -f compose.yaml down

logs:
	docker compose -f compose.yaml logs -f app

restart:
	docker compose -f compose.yaml restart app

tidy:
	go mod tidy

test:
	go test ./...

health:
	curl -s http://localhost:8080/health | jq .

demo:
	@echo "== list =="
	@curl -s http://localhost:8080/tasks | jq .
	@echo "\n== create =="
	@curl -s -X POST http://localhost:8080/tasks \
	  -H 'Content-Type: application/json' \
	  -d '{"title":"APIを叩いてみる","description":"curlでCRUDを確認"}' | jq .
	@echo "\n== list again =="
	@curl -s http://localhost:8080/tasks | jq .

# 各ステップのサンプルをローカルで実行（Docker外）
step1:
	go run ./steps/01-hello

step2:
	go run ./steps/02-http

step3:
	@echo "step3 は MySQL が必要です。先に make up してから実行してください"
	DB_HOST=127.0.0.1 go run ./steps/03-mysql

step4:
	@echo "step4 は MySQL が必要です。先に make up してから実行してください"
	DB_HOST=127.0.0.1 go run ./steps/04-layered
