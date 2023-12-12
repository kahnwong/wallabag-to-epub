import json

import xml2epub

# read data
filepath = "/Users/kahnwong/Downloads/Unread articles.json"
with open(filepath, "r") as f:
    data = json.load(f)

# sort from oldest to newest
data = data[::-1]
data = data[:50]  # debug

# main
book = xml2epub.Epub("My New E-book Name")

for article in data:
    title = article["title"]
    content = article["content"]

    print(f"Adding {title}...")
    book.add_chapter(
        xml2epub.create_chapter_from_string(html_string=content, title=title)
    )

## generate epub file
book.create_epub("output")
