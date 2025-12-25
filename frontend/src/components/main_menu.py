# Main menu component

import customtkinter as ctk

from src.config.settings import MAIN_MENU_CORNER_RADIUS

class MainMenu(ctk.CTkFrame):
    def __init__(self, parent):
        super().__init__(
            parent,
            corner_radius=MAIN_MENU_CORNER_RADIUS
        )

        # labels for hover info
        # title
        self.hover_title_main = ctk.CTkLabel(
            self,
            font=ctk.CTkFont(size=28, weight="bold"),
            anchor="w",
        )
        self.hover_title_main.pack(pady=20, padx=20, fill="x")

        # description
        self.hover_description_main = ctk.CTkLabel(
            self,
            font=ctk.CTkFont(size=24, weight="bold"),
            anchor="w"
        )
        self.hover_description_main.pack(pady=20, padx=20, fill="x")

        self.show_placeholder_text()

    # default placeholder
    def show_placeholder_text(self):
        self.hover_title_main.configure(text="Select cleaners")
        self.hover_description_main.configure(
            text="Hover over items to see details"
        )

    # show all info about certain cleaner
    def show_widget_info(self, widget):
        self.hover_title_main.configure(text=widget['name'])
        self.hover_description_main.configure(text=widget['description'])

    # TODO - fix this method
    def show_results(self, results):
        # clean everything from main menu
        for widget in self.winfo_children():
            widget.destroy()

        # show results
        title = ctk.CTkLabel(
            self,
            text="Results",
            font=ctk.CTkFont(size=24, weight="bold"),
        )
        title.pack(pady=20, padx=20, fill="x")

        for result in results:
            result_label = ctk.CTkLabel(
                self,
                text=f"{result['name']}: {result['size']}",
                font=ctk.CTkFont(size=14),
            )
            result_label.pack(pady=5, padx=20, fill="x")

        pass


