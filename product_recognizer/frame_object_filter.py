from typing import List

from frame_object import FrameObject


class FrameObjectFilter:
    def __init__(self, confidence_thresold: float):
        self.confidence_thresold = confidence_thresold

    def filter(self, objects: List[FrameObject]) -> List[FrameObject]:
        return list(
            (filter(lambda obj: obj.score >= self.confidence_thresold, objects))
        )
