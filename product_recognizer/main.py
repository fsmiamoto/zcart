#! /usr/bin/python3
import sys
import time
import argparse
from typing import List
from frame_object import FrameObject

from efficientdet_tpu import (
    EfficientDetFrameObjectDetector,
    EfficientDetFramePreprocessor,
)
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
    weight_sensor = WeightSensor()
    detector = EfficientDetFrameObjectDetector(
        model_path="./model/custom_whole_efficientdet_lite1_resizing_edgetpu.tflite"
    )

    height, width = detector.get_input_dimensions()

    preprocessor = EfficientDetFramePreprocessor(width, height)
    cart_service_client = CartServiceClient(args.cart_service)

    log.info("will start video stream")
    stream = VideoStream(resolution=(args.width, args.height)).start()
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

            input, _ = preprocessor.process(frame)

            objects = detector.get_objects(input)

            filtered_objects = object_filter.filter(objects)

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
    parser.add_argument("--cart_id", dest="cart_id", default="2", help="Cart ID")
    parser.add_argument(
        "--model_file",
        dest="model_file",
        type=str,
        default="./model/detect.tflite",
        help="Path for the TFLite model file",
    )
    parser.add_argument(
        "--label_file",
        dest="label_file",
        type=str,
        default="./model/labelmap.txt",
        help="Path for the label file",
    )
    parser.add_argument(
        "--confidence",
        dest="confidence_threshold",
        type=float,
        default=0.70,
        help="Confidence threshold",
    )
    parser.add_argument(
        "--width", dest="width", type=int, default=640, help="Frame width"
    )
    parser.add_argument(
        "--height", dest="height", type=int, default=480, help="Frame height"
    )
    parser.add_argument(
        "--with_window",
        dest="with_window",
        type=bool,
        default=True,
        help="Whether to show video window",
    )
    parser.add_argument(
        "--cart_service",
        dest="cart_service",
        default="http://localhost:3333",
        help="Cart Service Endpoint",
    )
    args = parser.parse_args()
    run(args)
