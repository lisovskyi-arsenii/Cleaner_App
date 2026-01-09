import json
from typing import Any, Dict

from src.config.settings import BASE_DIR, THEME_DIRECTORY

def load_theme_colors(theme_name: str) -> Dict[str, Any]:
    try:
        with open(f"{BASE_DIR}/{THEME_DIRECTORY}/theme_name", encoding="utf-8") as f:
            theme_colors = json.load(f)

            light_colors = {}
            dark_colors = {}

            light = 0
            dark = 1

            if "CTkLabel" in theme_colors:
                label_colors = theme_colors["CTkLabel"]
                if "text_color" in label_colors:
                    light_colors["text_color"] = label_colors["text_color"][light]
                    dark_colors["text_color"] = label_colors["text_color"][dark]

            if "CTkButton" in theme_colors:
                button_colors = theme_colors["CTkButton"]
                if "text_color" in button_colors:
                    light_colors["text_color"] = button_colors["text_color"][light]
                    dark_colors["text_color"] = button_colors["text_color"][dark]

    except FileNotFoundError:
        print(f"Theme name {theme_name} was not found.")
        return get_default_colors()
    except Exception as e:
        print(f"Error: {e}")
        return get_default_colors()


def get_default_colors() -> Dict[str, Any]:
    pass