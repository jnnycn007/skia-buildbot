generate_sql_schema:
	cd ./go/schema_generate; go run main.go

INDEX_DATE ?= 2026-04-20

build:
	bazelisk build --config=mayberemote -c dbg //rag/...

run-auth-proxy:
	bazelisk run //kube/cmd/auth-proxy -- \
		--prom-port=:20003 \
		--role=editor=google.com \
		--authtype=mocked \
		--mock_user=$(USER)@google.com \
		--port=:8003 \
		--passive \
		--target_port=:8002

run-history-api: build
	../_bazel_bin/rag/go/api/server/server_/server api \
		--config_filename=./configs/chrome-internal.json \
		--resources_dir=../_bazel_bin/rag/pages/development/ \
		--local

run-history-api-in-memory: build
	../_bazel_bin/rag/go/api/server/server_/server api \
		--config_filename=./configs/chrome-internal.json \
		--resources_dir=../_bazel_bin/rag/pages/development/ \
		--local \
		--use_in_memory_store \
		--index_date="$(INDEX_DATE)"

run-history-api-local: run-auth-proxy run-history-api