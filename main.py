from bs4 import BeautifulSoup
from functools import reduce
from datetime import datetime, timezone
import json

from request import request
from serialization import to_json, to_xml


def process_product(prod):
    title_a = prod.find(class_="item-title").find("a")
    link = title_a["href"]

    title = title_a.text
    if len(title) == 0:
        raise Exception("Product with empty title")

    price_span = (prod.find(class_="item-price")
                  .find(class_="special-price"))

    if price_span is None:
        price = None
    else:
        price_str = price_span.find(class_="price").text.replace(
            " ", "").replace("MDL", "")
        try:
            price = int(price_str)
        except ValueError:
            raise Exception(f'Price is not a number for product "{
                title}"')

    res = request(link)

    prod_page = BeautifulSoup(res, "html.parser")
    warranty_str = prod_page.find(
        class_="product-warranty").find("span").text.replace(" months", "")
    try:
        warranty = int(warranty_str)
    except ValueError:
        raise Exception(f'Warranty duration is not a number for product "{
            title}"')

    return {
        "name": title,
        "price": price,
        "warranty": warranty
    }


def add_eur(prod, rate):
    if prod["price"] is not None:
        prod["price_eur"] = prod["price"] * rate
    else:
        prod["price_eur"] = None
    return prod


def main():
    url = "https://nanoteh.md/en/computers-servers-parts/monitors-displays-screens/"
    res = request(url)

    s = BeautifulSoup(res, "html.parser")
    prods = s.find(class_="products-grid").find_all(class_="item-info")

    product_info = {
        "products": []
    }

    for prod in prods:
        try:
            product_info["products"].append(process_product(prod))
        except Exception as e:
            print(f"Error: {e}. Skipping...")

    # rate = json.loads(
        # request("https://open.er-api.com/v6/latest/MDL"))["rates"]["EUR"]
    rate = 1/20

    product_info["products"] = filter(
        lambda p: p["price"] is not None and p["price"] < 2000,
        product_info["products"])

    product_info["products"] = list(map(
        lambda p: add_eur(p, rate),
        product_info["products"]))

    product_info["sum"] = reduce(
        lambda col, cur: col + cur["price"],
        product_info["products"], 0)

    product_info["timestamp"] = datetime.now(timezone.utc)

    print("\nJSON:")
    to_json(product_info)
    print("\nXML:")
    to_xml(product_info)


if __name__ == "__main__":
    main()
