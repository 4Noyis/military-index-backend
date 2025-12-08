# Codecov Setup Guide - Step by Step

This guide will walk you through setting up Codecov for the Military Index Backend project.

---

## 📋 Prerequisites

- ✅ GitHub repository created and accessible
- ✅ CI pipeline files ready (`.github/workflows/ci.yml`)
- ✅ Tests written for services
- ⏳ Codecov account (we'll create this)

---

## 🎯 What is Codecov?

**Codecov** is a service that:
- **Visualizes test coverage** - Shows which lines of code are tested
- **Tracks coverage over time** - See if coverage improves or degrades
- **Comments on PRs** - Automatically tells you if a PR decreases coverage
- **Creates badges** - Show coverage percentage on README
- **Free for open source!** 🎉

---

## Step 1: Sign Up for Codecov

### 1.1 Go to Codecov Website
Open your browser and navigate to:
```
https://about.codecov.io/
```

### 1.2 Click "Sign Up" or "Log In"
- In the top right corner, click the **Sign Up** button
- Or go directly to: https://app.codecov.io/signup

### 1.3 Sign Up with GitHub
- Click **"Sign up with GitHub"** button
- This uses your existing GitHub account (no new password needed!)

### 1.4 Authorize Codecov
GitHub will ask you to authorize Codecov:
- ✅ Review the permissions (Codecov needs to read your repos and post comments)
- ✅ Click **"Authorize codecov"**

**✅ You now have a Codecov account!**

---

## Step 2: Add Your Repository to Codecov

### 2.1 On Codecov Dashboard
After signing up, you'll see the Codecov dashboard:
```
https://app.codecov.io/gh/[your-username]
```

### 2.2 Click "Add New Repository"
- Look for a **"+ Add Repository"** or **"Not yet set up"** button
- Or click on your profile → **"Not yet set up"** tab

### 2.3 Find Your Repository
You should see a list of your GitHub repositories:
```
□ 4noyis/military-index-backend
```

### 2.4 Enable the Repository
- ✅ Click the toggle or **"Setup repo"** next to `military-index-backend`
- Codecov will prepare the repository

### 2.5 Copy the Upload Token
After enabling, Codecov will show you a screen with:

```
Repository Upload Token
████████-████-████-████-████████████
```

**IMPORTANT: Copy this token!** You'll need it in the next step.

If you miss this screen:
1. Go to `https://app.codecov.io/gh/4noyis/military-index-backend`
2. Click **Settings** → **General**
3. Find **"Repository Upload Token"**
4. Click **"Copy"**

---

## Step 3: Add Token to GitHub Secrets

### 3.1 Go to Your GitHub Repository
Open your repository in GitHub:
```
https://github.com/4noyis/military-index-backend
```

### 3.2 Navigate to Secrets
1. Click **"Settings"** (top tab)
2. In the left sidebar, click **"Secrets and variables"**
3. Click **"Actions"**

### 3.3 Create New Secret
1. Click **"New repository secret"** (green button)
2. Fill in the form:
   - **Name**: `CODECOV_TOKEN`
   - **Value**: [Paste the token you copied from Codecov]
3. Click **"Add secret"**

**✅ Secret added!** GitHub Actions can now upload coverage to Codecov.

---

## Step 4: Push CI Files to GitHub

### 4.1 Check Current Git Status
```bash
git status
```

You should see:
- `.github/workflows/ci.yml` (new)
- `.golangci.yml` (new)
- `codecov.yml` (new)
- `docs/CI.md` (new)
- `docs/CODECOV_SETUP.md` (new)
- `docs/ROADMAP.md` (modified)
- `README.md` (modified)

### 4.2 Create a New Branch
```bash
# Make sure you're on the latest code
git checkout stage
git pull origin stage

# Create new branch for CI
git checkout -b ci/setup-github-actions
```

**Why a new branch?** So we can test CI on a pull request before merging.

### 4.3 Stage the CI Files
```bash
# Add the new CI files
git add .github/
git add .golangci.yml
git add codecov.yml
git add docs/CI.md
git add docs/CODECOV_SETUP.md
git add docs/ROADMAP.md
git add README.md
```

### 4.4 Commit
```bash
git commit -m "feat: Add CI/CD pipeline with GitHub Actions and Codecov

- Add GitHub Actions workflow for CI
- Configure golangci-lint with project standards
- Add Codecov integration for test coverage
- Add comprehensive CI documentation
- Update README with CI badges
- Mark Milestone 6.1 complete in ROADMAP

The CI pipeline includes:
- Code linting (golangci-lint)
- Unit tests with race detection
- Coverage reporting to Codecov
- Service builds verification
- Docker image builds

Coverage targets:
- Project: 70% minimum
- New code in PRs: 80% minimum"
```

### 4.5 Push to GitHub
```bash
git push -u origin ci/setup-github-actions
```

**✅ Files are now on GitHub!**

---

## Step 5: Watch CI Run (The Exciting Part!)

### 5.1 Go to GitHub Actions Tab
1. Open your repository on GitHub
2. Click the **"Actions"** tab at the top
3. You should see your workflow running:
   ```
   🟡 CI Pipeline
   ci/setup-github-actions
   Running... (feat: Add CI/CD pipeline...)
   ```

### 5.2 Watch the Jobs
Click on the workflow run to see details:

You'll see 4 jobs running in parallel:
- **Lint Code** - Checking code quality
- **Run Tests** - Running unit tests (2 matrix jobs)
- **Build Services** - Compiling code (3 matrix jobs)
- **Docker Build Test** - Building containers (3 matrix jobs)

**This takes about 3-5 minutes.**

### 5.3 Check for Errors

**If all jobs show ✅ green checkmarks:**
- Congratulations! Your CI is working!
- Skip to Step 6

**If any job shows ❌ red X:**
Click on the failed job to see the error. Common issues:

**Linter errors:**
```bash
# Run locally to fix
cd services/country-service
golangci-lint run
# Fix the issues it reports
```

**Test failures:**
```bash
# Run locally to debug
cd services/country-service
go test -v ./...
```

**Build errors:**
```bash
# Run locally
cd services/country-service
go build -v ./cmd/main.go
```

---

## Step 6: Check Codecov Report

### 6.1 Wait for Upload
After tests complete, GitHub Actions uploads coverage to Codecov.
This takes an extra 10-30 seconds.

### 6.2 View Coverage on Codecov
1. Go to `https://app.codecov.io/gh/4noyis/military-index-backend`
2. You should see your coverage report!

**What you'll see:**
- **Overall coverage percentage** (e.g., 86.2%)
- **Sunburst chart** - Visual breakdown of coverage
- **File tree** - Coverage per file
- **Trends graph** - Coverage over time

### 6.3 Explore the Report
Click around to see:
- **Files tab** - See which files have low coverage
- **Commits tab** - Coverage per commit
- **Branches tab** - Coverage per branch

**Cool features:**
- Red/yellow/green color coding
- Line-by-line coverage (click a file)
- Coverage diff (compare branches)

---

## Step 7: Create Pull Request

### 7.1 Create PR on GitHub
1. Go to your repository on GitHub
2. You should see a yellow banner:
   ```
   ci/setup-github-actions had recent pushes
   [Compare & pull request]
   ```
3. Click **"Compare & pull request"**

### 7.2 Fill in PR Details
- **Base branch**: `stage` (or `main` if you don't have stage)
- **Compare branch**: `ci/setup-github-actions`
- **Title**: `Add CI/CD pipeline with GitHub Actions`
- **Description**:
  ```markdown
  ## Changes
  - ✅ Add GitHub Actions CI pipeline
  - ✅ Configure golangci-lint
  - ✅ Integrate Codecov for coverage tracking
  - ✅ Add comprehensive documentation

  ## CI Pipeline
  - Lints code with golangci-lint
  - Runs all unit tests with race detection
  - Builds all services
  - Builds Docker images
  - Reports coverage to Codecov

  ## Coverage
  - Country Service: 87.8%
  - Technology Service: 84.6%
  - Overall: >85%

  Closes #[issue-number] (if you have one)
  ```

### 7.3 Create the PR
Click **"Create pull request"**

### 7.4 Watch CI Run Again
The CI will run again for the PR. This time you'll also see:

**Status checks at the bottom:**
```
✅ Lint Code
✅ Run Tests (country-service)
✅ Run Tests (technology-service)
✅ Build Services (country-service)
✅ Build Services (technology-service)
✅ Build Services (api-gateway)
✅ Docker Build Test (country-service)
✅ Docker Build Test (technology-service)
✅ Docker Build Test (api-gateway)
✅ CI Success
✅ codecov/patch - All good!
✅ codecov/project - All good!
```

**Codecov comment:**
Within 1-2 minutes, Codecov bot will comment on your PR:
```
📊 Coverage Report

Coverage: 86.2% (target: 70%)
Files changed: 12

+/- Coverage Δ
services/country-service  87.8%  +0.0%
services/technology-service  84.6%  +0.0%
```

**This is amazing!** Everyone can see coverage at a glance.

---

## Step 8: Merge the PR

### 8.1 Review the Changes
- Check that all CI checks passed ✅
- Review the code if you want
- Check Codecov report

### 8.2 Merge
1. Click **"Merge pull request"** (green button)
2. Confirm merge
3. Delete the branch (GitHub will offer)

**✅ CI is now on your main branch!**

---

## Step 9: Setup Branch Protection (Recommended)

This prevents anyone (including you) from pushing broken code.

### 9.1 Go to Branch Settings
1. GitHub repository → **Settings**
2. Left sidebar → **Branches**
3. Under "Branch protection rules" → **Add rule**

### 9.2 Configure Protection for `main`
Fill in:

**Branch name pattern:**
```
main
```

**Protection rules:**
- ✅ Require a pull request before merging
  - ✅ Require approvals: 1 (if you have team members)
  - ✅ Dismiss stale pull request approvals when new commits are pushed

- ✅ Require status checks to pass before merging
  - ✅ Require branches to be up to date before merging
  - **Search and select these checks:**
    - `CI Success`
    - `codecov/project`
    - `codecov/patch`

- ✅ Require conversation resolution before merging

- ❌ Require signed commits (optional, advanced)

- ✅ Require linear history (keeps git history clean)

- ✅ Include administrators (even you must follow rules)

Click **"Create"** or **"Save changes"**

### 9.3 Repeat for `stage` Branch
If you use a staging branch, repeat with:
- Branch name pattern: `stage`
- Same settings as main

**✅ Your branches are now protected!**

**What this means:**
- Can't push directly to main
- Must create PR
- CI must pass before merge
- Coverage must meet targets

---

## Step 10: Add Coverage Badge to README

### 10.1 Get Badge Markdown
We already added badges to README! But if you want to update them:

1. Go to Codecov dashboard for your repo
2. Click **Settings** → **Badge**
3. Copy the markdown:
   ```markdown
   [![codecov](https://codecov.io/gh/4noyis/military-index-backend/branch/main/graph/badge.svg?token=YOUR_TOKEN)](https://codecov.io/gh/4noyis/military-index-backend)
   ```

The badge shows real-time coverage!

---

## 🎉 You're Done! What You Now Have

### Automatic Quality Checks
Every push/PR now automatically:
- ✅ Lints code for quality issues
- ✅ Runs all unit tests
- ✅ Checks for race conditions
- ✅ Verifies builds work
- ✅ Verifies Docker images build
- ✅ Reports test coverage
- ✅ Prevents merging broken code

### Visual Coverage Reports
- Codecov dashboard with coverage trends
- PR comments showing coverage changes
- Badges on README
- Line-by-line coverage view

### Professional Workflow
- Branch protection prevents accidents
- CI catches bugs before production
- Team members can see test status
- Coverage goals are enforced

---

## 🔍 Troubleshooting

### "Codecov upload failed"
**Check:**
1. Is `CODECOV_TOKEN` set in GitHub Secrets?
2. Is the token correct? (copy again from Codecov)
3. Try regenerating the token on Codecov

### "Coverage report not showing"
**Check:**
1. Did tests actually run? (check CI logs)
2. Was `coverage.out` generated?
3. Wait a few minutes (Codecov can be slow)

### "Branch protection too strict"
If you're the only developer:
- Uncheck "Require approvals"
- But keep "Require status checks"

### "CI is too slow"
After it works:
- Consider removing some matrix jobs
- Cache is not working? Check `actions/setup-go` logs
- Docker builds slow? Ensure cache is working

---

## 📊 Viewing Coverage Locally

Want to see coverage without Codecov?

```bash
# Run tests with coverage
cd services/country-service
go test -coverprofile=coverage.out ./...

# Open in browser
go tool cover -html=coverage.out
```

This opens an HTML report showing:
- Green lines: Covered
- Red lines: Not covered
- Gray lines: Not executable

---

## 🚀 Next Steps

Now that CI is working:

1. **Keep coverage high**
   - Add tests for new features
   - Check coverage before merging PRs

2. **Add more checks** (optional)
   - Dependency vulnerability scanning
   - Docker image security scanning
   - Integration tests

3. **Build CD Pipeline** (Milestone 6.2)
   - Automatically deploy to staging
   - Automatically deploy to production

---

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Codecov Documentation](https://docs.codecov.com/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)

---

**Congratulations!** 🎉 You now have professional-grade CI/CD!

---

**Questions?** Check `docs/CI.md` or open an issue on GitHub.

**Last Updated**: December 6, 2025
