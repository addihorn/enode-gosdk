include .env.mk

PACT_GO_VERSION=2.2.0
PACT_DOWNLOAD_DIR=/tmp

export PACT_DIR = $(PWD)/pacts
export LOG_DIR = $(PWD)/log
export VERSION_COMMIT?=$(shell git rev-parse HEAD)
export VERSION_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)



ifeq ($(OS),Windows_NT)
	PACT_DOWNLOAD_DIR=$$TMP
endif

run-consumer:
	@go run cmd/main.go
  
run-provider:
	@go run provider-go/cmd/usersvc/main.go


install_cli:
	@if [ ! -d pact/bin ]; then\
		echo "--- Installing Pact CLI dependencies";\
		curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash;\
    fi

install:
	go install github.com/pact-foundation/pact-go/v2@v$(PACT_GO_VERSION)
	pact-go -l DEBUG install --libDir $(PACT_DOWNLOAD_DIR);

broker:
	docker compose up -d

publish_oapi:
	pact/bin/pactflow publish-provider-contract refs/2024-10-01.oapi.json \
	--provider "ENODE" \
	--content-type application/json \
	-a 2024-10-01 \
	-h main \
	--verification-exit-code=0 \
	--verification-results=test_result.txt \
	--verification-results-content-type=text/plain \
	--verifier go \
	-b https://abc-corp.pactflow.io \
	-k dMXolvw9KmnOpmu8FLNVAA

consumer: export PACT_TEST := true
consumer:
	@echo "--- 🔨Running Consumer Pact tests "
	go test github.com/addihorn/enode-gosdk/pkg/users -v


publish:
	@echo "--- 📝 Publishing Pacts"
	pact/bin/pact-broker publish ${PWD}/pacts --consumer-app-version ${VERSION_COMMIT} --branch ${VERSION_BRANCH} \
		-b $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL) -k ${PACT_BROKER_TOKEN} --skip-merge
	@echo
	@echo "Pact contract publishing complete!"
	@echo
	@echo "Head over to $(PACT_BROKER_PROTO)://$(PACT_BROKER_URL)"
	@echo "to see your published contracts.	"
