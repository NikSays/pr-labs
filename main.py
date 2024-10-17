from bs4 import BeautifulSoup
from functools import reduce
from datetime import datetime, timezone
import ssl
import json
import socket


def request(url):
    host, path = url.replace("https://", "").split("/", 1)
    port = 443

    context = ssl.create_default_context()

    sock = socket.create_connection((host, port))
    wrapped_socket = context.wrap_socket(sock, server_hostname=host)

    http_request = f"""GET /{path} HTTP/1.0\r
Host: {host}\r
Connection: close\r\n\r\n"""

    wrapped_socket.sendall(http_request.encode())

    response = b""
    while True:
        part = wrapped_socket.recv(4096)
        if not part:
            break
        response += part

    wrapped_socket.close()

    response_str = response.decode()
    headers, html_content = response_str.split("\r\n\r\n", 1)
    return html_content


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
        "title": title,
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

    product_info = []

    for prod in prods:
        try:
            product_info.append(process_product(prod))
        except Exception as e:
            print(f"Error: {e}. Skipping...")

    # rate = json.loads(
        # request("https://open.er-api.com/v6/latest/MDL"))["rates"]["EUR"]
    rate = 1/20
    product_info = filter(
        lambda p: p["price"] is not None and p["price"] < 2000,
        product_info)
    product_info = list(map(lambda p: add_eur(p, rate), product_info))
    price_sum = reduce(lambda col, cur: col + cur["price"], product_info, 0)
    timestamp = datetime.now(timezone.utc)
    for i in product_info:
        print(20*"=", "\n")

        print("Name: ", i["title"])
        print("Price (MDL): ", i["price"])
        print("Price (EUR): ", i["price_eur"])
        print("Warranty: ", i["warranty"])

        print("\n")
    print(f"Sum: {price_sum} MDL")
    print(f"Time: {timestamp}")


if __name__ == "__main__":
    main()
