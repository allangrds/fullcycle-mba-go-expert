# FullCycle | MBA | Go Expert

Repository created to store code, examples, and exercises completed during the **Go Expert** MBA in [FullCycle](https://fullcycle.com.br/)


## Table of contents
- [About](#about)
- [Go CLI Commands](#go-cli-commands)
  - [Information and Environment](#information-and-environment)
  - [Execution](#execution)
  - [Module Management](#module-management)
  - [Application Build](#application-build)
- [Useful links](#useful-links)

## About

This repository contains:

- Notes and summaries about Go
- Practical course exercises
- Experiment and test code
- Useful references for consultation


### Go CLI Commands

#### Information and Environment
- `go env` — Displays environment variables configured for Go
  ```bash
  go env GOOS
  ````

#### Execution
- `go run` — Compiles and runs the Go code
  ```bash
    go run main.go
  ```

#### Module Management
- `go mod init <name>` — Initializes a new Go module
- `go mod tidy `— Downloads dependencies and updates go.mod
- `go get <package>` — Downloads packages and dependencies
  ```bash
    go get github.com/google/uuid
  ```

#### Application Build
- `go build` — Compiles the current project
  ```bash
    go build main.go
    GOOS=windows go build main.go
    GOOS=linux go build main.go
  ```
  Reference: [Cross Compilation com Go](https://www.digitalocean.com/community/tutorials/building-go-applications-for-different-operating-systems-and-architectures)

#### Test
- `go test -v` — Execute tests inside folder
- `go test -cover` — Execute tests with coverage percent
- `go test -coverprofile=coverage.out` — Execute tests with coverage and putting the results in a file percent
- `go tool cover -html=coverage.out` — Print a coverage.html file using the coverage.out to show exactly covered/uncovered lines
- `go test -bench=.` - Execute tests and benchmark test
- `go test -bench=. -run=^#` - Execute only benchmark test
- `go test -fuzz=. -run=^# -fuzztime=3s` - Execute fuzz test with 3 seconds timeout

### Useful links

- [Official Go Documentation](https://go.dev/doc/)
- [Tour do Go](https://tour.golang.org/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
