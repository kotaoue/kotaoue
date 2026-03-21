# Workflows

## lint-markdown.yml

Lints all Markdown files (`**/*.md`) on pull requests.

### Triggers

- Pull request events that include changes to any `*.md` file.

### Steps

1. **Checkout** – Checks out the repository using `actions/checkout@v4`.
2. **Run cSpell** – Spell-checks Markdown files using [`streetsidesoftware/cspell-action@v6`](https://github.com/streetsidesoftware/cspell-action).
   - Configuration is loaded from [`cspell.json`](../../cspell.json) at the repository root.
   - Custom words accepted by cSpell are maintained in that file.
3. **Run markdownlint** – Lints Markdown style using [`DavidAnson/markdownlint-cli2-action@v17`](https://github.com/DavidAnson/markdownlint-cli2-action).

### Adding custom words to cSpell

Edit `cspell.json` at the repository root and add the word to the `words` array:

```json
{
    "words": [
        "yourNewWord"
    ]
}
```
