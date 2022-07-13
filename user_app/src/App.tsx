import { Layout, List, Avatar, Typography, Button, Statistic } from 'antd';
import "./App.css"
import React, { useState, useCallback, useEffect } from 'react';

type CartItem = {
    title: string
    price: number
    description? : string
}

const defaultCart:  CartItem[] = [
    {
        title: 'Coca-Cola 2L',
        price: 6.00,
    },
    {
        title: 'Batata Ruffles',
        price: 7.00,
    },
    {
        title: 'Pringles',
        price: 10.00
    },
];

const { Title } = Typography;
const { Header, Content, Footer } = Layout;

function App() {
    const [cartProducts, setCartProducts] = useState(defaultCart);
    const [subtotal, setSubtotal] = useState(0.0);

    useEffect(() => {
        setSubtotal(cartProducts.reduce((total, item) => total + item.price, 0.0))
    }, [cartProducts])

    const handleEmptyCart = useCallback(
        () => setCartProducts([]),
        []
    );

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Header>
                <Title level={3} style={{ color: 'white' }}>zCart</Title>
            </Header>
            <Content>
                <List
                    itemLayout="horizontal"
                    dataSource={cartProducts}
                    locale={{ emptyText: "Carrinho vazio" }}
                    renderItem={item => (
                        <List.Item>
                            <List.Item.Meta
                                avatar={<Avatar src="https://joeschmoe.io/api/v1/random" />}
                                title={item.title}
                                description={item.description ?? "Produto"}
                            />
                            <Statistic value={item.price} precision={2} prefix="R$" decimalSeparator="," />
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
