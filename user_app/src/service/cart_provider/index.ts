export interface Item {
  title: string;
  price: number;
  description?: string;
  image_url?: string;
}

export interface CartItem extends Item {
  quantity: number;
}

export interface CartProvider {
  ListCartItems(): Promise<CartItem[]>;
  OnAddProduct(handler: ItemHandler): void;
  OnRemoveProduct(handler: ItemHandler): void;
  Checkout(): Promise<void>;
}

export type ItemHandler = (item: CartItem) => void;
