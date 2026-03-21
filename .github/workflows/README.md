# Workflows

## lint-markdown.yml

Lints all Markdown files (`**/*.md`) on pull requests.

### Manual execution

#### Install

```sh
brew install cspell
brew install markdownlint-cli2
```

#### Usage

```sh
cspell "**/*.md" --config cspell.json
markdownlint-cli2 "**/*.md"
```
