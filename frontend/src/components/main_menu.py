# Main menu component

import customtkinter as ctk

from src.config.settings import MAIN_MENU_CORNER_RADIUS

class MainMenu(ctk.CTkFrame):
    def __init__(self, parent):
        super().__init__(
            parent,
            corner_radius=MAIN_MENU_CORNER_RADIUS
        )

        self.showing_results = False
        self._create_hover_widgets()


    # initialize all ui's
    def _create_hover_widgets(self):
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
        if self.showing_results:
            return

        self.hover_title_main.configure(text=widget['name'])
        self.hover_description_main.configure(text=widget['description'])

    # TODO - fix this method
    def show_results(self, results):
        self.showing_results = True

        # clean everything from main menu
        for widget in self.winfo_children():
            widget.destroy()

        scrollable = ctk.CTkScrollableFrame(self)

        # show results
        total_size_mb = results['total_size'] / 1024 / 1024
        title = ctk.CTkLabel(
            self,
            text="Results",
            font=ctk.CTkFont(size=28, weight="bold"),
        )
        title.pack(pady=10, padx=20, anchor="n")

        # General information
        total_info = ctk.CTkLabel(
            scrollable,
            text=f"Total: {total_size_mb:.2f} MB • {results['total_files']} files",
            font=ctk.CTkFont(size=20),
            text_color="gray70"
        )
        total_info.pack(pady=(0, 20), padx=20, anchor="w")

        for item in results['items']:
            size_mb = item['size'] / 1024 / 1024

            # Frame для кожного результату
            item_frame = ctk.CTkFrame(scrollable)
            item_frame.pack(pady=5, padx=20, fill="x")

            # Назва cleaner + option
            name_label = ctk.CTkLabel(
                item_frame,
                text=f"{item['cleaner_id']} - {item['option_id']}",
                font=ctk.CTkFont(size=16, weight="bold"),
                anchor="w"
            )
            name_label.pack(pady=(10, 5), padx=15, anchor="w")

            # Розмір та кількість файлів
            size_label = ctk.CTkLabel(
                item_frame,
                text=f"Size: {size_mb:.2f} MB • {item['file_count']} files",
                font=ctk.CTkFont(size=14),
                text_color="gray70",
                anchor="w"
            )
            size_label.pack(pady=(0, 5), padx=15, anchor="w")

        scrollable.pack(fill="both", expand=True, padx=10, pady=10)


    def restore_hover_view(self):
        """Повертає view до hover стану"""
        self.showing_results = False

        # Очистити результати
        for widget in self.winfo_children():
            widget.destroy()

        # Відтворити hover widgets
        self._create_hover_widgets()

