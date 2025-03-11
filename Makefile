SHELL := $(shell echo $$SHELL)

##@ Utility
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development
backend-run:  ## Run the backend server
	set -a && source .env.development && set +a && \
		cd backend && go run server.go

backend-test:  ## Run the backend tests
	set -a && source .env.development && set +a && \
		cd backend && go test -v ./...

##@ Terraform
_tf-cmd:
	@if [ -z "$(ENV)" ]; then \
		echo "ENV is not set"; \
		exit 1; \
	fi
	cd tf && tofu $(TF_CMD)

tf-init:  ## Initialize the Terraform project
	@TF_CMD=init $(MAKE) _tf-cmd

tf-plan:  ## Plan the Terraform changes
	@TF_CMD=plan $(MAKE) _tf-cmd

tf-apply:  ## Apply the Terraform changes
	@TF_CMD=apply $(MAKE) _tf-cmd

tf-destroy:  ## Destroy the Terraform resources
	@TF_CMD=destroy $(MAKE) _tf-cmd


