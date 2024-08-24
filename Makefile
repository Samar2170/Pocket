IP := $(shell hostname -I | awk '{print $$1}')
PATTERN := HOSTNAME=*
HOMEDIR := $(shell echo $$HOME)

BASEDIR := $(HOMEDIR)/Desktop/PROJECTS/fxb/fxb/
define HOSTVAR
HOSTNAME=$(IP)\nBASEDIR=$(BASEDIR)
endef

clear-existing-hostname:
	@grep -vE "^$(PATTERN)$@" < .env > .env.tmp
	@mv .env.tmp .env

clear-existing-basedir:
	@grep -vE "^BASEDIR=*" < .env > .env.tmp
	@mv .env.tmp .env

make-logs-dir:
	@mkdir -p logs

.PHONY:
run:
	@echo "Ip captured: $(IP)"
	@echo "HOME directory: $(HOMEDIR)"
	@echo "---------------------------------"
	@echo "Clearing existing hostname..."
	@make clear-existing-hostname
	@echo "Clearing existing base directory..."
	@echo "---------------------------------"
	@make clear-existing-basedir
	@echo "Overwriting the .env file..."
	@echo "$(HOSTVAR)" > .env
	@echo "---------------------------------"
	@echo "Creating logs directory..."
	@make make-logs-dir
	@echo "Running the application..."
	@go run .
	