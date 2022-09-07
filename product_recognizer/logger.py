import logging
import sys


class Logger(logging.Logger):
    def __init__(self, name="logger"):
        super().__init__(name)
        self.__handler = logging.StreamHandler(sys.stdout)
        self.__handler.setFormatter(
            logging.Formatter("%(asctime)s - %(levelname)s - %(message)s")
        )
        self.addHandler(self.__handler)
        self.setLevel(logging.DEBUG)
