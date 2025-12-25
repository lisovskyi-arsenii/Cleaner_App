# Main class App for all components

import customtkinter as ctk

from src.api import backend
from src.api.backend import get_cleaners
from src.components.left_menu import LeftMenu
from src.components.main_menu import MainMenu
from src.components.top_menu import TopMenu
from src.config.settings import *

# головний клас програми
class App(ctk.CTk):
    def __init__(self):
        super().__init__()

        self.cleaners_data = []

        # base config for window and widgets
        ctk.set_appearance_mode(APPEARANCE_MODE.value)
        ctk.set_default_color_theme(f"{THEME_DIRECTORY}/tokyonight.json")

        # window setup
        self.title(WINDOW_TITLE)
        self.geometry(f"{WINDOW_WIDTH}x{WINDOW_HEIGHT}")
        self.minsize(MIN_WIDTH, MIN_HEIGHT)
        self.resizable(True, True)

        # Configure grid
        self.grid_rowconfigure(1, weight=1)
        self.grid_columnconfigure(1, weight=1)

        # initialize all UI components
        self._initialize_ui()

        # load data
        self._load_data()



    def _initialize_ui(self):
        # top menu
        self.top_menu = TopMenu(
            self,
            on_analyze=self.on_analyze_clicked,
            on_clean=self.on_clean_clicked,
            on_settings=self.on_settings_clicked,
        )
        self.top_menu.grid(row=0, column=0, columnspan=2, sticky="ew")

        # left menu
        self.left_menu = LeftMenu(
            self,
            on_hover_callback=self.on_widget_hover,
        )
        self.left_menu.grid(row=1, column=0, sticky="ns")

        # main menu
        self.main_menu = MainMenu(self)
        self.main_menu.grid(row=1, column=1, sticky="nsew", padx=10, pady=10)



    def _load_data(self):
        cleaners = get_cleaners()

        if cleaners:
            self.left_menu.draw_cleaners(cleaners)
            # TODO - change to logging
            print(f"Loaded cleaners")
        else:
            print("❌ Failed to load cleaners")
            self.main_menu.hover_title_main.configure(
                text="❌ Cannot connect to backend",
                text_color="red"
            )

    def on_widget_hover(self, widget):
        self.main_menu.show_widget_info(widget)


    # TODO
    # after button `analyze` is pressed, this function will be invoked
    def on_analyze_clicked(self):
        selected_data = self.left_menu.get_selected()
        print(f"selected_data = {selected_data}")
        backend.analyze_cleaners(selected_data)
        # pass

    # TODO
    # after button `clean` is pressed, this function will be invoked
    def on_clean_clicked(self):
        # взяти вибрані чекбокси
        # та далі відправити на очищування на бекенд
        selected_data = self.left_menu.get_selected()
        print(f"selected_data = {selected_data}")

        # pass

    # TODO
    # after button `settings` is pressed, this function will be invoked
    def on_settings_clicked(self):
        pass
