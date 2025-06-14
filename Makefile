SHELL := /bin/bash

.PHONY: run local-run clean sqlc

help: ## This help dialog.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-30s %s\n" $$help_command $$help_info ; \
	done

run: ## Run server
	go run ./cmd/main.go

local-run: ## Run docker compose with local env file
	docker-compose --env-file .env.local up -d && docker-compose logs -f

clean: ## clean docker
	rm -rf ./data && docker-compose down

sqlc: ## Generate sqlc
	sqlc generate