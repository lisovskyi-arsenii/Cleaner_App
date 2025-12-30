# Main class App for all components
import threading

import customtkinter as ctk

from src.api import backend
from src.api.backend import get_cleaners
from src.components.settings_window import SettingsWindow
from src.components.left_menu import LeftMenu
from src.components.main_menu import MainMenu
from src.components.top_menu import TopMenu
from src.config.settings import *
from src.functions.wrapper_functions import async_action_clear_checkboxes



# головний клас програми
class App(ctk.CTk):
    def __init__(self):
        super().__init__()

        # data
        # all cleaners which were fetched from backend
        self.cleaners_data = []
        # settings window
        self.settings_window = None


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

    # initialize all ui components
    def _initialize_ui(self):
        # top menu
        self.top_menu = TopMenu(
            self,
            on_analyze=self.on_analyze_clicked,
            on_clean=self.on_clean_clicked,
            on_clear_options=self.on_clear_options_clicked,
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

    # load all data from backend
    def _load_data(self):
        def load():
            try:
                cleaners = get_cleaners()
                if cleaners:
                    self.after(0, lambda: self.left_menu.draw_cleaners(cleaners))
                else:
                    self.after(0, lambda: self._show_connection_error)
            except Exception as e:
                print(f"Error loading data: {e}")
                self.after(0, lambda: self._show_connection_error)

        threading.Thread(target=load, daemon=True).start()

    # show message `connection error` if failed to fetch data
    def _show_connection_error(self):
        if hasattr(self.main_menu, 'hover_title_main'):
            self.main_menu.hover_title_main.configure(
                text="❌ Cannot connect to backend", text_color="red"
            )

    # general function for all backend functions
    def _execute_backend_action(self, func):
        selected_data = self.left_menu.get_selected()
        if selected_data:
            func(selected_data)

    # after button `analyze` is pressed, this function will be invoked
    @async_action_clear_checkboxes
    def on_analyze_clicked(self):
        self._execute_backend_action(backend.analyze_cleaners)


    # TODO - change methods for backend request
    # after button `clean` is pressed, this function will be invoked
    @async_action_clear_checkboxes
    def on_clean_clicked(self):
        self._execute_backend_action(backend.clean_files)


    # after button `clear options` is pressed, this function will be invoked
    def on_clear_options_clicked(self):
        self.left_menu.clear_selected_checkboxes()

    # when user hovers one of the widgets, show some info about it in main component
    def on_widget_hover(self, widget):
        self.main_menu.show_widget_info(widget)

    # TODO
    # after button `settings` is pressed, this function will be invoked
    def on_settings_clicked(self):
        if self.settings_window is None or not self.settings_window.winfo_exists():
            self.settings_window = SettingsWindow(self)
        else:
            self.settings_window.focus()


