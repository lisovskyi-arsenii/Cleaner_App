# file for storing all wrapper functions
import functools
import threading


# wrapper function for clearing checkboxes after fetching data from backend
def async_action_clear_checkboxes(func):
    @functools.wraps(func)
    def wrapper(self, *args, **kwargs):
        def thread_func():
            func(self, *args, **kwargs)
            self.after(0, self.left_menu.clear_selected_checkboxes)
            print(f"DEBUG: {func.__name__} finished, checkboxes cleared")

        threading.Thread(target=thread_func, daemon=True).start()
    return wrapper