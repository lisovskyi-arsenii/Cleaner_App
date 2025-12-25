# Left component class
import customtkinter as ctk
from src.config.settings import LEFT_MENU_WIDTH, LEFT_MENU_CORNER_RADIUS

class LeftMenu(ctk.CTkScrollableFrame):
    def __init__(self, parent, on_hover_callback=None):
        super().__init__(
            parent,
            width=LEFT_MENU_WIDTH,
            corner_radius=LEFT_MENU_CORNER_RADIUS,
        )

        self.on_hover_callback = on_hover_callback
        self.checkboxes = []


    # draw all what can you delete in the left menu
    def draw_cleaners(self, cleaners_data) -> None:
        if not cleaners_data:
            print("No cleaners data to display.")
            return

        self.clear()

        for cleaner in cleaners_data:
            self.create_cleaner_section(cleaner)


    # create section for one of the cleaners
    def create_cleaner_section(self, cleaner):
        # label for cleaner
        cleaner_label = ctk.CTkLabel(
            self,
            text=cleaner["name"],
            font=ctk.CTkFont(size=20, weight="bold"),
            anchor="w"
        )
        cleaner_label.pack(pady=(15, 5), padx=10, fill="x")

        # callback for label
        if self.on_hover_callback:
            cleaner_label.bind(
                '<Enter>',
                lambda e, c=cleaner: self.on_hover_callback(widget=c)
            )

        # description for cleaner
        cleaner_description = ctk.CTkLabel(
            self,
            text=cleaner["description"],
            font=ctk.CTkFont(size=18, weight="bold"),
            text_color="gray",
            anchor="w",
        )
        cleaner_description.pack(pady=(0, 10), padx=10, fill="x")

        # callback for description
        if self.on_hover_callback:
            cleaner_description.bind(
                '<Enter>',
                lambda e, c=cleaner: self.on_hover_callback(widget=c)
            )

        # options checkboxes
        for option in cleaner["options"]:
            self.create_option_checkbox(cleaner, option)


        # separator between cleaners
        separator = ctk.CTkFrame(self, height=2, fg_color="gray30")
        separator.pack(fill="x", padx=10, pady=10)


    # create checkbox for one of the options
    def create_option_checkbox(self, cleaner, option):
        var = ctk.BooleanVar()

        option_checkbox = ctk.CTkCheckBox(
            self,
            text=option["label"],
            variable=var,
            font=ctk.CTkFont(size=14),
            cursor="hand2",
        )
        option_checkbox.pack(pady=3, padx=20, anchor="w")

        self.checkboxes.append({
            'variable': var,
            'cleaner_id': cleaner['id'],
            'option_id': option['id'],
        })


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

    # clear everything
    def clear(self):
        for widget in self.winfo_children():
            widget.destroy()
        self.checkboxes = []