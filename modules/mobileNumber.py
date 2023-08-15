import requests
import logging
import json


def mobile_number_notification():
    try:
        is_found = False
        tele2_url = "https://nnov.tele2.ru/api/shop/products/numbers/bundles/1/groups?query=9524596234&exclude&siteId=siteNNOV"
        headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0'}
        res_json = json.loads(requests.get(tele2_url, timeout=5, headers=headers).text)
        for group in res_json["data"]:
            for group_item in group['bundleGroups']:
                if len(group_item['bundles']) != 0:
                    is_found = True

        if is_found:
            return "Number was found in Tele2: https://nnov.tele2.ru/shop/number?pageParams=type%3Dchoose%26price%3D0%26search_num%3D9524596234"
        return ""
    except Exception as e:
        logging.exception("Couldn't get mobile number in Tele2", e)
        return ""
