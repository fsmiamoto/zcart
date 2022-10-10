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

enum CartEvent {
  ProductAdded = "product_added",
  ProductRemoved = "product_removed",
}

interface CartEventNotification {
  event: CartEvent;
  cart_product: CartServiceCartProduct;
}

export class CartServiceCartProvider implements CartProvider {
  private readonly axios: AxiosInstance;
  private readonly cartId: string;
  private readonly baseUrl: string;
  private websocket?: WebSocket;
  private addProductHandler?: ItemHandler;
  private removeProductHandler?: ItemHandler;

  constructor(url: string, cartId: string = "2") {
    this.cartId = cartId;
    this.baseUrl = url;
    this.axios = axios.create({
      baseURL: url,
    });
    this.setupWebsocket()
  }

  private setupWebsocket() {
    if (this.websocket) {
      this.websocket.close()
    }
    const webSocketUrl = `${this.baseUrl}/cart/${this.cartId}/ws`.replace("http", "ws");
    this.websocket = new WebSocket(webSocketUrl);
    this.websocket.onopen = (_event) => {
      console.log("opening websocket");
    };
    this.websocket.onclose = (_event) => {
      console.log("closing websocket");
      setTimeout(() => this.setupWebsocket(), 1000)
    };
    this.websocket.onmessage = (event) => {
      const payload = JSON.parse(event.data) as CartEventNotification;
      if (payload.event === CartEvent.ProductAdded) {
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

  async Checkout() {
    await this.axios.post(`/cart/${this.cartId}/checkout`)
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
