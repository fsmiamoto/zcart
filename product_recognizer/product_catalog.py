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
            "bottle": Product(product_id="1", weight_in_grams=24),
            "coke_soda": Product(product_id="7", weight_in_grams=490),
            "guarana_soda": Product(product_id="8", weight_in_grams=18),
            "blue_pens": Product(product_id="9", weight_in_grams=42),
            "post_it": Product(product_id="10", weight_in_grams=63),
            "card_deck": Product(product_id="11", weight_in_grams=60),
        }

    def get_product(self, label: str) -> Optional[Product]:
        return self.LABEL_TO_PRODUCT.get(label)
