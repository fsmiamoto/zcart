import numpy as np
import pdb
import tensorflow as tf
from typing import List, Tuple
from tflite_runtime.interpreter import Interpreter
from frame_object import FrameObject
import cv2


class EfficientDetFrameObjectDetector:
    __labels = ["post_it", "guarana_soda", "coke_soda", "card_deck", "blue_pens"]

    def __init__(self, model_path: str):
        self.__interpreter = Interpreter(model_path=model_path)
        self.__interpreter.allocate_tensors()
        self.__signature_runner = self.__interpreter.get_signature_runner()

    def get_input_dimensions(self) -> Tuple[int, int]:
        _, height, width, _ = self.__interpreter.get_input_details()[0]["shape"]
        return height, width

    def is_floating_model(self) -> bool:
        return True

    def get_objects(self, image) -> List[FrameObject]:
        output = self.__signature_runner(images=image)

        #count = int(np.squeeze(output['output_0']))
        scores = np.squeeze(output["output_1"])
        class_ids = np.squeeze(output["output_2"])
        boxes = np.squeeze(output["output_3"])

        #pdb.set_trace()
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
        frame_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        frame_resized = cv2.resize(frame_rgb, (self.__width, self.__height))
        #frame_resized  = tf.convert_to_tensor(frame_resized, dtype=tf.float32)
        frame_resized = tf.image.convert_image_dtype(frame_resized, tf.uint8)
        frame_resized = tf.image.resize(frame_resized, (320,320))

        preprocessed_frame = frame_resized[tf.newaxis, :]
        preprocessed_frame= tf.cast(preprocessed_frame, dtype=tf.uint8)

        return preprocessed_frame, frame_resized
