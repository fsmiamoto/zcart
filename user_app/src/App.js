import { Layout, List, Avatar, Typography, Button, Statistic } from 'antd';
import "./App.css"
import React, { useState, useCallback, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

const defaultCart = [
  {
    title: 'Coca-Cola 500ml',
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

const App = () => {
  const [socketUrl, _setSocketUrl] = useState('ws://localhost:3333/ws');
  const [cartProducts, setMessageHistory] = useState(defaultCart);

  const [subtotal, setSubtotal] = useState(0.0);

  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

  useEffect(() => {
    if (lastMessage !== null) {
      setMessageHistory((prev) => prev.concat(lastMessage));
    }
  }, [lastMessage, setMessageHistory]);

  useEffect(() => {
    setSubtotal(cartProducts.reduce((total, item) => total + item.price, 0.0))
  }, [cartProducts])

  const handleEmptyCart = useCallback(
    () => setMessageHistory([]),
    []
  );

  const handleClickSendMessage = useCallback(() => sendMessage('Hello'), []);

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Connecting',
    [ReadyState.OPEN]: 'Open',
    [ReadyState.CLOSING]: 'Closing',
    [ReadyState.CLOSED]: 'Closed',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

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
                description={item.description || "Produto"}
              />
              <Statistic value={item.price} precision={2} prefix="R$" decimalSeparator="," />
            </List.Item>
          )}
        />
        <div>Subtotal: <Statistic value={subtotal} prefix="R$" precision={2} decimalSeparator="," /></div>
        <Button onClick={handleEmptyCart}>
          Click Me to empty the cart
        </Button>
        <Button
          onClick={handleClickSendMessage}
          disabled={readyState !== ReadyState.OPEN}
        >
          Click Me to add a new product
        </Button>
      </Content>
      <Footer>
        <span>The WebSocket is currently {connectionStatus}</span>
        {lastMessage ? <span>Last message: {lastMessage.data}</span> : null}
      </Footer>
    </Layout>
  );
};

export default App; 
