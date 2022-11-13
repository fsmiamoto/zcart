import numpy as np
import tensorflow as tf
from typing import List, Tuple
from tflite_runtime.interpreter import Interpreter
from frame_object import FrameObject
import pycoral


class EfficientDetFrameObjectDetector:
    __labels = ["post_it", "coke_soda", "guarana_soda", "card_deck", "blue_pens"]

    def __init__(self, model_path: str):
        self.__interpreter = pycoral.make_interpreter(model_path)
        self.__interpreter.allocate_tensors()
        self.__signature_runner = self.__interpreter.get_signature_runner()

    def get_input_dimensions(self) -> Tuple[int, int]:
        _, height, width, _ = self.__interpreter.get_input_details()[0]["shape"]
        return height, width

    def is_floating_model(self) -> bool:
        return True

    def get_objects(self, image) -> List[FrameObject]:
        output = self.__signature_runner(images=image)

        scores = np.squeeze(output["output_1"])
        class_ids = np.squeeze(output["output_2"])
        boxes = np.squeeze(output["output_3"])

        labels = [self.__labels[id.astype("int")] for id in class_ids]

        return [
            FrameObject(label, score, box)
            for (score, label, box) in zip(
                scores,
                labels,
                boxes,
            )
        ]


class EfficientDetFramePreprocessor:
    def __init__(self, width, height):
        self.__width = width
        self.__height = height

    def process(self, frame):
        frame_resized = tf.image.resize(frame, (self.__height, self.__width))
        frame_resized = frame_resized[tf.newaxis, :]
        frame_resized = tf.cast(frame_resized, dtype=tf.uint8)

        return frame_resized
