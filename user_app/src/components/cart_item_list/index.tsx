import { Avatar, List, Statistic } from "antd";
import { CartItem } from "src/service/cart_provider";
import "./index.css"

function CartItemRow(props: { cartItem: CartItem }) {
    const item = props.cartItem;
    return (
        <List.Item className="cart-item">
            <List.Item.Meta
                className="cart-item-meta"
                avatar={<Avatar src={item.image_url} size="large" />}
                title={`${item.title}`}
                description={`${item.quantity} items`} />
            <Statistic
                className="cart-item-price"
                value={item.price}
                precision={2}
                prefix="R$"
                decimalSeparator=","
                groupSeparator="."
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
            locale={{ emptyText: "Empty cart" }}
            renderItem={(item) => <CartItemRow cartItem={item} />} />
    )

}
