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
import { CartItemList } from "./components/cart_item_list";

export interface Props {
    cartProvider: CartProvider;
}

const { Header, Content, Footer } = Layout;

function App(props: Props) {
    const [cartItems, setCartItems] = useState<CartItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [subtotal, setSubtotal] = useState(0.0);
    const [modalVisible, setModalVisible] = useState(false);

    useEffect(() => {
        if (!loading) return;
        props.cartProvider.ListCartItems().then((items) => {
            setCartItems(items);
            setLoading(false);
        });
    }, [props.cartProvider, loading]);

    useEffect(() => {
        props.cartProvider.OnAddProduct((item) => {
            message.info(
                <span>Produto {item.quantity}x {item.title} <b>adicionado</b> ao carrinho</span>
            );
            setLoading(true);
        });

        props.cartProvider.OnRemoveProduct((item) => {
            message.info(
                <span>Produto {item.quantity}x {item.title} <b>removido</b> do carrinho</span>
            );
            setLoading(true);
        });
    }, [props.cartProvider]);

    useEffect(() => {
        setSubtotal(
            cartItems.reduce(
                (total, item) => total + item.price * item.quantity,
                0.0
            )
        );
    }, [cartItems]);

    const handleFinalize = useCallback(() => {
        setModalVisible(false);
        setCartItems([]);
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
                    <div className="loading-spinner">
                        <LoadingSpinner fontSize={36} />
                    </div>
                ) : (
                    <CartItemList cartItems={cartItems} />
                )}
            </Content>
            <Footer className="footer">
                <div className="subtotal">
                    <Button
                        type="primary"
                        size="large"
                        onClick={() => setModalVisible(true)}
                        disabled={cartItems.length === 0}
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

            </Footer>
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
