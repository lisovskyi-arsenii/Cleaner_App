# Settings component

import customtkinter as ctk
from src.config.settings import (SETTINGS_WINDOW_TITLE,
                                 SETTINGS_WINDOW_WIDTH,
                                 SETTINGS_WINDOW_HEIGHT,
                                 SETTINGS_APPEARANCE_MODE_OPTION_MENU_WIDTH)
from src.enums.Enums import Appearance


class SettingsWindow(ctk.CTkToplevel):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        # window setup
        self.title(SETTINGS_WINDOW_TITLE)
        self.geometry(f"{SETTINGS_WINDOW_WIDTH}x{SETTINGS_WINDOW_HEIGHT}")
        self.attributes("-topmost", True)
        self.resizable(False, False)

        # main container for whole frame
        self.main_container = ctk.CTkFrame(self)
        self.main_container.pack(fill="both", expand=True, padx=20, pady=20)

        # header section
        self.header_label = ctk.CTkLabel(
            self.main_container,
            text="Settings",
            font=ctk.CTkFont(size=22, weight="bold"),
            anchor="w",
        )
        self.header_label.pack(fill="x", pady=(0, 20))

        # theme section
        self.appearance_frame = ctk.CTkFrame(self.main_container)
        self.appearance_frame.pack(fill="x", pady=10)

        self.appearance_label = ctk.CTkLabel(
            self.appearance_frame,
            text="Appearance",
            font=ctk.CTkFont(size=14),
            anchor="w",
        )
        self.appearance_label.pack(fill="x", padx=(0, 10))

        self.appearance_mode_option_menu = ctk.CTkOptionMenu(
            self,
            values=[item.value for item in Appearance],
            command=self.change_appearance_mode,
            width=SETTINGS_APPEARANCE_MODE_OPTION_MENU_WIDTH,
        )
        self.appearance_mode_option_menu.pack(side="right")

        self.appearance_mode_option_menu.set(ctk.get_appearance_mode())


    def change_appearance_mode(self, appearance_mode: str) -> None:
        ctk.set_appearance_mode(appearance_mode)
