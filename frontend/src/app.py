# Main class App for all components
import threading
from time import sleep

import customtkinter as ctk

from src.api import backend
from src.api.backend import get_cleaners
from src.components.menu_bar import MenuBar
from src.components.settings_window import SettingsWindow
from src.components.left_menu import LeftMenu
from src.components.main_menu import MainMenu
from src.components.top_menu import TopMenu
from src.config.settings import *
from src.functions.wrapper_functions import async_action_clear_checkboxes


# main app
class App(ctk.CTk):
    def __init__(self):
        super().__init__()

        # data
        # all cleaners which were fetched from backend
        self.cleaners_data = []
        # settings window
        self.settings_window = None

        # bool var for whether any request was sent to backend or not
        self.is_fetching_to_backend = False
        self.current_thread = None

        # base config for window and widgets
        self.current_color_theme = DEFAULT_THEME

        ctk.set_appearance_mode(APPEARANCE_MODE.value)
        ctk.set_default_color_theme(str(THEME_DIRECTORY / self.current_color_theme))

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
        # menu bar
        self.menu_bar = MenuBar(self)
        
        # top menu
        self.top_menu = TopMenu(self,
                                on_preview=self.on_preview_clicked,
                                on_clean=self.on_clean_clicked,
                                on_clear_options=self.on_clear_options_clicked,
                                on_abort=self.on_abort_clicked,
                                on_settings=self.on_settings_clicked
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
                    self.after(0, self._show_connection_error)

            except Exception as e:
                print(f"Error loading data: {e}")
                self.after(0, self._show_connection_error)

        threading.Thread(target=load, daemon=True).start()

    # show message `connection error` if failed to fetch data
    def _show_connection_error(self):
        if hasattr(self.main_menu, 'hover_title_main'):
            self.main_menu.hover_title_main.configure(
                text=STATUS_MESSAGES["connection_error"],
                text_color=STATUS_COLORS["error"]
            )

    # update button states depends on if is any request was sent or not
    def _update_button_states(self):
        if self.is_fetching_to_backend:
            self.top_menu.disable_preview()
            self.top_menu.disable_clean()
            self.top_menu.enable_abort()
        else:
            self.top_menu.enable_preview()
            self.top_menu.enable_clean()
            self.top_menu.disable_abort()

    # show processing status
    def _show_processing_status(self, message: str, color="blue"):
        try:
            if (hasattr(self.main_menu, 'hover_title_main') and
                    self.main_menu.hover_title_main is not None and
                    self.main_menu.hover_title_main.winfo_exists()):
                self.main_menu.hover_title_main.configure(
                    text=message,
                    text_color=color,
                )
            else:
                print(f"Status: {message}")
        except Exception as e:
            print(f"Error updating status: {e} - Message was: {message}")

    # general function for all backend functions
    def _execute_backend_action_async(self, func, action_name="Operation"):
        if self.is_fetching_to_backend:
            print(f"Already executing backend: {func}")
            return

        selected_data = self.left_menu.get_selected()
        if not selected_data:
            self._show_processing_status(
                message=STATUS_MESSAGES["no_selection"],
                color=STATUS_COLORS["warning"]
            )
            return

        self.is_fetching_to_backend = True
        self._update_button_states()
        self._show_processing_status(f"🔄 {action_name}...", "blue")

        def backend_thread():
            try:
                results = func(selected_data)
                self.after(0, lambda: self._handle_backend_complete(results, action_name))
            except Exception as e:
                print(f"{action_name} error: {e}")
                self.after(0, lambda: self._handle_backend_error(e, action_name))

        self.current_thread = threading.Thread(target=backend_thread, daemon=True)
        self.current_thread.start()

    # handle result from backend
    @async_action_clear_checkboxes
    def _handle_backend_complete(self, results, action_name):
        self.is_fetching_to_backend = False
        self._update_button_states()

        if results:
            if isinstance(results, dict) and results.get("partial"):
                self._show_processing_status(
                    f"⚠️ {action_name} cancelled",
                    STATUS_COLORS["warning"]
                )
                # Show partial results if available
                if results.get("data"):
                    self.main_menu.show_results(results["data"])
            else:
                self._show_processing_status(
                    f"✅ {action_name} complete",
                    STATUS_COLORS["success"]
                )
                self.main_menu.show_results(results)
        else:
            self._show_processing_status(
                f"⚠️ {action_name} returned no data",
                STATUS_COLORS["warning"]
            )

        self.after(1500, self.main_menu.show_placeholder_text)


    # show error whether something went wrong
    def _handle_backend_error(self, error, action_name):
        self.is_fetching_to_backend = False
        self._update_button_states()
        self._show_processing_status(f"❌ {action_name} error: {error}", "red")


    # after button `preview` is pressed, this function will be invoked
    def on_preview_clicked(self):
        self._execute_backend_action_async(
            func=backend.preview_cleaners,
            action_name="Preview",
        )


    # after button `clean` is pressed, this function will be invoked
    @async_action_clear_checkboxes
    def on_clean_clicked(self):
        self._execute_backend_action_async(
            func=backend.clean_files,
            action_name="Clean",
        )

    # after button `abort` is pressed, this function will be invoked
    def on_abort_clicked(self):
        if not self.is_fetching_to_backend:
            print("No operation to abort")
            return

        try:
            print(f"Sending abort request")
            response = backend.abort_request()
            print(f"abort: {response}")

            self.is_fetching_to_backend = False
            self._update_button_states()
            self._show_processing_status("⛔ Operation cancelled by user", "orange")

        except Exception as e:
            print(f"Error abort: {e}")
            self._show_processing_status(f"❌ Abort failed: {e}", "red")

    # after button `clear options` is pressed, this function will be invoked
    def on_clear_options_clicked(self):
        self.left_menu.clear_selected_checkboxes()

    # when user hovers one of the widgets, show some info about it in main component
    def on_widget_hover(self, widget):
        self.main_menu.show_widget_info(widget)

    # after button `settings` is pressed, this function will be invoked
    def on_settings_clicked(self):
        if self.settings_window is None or not self.settings_window.winfo_exists():
            self.settings_window = SettingsWindow(self)
        else:
            self.settings_window.focus()

