import json
from typing import Any, Dict, Optional
from pathlib import Path
from src.config.settings import (BASE_DIR, THEME_DIRECTORY,
                                 APPEARANCE_MODE, CURRENT_THEME, DEFAULT_THEME)
from src.enums.Enums import Appearance


def load_theme_colors(theme_name: str, appearance_mode: Optional[Appearance] = None) -> Dict[str, Any]:
    try:
        if appearance_mode is None:
            appearance_mode = APPEARANCE_MODE

        mode_index = 0 if appearance_mode == APPEARANCE_MODE.LIGHT else 1

        theme_path = Path(BASE_DIR) / THEME_DIRECTORY / theme_name

        with open(f"{BASE_DIR}/{THEME_DIRECTORY}/theme_name", encoding="utf-8") as f:
            theme_data = json.load(f)

        colors = {}

        if "CTk" in theme_data and "fg_color" in theme_data["CTk"]:
            colors["window_bg"] = theme_data["CTk"]["fg_color"][mode_index]

        if "CTkFrame" in theme_data:
            frame_data = theme_data["CTkFrame"]
            if "fg_color" in frame_data:
                colors["frame_bg"] = frame_data["fg_color"][mode_index]
            if "top_fg_color" in frame_data:
                colors["frame_top_bg"] = frame_data["top_fg_color"][mode_index]
            if "border_color" in frame_data:
                colors["border"] = frame_data["border_color"][mode_index]
            if "corner_radius" in frame_data:
                colors["corner_radius"] = frame_data["corner_radius"]
            if "border_width" in frame_data:
                colors["border_width"] = frame_data["border_width"]

        if "CTkLabel" in theme_data and "text_color" in theme_data["CTkLabel"]:
            colors["text"] = theme_data["CTkLabel"]["text_color"][mode_index]
            # Use text color for status colors that should match theme
            colors["ready"] = colors["text"]
            colors["cancelled"] = colors["text"]
            colors["info"] = colors["text"]

        if "CTkButton" in theme_data:
            button_data = theme_data["CTkButton"]
            if "fg_color" in button_data:
                colors["button"] = button_data["fg_color"][mode_index]
                colors["processing"] = colors["button"]  # Use button color for processing
            if "hover_color" in button_data:
                colors["button_hover"] = button_data["hover_color"][mode_index]
            if "text_color" in button_data:
                colors["button_text"] = button_data["text_color"][mode_index]
            if "text_color_disabled" in button_data:
                colors["button_text_disabled"] = button_data["text_color_disabled"][mode_index]

        if "CTkEntry" in theme_data:
            entry_data = theme_data["CTkEntry"]
            if "fg_color" in entry_data:
                colors["entry_bg"] = entry_data["fg_color"][mode_index]
            if "border_color" in entry_data:
                colors["entry_border"] = entry_data["border_color"][mode_index]
            if "text_color" in entry_data:
                colors["entry_text"] = entry_data["text_color"][mode_index]
            if "placeholder_text_color" in entry_data:
                colors["entry_placeholder"] = entry_data["placeholder_text_color"][mode_index]

        if "CTkCheckBox" in theme_data:
            checkbox_data = theme_data["CTkCheckBox"]
            if "fg_color" in checkbox_data:
                colors["checkbox"] = checkbox_data["fg_color"][mode_index]
            if "hover_color" in checkbox_data:
                colors["checkbox_hover"] = checkbox_data["hover_color"][mode_index]
            if "checkmark_color" in checkbox_data:
                colors["checkmark"] = checkbox_data["checkmark_color"][mode_index]

        if "CTkScrollbar" in theme_data:
            scrollbar_data = theme_data["CTkScrollbar"]
            if "button_color" in scrollbar_data:
                colors["scrollbar"] = scrollbar_data["button_color"][mode_index]
            if "button_hover_color" in scrollbar_data:
                colors["scrollbar_hover"] = scrollbar_data["button_hover_color"][mode_index]

        colors["success"] = "#2fa572"
        colors["warning"] = "#ff9500"
        colors["error"] = "#d62828"

        return colors

    except FileNotFoundError:
        print(f"Theme name {theme_name} was not found.")
        return get_default_colors()
    except KeyError as e:
        print(f"Missing key in theme file.")
        return get_default_colors()
    except Exception as e:
        print(f"Error loading theme colors: {e}")
        return get_default_colors()


def get_default_colors() -> Dict[str, Any]:
    return {
        # Status colors
        "ready": "gray",
        "processing": "#3b8ed0",
        "success": "#2fa572",
        "warning": "#ff9500",
        "error": "#d62828",
        "cancelled": "gray",
        "info": "#6c757d",

        # UI colors
        "text": "#c0caf5",
        "window_bg": "#1a1b26",
        "frame_bg": "#24283b",
        "frame_top_bg": "#1f2335",
        "border": "#565f89",

        # Button colors
        "button": "#7aa2f7",
        "button_hover": "#5a7abf",
        "button_text": "#c0caf5",
        "button_text_disabled": "#565f89",

        # Entry colors
        "entry_bg": "#1a1b26",
        "entry_border": "#565f89",
        "entry_text": "#c0caf5",
        "entry_placeholder": "#565f89",

        # Checkbox colors
        "checkbox": "#7aa2f7",
        "checkbox_hover": "#5a7abf",
        "checkmark": "#1a1b26",

        # Scrollbar colors
        "scrollbar": "#565f89",
        "scrollbar_hover": "#7aa2f7",

        # Geometry
        "corner_radius": 10,
        "border_width": 0,
    }

def get_current_theme_color(color_key: str, default: Optional[str]=None) -> str:
    if default is None:
        default = "gray"

    try:
        theme_name = CURRENT_THEME if CURRENT_THEME else DEFAULT_THEME

        colors = load_theme_colors(theme_name)

        return colors.get(color_key, default)

    except Exception as e:
        print(f"Error loading theme colors: {e}")
        return default

def get_all_theme_colors() -> Dict[str, Any]:
    theme_name = CURRENT_THEME if CURRENT_THEME else DEFAULT_THEME
    return load_theme_colors(theme_name, APPEARANCE_MODE)

def reload_theme_colors() -> Dict[str, Any]:
    return get_all_theme_colors()