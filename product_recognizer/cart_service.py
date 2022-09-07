import requests
from enum import Enum


class UpdateCartRequestAction(Enum):
    ADD = "add"
    REMOVE = "remove"


class UpdateCartRequest:
    def __init__(self, product_id: str, quantity: int, action: UpdateCartRequestAction):
        self.__product_id = product_id
        self.__quantity = quantity
        self.__action = action

    def to_json(self):
        return {
            "product_id": self.__product_id,
            "quantity": self.__quantity,
            "action": self.__action.value,
        }


class CartServiceClient:
    def __init__(self, base_url="http://localhost:3333"):
        self.__base_url = base_url

    def execute(self, cart_id: str, request: UpdateCartRequest):
        url = f"{self.__base_url}/cart/{cart_id}/products"
        return requests.post(url, json=request.to_json())
