from typing import Tuple, List
import numpy as np
from frame_object import FrameObject
from tflite_runtime.interpreter import Interpreter


class MobileNetFrameObjectDetector:
    def __init__(self, model_path: str, labelmap_path: str):
        self.__interpreter = Interpreter(model_path=model_path)
        self.__interpreter.allocate_tensors()

        self.__input_details = self.__interpreter.get_input_details()
        self.__output_details = self.__interpreter.get_output_details()

        self.__read_labels_from_file(labelmap_path)

    def get_input_dimensions(self) -> Tuple[int, int]:
        """Get height and width (in that order) of the input"""
        return self.__input_details[0]["shape"][1], self.__input_details[0]["shape"][2]

    def is_floating_model(self) -> bool:
        return self.__input_details[0]["dtype"] == np.float32

    def infer(self, input):
        self.__interpreter.set_tensor(self.__input_details[0]["index"], input)
        self.__interpreter.invoke()

    def get_objects(self) -> List[FrameObject]:
        return [
            FrameObject(label, score, box)
            for (score, label, box) in zip(
                self.get_scores(), self.get_labels(), self.get_boxes()
            )
        ]

    def get_boxes(self):
        return self.__interpreter.get_tensor(self.__output_details[0]["index"])[0]

    def get_scores(self):
        return [
            float(score)
            for score in self.__interpreter.get_tensor(
                self.__output_details[2]["index"]
            )[0]
        ]

    def get_labels(self):
        return [
            self.__get_label(int(class_index)) for class_index in self.__get_classes()
        ]

    def __get_classes(self):
        return self.__interpreter.get_tensor(self.__output_details[1]["index"])[0]

    def __get_label(self, class_index: int) -> str:
        return self.__labels[class_index]

    def __read_labels_from_file(self, path: str):
        with open(path, "r") as file:
            self.__labels = [line.strip() for line in file.readlines()]

        # Have to do a weird fix for label map if using the COCO "starter model" from
        # https://www.tensorflow.org/lite/models/object_detection/overview
        # First label is '???', which has to be removed.
        if self.__labels[0] == "???":
            del self.__labels[0]
