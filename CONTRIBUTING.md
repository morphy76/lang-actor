# Contributing to Lang-Actor

Thank you for your interest in contributing to Lang-Actor! This document provides guidelines and instructions to help you contribute effectively to this project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please read it to understand what behavior will and will not be tolerated.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report. Following these guidelines helps maintainers understand your report, reproduce the behavior, and find related reports.

- **Use the GitHub issue tracker**: Submit bugs as GitHub issues.
- **Use a clear and descriptive title** for the issue to identify the problem.
- **Describe the exact steps to reproduce the problem** in as much detail as possible.
- **Provide specific examples** to demonstrate the steps.
- **Describe the behavior you observed** after following the steps and explain which behavior you expected to see instead and why.
- **Include screenshots and animated GIFs** if possible.
- **If you're reporting a crash**, include a crash report with a stack trace.
- **Include details about your environment**: OS, Go version, etc.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion, including new features or improvements to existing functionality.

- **Use the GitHub issue tracker**: Submit enhancement suggestions as GitHub issues.
- **Use a clear and descriptive title** for the issue.
- **Provide a step-by-step description of the suggested enhancement** in as much detail as possible.
- **Provide specific examples to demonstrate the steps** or mockups to illustrate the suggestion.
- **Explain why this enhancement would be useful** to most Lang-Actor users.

### Pull Requests

- **Fill in the required pull request template**.
- **Do not include issue numbers in the PR title**.
- **Follow Go style conventions**: Run `gofmt`.
- **Include appropriate test cases**.
- **Document new code based on the Documentation Styleguide**.
- **End all files with a newline**.

## Development Setup

### Prerequisites

- Go 1.18 or higher
- Git

### Setup Steps

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```
   git clone https://github.com/your-username/lang-actor.git
   cd lang-actor
   ```
3. Add the original repository as an upstream remote to keep your fork in sync:
   ```
   git remote add upstream https://github.com/original-owner/lang-actor.git
   ```
4. Install dependencies:
   ```
   go mod download
   ```

## Testing

- Run tests with `make test` or `go test ./...`
- Ensure that your changes pass all tests
- Add new tests for new functionality

## Styleguides

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

### Go Styleguide

- Follow standard Go conventions and `gofmt`
- Use meaningful variable and function names
- Include comments for exported functions and packages
- Keep functions focused and small

### Documentation Styleguide

- Use [Markdown](https://guides.github.com/features/mastering-markdown/) for documentation.
- Keep API documentation with the code.
- For more complex topics, add documentation to the `design/` directory.

## Additional Notes

### Issue and Pull Request Labels

This section lists the labels we use to help us track and manage issues and pull requests.

* `bug` - Issues that are bugs
* `documentation` - Issues or PRs related to documentation
* `enhancement` - Issues for enhancements or new features
* `good first issue` - Good for newcomers
* `help wanted` - Extra attention is needed
* `question` - Further information is requested

Thank you for contributing to Lang-Actor!

## Disclaimer

This document has been generated with Full Vibes (Github Copilot using Claude 3.7 Sonnet).

Prompt:

```text
create a github standard contribuition guide
```
