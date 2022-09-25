from abc import ABC, abstractmethod
from typing import Optional
from product import Product


class ProductCatalog(ABC):
    @abstractmethod
    def get_product(self, label: str) -> Optional[Product]:
        pass


@ProductCatalog.register
class StubProductCatalog(ProductCatalog):
    def __init__(self):
        self.LABEL_TO_PRODUCT = {
            "bottle": Product(product_id="1", weight_in_grams=18),
            "mouse": Product(product_id="2", weight_in_grams=100),
            "scissors": Product(product_id="3", weight_in_grams=10),
        }

    def get_product(self, label: str) -> Optional[Product]:
        return self.LABEL_TO_PRODUCT.get(label)
