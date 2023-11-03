# -*- coding: utf-8 -*-
import json



with open('products_export.json', 'r') as f:
    data = json.load(f)

print("const data =", json.dumps(data, indent=2, ensure_ascii=False) + ";")
