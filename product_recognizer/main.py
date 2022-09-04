#! /usr/bin/python3

import sys
import time
from typing import Dict
import cv2
from collections import defaultdict

from frame_preprocessor import FramePreprocessor
from object_detector import ObjectDetector
from weight_sensor import WeightSensor
from cart_service import CartServiceClient, UpdateCartRequest, UpdateCartRequestAction
from logger import Logger
from video_stream import VideoStream

WEIGHT_THRESHOLD = 30
CONFIDENCE_THRESHOLD = 0.55

MODEL_FILE = './model/detect.tflite'
LABELMAP_FILE = './model/labelmap.txt'

# Where do these come from??
INPUT_MEAN = 127.5
INPUT_STD = 127.5

FRAME_WIDTH = 640
FRAME_HEIGHT = 480

CV2_WINDOW_NAME='Object Detection'

def main():
    log = Logger()
    weight_sensor = WeightSensor()
    detector = ObjectDetector(model_path=MODEL_FILE,labelmap_path=LABELMAP_FILE)

    height, width = detector.get_input_dimensions()

    preprocessor = FramePreprocessor(height,width, INPUT_MEAN, INPUT_STD, detector.is_floating_model())
    cart_service_client = CartServiceClient("http://tokyo:3333")

    log.info("will start video stream")
    stream = VideoStream(resolution=(FRAME_WIDTH, FRAME_HEIGHT)).start()
    time.sleep(1)
    log.info("done starting video stream")

    log.info("will tare sensor")
    weight_sensor.tare()
    log.info("done taring")

    cv2.namedWindow(CV2_WINDOW_NAME, cv2.WINDOW_NORMAL)

    frame_rate = 0
    tick_frequency = cv2.getTickFrequency()

    last_weight_reading = 0.0
    last_frame_objects = {}

    # TODO: Move to module
    def build_object_dict(objects) -> Dict[str,int]:
        result = defaultdict(int)
        for (_, _, class_index) in objects:
            result[detector.get_label(int(class_index))] += 1
        return result

    def get_difference(current_frame_objects: Dict[str,int], last_frame_objects: Dict[str,int]) -> Dict[str, int]:
        result = {}
        for label in current_frame_objects:
            result[label] = current_frame_objects[label]
        for label in last_frame_objects:
            if label in result:
                # Might be negative
                result[label] -= last_frame_objects[label]
                if result[label] == 0:
                    del result[label]
            else:
                result[label] = -last_frame_objects[label]

        return result

    while True:
        try:
            start_tick = cv2.getTickCount()
            frame = stream.read_frame().copy()

            input = preprocessor.process(frame)
            detector.infer(input)

            boxes, classes, scores = detector.get_boxes(), detector.get_classes(), detector.get_scores()

            filtered_objects = list(
                filter(lambda tuple: tuple[0] >= CONFIDENCE_THRESHOLD, zip(scores,boxes,classes))
            )

            current_frame_objects = build_object_dict(filtered_objects)

            diff = get_difference(current_frame_objects, last_frame_objects)

            # TODO: This should be a thread
            for label, count in diff.items():
                if label != 'bottle':
                    continue

                current_weight_reading = weight_sensor.get_reading()
                if abs(current_weight_reading-last_weight_reading) < WEIGHT_THRESHOLD:
                    log.info("ignoring, not reached weight threshold")
                    continue

                log.info("will call cart service")
                request = UpdateCartRequest("1", abs(count), UpdateCartRequestAction.ADD if count > 0 else UpdateCartRequestAction.REMOVE)
                response = cart_service_client.execute("2", request)
                log.info(f"got status {response.status_code}")
                last_weight_reading = current_weight_reading

            log.info(f"diff: {diff}")

            # Display rectangle
            # for i, (score, box, class_index) in enumerate(filtered_objects):
            #     ymin = int(max(1,(box[0] * FRAME_HEIGHT)))
            #     xmin = int(max(1,(box[1] * FRAME_WIDTH)))
            #     ymax = int(min(FRAME_HEIGHT,(box[2] * FRAME_HEIGHT)))
            #     xmax = int(min(FRAME_WIDTH,(box[3] * FRAME_WIDTH)))
            #
            #     cv2.rectangle(frame, (xmin,ymin), (xmax,ymax), (10, 255, 0), 2)
            #
            #     object_name = detector.get_label(int(class_index))
            #
            #     # Example: bottle: 75%
            #     label = f'{object_name}: {int(score*100)}%'
            #
            #     labelSize, _baseLine = cv2.getTextSize(label, cv2.FONT_HERSHEY_SIMPLEX, 0.7, 2)
            #     label_ymin = max(ymin, labelSize[1],  10) # Make sure not to draw label too close to top of window
            #     cv2.putText(frame, label, (xmin, label_ymin-7), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 0, 0), 2) # Draw label text
            #
            #     # Draw circle in center
            #     xcenter = xmin + (int(round((xmax - xmin) / 2)))
            #     ycenter = ymin + (int(round((ymax - ymin) / 2)))
            #     cv2.circle(frame, (xcenter, ycenter), 5, (0,0,255), thickness=-1)
            #
            #     log.info(f"Object {i}: {object_name} at ({xcenter},{ycenter}) with confidence {scores[i]*100}")

            cv2.putText(frame,'FPS: {0:.2f}'.format(frame_rate),(30,50),cv2.FONT_HERSHEY_SIMPLEX,1,(255,255,0),2,cv2.LINE_AA)
            cv2.imshow(CV2_WINDOW_NAME, frame)

            end_tick = cv2.getTickCount()
            frame_rate = tick_frequency/(end_tick-start_tick)

            last_frame_objects = current_frame_objects

            # Need to call waitKey to display the frame
            cv2.waitKey(1)


        except (KeyboardInterrupt, SystemExit):
            log.info("received exit signal, cleaning up")
            weight_sensor.cleanup()
            cv2.destroyAllWindows()
            stream.stop()
            sys.exit()


if __name__ == '__main__':
    main()
