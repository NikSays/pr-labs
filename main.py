from bs4 import BeautifulSoup
import requests


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

    res = requests.get(link)
    if res.status_code != 200:
        raise Exception("Error getting product information")

    prod_page = BeautifulSoup(res.text, "html.parser")
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


def main():
    url = "https://nanoteh.md/en/computers-servers-parts/monitors-displays-screens/"
    res = requests.get(url)

    if res.status_code != 200:
        raise Exception(f"The request failed with code {res.status_code}")

    s = BeautifulSoup(res.text, "html.parser")
    prods = s.find(class_="products-grid").find_all(class_="item-info")

    product_info = []

    for prod in prods:
        try:
            product_info.append(process_product(prod))
        except Exception as e:
            print(f"Error: {e}. Skipping...")

    for i in product_info:
        print(20*"=", "\n")

        print("Name: ", i["title"])
        print("Price: ", i["price"])
        print("Warranty: ", i["warranty"])

        print("\n")


if __name__ == "__main__":
    main()
