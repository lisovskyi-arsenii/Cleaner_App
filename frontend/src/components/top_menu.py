# Top menu component

import customtkinter as ctk
from PIL import Image
from src.config.settings import *

class TopMenu(ctk.CTkFrame):
    def __init__(self, parent, on_preview=None, on_clean=None,
                 on_clear_options=None, on_abort=None, on_settings=None):
        super().__init__(
            parent,
            height=TOP_MENU_HEIGHT,
            corner_radius=TOP_MENU_CORNER_RADIUS
        )

        self.grid_propagate(False)

        # store callbacks
        self.on_preview = on_preview
        self.on_clean = on_clean
        self.on_clear_options = on_clear_options
        self.on_settings = on_settings
        self.on_abort = on_abort

        self._create_ui_buttons()


    def _create_ui_buttons(self):
        TOP_ICONS_DIRECTORY = ICONS_DIRECTORY.joinpath("top_menu_icons")
        # icons
        # clean icon
        self.btn_clean_icon = ctk.CTkImage(
            light_image=Image.open(TOP_ICONS_DIRECTORY / "clean.png"),
            dark_image=Image.open(TOP_ICONS_DIRECTORY / "clean.png"),
        )

        # preview icon
        self.btn_preview_icon = ctk.CTkImage(
            light_image=Image.open(TOP_ICONS_DIRECTORY / "preview.png"),
            dark_image=Image.open(TOP_ICONS_DIRECTORY / "preview.png"),
        )

        # unselect icon
        self.btn_unselect_icon = ctk.CTkImage(
            light_image=Image.open(TOP_ICONS_DIRECTORY / "unselect.png"),
            dark_image=Image.open(TOP_ICONS_DIRECTORY / "unselect.png"),
        )

        # abort icon
        self.btn_abort_icon = ctk.CTkImage(
            light_image=Image.open(TOP_ICONS_DIRECTORY / "abort.png"),
            dark_image=Image.open(TOP_ICONS_DIRECTORY / "abort.png"),
        )

        # settings icon
        self.btn_settings_icon = ctk.CTkImage(
            light_image=Image.open(TOP_ICONS_DIRECTORY / "settings.png"),
            dark_image=Image.open(TOP_ICONS_DIRECTORY / "settings.png"),
        )

        # button `clean`
        self.btn_clean = ctk.CTkButton(
            self,
            text="Clean",
            image=self.btn_clean_icon,
            width=TOP_MENU_BUTTON_CLEAN_WIDTH,
            height=TOP_MENU_BUTTON_CLEAN_HEIGHT,
            command=self.on_clean
        )
        self.btn_clean.pack(side="left", padx=10, pady=10)

        # button `preview`
        self.btn_preview = ctk.CTkButton(
            self,
            text="Preview",
            image=self.btn_preview_icon,
            width=TOP_MENU_BUTTON_PREVIEW_WIDTH,
            height=TOP_MENU_BUTTON_PREVIEW_HEIGHT,
            command=self.on_preview
        )
        self.btn_preview.pack(side="left", padx=10, pady=10)

        # button `clear_options`
        self.btn_unselect = ctk.CTkButton(
            self,
            text="Unselect",
            image=self.btn_unselect_icon,
            width=TOP_MENU_BUTTON_UNSELECT_WIDTH,
            height=TOP_MENU_BUTTON_UNSELECT_HEIGHT,
            command=self.on_clear_options
        )
        self.btn_unselect.pack(side="left", padx=10, pady=10)

        # button `abort`
        self.btn_abort = ctk.CTkButton(
            self,
            text="Abort",
            image=self.btn_abort_icon,
            width=TOP_MENU_BUTTON_ABORT_WIDTH,
            height=TOP_MENU_BUTTON_ABORT_HEIGHT,
            command=self.on_abort
        )
        self.btn_abort.pack(side="left", padx=10, pady=10)

        # button `settings`
        self.btn_settings = ctk.CTkButton(
            self,
            text="Settings",
            image=self.btn_settings_icon,
            width=TOP_MENU_BUTTON_SETTINGS_WIDTH,
            height=TOP_MENU_BUTTON_SETTINGS_HEIGHT,
            command=self.on_settings
        )
        self.btn_settings.pack(side="right", padx=10, pady=10)


    def enable_preview(self):
        self.btn_preview.configure(state="normal")

    def disable_preview(self):
        self.btn_preview.configure(state="disabled")

    def enable_clean(self):
        self.btn_clean.configure(state="normal")

    def disable_clean(self):
        self.btn_clean.configure(state="disabled")

    def enable_abort(self):
        self.btn_abort.configure(state="normal")

    def disable_abort(self):
        self.btn_abort.configure(state="disabled")