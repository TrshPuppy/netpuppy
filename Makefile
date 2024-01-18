mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
package_path := ${shell dirname ${mkfile_path}}

# Define variables
base_python ?= python3
build_path = ${package_path}/build
env_path = ${package_path}/.venv
python = $(env_path)/bin/python
library_name = netpuppy
library_path = ${package_path}/${library_name}

# Define source files
source_files = $(shell find $(library_path) -name "*.py")
test_files = $(shell find $(package_path)/tests -name "*.py")

# Help command
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build:       Build the project"
	@echo "  clean:       Clean the project"
	@echo "  test:        Run the tests"
	@echo "  format:      Format the code using black"
	@echo "  check_types: Check types using mypy"
	@echo "  env/venv:    Create a virtual environment"
	@echo "  help:        Show this help message"

# Define header function
define header
	@printf "\033[0;32m==> $(1) \033[0m\n"
endef

# Define virtual environment
venv := $(env_path)/venv.lock
${venv}: pyproject.toml
	$(call header, Creating virtual environment)
	@$(base_python) -m venv --prompt="netpuppy" $(env_path)
	@${python} -m pip install --upgrade pip
	@${python} -m pip install -e .[dev]
	@touch ${venv}
	$(call header, Virtual environment created)
	@echo "Run 'source .venv/bin/activate' to activate the virtual environment."

.PHONY: venv
venv: ${venv}

.PHONY: env
env: venv

# Build the project
.PHONY: build
build: ${venv} ${source_files}
	$(call header, Building the project)
	@$(python) -m build
	@[ -d ${build_path} ] || mkdir -p ${build_path}
	@[ -d ${build_path}/package ] || mkdir -p ${build_path}/package
	@rm -rf ${build_path}/package/*
	@mv dist ${build_path}/package
	$(call header, Project built at ${build_path}/package)

check_format := ${build_path}/black.diff
${check_format}: ${venv} ${source_files} ${test_files}
	$(call header, Checking the code formatting)
	@mkdir -p ${build_path}
	@$(python) -m black --check --diff $(source_files) ${test_files} > ${check_format}

# Format the code using black
.PHONY: format
format: ${venv}
	$(call header, Formatting the code)
	@$(python) -m black $(source_files)

pytest_passed := ${build_path}/pytest_passed.lock
${pytest_passed}: ${venv} ${source_files} ${test_files}
	$(call header, Running tests)
	@[ -d ${build_path} ] || mkdir -p ${build_path}
	${python} -m pytest ${package_path} \
		--junitxml=${build_path}/report.xml \
		--cov=${library_name} \
		--cov-report=xml:${build_path}/coverage.xml
	@touch ${pytest_passed}

mypy_report := ${build_path}/coverage.json
${mypy_report}: ${venv} ${source_files} ${test_files}
	$(call header, Checking types)
	@${python} -m mypy ${source_files} ${test_files} \
		--incremental \
		--cache-dir=${build_path}/mypy_cache \
		--linecoverage-report ${build_path}
	
check_types: ${mypy_report}

# Run the tests
.PHONY: test
test: ${check_format} ${pytest_passed}

# Clean the project
.PHONY: clean
clean:
	$(call header, Cleaning the project)
	@rm -rf ${build_path}
	@rm -rf ${env_path}
	@rm -rf ${library_path}.egg-info
	@rm -rf ${library_path}/__pycache__
	@rm -rf ${package_path}/.pytest_cache
	@rm -rf ${package_path}/tests/__pycache__
	@rm -rf ${package_path}/.coverage