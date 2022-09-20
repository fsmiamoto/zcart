from typing import List, Optional


class FrameObject:
    def __init__(self, label: str, score: float, box: List[float]):
        self.label = label
        self.score = score
        self.box = box
