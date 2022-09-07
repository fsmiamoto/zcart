import cv2
from threading import Thread

# TODO: Type methods
class VideoStream:
    def __init__(self, resolution=(640, 480)):
        self.__stream = cv2.VideoCapture(0)

        if not self.__stream.isOpened():
            raise Exception("failed to open video camera")

        self.__stream.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter_fourcc(*"MJPG"))
        self.__stream.set(3, resolution[0])
        self.__stream.set(4, resolution[1])

        (_, self.__frame) = self.__stream.read()

        self.__stopped = False

    def start(self):
        Thread(target=self.__frame_reader, args=()).start()
        return self

    def stop(self):
        self.__stopped = True

    def read_frame(self):
        return self.__frame

    def __frame_reader(self):
        while True:
            if self.__stopped:
                self.__stream.release()
                return

            (_, self.__frame) = self.__stream.read()
