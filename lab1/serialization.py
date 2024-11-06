def to_json(prod_info):
    print(f"""{{
    "sum": {prod_info['sum']},
    "timestamp": "{prod_info['timestamp']}",
    "products": [
        {','.join([f"""
        {{
            "name": "{p['name'].replace('"', '\\"')}",
            "price_mdl": {p['price']},
            "price_eur": {p['price_eur']},
            "warranty": {p['warranty']}
        }}"""
          for p in prod_info['products']
                   ])}
    ]
}}""")


def to_xml(prod_info):
    print(f"""<product_info>
    <sum> {prod_info['sum']} </sum>
    <timestamp> {prod_info['timestamp']} </timestamp>
    <products>
        {'\n'.join([f"""
        <product>
            <name> {p['name']} </name>
            <price_mdl> {p['price']} </price_mdl>
            <price_eur> {p['price_eur']} </price_eur>
            <warranty> {p['warranty']} </warranty>
        </product>"""
          for p in prod_info['products']
                    ])}
    </products>
</product_info>""")


