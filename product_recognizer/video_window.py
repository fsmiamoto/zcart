import cv2


class VideoWindow:
    WINDOW_NAME = "Object Detection"

    def __init__(self):
        self.tick_frequency = cv2.getTickFrequency()
        self.__start_tick = 0

    def start(self):
        cv2.namedWindow(self.WINDOW_NAME, cv2.WINDOW_NORMAL)

    def stop(self):
        cv2.destroyAllWindows()

    def start_tick(self):
        self.__start_tick = cv2.getTickCount()

    def display(self, frame):
        end_tick = cv2.getTickCount()
        frame_rate = self.tick_frequency / (end_tick - self.__start_tick)

        cv2.putText(
            frame,
            "FPS: {0:.2f}".format(frame_rate),
            (30, 50),
            cv2.FONT_HERSHEY_SIMPLEX,
            1,
            (255, 255, 0),
            2,
            cv2.LINE_AA,
        )
        cv2.imshow(self.WINDOW_NAME, frame)
        cv2.waitKey(1)
