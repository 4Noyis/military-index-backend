# Continuous Integration (CI) Guide

## Overview

This project uses **GitHub Actions** for Continuous Integration. Every time you push code or create a pull request, automated checks run to ensure code quality and prevent bugs.

## What Gets Checked?

Our CI pipeline runs 4 parallel jobs:

### 1. **Lint** - Code Quality Checks
**Purpose**: Ensures code follows Go best practices and style guidelines

**What it does**:
- Runs `golangci-lint` on all three services
- Checks for:
  - Code formatting issues
  - Unused code
  - Potential bugs
  - Security vulnerabilities
  - Style violations

**How to run locally**:
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run on a service
cd services/country-service
golangci-lint run

# Run on all services
golangci-lint run ./services/...
```

**Common issues**:
- "File is not gofmt-ed" → Run `go fmt ./...`
- "Unused variable" → Remove or use the variable
- "Error return value not checked" → Add error handling

### 2. **Test** - Unit Tests
**Purpose**: Runs all unit tests and measures code coverage

**What it does**:
- Runs tests for Country Service and Technology Service
- Enables race detection (`-race` flag)
- Generates coverage reports
- Uploads coverage to Codecov

**How to run locally**:
```bash
# Test a single service
cd services/country-service
go test -v -race ./...

# With coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View in browser
```

**Coverage goals**:
- Country Service: >85% ✅
- Technology Service: >85% ✅

### 3. **Build** - Compilation Check
**Purpose**: Verifies all services compile successfully

**What it does**:
- Compiles each service binary
- Saves binaries as artifacts (for debugging)

**How to run locally**:
```bash
# Build a service
cd services/country-service
go build -v -o ../../bin/country-service ./cmd/main.go

# Run the binary
../../bin/country-service
```

### 4. **Docker Build** - Container Verification
**Purpose**: Ensures Docker images build correctly

**What it does**:
- Builds Docker images for all services
- Uses layer caching for speed
- Doesn't push (just verification)

**How to run locally**:
```bash
# Build Docker image
docker build -t military-index/country-service:test services/country-service

# Verify it runs
docker run --rm military-index/country-service:test
```

---

## When Does CI Run?

The pipeline triggers on:
- ✅ Every push to `stage` branch
- ✅ Every push to `main` branch
- ✅ Every pull request to `stage` or `main`

It does NOT run on:
- Feature branches (unless you open a PR)
- Draft pull requests (until marked ready)

---

## Understanding CI Status

### ✅ All Checks Passed
Your code is good! The PR can be merged.

### ❌ Some Checks Failed
Click "Details" to see which job failed:

**Lint failed**:
- Fix code quality issues
- Run `golangci-lint run` locally first

**Test failed**:
- Fix failing tests
- Run `go test ./...` locally first

**Build failed**:
- Fix compilation errors
- Run `go build ./...` locally first

**Docker Build failed**:
- Fix Dockerfile issues
- Test `docker build` locally first

---

## CI Workflow File Breakdown

Our CI is defined in `.github/workflows/ci.yml`:

```yaml
on:
  push:
    branches: [ stage, main ]
  pull_request:
    branches: [ stage, main ]
```
**EXPLANATION**: Trigger conditions. Runs on push to stage/main or PR to those branches.

```yaml
strategy:
  matrix:
    service: [country-service, technology-service]
```
**EXPLANATION**: Matrix strategy runs the same job for multiple services. This creates 2 parallel jobs from 1 definition.

```yaml
uses: actions/checkout@v4
```
**EXPLANATION**: Pre-built action that checks out your code. Almost every workflow needs this.

```yaml
uses: actions/setup-go@v5
with:
  go-version: ${{ env.GO_VERSION }}
  cache: true
```
**EXPLANATION**: Installs Go and caches dependencies for faster builds.

```yaml
go test -v -race -coverprofile=coverage.out ./...
```
**EXPLANATION**:
- `-v`: Verbose output
- `-race`: Detect race conditions
- `-coverprofile`: Generate coverage report
- `./...`: Test all packages recursively

```yaml
needs: [lint, test, build, docker-build]
```
**EXPLANATION**: `ci-success` job only runs if all 4 previous jobs pass. Used for branch protection.

---

## Branch Protection Rules

Once CI is working, configure branch protection on GitHub:

**Settings → Branches → Add rule for `main`**:
- ✅ Require pull request reviews
- ✅ Require status checks to pass
  - Select: `CI Success`
- ✅ Require branches to be up to date
- ✅ Include administrators

This prevents broken code from reaching main.

---

## Code Coverage with Codecov

### Setup
1. Go to [codecov.io](https://codecov.io)
2. Sign in with GitHub
3. Add `military-index-backend` repository
4. Copy the token
5. Add to GitHub Secrets:
   - Go to repo **Settings → Secrets → Actions**
   - Click **New repository secret**
   - Name: `CODECOV_TOKEN`
   - Value: [paste token]

### Benefits
- Visual coverage reports
- Coverage trends over time
- PR comments with coverage changes
- Coverage badges for README

### Viewing Coverage
- **On PR**: Codecov bot comments with coverage changes
- **On main**: Visit `https://codecov.io/gh/your-org/military-index-backend`

---

## Performance Optimization

### Caching
We use caching to speed up CI:

```yaml
cache: true  # In setup-go action
```
Caches Go modules between runs.

```yaml
cache-from: type=gha
cache-to: type=gha,mode=max
```
Caches Docker layers between runs.

**Impact**: First run ~5 minutes, subsequent runs ~2 minutes

### Parallel Jobs
Jobs run in parallel by default:
- Lint (1 min)
- Test (2 min)
- Build (1 min)
- Docker (3 min)

**Total time**: ~3 minutes (not 7 minutes!)

---

## Troubleshooting

### "golangci-lint not found"
The workflow uses the official golangci-lint action, so this shouldn't happen. If it does:
```yaml
- name: Install golangci-lint
  run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### "Tests timeout"
Increase timeout in workflow:
```yaml
run: go test -timeout 10m ./...
```

### "Out of disk space"
Clean up Docker layers:
```yaml
- name: Clean up
  run: docker system prune -f
```

### "Rate limit exceeded"
GitHub has rate limits for actions. Use caching to reduce API calls.

---

## Local Development Workflow

**Before pushing**:
```bash
# 1. Format code
go fmt ./...

# 2. Run linter
golangci-lint run ./...

# 3. Run tests
go test -v ./...

# 4. Build
go build ./...
```

This catches issues before CI runs, saving time.

---

## CI vs CD

**CI (what we just built)**:
- Runs on every push/PR
- Tests and validates code
- Catches bugs early
- Doesn't deploy anything

**CD (next phase)**:
- Runs when CI passes
- Deploys to staging/production
- Requires successful CI first

---

## Next Steps

After CI is working:

1. **Add badges to README**:
   ```markdown
   ![CI](https://github.com/your-org/military-index-backend/workflows/CI%20Pipeline/badge.svg)
   [![codecov](https://codecov.io/gh/your-org/military-index-backend/branch/main/graph/badge.svg)](https://codecov.io/gh/your-org/military-index-backend)
   ```

2. **Setup CD Pipeline** (Milestone 6.2):
   - Deploy to staging on merge to `stage`
   - Deploy to production on merge to `main`

3. **Add more checks**:
   - Dependency vulnerability scanning
   - Docker image scanning
   - Integration tests

---

## Reference

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [golangci-lint Linters](https://golangci-lint.run/usage/linters/)
- [Codecov Documentation](https://docs.codecov.com/)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)

---

**Last Updated**: December 6, 2025
**Maintained By**: DevOps Team
