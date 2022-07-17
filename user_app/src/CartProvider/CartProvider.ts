
export interface Item {
  title: string
  price: number
  description?: string
  image_url?: string
}

export interface CartItem extends Item{
    quantity: number
}

export interface CartProvider {
  ListCartItems() : CartItem[]
  OnAddProduct(handler: ItemHandler): void
  OnRemoveProduct(handler: ItemHandler): void
}

export type ItemHandler = (item: Item) => void
