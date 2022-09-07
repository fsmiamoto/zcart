from queue import Queue, Empty
from typing import Callable, Dict
from collections import defaultdict
from threading import Thread
from logger import Logger

from cart_service import CartServiceClient, UpdateCartRequest, UpdateCartRequestAction
from weight_sensor import WeightSensor

WEIGHT_THRESHOLD = 20

# TODO: Add local database
# Label -> ID
RECOGNIZED_OBJECTS = {"bottle": "1"}


class ObjectDiff:
    def __init__(
        self,
        queue: Queue,
        label_getter: Callable[[int], str],
        weight_sensor: WeightSensor,
        logger: Logger,
        cart_service_client: CartServiceClient,
    ):
        self.queue = queue
        self.last_frame_objects = {}
        self.last_weight_reading = 0.0
        self.label_getter = label_getter
        self.weight_sensor = weight_sensor
        self.cart_service_client = cart_service_client
        self.log = logger

    def start(self):
        self.__stopped = False
        Thread(target=self.__worker, args=[]).start()

    def stop(self):
        self.__stopped = True

    def __worker(self):
        while True:
            try:
                objects = self.queue.get_nowait()
                current_frame_objects = self.__build_object_dict(objects)
                diff = self.__get_difference(
                    current_frame_objects, self.last_frame_objects
                )

                self.log.debug(
                    f"current: {current_frame_objects} - last {self.last_frame_objects}"
                )

                for label, count in diff.items():
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

                self.log.debug(f"diff: {diff}")
            except Empty:
                if self.__stopped:
                    return

    def __call_cart_service(self, label: str, count: int):
        self.log.info("will call cart service")
        request = UpdateCartRequest(
            RECOGNIZED_OBJECTS[label],
            abs(count),
            UpdateCartRequestAction.ADD
            if count > 0
            else UpdateCartRequestAction.REMOVE,
        )
        try:
            response = self.cart_service_client.execute("2", request)
            self.log.info(f"got status {response.status_code}")
        except:
            self.log.error("exception while calling cart service")

    def __build_object_dict(self, objects) -> Dict[str, int]:
        result = defaultdict(int)
        for (_, _, class_index) in objects:
            result[self.label_getter(int(class_index))] += 1
        return result

    def __get_difference(
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
