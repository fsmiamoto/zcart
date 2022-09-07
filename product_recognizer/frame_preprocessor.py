import numpy as np
import cv2


class FramePreprocessor:
    def __init__(
        self, height: int, width: int, mean: float, std: float, is_floating_model: bool
    ):
        self.__height = height
        self.__width = width
        self.__mean = mean
        self.__std = std
        self.__is_floating_model = is_floating_model

        pass

    def process(self, frame):
        frame_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        frame_resized = cv2.resize(frame_rgb, (self.__width, self.__height))

        result = np.expand_dims(frame_resized, axis=0)

        if self.__is_floating_model:
            result = (np.float32(result) - self.__mean) / self.__std

        return result
