# Check for python or python3
PYTHON_EXISTS := $(shell python -c "print('exists')" 2>/dev/null)
PYTHON3_EXISTS := $(shell python3 -c "print('exists')" 2>/dev/null)

ifeq ($(PYTHON_EXISTS),exists)
BASE_PYTHON := python
else ifeq ($(PYTHON3_EXISTS),exists)
BASE_PYTHON := python3
else
$(error "No python or python3 found. Please install Python.")
endif

# Define variables
VENV_DIR := .venv
PIP := $(VENV_DIR)/bin/pip
PYTEST := $(VENV_DIR)/bin/pytest
BLACK := $(VENV_DIR)/bin/black
MYPY := $(VENV_DIR)/bin/mypy
PYTHON := $(VENV_DIR)/bin/python
DEL_COMMAND := rm -rf

# Define phony targets
.PHONY: all clean test format env

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

# Build the project
build:
	@$(PYTHON) -m build

# Format the code using black
format:
	@if [ ! -f "$(VENV_DIR)/bin/python" ] && [ ! -f "$(VENV_DIR)/bin/python.exe" ]; then \
		echo "Error: Virtual environment not set up. Run 'make venv' first."; \
		exit 1; \
	fi
	@$(BLACK) .

# Run the tests
test:
	@if [ ! -f "$(VENV_DIR)/bin/python" ] && [ ! -f "$(VENV_DIR)/bin/python.exe" ]; then \
		echo "Error: Virtual environment not set up. Run 'make venv' first."; \
		exit 1; \
	fi
	@$(PYTEST)
	@$(MYPY) --config-file=pyproject.toml netpuppy

# Create a virtual environment
env:
	@if [ -d "$(VENV_DIR)" ]; then \
		echo "Virtual environment already exists."; \
		echo "To activate the virtual environment, run:"; \
		echo "source $(VENV_DIR)/bin/activate"; \
	else \
		echo "Creating a virtual environment..."; \
		$(BASE_PYTHON) -m venv --prompt netpuppy $(VENV_DIR); \
		$(PIP) install -e .[dev]; \
		echo "To activate the virtual environment, run:"; \
		echo "source $(VENV_DIR)/bin/activate"; \
	fi
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