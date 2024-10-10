import requests


def main():
    url = "https://nanoteh.md/en/computers-servers-parts/monitors-displays-screens/"
    res = requests.get(url)

    if res.status_code != 200:
        raise Exception(f"The request failed with code {res.status_code}")

    print(res.text)


if __name__ == "__main__":
    main()
