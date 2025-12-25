import requests
from requests.exceptions import Timeout, ConnectionError
import customtkinter as ctk

# константи
BACKEND_URL         = "http://localhost:8080"
REQUEST_TIMEOUT     = 3  # seconds


# будь-який API запит до бекенду
def api_request(endpoint, method="GET", data=None):
    url = f"{BACKEND_URL}{endpoint}"

    try:
        if method == "GET":
            response = requests.get(url, timeout=REQUEST_TIMEOUT)
        elif method == "POST":
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




# головний клас програми
class App(ctk.CTk):
    def __init__(self):
        super().__init__()

        ctk.set_appearance_mode("dark")
        ctk.set_default_color_theme("dark-blue")

        # змінні для налаштування програми
        self.cleaners_data = []

        self.check_var = ctk.StringVar(value="off")


        # налаштування вікна
        self.title("CLEANER APP")
        self.geometry("1200x700")

        self.minsize(800,600)
        self.resizable(False, False)


        # налаштування додаткових віджетів
        self.checkbox = None
        self.checkbox2 = None



    def draw_cleaners(self):
        # self.checkbox = ctk.CTkCheckBox(self, text="Checkbox", onvalue="on", offvalue="off", variable=self.check_var)
        # self.checkbox.pack(pady=20, padx=20)

        if not self.cleaners_data:
            print("No cleaners data to display.")
            return

        frame = ctk.CTkFrame(self)
        frame.pack(pady=20, padx=20, fill="both", expand=True)

        for cleaner in self.cleaners_data:
            cleaner_name = cleaner["name"]
            label = ctk.CTkLabel(frame, text=cleaner_name, font=ctk.CTkFont(size=20, weight="bold"))
            label.pack(pady=10)

            for option in cleaner["options"]:
                option_name = option["label"]
                option_label = ctk.CTkLabel(frame, text=option_name, font=ctk.CTkFont(size=16))
                option_label.pack(pady=5)

                for action in option["actions"]:
                    action_name = action["command"]
                    action_button = ctk.CTkButton(frame, text=action_name)
                    action_button.pack(pady=2)







    def set_cleaners(self):
        self.cleaners_data = api_request("/api/cleaners", method="GET")
        for cleaner in self.cleaners_data:
            print(f"Cleaner: {cleaner['name']}, description: {cleaner['description']}")
            for option in cleaner['options']:
                print(f" - Option: {option['actions']}")
                for action in option['actions']:
                    print(f"   - Action: {action['command']}")
