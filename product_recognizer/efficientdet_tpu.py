import numpy as np
from typing import List, Tuple
from frame_object import FrameObject
from pycoral.utils.edgetpu import make_interpreter
from pycoral.adapters.common import set_input, set_resized_input
from pycoral.adapters.detect import get_objects
from PIL import Image
import cv2


class EfficientDetFrameObjectDetector:
    __labels = ["post_it", "guarana_soda", "coke_soda", "card_deck", "blue_pens"]

    def __init__(self, model_path: str):
        self.__interpreter = make_interpreter(model_path)
        self.__interpreter.allocate_tensors()

    def get_input_dimensions(self) -> Tuple[int, int]:
        _, height, width, _ = self.__interpreter.get_input_details()[0]["shape"]
        return height, width

    def is_floating_model(self) -> bool:
        return True

    def get_objects(self, image) -> List[FrameObject]:
        _, scale = set_resized_input(
            self.__interpreter,
            image.size,
            lambda size: image.resize(size, Image.ANTIALIAS),
        )
        self.__interpreter.invoke()
        objs = get_objects(self.__interpreter, image_scale=scale)

        return [FrameObject(self.__labels[obj.id], obj.score, obj.bbox) for obj in objs]


class EfficientDetFramePreprocessor:
    def __init__(self, width, height):
        self.__width = width
        self.__height = height

    def process(self, frame):
        frame_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        frame_resized = cv2.resize(frame_rgb, (self.__width, self.__height))

        return Image.fromarray(frame_resized), frame_resized
