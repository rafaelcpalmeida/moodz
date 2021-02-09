ENV_FILE := .env

.PHONY: help run build-and-run

help:	# The following lines will print the available commands when entering just 'make', as long as a comment with three cardinals is found after the recipe name. Look at the examples below.
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

build:			### Builds all of the containers required by the application, without starting them.
	@docker-compose up --no-start

run:			### Starts all of the containers required by the application.
ifeq ("$(wildcard $(ENV_FILE))","")
	$(error Missing .env file. Make sure to clone it from .env.example file)
endif
	@docker-compose up

build-and-run:	### Starts all of the containers required by the application, rebuilding the images first.
ifeq ("$(wildcard $(ENV_FILE))","")
	$(error Missing .env file. Make sure to clone it from .env.example file)
endif
	@docker-compose up --build
