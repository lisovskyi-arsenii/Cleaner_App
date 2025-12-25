# file for communication with backend
from typing import Any

# будь-який API запит до бекенду
import requests
from requests import Timeout

from src.config.settings import BACKEND_URL, REQUEST_TIMEOUT
from src.enums.Enums import APIMethods


def api_request(endpoint: str, method: APIMethods, data=None) -> Any:
    url = f"{BACKEND_URL}{endpoint}"

    try:
        if method == APIMethods.GET:
            response = requests.get(url, timeout=REQUEST_TIMEOUT)
        elif method == APIMethods.POST:
            response = requests.post(url, json=data, timeout=REQUEST_TIMEOUT)
        else:
            raise ValueError("Unsupported HTTP method")

        response.raise_for_status()
        return response.json()

    except Timeout:
        print("Backend не відповідає (timeout)")
        return None

    except ConnectionError:
        print("Не можу підключитися до backend")
        return None

    except Exception as e:
        print(f"Інша помилка: {e}")
        return None

# get list of cleaners
def get_cleaners() -> Any:
    return api_request("/api/cleaners", method=APIMethods.GET)


# analyze cleaners
# def analyze_cleaners() -> Any:
#     return api_request("/api/analyzer", method=APIMethods.GET)


# send request to backend for cleaning all selected files
def clean_files(selected_options) -> Any:
    return api_request("/api/clean", method=APIMethods.POST, data=selected_options)
