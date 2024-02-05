import os
from typing import Any
from typing import Dict
from typing import List

import requests
import xml2epub
from dotenv import load_dotenv

load_dotenv()


def get_access_token():
    url = "https://wallabag.karnwong.me/oauth/v2/token"
    data = {
        "grant_type": "password",
        "client_id": os.getenv("WALLABAG_CLIENT_ID"),
        "client_secret": os.getenv("WALLABAG_CLIENT_SECRET"),
        "username": os.getenv("WALLABAG_USERNAME"),
        "password": os.getenv("WALLABAG_PASSWORD"),
    }

    r = requests.post(url, data=data)

    return r.json()["access_token"]


def get_articles(access_token: str) -> List[Dict[str, Any]]:
    headers = {
        "Authorization": f"Bearer {access_token}",
        "Content-Type": "application/json",  # Adjust content type if needed
    }

    url = "https://wallabag.karnwong.me/api/entries?perPage=160&starred=0&archived=0"

    r = requests.get(url=url, headers=headers)

    return r.json()["_embedded"]["items"]


def parse_data(data: List[Dict[str, Any]]) -> Any:
    # sort from oldest to newest
    data = data[::-1]
    # data = data[:140]  # debug

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
    access_token = get_access_token()
    articles = get_articles(access_token)

    data = parse_data(articles)
    create_epub(data)
