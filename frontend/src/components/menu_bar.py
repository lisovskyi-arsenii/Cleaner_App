# menu bar component

import tkinter as tk


class MenuBar:
    def __init__(self, parent):
        self.parent = parent
        self.menubar = tk.Menu(parent)
        parent.config(menu=self.menubar)

        # file menu
        self._create_file_menu()
        # edit menu
        self._create_edit_menu()
        # view menu
        self._create_view_menu()
        # tools
        self._create_tools_menu()
        # help menu
        self._create_help_menu()


        # TODO decide should I change theme here or in settings
        # theme menu



    # util functions for creation menus
    def _create_file_menu(self):
        file_menu = tk.Menu(self.menubar, tearoff=0)
        self.menubar.add_cascade(label="File", menu=file_menu)
        self.menubar.add_command(label="Exit", command=lambda: self.exit)

    def _create_edit_menu(self):
        edit_menu = tk.Menu(self.menubar, tearoff=0)
        self.menubar.add_cascade(label="Edit", menu=edit_menu)
        self.menubar.add_command(label="Unselect options", command=self.unselect_options)

    def _create_view_menu(self):
        view_menu = tk.Menu(self.menubar, tearoff=0)
        self.menubar.add_cascade(label="View", menu=view_menu)

    def _create_tools_menu(self):
        tools_menu = tk.Menu(self.menubar, tearoff=0)
        self.menubar.add_cascade(label="Tools", menu=tools_menu)

    def _create_help_menu(self):
        help_menu = tk.Menu(self.menubar, tearoff=0)
        self.menubar.add_cascade(label="Help", menu=help_menu)
        help_menu.add_command(label="Documentation", command=self.show_documentation)
        help_menu.add_command(label="About", command=self.show_about)



    # commands
    def unselect_options(self):
        pass

    def exit(self):
        pass

    def show_documentation(self):
        pass

    def show_about(self):
        pass


