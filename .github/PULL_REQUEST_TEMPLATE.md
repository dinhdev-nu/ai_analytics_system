# Template for GitHub Pull Request

## Description

<!-- Provide a brief description of the changes in this PR -->

## Type of Change

<!-- Mark the relevant option with an 'x' -->

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring
- [ ] Tests
- [ ] CI/CD changes

## Related Issues

<!-- Link to related issues using #issue_number -->

Closes #
Related to #

## Changes Made

<!-- List the main changes made in this PR -->

- 
- 
- 

## Testing

<!-- Describe the tests you ran and how to reproduce them -->

### Test Environment

- OS: <!-- e.g., Ubuntu 22.04, macOS 13, Windows 11 -->
- Go version: <!-- e.g., 1.21.0 -->
- Python version: <!-- e.g., 3.11.0 -->
- Node version: <!-- e.g., 18.0.0 -->
- Docker version: <!-- e.g., 24.0.0 -->

### Test Steps

1. 
2. 
3. 

### Test Results

<!-- Paste test output or describe results -->

```
# Test output here
```

## Checklist

<!-- Mark completed items with an 'x' -->

### Code Quality

- [ ] My code follows the project's coding standards
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have removed all debug statements and console logs
- [ ] My changes generate no new warnings or errors

### Testing

- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have tested the changes in a development environment
- [ ] I have tested edge cases and error scenarios

### Documentation

- [ ] I have updated the relevant documentation (README, API docs, etc.)
- [ ] I have added/updated code comments where necessary
- [ ] I have updated the CHANGELOG.md file
- [ ] I have updated type definitions if applicable

### Dependencies

- [ ] I have checked that dependencies are up to date
- [ ] I have not introduced any unnecessary dependencies
- [ ] I have updated `go.mod`, `requirements.txt`, or `package.json` if needed
- [ ] I have run dependency vulnerability scans

### Security

- [ ] My code does not introduce any security vulnerabilities
- [ ] I have not exposed any sensitive information (API keys, passwords, etc.)
- [ ] I have validated all user inputs
- [ ] I have followed secure coding practices

### Performance

- [ ] I have considered the performance impact of my changes
- [ ] I have optimized database queries if applicable
- [ ] I have not introduced any N+1 query problems
- [ ] I have tested with realistic data volumes

### Database

- [ ] Database migrations are included (if applicable)
- [ ] Database migrations are reversible
- [ ] I have tested migrations on a copy of production data
- [ ] Database indexes are added for new queries

## Screenshots (if applicable)

<!-- Add screenshots to help explain your changes -->

### Before

<!-- Screenshot before changes -->

### After

<!-- Screenshot after changes -->

## Breaking Changes

<!-- List any breaking changes and migration steps -->

- 

## Deployment Notes

<!-- Special instructions for deployment -->

- 
- 

## Additional Context

<!-- Add any other context about the PR here -->

## Reviewer Notes

<!-- Anything specific you want reviewers to focus on -->

- 
- 

---

**PR Author Checklist:**

- [ ] I have read the [CONTRIBUTING.md](../CONTRIBUTING.md) guidelines
- [ ] I have assigned appropriate labels to this PR
- [ ] I have added reviewers
- [ ] I have linked related issues
- [ ] I am ready for review
