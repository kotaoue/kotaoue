import html
import json
import random
import sys
import urllib.request

WISH_JSON_URL = "https://raw.githubusercontent.com/kotaoue/readingLog/main/wish.json"
README_PATH = "README.md"
START_MARKER = "<!-- WISH_BOOK_START -->"
END_MARKER = "<!-- WISH_BOOK_END -->"


def fetch_books(url: str) -> list:
    with urllib.request.urlopen(url) as res:
        return json.load(res)


def filter_valid_books(books: list) -> list:
    return [b for b in books if all(k in b for k in ("url", "image", "title"))]


def build_book_html(book: dict) -> str:
    return (
        f'<a href="{html.escape(book["url"])}">'
        f'<img src="{html.escape(book["image"])}" alt="{html.escape(book["title"])}">'
        f"</a>"
    )


def replace_between_markers(content: str, start: str, end: str, replacement: str) -> str:
    if start not in content or end not in content:
        raise ValueError(f"Markers not found in content: {start!r}, {end!r}")
    start_idx = content.index(start)
    end_idx = content.index(end)
    if start_idx >= end_idx:
        raise ValueError(f"Start marker must appear before end marker")
    return content[:start_idx + len(start)] + replacement + content[end_idx:]


def main():
    books = fetch_books(WISH_JSON_URL)

    if not books:
        print("wish.json is empty, skipping update.")
        sys.exit(0)

    valid = filter_valid_books(books)
    if not valid:
        print("No valid book entries found in wish.json, skipping update.")
        sys.exit(0)

    book = random.choice(valid)
    book_html = build_book_html(book)

    with open(README_PATH, "r", encoding="utf-8") as f:
        content = f.read()

    try:
        new_content = replace_between_markers(content, START_MARKER, END_MARKER, book_html)
    except ValueError as e:
        print(e)
        sys.exit(1)

    with open(README_PATH, "w", encoding="utf-8") as f:
        f.write(new_content)

    print(f"Updated README.md with: {book_html}")
