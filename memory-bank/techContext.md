# Technical Context

## Technologies Used

### Programming Language
- **Go** (version 1.23)
  - Standard library focused approach
  - Minimal external dependencies

### Database
- **SQLite**
  - Embedded database
  - File-based storage
  - No separate server required

### Development Environment
- **VSCode** with Dev Containers
  - Consistent development environment
  - Isolated from host system
  - Pre-configured tools and extensions

## Development Setup

### Prerequisites
- Docker and Docker Compose
- VSCode with Remote Containers extension

### Setup Process
1. Clone the repository
2. Open in VSCode
3. Use the "Remote-Containers: Open Folder in Container" command
4. Development container will be built and started automatically

## Build Process

### Linux Build
```bash
go build -o sqlite-playground main.go
```

### Windows Build
```bash
GOOS=windows GOARCH=amd64 go build -o sqlite-playground.exe main.go
```

## Testing Framework

### Test Command
```bash
go test -v -cover ./...
```

### Testing Approach
- Unit tests for business logic (usecase package)
- Tests follow Arrange-Act-Assert pattern
- Test files named with _test.go suffix
- Test functions named TestXxx where Xxx describes the business rule being tested

## Technical Constraints

### Code Conventions
- Follow Go standard library patterns
- Clear naming for functions, arguments, variables, and constants
- Comments explaining function purpose, parameters, and return values
- Complex logic broken down into smaller, named operations
- Literal values defined as constants at the top of files

### Architecture Constraints
- Domain models defined in domain/models package
- Infrastructure functions only convert between models and external formats
- Business logic in usecase package as pure functions
- Controller functions coordinate calls to infrastructure and usecase functions
- Handler functions process input parameters and call controller functions
