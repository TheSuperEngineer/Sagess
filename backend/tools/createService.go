package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	serviceName := flag.String("name", "", "Name of the service (e.g., user, payment)")
	language := flag.String("language", "", "Service language: go or python")
	flag.Parse()

	args := flag.Args()
	if *serviceName == "" && len(args) > 0 {
		*serviceName = args[0]
	}
	if *language == "" && len(args) > 1 {
		*language = args[1]
	}

	if *serviceName == "" || *language == "" {
		scanner := bufio.NewScanner(os.Stdin)
		if *serviceName == "" {
			*serviceName = prompt(scanner, "Service name (e.g., user, payment): ")
		}
		if *language == "" {
			*language = prompt(scanner, "Service language (go/python): ")
		}
	}

	normalizedLanguage := normalizeLanguage(*language)
	if *serviceName == "" {
		fmt.Println("Please provide a service name using -name flag or as the first argument")
		os.Exit(1)
	}
	if normalizedLanguage == "" {
		fmt.Println("Please provide a valid service language: go or python")
		os.Exit(1)
	}

	switch normalizedLanguage {
	case "go":
		createGoService(*serviceName)
	case "python":
		createPythonService(*serviceName)
	}
}

func prompt(scanner *bufio.Scanner, label string) string {
	fmt.Print(label)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func normalizeLanguage(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "go", "golang":
		return "go"
	case "py", "python", "python3":
		return "python"
	default:
		return ""
	}
}

func createGoService(serviceName string) {
	// Create service directory structure
	basePath := filepath.Join("services", serviceName+"-service")
	dirs := []string{
		"cmd",
		"internal/domain",
		"internal/service",
		"internal/infrastructure/events",
		"internal/infrastructure/grpc",
		"internal/infrastructure/repository",
		"pkg/types",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Create an empty README.md
	readmePath := filepath.Join(basePath, "README.md")
	readmeContent := fmt.Sprintf(`# %s service

This service handles all %s-related operations in the system.

## Architecture

The service follows Clean Architecture principles with the following structure:

`+"```"+`
services/%s-service/
├── cmd/                    # Application entry points
│   └── main.go            # Main application setup
├── internal/              # Private application code
│   ├── domain/           # Business domain models and interfaces
│   ├── service/          # Business logic implementation
│   │   └── service.go    # Service implementations
│   └── infrastructure/   # External dependencies implementations (abstractions)
│       ├── events/       # Event handling (RabbitMQ)
│       ├── grpc/         # gRPC server handlers
│       └── repository/   # Data persistence
├── pkg/                  # Public packages
│   └── types/           # Shared types and models
└── README.md            # This file
`+"```"+`

### Layer Responsibilities

1. **Domain Layer** (`+"`internal/domain/`"+`)
   - Contains business domain interfaces
   - Defines contracts for repositories and services
   - Pure business logic, no implementation details

2. **Service Layer** (`+"`internal/service/`"+`)
   - Implements business logic
   - Uses repository interfaces
   - Coordinates between different parts of the system

3. **Infrastructure Layer** (`+"`internal/infrastructure/`"+`)
   - `+"`repository/`"+`: Implements data persistence
   - `+"`events/`"+`: Handles event publishing and consuming
   - `+"`grpc/`"+`: Handles gRPC communication

4. **Public Types** (`+"`pkg/types/`"+`)
   - Contains shared types and models
   - Can be imported by other services

## Key Benefits

1. **Dependency Inversion**: Services depend on interfaces, not implementations
2. **Separation of Concerns**: Each layer has a specific responsibility
3. **Testability**: Easy to mock dependencies for testing
4. **Maintainability**: Clear boundaries between components
5. **Flexibility**: Easy to swap implementations without affecting business logic
`, serviceName, serviceName, serviceName)

	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		fmt.Printf("Error creating README.md: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s Go service structure in %s\n", serviceName, basePath)
	fmt.Println("\nDirectory structure created:")
	fmt.Printf(`
services/%s-service/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── domain/           # Business domain models and interfaces
│   │   └── %s.go         # Core domain interfaces
│   ├── service/          # Business logic implementation
│   │   └── service.go    # Service implementations
│   └── infrastructure/   # External dependencies implementations (abstractions)
│       ├── events/       # Event handling (RabbitMQ)
│       ├── grpc/         # gRPC server handlers
│       └── repository/   # Data persistence
├── pkg/                  # Public packages
│   └── types/           # Shared types and models
└── README.md            # This file
`, serviceName, serviceName)
}

func createPythonService(serviceName string) {
	basePath := filepath.Join("services", serviceName+"-service")
	packageName := pythonPackageName(serviceName)

	dirs := []string{
		filepath.Join("src", packageName),
		"tests",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	files := map[string]string{
		filepath.Join(basePath, ".python-version"): "3.12\n",
		filepath.Join(basePath, "pyproject.toml"):  pythonPyproject(serviceName),
		filepath.Join(basePath, "README.md"):       pythonReadme(serviceName, packageName),
		filepath.Join(basePath, "src", packageName, "__init__.py"): fmt.Sprintf(`"""%s service."""
`, serviceName),
		filepath.Join(basePath, "src", packageName, "main.py"): fmt.Sprintf(`def main() -> None:
    print("Starting %s service")


if __name__ == "__main__":
    main()
`, serviceName),
		filepath.Join(basePath, "tests", "__init__.py"): "",
		filepath.Join(basePath, "tests", "test_main.py"): fmt.Sprintf(`from %s.main import main


def test_main_runs(capsys):
    main()
    assert "Starting %s service" in capsys.readouterr().out
`, packageName, serviceName),
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Printf("Error creating %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully created %s Python service structure in %s\n", serviceName, basePath)
	fmt.Println("\nDirectory structure created:")
	fmt.Printf(`
services/%s-service/
├── .python-version       # Python version used by uv
├── pyproject.toml        # uv project configuration
├── README.md            # This file
├── src/
│   └── %s/
│       ├── __init__.py
│       └── main.py      # Application entry point
└── tests/
    ├── __init__.py
    └── test_main.py

Next steps:
  cd %s
  uv sync --dev
  uv run python -m %s.main
`, serviceName, packageName, basePath, packageName)
}

func pythonPackageName(serviceName string) string {
	name := strings.ToLower(serviceName)
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.Trim(name, "_")
	if name == "" {
		return "service"
	}
	if strings.HasSuffix(name, "_service") {
		return name
	}
	return name + "_service"
}

func pythonPyproject(serviceName string) string {
	return fmt.Sprintf(`[project]
name = "%s-service"
version = "0.1.0"
description = "%s service"
readme = "README.md"
requires-python = ">=3.12"
dependencies = []

[dependency-groups]
dev = [
    "pytest>=8.0.0",
]
`, serviceName, serviceName)
}

func pythonReadme(serviceName string, packageName string) string {
	return fmt.Sprintf(`# %s service

This service handles all %s-related operations in the system.

## Project management

This Python service uses uv for project and dependency management.

`+"```"+`bash
uv sync --dev
uv run python -m %s.main
uv run pytest
`+"```"+`

## Structure

`+"```"+`
services/%s-service/
├── .python-version
├── pyproject.toml
├── README.md
├── src/
│   └── %s/
│       ├── __init__.py
│       └── main.py
└── tests/
    ├── __init__.py
    └── test_main.py
`+"```"+`
`, serviceName, serviceName, packageName, serviceName, packageName)
}
