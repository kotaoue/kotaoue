import random
import unittest

from tools.python.update_wish_book import build_book_html, filter_valid_books, replace_between_markers


class TestWishBookUpdateSimple(unittest.TestCase):
    START = "<!-- WISH_BOOK_START -->"
    END = "<!-- WISH_BOOK_END -->"

    def _render_readme(self, books: list[dict]) -> str:
        valid_books = filter_valid_books(books)
        chosen = random.choice(valid_books)
        book_html = build_book_html(chosen)
        readme = f"before\n{self.START}<old>{self.END}\nafter"
        return replace_between_markers(readme, self.START, self.END, book_html)

    def test_single_book_json_creates_expected_readme_string(self):
        books = [
            {
                "url": "https://bookmeter.com/books/1",
                "image": "https://img.example.com/1.jpg",
                "title": "Book A",
            }
        ]

        updated = self._render_readme(books)

        expected = (
            "before\n"
            "<!-- WISH_BOOK_START -->"
            '<a href="https://bookmeter.com/books/1"><img src="https://img.example.com/1.jpg" alt="Book A"></a>'
            "<!-- WISH_BOOK_END -->\n"
            "after"
        )
        self.assertEqual(updated, expected)

    def test_multiple_books_json_selects_one_of_them_when_repeated(self):
        books = [
            {
                "url": "https://bookmeter.com/books/1",
                "image": "https://img.example.com/1.jpg",
                "title": "Book A",
            },
            {
                "url": "https://bookmeter.com/books/2",
                "image": "https://img.example.com/2.jpg",
                "title": "Book B",
            },
        ]

        candidate_a = (
            "before\n"
            "<!-- WISH_BOOK_START -->"
            '<a href="https://bookmeter.com/books/1"><img src="https://img.example.com/1.jpg" alt="Book A"></a>'
            "<!-- WISH_BOOK_END -->\n"
            "after"
        )
        candidate_b = (
            "before\n"
            "<!-- WISH_BOOK_START -->"
            '<a href="https://bookmeter.com/books/2"><img src="https://img.example.com/2.jpg" alt="Book B"></a>'
            "<!-- WISH_BOOK_END -->\n"
            "after"
        )
        allowed = {candidate_a, candidate_b}

        results = {self._render_readme(books) for _ in range(100)}

        self.assertTrue(results.issubset(allowed))
        self.assertEqual(results, allowed)


if __name__ == "__main__":
    unittest.main()
