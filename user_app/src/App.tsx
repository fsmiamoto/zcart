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
    const [subtotal, setSubtotal] = useState(0.5);

    useEffect(() => {
        setSubtotal(cartProducts.reduce((total, item) => total + item.price * item.quantity, 0.0))
    }, [cartProducts])

    const handleEmptyCart = useCallback(() => setCartProducts([]), []);

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
                                avatar={<Avatar src={item.image_url} />}
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
