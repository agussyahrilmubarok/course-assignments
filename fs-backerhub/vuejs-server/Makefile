# Makefile for Vue + PNPM project
# Commands prefixed with 'pnpm/' to manage the project

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make pnpm/install               # Install project dependencies"
	@echo "  make pnpm/dev                   # Start development server"
	@echo "  make pnpm/build                 # Build the project for production"
	@echo "  make pnpm/preview               # Preview production build"
	@echo "  make pnpm/update                # Update all dependencies (with backup)"
	@echo "  make pnpm/update PACKAGES=\"vue axios\"  # Update specific dependencies"
	@echo "  make pnpm/outdated              # Check outdated dependencies"

# Install dependencies
.PHONY: pnpm/install
pnpm/install:
	@echo "Installing dependencies..."
	pnpm install

# Start development server
.PHONY: pnpm/dev
pnpm/dev:
	@echo "Starting dev server..."
	pnpm run dev

# Build for production
.PHONY: pnpm/build
pnpm/build:
	@echo "Building project..."
	pnpm run build

# Preview production build
.PHONY: pnpm/preview
pnpm/preview:
	@echo "Previewing production build..."
	pnpm run preview

# Check outdated dependencies
.PHONY: pnpm/outdated
pnpm/outdated:
	@echo "Checking outdated dependencies..."
	pnpm outdated

# Update dependencies (all or specific)
# pnpm up axios
# Usage: make pnpm/update             # updates all dependencies
#        make pnpm/update PACKAGES="vue axios"  # updates specific packages
.PHONY: pnpm/update
pnpm/update:
	@echo "Backing up package.json and pnpm-lock.yaml..."
	cp package.json package.json.bak
	cp pnpm-lock.yaml pnpm-lock.yaml.bak
ifndef PACKAGES
	@echo "Updating all dependencies..."
	pnpm up
else
	@echo "Updating specified dependencies: $(PACKAGES)"
	pnpm up $(PACKAGES)
endif
	@echo "Checking for outdated packages..."
	pnpm outdated
	@echo "Update complete. Test your project with 'make pnpm/dev'"
