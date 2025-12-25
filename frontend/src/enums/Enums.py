from enum import Enum

# appearance enum
class Appearance(Enum):
    DARK = "dark"
    LIGHT = "light"
    SYSTEM = "system"


# API_methods enum
class APIMethods(Enum):
    GET = "GET"
    POST = "POST"
    PUT = "PUT"
    UPDATE = "UPDATE"
    DELETE = "DELETE"
