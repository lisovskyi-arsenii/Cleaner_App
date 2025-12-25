from App import App


def setup_window():
    app = App()
    app.set_cleaners()
    app.draw_cleaners()
    app.mainloop()



def main():
    setup_window()





if __name__ == '__main__':
    main()