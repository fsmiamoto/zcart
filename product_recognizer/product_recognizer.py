from queue import Queue, Empty
from typing import Dict, List
from collections import defaultdict
from threading import Thread
from logger import Logger

from cart_service import CartServiceClient, UpdateCartRequest, UpdateCartRequestAction
from frame_object import FrameObject
from weight_sensor import WeightSensor
from product_catalog import ProductCatalog, StubProductCatalog


class ProductRecognizer:
    def __init__(
        self,
        queue: "Queue[List[FrameObject]]",
        weight_sensor: WeightSensor,
        logger: Logger,
        cart_service_client: CartServiceClient,
        cart_id: str,
        catalog: ProductCatalog = StubProductCatalog(),
        weight_tolerance: float = 0.15,
    ):
        self.queue = queue
        self.weight_sensor = weight_sensor
        self.cart_service_client = cart_service_client
        self.log = logger
        self.catalog = catalog

        self.weight_tolerance = weight_tolerance
        self.cart_id = cart_id

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

                if len(frame_diff) == 0:
                    self.log.debug("empty diff")
                    continue

                self.log.info(f"object diff in frame: {frame_diff}")

                weight_reading = self.weight_sensor.get_reading(samples=5)
                self.log.debug(
                    f"weight: {weight_reading} - last weight: {self.last_weight_reading}"
                )

                for label, count in frame_diff.items():
                    if not self.__valid_weight_difference(label, count, weight_reading):
                        self.log.info("ignoring, not valid weight difference")
                        continue

                    self.__call_cart_service(label, count)
                    self.last_weight_reading = weight_reading
                    self.last_frame_objects[label] = current_frame_objects[label]
                    if self.last_frame_objects[label] == 0:
                        del self.last_frame_objects[label]

            except Empty:
                if self.__stopped:
                    return

    def __valid_weight_difference(self, label: str, count: int, reading: float) -> bool:
        if reading > self.last_weight_reading and count < 0:
            self.log.info(f"weight increased but count is negative ({count}), ignoring")
            return False

        if reading < self.last_weight_reading and count > 0:
            self.log.info(f"weight decreased but count is positive ({count}), ignoring")
            return False

        product = self.catalog.get_product(label)
        if not product:
            # TODO: Review this
            return False

        weight_difference = reading - self.last_weight_reading
        expected_difference = count * product.weight_in_grams

        self.log.info(
            f"expected_difference: {expected_difference}, actual: {weight_difference}"
        )

        return self.__is_weight_diff_withing_tolerance(
            expected_difference, weight_difference
        )

    def __is_weight_diff_withing_tolerance(
        self, expected: float, actual: float
    ) -> bool:
        # Use absolute value since with negative values the comparisons would invert
        expected = abs(expected)
        actual = abs(actual)
        return (1 - self.weight_tolerance) * expected <= actual and actual <= (
            1 + self.weight_tolerance
        ) * expected

    def __call_cart_service(self, label: str, count: int):
        self.log.info("will call cart service")

        request = self.__build_cart_service_request(label, count)

        try:
            response = self.cart_service_client.execute(self.cart_id, request)
            self.log.info(f"got status {response.status_code}")
        except:
            self.log.error("exception while calling cart service")

    def __build_cart_service_request(self, label: str, count: int):
        product = self.catalog.get_product(label)
        if not product:
            raise Exception(f"product not found for label {label}")

        return UpdateCartRequest(
            product.product_id,
            abs(count),
            UpdateCartRequestAction.ADD
            if count > 0
            else UpdateCartRequestAction.REMOVE,
        )

    def __build_object_dict(self, objects: List[FrameObject]) -> Dict[str, int]:
        result = defaultdict(int)
        for object in objects:
            result[object.label] += 1
        return result

    def __get_frame_diff(
        self, current_frame_objects: Dict[str, int], last_frame_objects: Dict[str, int]
    ) -> Dict[str, int]:
        result = {}

        for label in current_frame_objects:
            self.log.debug(f"label: {label}")
            if not self.catalog.get_product(label):
                continue
            result[label] = current_frame_objects[label]

        for label in last_frame_objects:
            if not self.catalog.get_product(label):
                continue
            if label in result:
                # Might be negative
                result[label] -= last_frame_objects[label]
                if result[label] == 0:
                    del result[label]
            else:
                result[label] = -last_frame_objects[label]

        return result
