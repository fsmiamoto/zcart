import {
  CartProvider,
  Item,
  CartItem,
  ItemHandler,
} from "src/service/cart_provider";

const IMAGE_BASE_URL = "https://zcart-test-images.s3.amazonaws.com";

export class StubCartProvider implements CartProvider {
  private ITEMS: Item[] = [
    {
      title: "Coca-Cola 2L",
      price: 9.0,
      image_url: `${IMAGE_BASE_URL}/coca2l.png`,
    },
    {
      title: "Batata Ruffles",
      price: 7.0,
      image_url: `${IMAGE_BASE_URL}/ruffles.png`,
    },
    {
      title: "Chamyto",
      price: 5.99,
      image_url: `${IMAGE_BASE_URL}/chamyto.png`,
    },
    {
      title: "Bombril",
      price: 2.99,
      image_url: `${IMAGE_BASE_URL}/bombril.png`,
    },
    {
      title: "Café ",
      price: 9.99,
      image_url: `${IMAGE_BASE_URL}/cafe.png`,
    },
    {
      title: "Nissin Lámen",
      price: 3.99,
      image_url: `${IMAGE_BASE_URL}/lamen.png`,
    },
    {
      title: "Leite Longa Vida 1L",
      price: 3.99,
      image_url: `${IMAGE_BASE_URL}/leite.png`,
    },
    {
      title: "Nuggets",
      price: 7.99,
      image_url: `${IMAGE_BASE_URL}/nuggets.png`,
    },
    {
      title: "Pão de Alho",
      price: 4.99,
      image_url: `${IMAGE_BASE_URL}/paodealho.png`,
    },
    {
      title: "Tang",
      price: 1.99,
      image_url: `${IMAGE_BASE_URL}/tang.png`,
    },
  ];

  private cartItems: CartItem[];
  private addHandler?: ItemHandler;
  private removeHandler?: ItemHandler;

  private readonly interval: number;
  private readonly delay: number;

  constructor(
    delay: number = 2000,
    interval: number = 60000,
    maxQuantity: number = 4
  ) {
    this.interval = interval;
    this.delay = delay;
    this.cartItems = this.generateRandomCartItemList();
    setInterval(() => {
      const randomItem = this.ITEMS[this.randomIndex()];
      const shouldAdd = Math.round(Math.random());
      if (shouldAdd) {
        this.addHandler &&
          this.addHandler({
            quantity: Math.round((maxQuantity - 1) * Math.random()) + 1,
            ...randomItem,
          });
      } else {
        this.removeHandler &&
          this.removeHandler({
            quantity: Math.round((maxQuantity - 1) * Math.random()) + 1,
            ...randomItem,
          });
      }
    }, this.interval);
  }

  async ListCartItems(): Promise<CartItem[]> {
    await new Promise((resolve) => {
      setTimeout(resolve, this.delay);
    });
    return this.cartItems;
  }

  OnAddProduct(handler: ItemHandler) {
    this.addHandler = handler;
  }

  OnRemoveProduct(handler: ItemHandler) {
    this.removeHandler = handler;
  }

  AddItem() {
    const randomItem = this.ITEMS[this.randomIndex()];

    const inCart = this.cartItems.find(
      (item) => item.title === randomItem.title
    );
    if (!inCart) {
      this.cartItems.push({ ...randomItem, quantity: this.randomQuantity() });
    } else {
      inCart.quantity += 1;
    }
  }

  RemoveItem() {
    const randomItem = this.ITEMS[this.randomIndex()];

    const inCart = this.cartItems.find(
      (item) => item.title === randomItem.title
    );
    if (!inCart) {
      // Should not happen on real implementation
      return;
    }

    inCart.quantity -= 1;
    if (inCart.quantity === 0) {
      // Good enough for small arrays
      this.cartItems = this.cartItems.filter(
        (item) => item.title === randomItem.title
      );
    }
  }

  private randomQuantity(): number {
    return 1 + Math.floor(Math.random() * 9);
  }

  private randomCartLenth(): number {
    return 1 + Math.floor(Math.random() * this.ITEMS.length);
  }

  private randomIndex(): number {
    return Math.floor(Math.random() * this.ITEMS.length);
  }

  private generateRandomCartItemList(): CartItem[] {
    const randomOrder = this.ITEMS.sort(() => 0.5 - Math.random());

    return randomOrder
      .map((item) => ({ ...item, quantity: this.randomQuantity() }))
      .slice(0, this.randomCartLenth());
  }
}
