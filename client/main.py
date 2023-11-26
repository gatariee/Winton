#!/usr/bin/env python3

from UserInterface.widgets.winton import Winton
from Utils.beacon import dispatch

def main():
    app = Winton(dispatch)
    app.mainloop()


if __name__ == "__main__":
    main()
