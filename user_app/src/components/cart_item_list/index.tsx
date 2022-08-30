import { Avatar, List, Statistic } from "antd";
import { CartItem } from "src/service/cart_provider";
import "./index.css"

function CartItemRow(props: { cartItem: CartItem }) {
    const item = props.cartItem;
    return (
        <List.Item className="cart-item">
            <List.Item.Meta
                className="cart-item-meta"
                avatar={<Avatar src={item.image_url} />}
                title={`${item.quantity}x ${item.title}`}
                description={`R$ ${item.price
                    .toFixed(2)
                    .replace(".", ",")}`}
            />
            <Statistic
                className="cart-item-price"
                value={item.price * item.quantity}
                precision={2}
                prefix="R$"
                decimalSeparator=","
            />
        </List.Item>
    )

}

export function CartItemList(props: { cartItems: CartItem[] }) {
    const { cartItems } = props;
    return (
        <List
            className="cart-items"
            itemLayout="horizontal"
            dataSource={cartItems}
            locale={{ emptyText: "Carrinho vazio" }}
            renderItem={(item) => <CartItemRow cartItem={item} />} />
    )

}
