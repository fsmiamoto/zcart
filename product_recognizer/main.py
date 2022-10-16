#! /usr/bin/python3
import sys
import time
import pdb
import argparse
import logging
from typing import List
from frame_object import FrameObject

from frame_preprocessor import FramePreprocessor
from frame_object_detector import FrameObjectDetector, PytorchFrameObjectDetector, TfLiteFrameObjectDetector
from frame_object_filter import FrameObjectFilter
from weight_sensor import WeightSensor
from cart_service import CartServiceClient
from video_window import VideoWindow
from product_recognizer import ProductRecognizer
from queue import Queue
from logger import Logger
from video_stream import VideoStream


def run(args):
    log = Logger()
    log.setLevel(logging.DEBUG)
    weight_sensor = WeightSensor()
    # detector: FrameObjectDetector = TfLiteFrameObjectDetector(
    #     model_path=args.model_file, labelmap_path=args.label_file
    # )
    detector: FrameObjectDetector = PytorchFrameObjectDetector('./model/mobilenet_classes.txt', logger=log)

    height, width = detector.get_input_dimensions()

    # preprocessor = FramePreprocessor(
    #     height, width, is_floating_model=detector.is_floating_model()
    # )
    cart_service_client = CartServiceClient(args.cart_service)

    log.info("will start video stream")
    stream = VideoStream(resolution=(width, height)).start()
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
        cart_id=args.cart_id,
    )

    object_filter = FrameObjectFilter(confidence_thresold=args.confidence_threshold)
    product_recognizer.start()

    video_window = VideoWindow() if args.with_window else None

    while True:
        try:
            if video_window:
                # Used for FPS calculation
                video_window.start_tick()

            frame = stream.read_frame()

            # input = preprocessor.process(frame)
            # log.info("will preprocess")
            input = frame[:,:, [2,1,0]]
            detector.infer(input)
            objects = detector.get_objects()
            log.info("got objects")
            # pdb.set_trace()
            filtered_objects = object_filter.filter(objects)

            log.info(f"Survivors: {len(filtered_objects)}")

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
    parser = argparse.ArgumentParser("zCart Product Recognizer Application")
    parser.add_argument("--cart_id", dest="cart_id", default="1", help="Cart ID")
    parser.add_argument(
        "--model_file",
        dest="model_file",
        default="./model/detect.tflite",
        help="Path for the TFLite model file",
    )
    parser.add_argument(
        "--label_file",
        dest="label_file",
        default="./model/labelmap.txt",
        help="Path for the label file",
    )
    parser.add_argument(
        "--confidence",
        dest="confidence_threshold",
        default=0.55,
        type=float,
        help="Confidence threshold",
    )
    parser.add_argument("--width", dest="width", default=640, help="Frame width")
    parser.add_argument("--height", dest="height", default=480, help="Frame height")
    parser.add_argument(
        "--with_window",
        dest="with_window",
        default=True,
        help="Whether to show video window",
    )
    parser.add_argument(
        "--cart_service",
        dest="cart_service",
        default=True,
        help="Cart Service Endpoint",
    )
    args = parser.parse_args()
    run(args)
