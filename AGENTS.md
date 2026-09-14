# Agent operating guide

## Project purpose

This project will eventually support a WhatsApp-based lost-dog matching service that accepts dog information and pictures, catalogs reports, and suggests possible matches. The current implementation intentionally contains only a subset of that product: a minimal Go API with `GET /health`.

## Engineering principles

- Prefer simple, explicit, idiomatic Go.
- Favor the standard library when it is sufficient.
- Do not add a dependency without a clear reason.
- Do not introduce architecture for hypothetical future requirements.
- Keep domain and business logic separate from transport concerns as the application grows.
- Make small, scoped changes and avoid unrelated refactoring.
- Preserve backward compatibility unless the task explicitly permits a breaking change.
- Add or update tests for behavioral changes.
- Do not weaken or remove tests merely to make checks pass.
- Never claim verification succeeded unless the command was actually executed successfully.

## Repository structure

- `cmd/api/`: the executable entry point and process wiring.
- `internal/httpapi/`: HTTP routes and handlers. Keep transport-specific code here.
- `docs/architecture.md`: architecture that exists today, not speculative design.
- `docs/adr/`: accepted or proposed architecture decisions and their tradeoffs.
- `devenv.nix` and `devenv.yaml`: canonical development environment, tasks, and processes.

Add new packages only when current behavior gives them a clear responsibility. As domain behavior appears, keep it independent of HTTP packages.

## Development environment

devenv is the canonical environment. Enter it before using direct Go commands:

```bash
devenv shell
```

Common commands:

| Purpose | Direct command | devenv task |
| --- | --- | --- |
| Run the API | `go run ./cmd/api` | `devenv tasks run app:run` |
| Format Go code | `go fmt ./...` | `devenv tasks run app:format` |
| Run tests | `go test ./...` | `devenv tasks run app:test` |
| Vet code | `go vet ./...` | `devenv tasks run app:vet` |
| Complete verification | — | `devenv tasks run app:check` |

`devenv up` starts the API as a supervised development process.

## Agent workflow

Use this loop:

```text
Understand
    ↓
Plan
    ↓
Implement
    ↓
Test
    ↓
Review diff
    ↓
Fix issues
    ↓
Verify again
    ↓
Report evidence
```

Before editing, inspect the relevant code and repository guidance, understand existing conventions, identify affected behavior, and confirm assumptions from repository evidence where possible.

During implementation, keep changes scoped, follow the established architecture, update tests with behavioral changes, and avoid unrelated refactoring.

Before finishing:

1. Run formatting checks.
2. Run tests.
3. Run vet or other configured static checks.
4. Inspect `git diff`.
5. Check for unrelated changes and generated artifacts.
6. Mention any unresolved uncertainty in the final report.

## Definition of Done

A task is complete only when:

- The requested behavior exists.
- Relevant tests exist and pass.
- Formatting passes.
- Static checks pass.
- Documentation affected by the change is updated.
- The final diff has been reviewed.
- Remaining limitations or uncertainties are explicitly reported.
