# Lost Dog Matcher

Lost Dog Matcher is the foundation for a future WhatsApp-based service that will help reunite lost dogs with their owners. The project currently provides only a small Go HTTP API with a health endpoint; reporting, searching, matching, storage, and WhatsApp integration have not been implemented.

## Prerequisites

- [Nix](https://nixos.org/download/)
- [devenv](https://devenv.sh/getting-started/)

Go is supplied by devenv, so a separate Go installation is not required.

## Development

Enter the canonical development environment from the project root:

```bash
devenv shell
```

Entering the environment installs the repository's pre-commit hook. Before each
commit, it formats staged Go files with `gofmt` and runs `go vet` for their
packages. If formatting changes a file, stage the updated file and commit again.

Run the API:

```bash
go run ./cmd/api
```

The same command is available as a task:

```bash
devenv tasks run app:run
```

Alternatively, start the configured development process:

```bash
devenv up
```

The API listens on port `8080`. In another terminal, call its health endpoint:

```bash
curl -i http://localhost:8080/health
```

The response is `200 OK`, has an `application/json` content type, and contains:

```json
{"status":"ok"}
```

## Checks

Run individual commands inside `devenv shell`:

```bash
go fmt ./...
go test ./...
go vet ./...
```

The equivalent devenv tasks are:

```bash
devenv tasks run app:format
devenv tasks run app:test
devenv tasks run app:vet
```

Run the complete, non-mutating verification suite before finishing a change:

```bash
devenv tasks run app:check
```

This checks formatting and runs the tests and vet analysis.

## Project structure

```text
cmd/api/           Application entry point
internal/httpapi/  HTTP routing and handlers
docs/              Current architecture documentation
docs/adr/          Architecture decision records
AGENTS.md          Operating guide for coding agents
devenv.nix         Development tasks and process configuration
devenv.yaml        Development environment inputs
```
