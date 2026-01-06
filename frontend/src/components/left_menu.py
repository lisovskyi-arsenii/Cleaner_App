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

        self.cleaners_widgets = {}
        self.on_hover_callback = on_hover_callback
        self.checkboxes = []

        # expandable
        # {cleaner_id: True/False}
        self.expanded_state = {}

    # draw all what can you delete in the left menu
    def draw_cleaners(self, cleaners_data) -> None:
        if not cleaners_data:
            print("No cleaners data to display.")
            return

        self.clear()

        for cleaner in cleaners_data:
            cleaner_id = cleaner["id"]
            self.expanded_state[cleaner_id] = True
            self.create_cleaner_section(cleaner)

    # create section for one of the cleaners
    def create_cleaner_section(self, cleaner) -> None:
        cleaner_id = cleaner["id"]

        main_frame = ctk.CTkFrame(
            self,
            fg_color="transparent",
        )
        main_frame.pack(fill="x", padx=5, pady=2)

        header_frame = ctk.CTkFrame(
            main_frame,
            fg_color=("gray85", "gray25")
        )
        header_frame.pack(fill="x", padx=2, pady=2)

        toggle_btn = ctk.CTkButton(
            header_frame,
            text="▼",
            width=30,
            height=30,
            font=ctk.CTkFont(size=20, weight="bold"),
            fg_color="transparent",
            hover_color=("gray75", "gray35"),
            command=lambda: self._toggle_cleaner(cleaner_id),
        )
        toggle_btn.pack(side="left", padx=5)

        # label for cleaner
        cleaner_label = ctk.CTkLabel(
            header_frame,
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
            main_frame,
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

        options_frame = ctk.CTkFrame(
            main_frame,
            fg_color="transparent",
        )
        options_frame.pack(fill="x", padx=(10, 2), pady=(0, 5))

        # options checkboxes
        for option in cleaner["options"]:
            self.create_option_checkbox(options_frame, cleaner, option)

        # separator between cleaners
        separator = ctk.CTkFrame(self, height=2, fg_color="gray30")
        separator.pack(fill="x", padx=10, pady=10)

        self.cleaners_widgets[cleaner_id] = {
            "main_frame": main_frame,
            "toggle_btn": toggle_btn,
            "cleaner_description": cleaner_description,
            "options_frame": options_frame,
            "separator": separator,
            "expanded_state": True,
        }


    def _toggle_cleaner(self, cleaner_id) -> None:
        if cleaner_id not in self.cleaners_widgets:
            return

        widgets = self.cleaners_widgets[cleaner_id]
        current_state = widgets["expanded_state"]

        new_state = not current_state
        widgets["expanded_state"] = new_state
        self.expanded_state[cleaner_id] = new_state

        widgets["toggle_btn"].configure(text="▼" if new_state else "▶")

        if new_state:
            # Розгорнути
            widgets["cleaner_description"].pack(fill="x", padx=10, pady=(5, 0))
            widgets["options_frame"].pack(fill="x", padx=(15, 2), pady=(5, 5))
        else:
            # Згорнути
            widgets["cleaner_description"].pack_forget()
            widgets["options_frame"].pack_forget()
            

    # create checkbox for one of the options
    def create_option_checkbox(self, parent, cleaner, option):
        var = ctk.BooleanVar()

        option_checkbox = ctk.CTkCheckBox(
            parent,
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
            'widget': option_checkbox,
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

    # clear all selected options
    def clear_selected_checkboxes(self):
        for checkbox in self.checkboxes:
            if 'variable' in checkbox:
                checkbox['variable'].set(False)
            elif 'widget' in checkbox:
                checkbox['widget'].deselect()

    # clear everything
    def clear(self):
        for widget in self.winfo_children():
            widget.destroy()
        self.checkboxes = []
