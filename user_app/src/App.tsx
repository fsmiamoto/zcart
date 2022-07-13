import { Layout, List, Avatar, Button, Statistic } from 'antd';
import "./App.css"
import React, { useState, useCallback, useEffect } from 'react';

type CartItem = {
    title: string
    price: number
    quantity: number
    description?: string
    image_url?: string
}

const PLACEHOLDER_PRODUCT_URL = "./images/placeholder.png"

const defaultCart: CartItem[] = [
    {
        title: 'Coca-Cola 2L',
        price: 6.00,
        quantity: 1,
    },
    {
        title: 'Batata Ruffles',
        price: 7.00,
        quantity: 2,
    },
    {
        title: 'Pringles',
        price: 10.00,
        quantity: 3,
    },
];

const { Header, Content, Footer } = Layout;

function App() {
    const [cartProducts, setCartProducts] = useState(defaultCart);
    const [subtotal, setSubtotal] = useState(1.0);

    useEffect(() => {
        setSubtotal(cartProducts.reduce((total, item) => total + item.price * item.quantity, 0.0))
    }, [cartProducts])

    const handleEmptyCart = useCallback(
        () => setCartProducts([]),
        []
    );

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Header>
                <h3 style={{ color: 'white' }}>zCart</h3>
            </Header>
            <Content>
                <List
                    itemLayout="horizontal"
                    dataSource={cartProducts}
                    locale={{ emptyText: "Carrinho vazio" }}
                    renderItem={item => (
                        <List.Item>
                            <List.Item.Meta
                                avatar={<Avatar src={item.image_url ?? PLACEHOLDER_PRODUCT_URL} />}
                                title={`${item.quantity}x ${item.title}`}
                                description={item.description ?? "Produto"}
                            />
                            <Statistic value={item.price * item.quantity} precision={2} prefix="R$" decimalSeparator="," />
                        </List.Item>
                    )}
                />
                <div>Subtotal: <Statistic value={subtotal} prefix="R$" precision={2} decimalSeparator="," /></div>
                <Button onClick={handleEmptyCart}>
                    Click Me to empty the cart
                </Button>
                <Button>
                    Click Me to add a new product
                </Button>
            </Content>
            <Footer>
            </Footer>
        </Layout>
    );
};

export default App; 
