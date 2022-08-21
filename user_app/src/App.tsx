import { useState, useCallback, useEffect } from "react";
import {
    Layout,
    List,
    message,
    Avatar,
    Button,
    Statistic,
    Modal,
} from "antd";
import { ShoppingCartOutlined } from "@ant-design/icons";
import { CartProvider, CartItem } from "src/service/cart_provider";
import { LoadingSpinner } from "src/components/loading_spinner";
import "./App.css";

export interface Props {
    cartProvider: CartProvider;
}

const { Header, Content, Footer } = Layout;

function App(props: Props) {
    const [cartProducts, setCartProducts] = useState<CartItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [subtotal, setSubtotal] = useState(0.0);
    const [modalVisible, setModalVisible] = useState(false);

    useEffect(() => {
        if (!loading) return;
        props.cartProvider.ListCartItems().then((items) => {
            setCartProducts(items);
            setLoading(false);
        });
    }, [props.cartProvider, loading]);

    useEffect(() => {
        props.cartProvider.OnAddProduct((item) => {
            message.success(
                `Produto ${item.quantity}x ${item.title} adicionado ao carrinho`
            );
            setLoading(true);
        });

        props.cartProvider.OnRemoveProduct((item) => {
            message.success(
                `Produto ${item.quantity}x ${item.title} removido do carrinho`
            );
            setLoading(true);
        });
    }, [props.cartProvider]);

    useEffect(() => {
        setSubtotal(
            cartProducts.reduce(
                (total, item) => total + item.price * item.quantity,
                0.0
            )
        );
    }, [cartProducts]);

    const handleFinalize = useCallback(() => {
        setModalVisible(false);
        setCartProducts([]);
        message.success("Obrigado!");
    }, []);

    return (
        <Layout className="app">
            <Header className="header">
                <ShoppingCartOutlined id="icon" />
                <span id="title">zCart</span>
            </Header>
            <Content className="content">
                {loading ? (
                    <LoadingSpinner />
                ) : (
                    <>
                        <List
                            className="products"
                            itemLayout="horizontal"
                            dataSource={cartProducts}
                            locale={{ emptyText: "Carrinho vazio" }}
                            renderItem={(item) => (
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
                            )}
                        />
                        <div className="subtotal">
                            <Button
                                type="primary"
                                size="large"
                                onClick={() => setModalVisible(true)}
                                disabled={cartProducts.length === 0}
                            >
                                Finalizar compra
                            </Button>
                            <span>
                                Subtotal:{" "}
                                <Statistic
                                    value={subtotal}
                                    prefix="R$"
                                    precision={2}
                                    decimalSeparator=","
                                />
                            </span>
                        </div>
                    </>
                )}
            </Content>
            <Footer></Footer>
            <Modal
                visible={modalVisible}
                onOk={handleFinalize}
                onCancel={() => setModalVisible(false)}
                okText={"Finalizar"}
                cancelText={"Cancelar"}
            >
                Deseja finalizar a compra?
            </Modal>
        </Layout>
    );
}

export default App;
