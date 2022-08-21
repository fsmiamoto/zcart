import axios, { AxiosInstance } from "axios";
import { CartProvider, CartItem, ItemHandler } from "src/service/cart_provider";

interface CartServiceResponse {
  id: string;
  products: CartServiceCartProduct[];
}

interface CartServiceCartProduct {
  cart_id: string;
  product_id: string;
  quantity: number;
  product: {
    id: string;
    name: string;
    price: number;
    image_url?: string;
    description?: string;
  };
}

enum actionType {
  ADD = 0,
  REMOVE = 1,
}

interface cartServiceAction {
  type: number;
  cart_product: CartServiceCartProduct;
}

export class CartServiceCartProvider implements CartProvider {
  private readonly axios: AxiosInstance;
  private readonly cartId: string;
  private readonly websocket: WebSocket;
  private addProductHandler?: ItemHandler;
  private removeProductHandler?: ItemHandler;

  constructor(url: string, cartId: string = "2") {
    this.cartId = cartId;
    this.axios = axios.create({
      baseURL: url,
    });
    const webSocketUrl = `${url}/cart/${cartId}/ws`.replace("http", "ws");
    console.log(webSocketUrl);
    this.websocket = new WebSocket(webSocketUrl);
    this.websocket.onopen = (_event) => {
      console.log("opening websocket");
    };
    this.websocket.onclose = (_event) => {
      console.log("closing websocket");
    };
    this.websocket.onmessage = (event) => {
      const payload = JSON.parse(event.data) as cartServiceAction;
      console.log(payload);
      if (payload.type === actionType.ADD) {
        this.addProductHandler &&
          this.addProductHandler(this.adapter(payload.cart_product));
      } else {
        this.removeProductHandler &&
          this.removeProductHandler(this.adapter(payload.cart_product));
      }
    };
  }

  async ListCartItems(): Promise<CartItem[]> {
    const items = (await this.axios.get(`/cart/${this.cartId}`))
      .data as CartServiceResponse;
    return items.products.map(this.adapter);
  }

  OnAddProduct(handler: ItemHandler) {
    this.addProductHandler = handler;
  }

  OnRemoveProduct(handler: ItemHandler) {
    this.removeProductHandler = handler;
  }

  private adapter(cartProduct: CartServiceCartProduct): CartItem {
    return {
      quantity: cartProduct.quantity,
      title: cartProduct.product.name,
      price: cartProduct.product.price,
      image_url: cartProduct.product.image_url,
      description: cartProduct.product.description,
    };
  }
}
