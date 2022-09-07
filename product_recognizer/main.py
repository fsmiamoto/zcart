#! /usr/bin/python3

import sys
import time

from frame_preprocessor import FramePreprocessor
from object_detector import ObjectDetector
from weight_sensor import WeightSensor
from cart_service import CartServiceClient
from video_window import VideoWindow
from object_diff import ObjectDiff
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
    detector = ObjectDetector(model_path=MODEL_FILE, labelmap_path=LABELMAP_FILE)

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

    object_diff = ObjectDiff(
        queue=objects_queue,
        label_getter=detector.get_label,
        weight_sensor=weight_sensor,
        logger=log,
        cart_service_client=cart_service_client,
    )

    object_diff.start()

    video_window = VideoWindow()

    while True:
        try:
            # Use for FPS calculation
            video_window.start_tick()

            frame = stream.read_frame().copy()

            input = preprocessor.process(frame)
            detector.infer(input)

            boxes, classes, scores = (
                detector.get_boxes(),
                detector.get_classes(),
                detector.get_scores(),
            )

            filtered_objects = list(
                filter(
                    lambda tuple: tuple[0] >= CONFIDENCE_THRESHOLD,
                    zip(scores, boxes, classes),
                )
            )

            objects_queue.put(filtered_objects)

            video_window.display(frame)

        except (KeyboardInterrupt, SystemExit):
            log.info("received exit signal, cleaning up")
            weight_sensor.cleanup()
            video_window.stop()
            object_diff.stop()
            stream.stop()
            sys.exit()


if __name__ == "__main__":
    main()
