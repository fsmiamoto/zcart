from queue import Queue, Empty
from typing import Dict, List
from collections import defaultdict
from threading import Thread
from logger import Logger

from cart_service import CartServiceClient, UpdateCartRequest, UpdateCartRequestAction
from object import Object
from weight_sensor import WeightSensor

WEIGHT_THRESHOLD = 20
CART_ID = "2"

# TODO: Add local database
# Label -> ID
RECOGNIZED_OBJECTS = {"bottle": "1"}


class ProductRecognizer:
    def __init__(
        self,
        queue: Queue,
        weight_sensor: WeightSensor,
        logger: Logger,
        cart_service_client: CartServiceClient,
    ):
        self.queue = queue
        self.weight_sensor = weight_sensor
        self.cart_service_client = cart_service_client
        self.log = logger

        self.last_frame_objects = {}
        self.last_weight_reading = 0.0

    def start(self):
        self.__stopped = False
        Thread(target=self.__worker, args=[]).start()

    def stop(self):
        self.__stopped = True

    def __worker(self):
        while True:
            try:
                # Non-blocking read to avoid getting stuck here
                objects = self.queue.get_nowait()

                current_frame_objects = self.__build_object_dict(objects)
                frame_diff = self.__get_frame_diff(
                    current_frame_objects, self.last_frame_objects
                )

                for label, count in frame_diff.items():
                    current_weight_reading = self.weight_sensor.get_reading(samples=5)
                    self.log.debug(
                        f"weight: {current_weight_reading} - last weight: {self.last_weight_reading}"
                    )

                    if current_weight_reading < 0:
                        self.log.info(
                            f"ignoring bad weight reading: {current_weight_reading}"
                        )

                    if (
                        abs(current_weight_reading - self.last_weight_reading)
                        < WEIGHT_THRESHOLD
                    ):
                        self.log.info("ignoring, not reached weight threshold")
                        continue

                    self.__call_cart_service(label, count)
                    self.last_weight_reading = current_weight_reading
                    self.last_frame_objects[label] = current_frame_objects[label]
                    if self.last_frame_objects[label] == 0:
                        del self.last_frame_objects[label]

                self.log.debug(f"diff: {frame_diff}")
            except Empty:
                if self.__stopped:
                    return

    def __call_cart_service(self, label: str, count: int):
        self.log.info("will call cart service")

        request = self.__build_cart_service_request(label, count)

        try:
            response = self.cart_service_client.execute(CART_ID, request)
            self.log.info(f"got status {response.status_code}")
        except:
            self.log.error("exception while calling cart service")

    def __build_cart_service_request(self, label: str, count: int):
        return UpdateCartRequest(
            RECOGNIZED_OBJECTS[label],
            abs(count),
            UpdateCartRequestAction.ADD
            if count > 0
            else UpdateCartRequestAction.REMOVE,
        )

    def __build_object_dict(self, objects: List[Object]) -> Dict[str, int]:
        result = defaultdict(int)
        for object in objects:
            result[object.label] += 1
        return result

    def __get_frame_diff(
        self, current_frame_objects: Dict[str, int], last_frame_objects: Dict[str, int]
    ) -> Dict[str, int]:
        result = {}

        for label in current_frame_objects:
            if label not in RECOGNIZED_OBJECTS:
                continue
            result[label] = current_frame_objects[label]

        for label in last_frame_objects:
            if label not in RECOGNIZED_OBJECTS:
                continue
            if label in result:
                # Might be negative
                result[label] -= last_frame_objects[label]
                if result[label] == 0:
                    del result[label]
            else:
                result[label] = -last_frame_objects[label]

        return result
