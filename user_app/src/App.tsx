import React, { useState, useCallback, useEffect } from 'react';
import { Layout, List, Avatar, Button, Statistic } from 'antd';
import { ShoppingCartOutlined } from "@ant-design/icons"
import "./App.css"

type CartItem = {
    title: string
    price: number
    quantity: number
    description?: string
    image_url?: string
}

const IMAGE_BASE_URL = "https://zcart-test-images.s3.amazonaws.com"

const defaultCart: CartItem[] = [
    {
        title: 'Coca-Cola 2L',
        price: 6.00,
        quantity: 1,
        image_url: `${IMAGE_BASE_URL}/coca2l.png`
    },
    {
        title: 'Batata Ruffles',
        price: 7.00,
        quantity: 2,
        image_url: `${IMAGE_BASE_URL}/ruffles.png`
    },
    {
        title: 'Chamyto',
        price: 5.99,
        quantity: 5,
        image_url: `${IMAGE_BASE_URL}/chamyto.png`
    },
];

const { Header, Content, Footer } = Layout;

function App() {
    const [cartProducts, setCartProducts] = useState(defaultCart);
    const [subtotal, setSubtotal] = useState(0.0);

    useEffect(() => {
        setSubtotal(cartProducts.reduce((total, item) => total + item.price * item.quantity, 0.0))
    }, [cartProducts])

    const handleEmptyCart = useCallback(() => setCartProducts([]), []);

    return (
        <Layout className="app">
            <Header className="header">
                <ShoppingCartOutlined id="icon" />
                <span id="title">zCart</span>
            </Header>
            <Content className="content">
                <List
                    className="products"
                    itemLayout="horizontal"
                    dataSource={cartProducts}
                    locale={{ emptyText: "Carrinho vazio" }}
                    renderItem={item => (
                        <List.Item className="cart-item">
                            <List.Item.Meta
                                className="cart-item-meta"
                                avatar={<Avatar src={item.image_url} />}
                                title={`${item.quantity}x ${item.title}`}
                                description={`R$ ${item.price.toFixed(2).replace(".", ",")}`}
                            />
                            <Statistic className="cart-item-price" value={item.price * item.quantity} precision={2} prefix="R$" decimalSeparator="," />
                        </List.Item>
                    )}
                />
                <div className="subtotal">
                    <Button type="primary" className="buy-button">
                        Finalizar compra
                    </Button>
                    <span>Subtotal: <Statistic value={subtotal} prefix="R$" precision={2} decimalSeparator="," /></span>
                </div>
            </Content>
            <Footer>
                <div className="buttons">
                    <Button onClick={handleEmptyCart}>
                        Click Me to empty the cart
                    </Button>
                    <Button>
                        Click Me to add a new product
                    </Button>
                </div>
            </Footer>
        </Layout>
    );
};

export default App; 
