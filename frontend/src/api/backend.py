# file for communication with backend
from typing import Any
import json

import requests
from requests import Timeout

from src.config.settings import BACKEND_URL, REQUEST_TIMEOUT
from src.enums.Enums import APIMethods

# any request to backend
def api_request(endpoint: str, method: APIMethods, data=None) -> Any:
    url = f"{BACKEND_URL}{endpoint}"

    try:
        if method == APIMethods.GET:
            response = requests.get(url, timeout=REQUEST_TIMEOUT)
        elif method == APIMethods.POST:
            print(f"DEBUG: Sending POST request to {url}")
            print(f"DEBUG: Data = {data}")
            response = requests.post(url, json=data, timeout=REQUEST_TIMEOUT)
        else:
            raise ValueError("Unsupported HTTP method")

        response.raise_for_status()

        if not response.text or response.text.strip() == "":
            print("Backend is not answering (timeout)")
            return None

        try:
            result = response.json()
            print(f"DEBUG: ✅ JSON parsed successfully: {result}")
            return result
        except json.decoder.JSONDecodeError as e:
            print(f"❌ JSON decode error: {e}")
            print(f"❌ Trying to parse: '{response.text}'")
            return None

    except Timeout:
        print("Backend is not answering (timeout)")
        return None

    except ConnectionError:
        print("Cannot connect to backend")
        return None

    except Exception as e:
        print(f"Error: {e}")
        return None

# certain requests
# get list of cleaners
def get_cleaners() -> Any:
    return api_request("/api/cleaners", method=APIMethods.GET)


# analyze cleaners
def preview_cleaners(selected_options) -> Any:
    return api_request("/api/analyze", method=APIMethods.POST, data=selected_options)


# send request to backend for cleaning all selected files
def clean_files(selected_options) -> Any:
    return api_request("/api/clean", method=APIMethods.POST, data=selected_options)

