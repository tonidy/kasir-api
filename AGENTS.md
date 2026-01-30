# Agent Operational Guide - Kasir API

This document provides a focused overview for AI agents working on the Kasir API project. It is derived from [GEMINI.md](GEMINI.md).

## 1. Project Identity
*   **Name:** Kasir API
*   **Type:** RESTful API (Go/Golang)
*   **Core Logic:** In-memory storage, Standard Library HTTP (no external web frameworks like Gin/Echo).

## 2. Critical Commands
Agents should prioritize these commands for verification and maintenance.

*   **Build Verification:** `make build`
    *   *Goal:* Ensure code compiles without errors.
*   **Test Execution:** `make test`
    *   *Goal:* Verify functionality. Run this after any code modification.
*   **Code Quality:** `make audit`
    *   *Goal:* Format, Vet, and Staticcheck. Run this before finalizing any task.
*   **Documentation:** `make docs`
    *   *Goal:* Regenerate Swagger docs. Run this if you modify API handlers or models.

## 3. Code Modification Rules
*   **Style:** Strict adherence to `go fmt`.
*   **Handlers:** located in `handler.go`. Ensure Swagger comments (`// @Summary`, etc.) are updated if logic changes.
*   **Models:** located in `model.go`. If struct fields change, update JSON tags.
*   **Tests:** located in `*_test.go` files. New features *must* include corresponding tests.

## 4. Environment
*   **Language:** Go 1.25.6
*   **Config:** Environment variables `SERVER_HOST` and `SERVER_PORT`.
*   **Local URL:** `http://localhost:8300` (default)

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
