from typing import List

from object import Object


class ObjectFilter:
    def __init__(self, confidence_thresold: float):
        self.confidence_thresold = confidence_thresold

    def filter(self, objects: List[Object]) -> List[Object]:
        return list(
            (filter(lambda obj: obj.score >= self.confidence_thresold, objects))
        )
