# configuration file
from pathlib import Path

from src.enums.Enums import Appearance

# API settings
BACKEND_URL         = "http://localhost:8080"
BACKEND_URL_HTTPS   = "https://localhost:8080"
REQUEST_TIMEOUT     = 3  # seconds

# WINDOW settings
WINDOW_TITLE = "CLEANER APP"
WINDOW_WIDTH = 1200
WINDOW_HEIGHT = 700
MIN_WIDTH = 800
MIN_HEIGHT = 600

# FILES SETTINGS
BASE_DIR = Path(__file__).resolve().parent.parent.parent
THEME_DIRECTORY = BASE_DIR / "resources" / "themes"
ICONS_DIRECTORY = BASE_DIR / "resources" / "icons"

# TOP menu settings
TOP_MENU_HEIGHT = 60
TOP_MENU_CORNER_RADIUS = 0

# button `clean`
TOP_MENU_BUTTON_CLEAN_WIDTH = 120
TOP_MENU_BUTTON_CLEAN_HEIGHT = 50

# button `preview`
TOP_MENU_BUTTON_PREVIEW_WIDTH = 110
TOP_MENU_BUTTON_PREVIEW_HEIGHT = 50

# button `unselect`
TOP_MENU_BUTTON_UNSELECT_WIDTH = 120
TOP_MENU_BUTTON_UNSELECT_HEIGHT = 50

# button `abort`
TOP_MENU_BUTTON_ABORT_WIDTH = 120
TOP_MENU_BUTTON_ABORT_HEIGHT = 50

# button `settings`
TOP_MENU_BUTTON_SETTINGS_WIDTH = 100
TOP_MENU_BUTTON_SETTINGS_HEIGHT = 50

# LEFT menu settings
LEFT_MENU_WIDTH = 250
LEFT_MENU_CORNER_RADIUS = 0

# MAIN menu settings
MAIN_MENU_CORNER_RADIUS = 0

# SETTINGS window settings
SETTINGS_WINDOW_TITLE = "Settings"
SETTINGS_WINDOW_WIDTH = 800
SETTINGS_WINDOW_HEIGHT = 600
SETTINGS_APPEARANCE_MODE_OPTION_MENU_WIDTH = 140

# THEME settings
APPEARANCE_MODE = Appearance.DARK
DEFAULT_THEME = "tokyonight.json"
CURRENT_THEME = "tokyonight.json"

# Status colors for different states
STATUS_COLORS = {
    "ready": "gray",
    "processing": "#3b8ed0",  # blue
    "success": "#2fa572",     # green
    "warning": "#ff9500",     # orange
    "error": "#d62828",       # red
    "info": "#6c757d",        # gray-blue
}

# Default text for different states
STATUS_MESSAGES = {
    "ready": "✓ Ready",
    "no_selection": "⚠️ No items selected",
    "connection_error": "❌ Cannot connect to backend",
}
