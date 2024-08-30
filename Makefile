OAPI_CODEGEN_VERSION := 2.3.0
SWAGGER_CLI_CONTAINER_VERSION := 1.1.0

.PHONY: bootstrap
bootstrap: bootstrap_oapi bootstrap_swagger-cli

.PHONY: bootstrap_swagger-cli
bootstrap_swagger-cli:
	docker build -f ./docker/swagger_cli/Dockerfile -t swagger_cli:$(SWAGGER_CLI_CONTAINER_VERSION) .

.PHONY: bootstrap_oapi
bootstrap_oapi:
	mkdir -p ./bin
	GOBIN=${PWD}/bin/ go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v$(OAPI_CODEGEN_VERSION)

.PHONY: server_generate
server_generate:
	bin/oapi-codegen -config docs/openapi/config/oapi.yaml docs/openapi/openapi.yaml

.PHONY: open_api_bundle
open_api_bundle:
	docker run \
		--mount "type=bind,source=${PWD}/docs/openapi/,target=/openapi/" \
		--mount "type=bind,source=${PWD}/app/general/internal/initialize/,target=/openapi_bundle/" \
		swagger_cli:$(SWAGGER_CLI_CONTAINER_VERSION)

.PHONY: deploy
deploy:
	gcloud builds submit --config=cloudbuild.yaml
