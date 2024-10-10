from bs4 import BeautifulSoup
import requests


def main():
    url = "https://nanoteh.md/en/computers-servers-parts/monitors-displays-screens/"
    res = requests.get(url)

    if res.status_code != 200:
        raise Exception(f"The request failed with code {res.status_code}")

    s = BeautifulSoup(res.text, "html.parser")
    prods = s.find(class_="products-grid").find_all(class_="item-info")
    for prod in prods:
        print("\n", 20 * "=", "\n")

        title = prod.find(class_="item-title").find("a")
        print(f"Title: {title.text}")

        price_span = (prod.find(class_="item-price")
                      .find(class_="special-price"))
        if price_span is None:
            print("Price: None")
            continue

        price = price_span.find(class_="price")
        print(f"Price: {price.text}")


if __name__ == "__main__":
    main()
