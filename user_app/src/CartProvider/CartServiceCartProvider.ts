import axios, { AxiosInstance } from "axios";
import { CartProvider, CartItem, ItemHandler } from "./CartProvider";

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

export class CartServiceCartProvider implements CartProvider {
  private axios: AxiosInstance;

  constructor(url: string) {
    this.axios = axios.create({
      baseURL: url,
    });
  }

  async ListCartItems(): Promise<CartItem[]> {
    const items = (await this.axios.get("/cart/2")).data as CartServiceResponse;
    return items.products.map(this.adapter);
  }

  OnAddProduct(handler: ItemHandler) {}

  OnRemoveProduct(handler: ItemHandler) {}

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
