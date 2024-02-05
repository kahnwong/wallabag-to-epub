import json
from typing import Any

import xml2epub


def parse_data(filename: str) -> Any:
    with open(filename, "r") as f:
        data = json.load(f)

    # sort from oldest to newest
    data = data[::-1]
    data = data[:140]  # debug

    chunksize = 20
    data = [data[x : x + chunksize] for x in range(0, len(data), chunksize)]

    return data


def create_epub(data: Any):
    for index, chunk in enumerate(data):
        book = xml2epub.Epub(f"chunk {index}")

        for article in chunk:
            title = article["title"]
            content = f"<h1>{title}</h1>" + article["content"]

            print(f"Adding {title}...")

            try:
                book.add_chapter(
                    xml2epub.create_chapter_from_string(
                        html_string=content, title=title
                    )
                )
            except ValueError:
                pass

        ## generate epub file
        book.create_epub("output")


if __name__ == "__main__":
    data = parse_data("/Users/kahnwong/Downloads/Unread articles.json")
    create_epub(data)
