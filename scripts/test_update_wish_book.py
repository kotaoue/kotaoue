import unittest

from update_wish_book import (
    build_book_html,
    filter_valid_books,
    replace_between_markers,
)


class TestFilterValidBooks(unittest.TestCase):
    def test_returns_only_books_with_all_required_keys(self):
        books = [
            {"url": "https://example.com", "image": "https://example.com/img.jpg", "title": "Book A"},
            {"url": "https://example.com", "title": "No image"},
            {"image": "https://example.com/img.jpg", "title": "No url"},
        ]
        result = filter_valid_books(books)
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]["title"], "Book A")

    def test_returns_empty_list_when_all_books_invalid(self):
        books = [{"title": "Only title"}, {}]
        self.assertEqual(filter_valid_books(books), [])

    def test_returns_all_books_when_all_valid(self):
        books = [
            {"url": "u1", "image": "i1", "title": "t1"},
            {"url": "u2", "image": "i2", "title": "t2"},
        ]
        self.assertEqual(len(filter_valid_books(books)), 2)

    def test_returns_empty_list_for_empty_input(self):
        self.assertEqual(filter_valid_books([]), [])


class TestBuildBookHtml(unittest.TestCase):
    def test_builds_anchor_with_img(self):
        book = {"url": "https://bookmeter.com/books/1", "image": "https://img.example.com/1.jpg", "title": "My Book"}
        result = build_book_html(book)
        self.assertIn('<a href="https://bookmeter.com/books/1">', result)
        self.assertIn('<img src="https://img.example.com/1.jpg" alt="My Book">', result)
        self.assertIn("</a>", result)

    def test_escapes_special_characters_in_title(self):
        book = {"url": "https://example.com", "image": "https://example.com/img.jpg", "title": 'A "quoted" <title>'}
        result = build_book_html(book)
        self.assertIn("A &quot;quoted&quot; &lt;title&gt;", result)

    def test_escapes_special_characters_in_url(self):
        book = {"url": "https://example.com/?a=1&b=2", "image": "https://example.com/img.jpg", "title": "T"}
        result = build_book_html(book)
        self.assertIn("&amp;", result)


class TestReplaceBetweenMarkers(unittest.TestCase):
    START = "<!-- WISH_BOOK_START -->"
    END = "<!-- WISH_BOOK_END -->"

    def _wrap(self, inner: str) -> str:
        return f"before\n{self.START}{inner}{self.END}\nafter"

    def test_replaces_content_between_markers(self):
        content = self._wrap("<old>")
        result = replace_between_markers(content, self.START, self.END, "<new>")
        self.assertEqual(result, self._wrap("<new>"))

    def test_replaces_empty_content(self):
        content = self._wrap("")
        result = replace_between_markers(content, self.START, self.END, "<new>")
        self.assertEqual(result, self._wrap("<new>"))

    def test_raises_when_start_marker_missing(self):
        content = f"no markers here {self.END}"
        with self.assertRaises(ValueError):
            replace_between_markers(content, self.START, self.END, "<new>")

    def test_raises_when_end_marker_missing(self):
        content = f"{self.START} no end marker"
        with self.assertRaises(ValueError):
            replace_between_markers(content, self.START, self.END, "<new>")

    def test_raises_when_both_markers_missing(self):
        with self.assertRaises(ValueError):
            replace_between_markers("no markers", self.START, self.END, "<new>")

    def test_raises_when_end_marker_before_start_marker(self):
        content = f"{self.END} some content {self.START}"
        with self.assertRaises(ValueError):
            replace_between_markers(content, self.START, self.END, "<new>")

    def test_preserves_content_outside_markers(self):
        content = self._wrap("<old>")
        result = replace_between_markers(content, self.START, self.END, "<new>")
        self.assertTrue(result.startswith("before\n"))
        self.assertTrue(result.endswith("\nafter"))


if __name__ == "__main__":
    unittest.main()
