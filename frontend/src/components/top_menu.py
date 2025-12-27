# Top menu component

import customtkinter as ctk

from src.config.settings import *

class TopMenu(ctk.CTkFrame):
    def __init__(self, parent, on_analyze=None, on_clean=None,
                 on_clear_options=None, on_settings=None):
        super().__init__(
            parent,
            height=TOP_MENU_HEIGHT,
            corner_radius=TOP_MENU_CORNER_RADIUS
        )

        self.grid_propagate(False)

        # button `clean`
        self.btn_clean = ctk.CTkButton(
            self,
            text="clean",
            width=TOP_MENU_BUTTON_CLEAN_WIDTH,
            height=TOP_MENU_BUTTON_CLEAN_HEIGHT,
            command=on_clean
        )
        self.btn_clean.pack(side="left", padx=10, pady=10)


        # TODO
        # button `analyze`
        self.btn_analyze = ctk.CTkButton(
            self,
            text="analyze",
            width=TOP_MENU_BUTTON_ANALYZE_WIDTH,
            height=TOP_MENU_BUTTON_ANALYZE_HEIGHT,
            command=on_analyze
        )
        self.btn_analyze.pack(side="left", padx=10, pady=10)

        # button `clear_options`
        self.btn_clear_option = ctk.CTkButton(
            self,
            text="clear options",
            width=TOP_MENU_BUTTON_CLEAR_OPTIONS_WIDTH,
            height=TOP_MENU_BUTTON_CLEAR_OPTIONS_HEIGHT,
            command=on_clear_options
        )
        self.btn_clear_option.pack(side="left", padx=10, pady=10)

        # button `settings`
        self.btn_settings = ctk.CTkButton(
            self,
            text="settings",
            width=TOP_MENU_BUTTON_SETTINGS_WIDTH,
            height=TOP_MENU_BUTTON_SETTINGS_HEIGHT,
            command=on_settings
        )
        self.btn_settings.pack(side="right", padx=10, pady=10)
