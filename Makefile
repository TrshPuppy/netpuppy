mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
package_path := ${shell dirname ${mkfile_path}}

# Define variables
base_python = python
build_path = ${package_path}/build
env_path = ${package_path}/.venv
PYTHON = $(env_path)/bin/python
library_name = netpuppy
library_path = ${package_path}/${library_name}

source_files = $(shell find $(library_path) -name "*.py")

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

# Define header function
define header
	@printf "\033[0;32m==> $(1) \033[0m\n"
endef

# Define virtual environment
venv := $(env_path)/venv.lock
${venv}: pyproject.toml
	$(call header, Creating virtual environment)
	@$(base_python) -m venv --prompt="netpuppy" $(env_path)
	@${PYTHON} -m pip install --upgrade pip
	@${PYTHON} -m pip install -e .[dev]
	@touch ${venv}
	$(call header, Virtual environment created)
	@echo "Run 'source ${env_path}/bin/activate' to activate the virtual environment."

.PHONY: venv
venv: ${venv}

# Build the project
build:
	$(call header, Building the project)
	@$(PYTHON) -m build
	@mv dist ${build_path}

# Format the code using black
.PHONY: format
format: ${venv}
	$(call header, Formatting the code)
	@$(PYTHON) -m black $(source_files)

check_format: ${build_path}/black.diff
${check_format}: ${venv} ${source_files}
	$(call header, Checking the code formatting)
	@mkdir -p ${build_path}
	@$(PYTHON) -m black --check ---diff $(source_files) > ${check_format}

# Run the tests
.PHONY: test
test: ${check_format}

# Clean the project
.PHONY: clean
clean:
	$(call header, Cleaning the project)
	@rm -rf ${build_path}
	@rm -rf ${env_path}