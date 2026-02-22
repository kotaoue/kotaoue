# update_wish_book

## Usage

```bash
PIPENV_PIPFILE=tools/python/Pipfile pipenv --python /opt/homebrew/bin/python
PIPENV_PIPFILE=tools/python/Pipfile pipenv install

# Test
PIPENV_PIPFILE=tools/python/Pipfile pipenv run python -m unittest tools/python/tests/update_wish_book/test_update_wish_book.py

# 実行
PIPENV_PIPFILE=tools/python/Pipfile pipenv run python -m tools.python.update_wish_book
```
