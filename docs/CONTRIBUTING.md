# Contributing Guide

## 🌟 Welcome!

Thank you for contributing to Military Index Backend!

## 🔀 Git Workflow

### Branch Strategy
```
main (production)
  ↑
stage (staging)
  ↑
feature/your-feature (development)
```

### Creating a Feature
```bash
# 1. Start from stage
git checkout stage
git pull origin stage

# 2. Create feature branch
git checkout -b feature/add-tech-filtering

# 3. Make changes and commit
git add .
git commit -m "feat: add filtering by year to technology endpoint"

# 4. Push to remote
git push origin feature/add-tech-filtering

# 5. Create Pull Request on GitHub
# Base: stage ← Compare: feature/add-tech-filtering
```

### Commit Message Convention

Use [Conventional Commits](https://www.conventionalcommits.org/):
```
feat: add new feature
fix: bug fix
docs: documentation changes
style: formatting, missing semicolons, etc.
refactor: code refactoring
test: adding tests
chore: maintenance tasks
perf: performance improvements
```

Examples:
```bash
git commit -m "feat: add pagination to countries endpoint"
git commit -m "fix: resolve database connection timeout"
git commit -m "docs: update API documentation for tech endpoint"
git commit -m "refactor: improve error handling in service layer"
```

## 📝 Pull Request Guidelines

### Before Creating PR

- [ ] Code follows Go conventions
- [ ] All tests pass: `make test`
- [ ] No linting errors: `make lint`
- [ ] Documentation updated if needed
- [ ] `.env.example` updated if new vars added
- [ ] Migrations created if database changed

### PR Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How to test these changes

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes
- [ ] Follows code style
```

## 🧪 Testing
```bash
# Run all tests
make test

# Run specific package tests
go test ./services/technology-service/...

# Run with coverage
make test-coverage

# Run integration tests
make test-integration
```

## 📏 Code Style

### Go Style Guide

Follow [Uber's Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

Key points:
- Use `gofmt` for formatting
- Use meaningful variable names
- Add comments for exported functions
- Handle errors properly
- Use context for cancellation

Example:
```go
// GetTechnologyByID retrieves a technology by its ID
func (s *techService) GetTechnologyByID(ctx context.Context, id uint) (*models.Technology, error) {
    if id == 0 {
        return nil, errors.New("invalid technology ID")
    }

    tech, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get technology: %w", err)
    }

    return tech, nil
}
```

## 🗄️ Database Changes

### Creating Migrations
```bash
# Create new migration
make db-migrate-create name=add_tech_specs_column

# This creates:
# migrations/000002_add_tech_specs_column.up.sql
# migrations/000002_add_tech_specs_column.down.sql
```

### Migration Best Practices

1. **Always create both up and down migrations**
2. **Test migrations before PR**
3. **Keep migrations atomic** (one change per migration)
4. **Never modify existing migrations** (create new ones)

Example:
```sql
-- 000002_add_tech_specs_column.up.sql
ALTER TABLE technologies
ADD COLUMN specifications JSONB DEFAULT '{}';

-- 000002_add_tech_specs_column.down.sql
ALTER TABLE technologies
DROP COLUMN specifications;
```

## 🚀 Deployment Process

### To Staging
```bash
# 1. Merge PR to stage
# 2. CI/CD automatically deploys to staging
# 3. Test on https://api-staging.military-index.com
```

### To Production
```bash
# 1. Create PR from stage to main
# 2. Get approval from 2+ reviewers
# 3. Merge to main
# 4. CI/CD automatically deploys to production
# 5. Tag release: git tag v1.0.0 && git push origin v1.0.0
```

## ❓ Questions?

- Create an issue on GitHub
- Ask in team Slack/Discord
- Email: dev@military-index.com
