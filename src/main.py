import json

import xml2epub

# read data
filepath = "/Users/kahnwong/Downloads/Unread articles.json"
with open(filepath, "r") as f:
    data = json.load(f)

# sort from oldest to newest
data = data[::-1]
data = data[:50]  # debug

chunksize = 20
data = [data[x : x + chunksize] for x in range(0, len(data), chunksize)]

# main
for index, chunk in enumerate(data):
    book = xml2epub.Epub(f"chunk {index}")

    for article in chunk:
        title = article["title"]
        content = f"<h1>{title}</h1>" + article["content"]

        print(f"Adding {title}...")
        book.add_chapter(
            xml2epub.create_chapter_from_string(html_string=content, title=title)
        )

    ## generate epub file
    book.create_epub("output")
