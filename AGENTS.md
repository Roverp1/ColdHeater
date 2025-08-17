# AGENTS.md - Development Guidelines

## Build/Test Commands
- `go run main.go` - Run the application
- `go build` - Build the binary
- `go test ./...` - Run all tests  
- `go test ./internal/database` - Run tests for specific package
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run static analysis

## Code Style Guidelines
- **Imports**: Standard library first, then external packages, then internal packages
- **Naming**: Use camelCase for variables/functions, PascalCase for exported types
- **Error handling**: Always check errors explicitly with `if err != nil`
- **Pointers**: Use `*string` for nullable database fields, direct scanning preferred
- **Database**: Use `defer rows.Close()` and proper resource cleanup
- **Config**: Load from `configs/app.yaml` using struct tags
- **Struct tags**: Use `yaml:"field_name"` for config, `db:"column_name"` for database

## Architecture
- Store business logic in `internal/` packages
- Use dependency injection for config and database connections
- Follow existing patterns in `internal/database/bot_operations.go`
- CLI interface in `internal/ui/cli/`