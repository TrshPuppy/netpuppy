# check if the OS is Windows_NT
ifeq ($(OS),Windows_NT)
	ENV_NAME := Scripts
else
	ENV_NAME := bin
endif

# Define variables
VENV_DIR := .venv
PYTHON := python
PIP := $(VENV_DIR)/${ENV_NAME}/pip
PYTEST := $(VENV_DIR)/${ENV_NAME}/pytest
BLACK := $(VENV_DIR)/${ENV_NAME}/black
MYPY := $(VENV_DIR)/${ENV_NAME}/mypy

# Define phony targets
.PHONY: all clean test format env

# Detect the shell
ifdef MSYSTEM
    DEL_COMMAND := rm -rf
else
    # Windows Command Prompt or PowerShell
    DEL_COMMAND := del /s /q
endif

print-del-command:
	@echo $(DEL_COMMAND)

# Build the project
build:
	@$(PYTHON) -m build

# Format the code using black
format:
	@$(BLACK) .

# Run the tests
test:
	@$(PYTEST)
	@$(MYPY) --config-file=pyproject.toml netpuppy

# Create a virtual environment
env:
	@$(PYTHON) -m venv $(VENV_DIR)
	@if [ "${MSYSTEM}" ]; then \
		source $(VENV_DIR)/Scripts/activate && $(PIP) install -e .[dev]; \
	else \
		source $(VENV_DIR)/bin/activate && $(PIP) install -e .[dev]; \
	fi
	@echo "To activate the virtual environment, run:"
	@echo "source $(VENV_DIR)/${ENV_NAME}/activate"
venv: env

# Clean the project
clean:
	-@$(DEL_COMMAND) build/
	-@$(DEL_COMMAND) dist/
	-@$(DEL_COMMAND) *.egg-info/
	-@$(DEL_COMMAND) __pycache__/
	-@$(DEL_COMMAND) $(VENV_DIR)
	-@$(DEL_COMMAND) .mypy_cache
	-@$(DEL_COMMAND) .pytest_cache

# Help command
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build:    Build the project"
	@echo "  clean:    Clean the project"
	@echo "  test:     Run the tests"
	@echo "  format:   Format the code using black"
	@echo "  env/venv: Create a virtual environment"
	@echo "  help:     Show this help message"