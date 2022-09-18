#! /usr/bin/python3

import sys
import time
from os import environ
from typing import List
from frame_object import FrameObject

from frame_preprocessor import FramePreprocessor
from frame_object_detector import FrameObjectDetector
from frame_object_filter import FrameObjectFilter
from weight_sensor import WeightSensor
from cart_service import CartServiceClient
from video_window import VideoWindow
from product_recognizer import ProductRecognizer
from queue import Queue
from logger import Logger
from video_stream import VideoStream

CONFIDENCE_THRESHOLD = 0.55

MODEL_FILE = "./model/detect.tflite"
LABELMAP_FILE = "./model/labelmap.txt"

# Where do these come from??
INPUT_MEAN = 127.5
INPUT_STD = 127.5

FRAME_WIDTH = 640
FRAME_HEIGHT = 480

WITH_VIDEO_WINDOW = environ.get("WITH_VIDEO_WINDOW")


def main():
    log = Logger()
    weight_sensor = WeightSensor()
    detector = FrameObjectDetector(model_path=MODEL_FILE, labelmap_path=LABELMAP_FILE)

    height, width = detector.get_input_dimensions()

    preprocessor = FramePreprocessor(
        height, width, INPUT_MEAN, INPUT_STD, detector.is_floating_model()
    )
    cart_service_client = CartServiceClient("http://localhost:3333")

    log.info("will start video stream")
    stream = VideoStream(resolution=(FRAME_WIDTH, FRAME_HEIGHT)).start()
    time.sleep(1)
    log.info("done starting video stream")

    log.info("will tare sensor")
    weight_sensor.tare()
    log.info("done taring")

    frame_objects_queue: "Queue[List[FrameObject]]" = Queue()

    product_recognizer = ProductRecognizer(
        queue=frame_objects_queue,
        weight_sensor=weight_sensor,
        logger=log,
        cart_service_client=cart_service_client,
    )

    object_filter = FrameObjectFilter(confidence_thresold=CONFIDENCE_THRESHOLD)

    product_recognizer.start()

    video_window = VideoWindow() if WITH_VIDEO_WINDOW else None

    while True:
        try:
            if video_window:
                # Used for FPS calculation
                video_window.start_tick()

            frame = stream.read_frame()

            input = preprocessor.process(frame)
            detector.infer(input)

            filtered_objects = object_filter.filter(detector.get_objects())

            frame_objects_queue.put(filtered_objects)

            if video_window:
                video_window.display(frame)

        except (KeyboardInterrupt, SystemExit):
            log.info("received exit signal, cleaning up")
            weight_sensor.cleanup()
            if video_window:
                video_window.stop()
            product_recognizer.stop()
            stream.stop()
            sys.exit()


if __name__ == "__main__":
    main()
