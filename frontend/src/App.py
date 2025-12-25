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

        # base config for window and widgets
        ctk.set_appearance_mode("dark")
        ctk.set_default_color_theme("dark-blue")

        # змінні для налаштування програми
        self.cleaners_data = []

        self.check_var = ctk.StringVar(value="off")


        # налаштування вікна
        self.title("CLEANER APP")
        self.geometry("1200x700")

        self.minsize(800,600)
        self.resizable(True, True)


        # налаштування додаткових віджетів
        self.checkboxes = []

        # Configure grid weights
        self.grid_rowconfigure(1, weight=1)  # main content розтягується
        self.grid_columnconfigure(1, weight=1)  # main content розтягується

        # 1. Top menu (row=0, займає 2 колонки)
        self.top_menu = ctk.CTkFrame(
            self,
            height=60,
            corner_radius=0,
            fg_color=("#4A90E2", "#2E5C8A")
        )
        self.top_menu.grid(row=0, column=0, columnspan=2, sticky="ew")
        self.top_menu.grid_propagate(False)

        # 2. Left menu (row=1, column=0)
        self.left_menu = ctk.CTkScrollableFrame(
            self,
            width=250,
            corner_radius=0,
            fg_color=("#DBDBDB", "#2B2B2B")
        )
        self.left_menu.grid(row=1, column=0, sticky="ns")
        self.left_menu.grid_propagate()

        # 3. Main content (row=1, column=1)
        self.main_content_menu = ctk.CTkFrame(
            self,
            corner_radius=0
        )
        self.main_content_menu.grid(row=1, column=1, sticky="nsew", padx=10, pady=10)


    def draw_cleaners(self):
        if not self.cleaners_data:
            print("No cleaners data to display.")
            return

        for cleaner in self.cleaners_data:
            cleaner_label = ctk.CTkLabel(
                self.left_menu,
                text=cleaner["name"],
                font=ctk.CTkFont(size=20, weight="bold"),
                anchor="w"
            )
            cleaner_label.pack(pady=(15,5), padx=10, fill="x")

            cleaner_description = ctk.CTkLabel(
                self.left_menu,
                text=cleaner["description"],
                font=ctk.CTkFont(size=18, weight="bold"),
                text_color="gray",
                anchor="w"
            )
            cleaner_description.pack(pady=(0,10), padx=10, fill="x")

            # options checkboxes
            for option in cleaner["options"]:
                var = ctk.BooleanVar()

                option_checkbox = ctk.CTkCheckBox(
                    self.left_menu,
                    text=option["label"],
                    variable=var,
                    font=ctk.CTkFont(size=14)
                )
                option_checkbox.pack(pady=3, padx=20, anchor="w")

                self.checkboxes.append({
                    'variable': var,
                    'cleaner_id': cleaner['id'],
                    'options_id': option['id'],
                })

            # separator between cleaners
            separator = ctk.CTkFrame(self.left_menu, height=2, fg_color="gray30")
            separator.pack(fill="x", padx=10, pady=10)



    # fetch data from backend
    def set_cleaners(self):
        self.cleaners_data = api_request("/api/cleaners", method="GET")
        for cleaner in self.cleaners_data:
            print(f"Cleaner: {cleaner['name']}, description: {cleaner['description']}")
            for option in cleaner['options']:
                print(f" - Option: {option['actions']}")
                for action in option['actions']:
                    print(f"   - Action: {action['command']}")

    # choose selected checkboxes from list
    def get_selected(self):
        selected = []
        for item in self.checkboxes:
            if item['variable'].get():
                selected.append({
                    'cleaner_id': item['cleaner_id'],
                    'option_id': item['option_id'],
                })

        return selected

    # after button `Clean` press, do this function
    def clean_clicked(self):
        selected = self.get_selected()
        print(f"Selected options: {selected}")
        # далі відправка на очищування на бекенд