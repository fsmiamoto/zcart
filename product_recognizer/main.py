#! /usr/bin/python3

import sys
import time

from frame_preprocessor import FramePreprocessor
from object import Object
from frame_object_detector import FrameObjectDetector
from object_filter import ObjectFilter
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


def main():
    log = Logger()
    weight_sensor = WeightSensor()
    detector = FrameObjectDetector(model_path=MODEL_FILE, labelmap_path=LABELMAP_FILE)

    height, width = detector.get_input_dimensions()

    preprocessor = FramePreprocessor(
        height, width, INPUT_MEAN, INPUT_STD, detector.is_floating_model()
    )
    cart_service_client = CartServiceClient("http://tokyo:3333")

    log.info("will start video stream")
    stream = VideoStream(resolution=(FRAME_WIDTH, FRAME_HEIGHT)).start()
    time.sleep(1)
    log.info("done starting video stream")

    log.info("will tare sensor")
    weight_sensor.tare()
    log.info("done taring")

    objects_queue = Queue()

    product_recognizer = ProductRecognizer(
        queue=objects_queue,
        weight_sensor=weight_sensor,
        logger=log,
        cart_service_client=cart_service_client,
    )

    object_filter = ObjectFilter(confidence_thresold=CONFIDENCE_THRESHOLD)

    product_recognizer.start()

    video_window = VideoWindow()

    while True:
        try:
            # Use for FPS calculation
            video_window.start_tick()

            frame = stream.read_frame()

            input = preprocessor.process(frame)
            detector.infer(input)

            objects = [
                Object(label=label, score=score)
                for (score, label) in zip(detector.get_scores(), detector.get_labels())
            ]

            filtered_objects = object_filter.filter(objects)

            objects_queue.put(filtered_objects)

            video_window.display(frame)

        except (KeyboardInterrupt, SystemExit):
            log.info("received exit signal, cleaning up")
            weight_sensor.cleanup()
            video_window.stop()
            product_recognizer.stop()
            stream.stop()
            sys.exit()


if __name__ == "__main__":
    main()
