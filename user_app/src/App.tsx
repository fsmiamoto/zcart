import { useState, useCallback, useEffect } from "react";
import {
    Layout,
    message,
    Button,
    Statistic,
    Modal,
    Result,
} from "antd";
import { DollarCircleFilled, DownCircleFilled, ShoppingCartOutlined, UpCircleFilled } from "@ant-design/icons";
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
    const [checkedout, setCheckedout] = useState(false);

    useEffect(() => {
        if (!loading) return;
        props.cartProvider.ListCartItems().then((items) => {
            setCartItems(items);
            setLoading(false);
        });
    }, [props.cartProvider, loading]);

    useEffect(() => {
        props.cartProvider.OnAddProduct((item) => {
            message.info({
                icon: <UpCircleFilled style={{ fontSize: "1.2rem" }} />,
                content: <span>{item.quantity}x {item.title} <b>added</b> to the cart</span>,
                style: { fontSize: "1.2rem", marginTop: "5vh" }
            });
            setLoading(true);
        });

        props.cartProvider.OnRemoveProduct((item) => {
            message.info({
                icon: <DownCircleFilled style={{ fontSize: "1.2rem" }} />,
                content: <span>{item.quantity}x {item.title} <b>removed</b> from the cart</span>,
                style: { fontSize: "1.2rem", marginTop: "5vh" }
            });
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

    useEffect(() => {
        if (checkedout) {
            props.cartProvider.Checkout()
            return
        };
        setCartItems([])
        setLoading(true)
    }, [checkedout, props.cartProvider])

    const handleFinalize = useCallback(() => {
        setModalVisible(false);
        setCheckedout(true);
    }, []);

    if (checkedout) {
        return (
            <Result
                status="success"
                title="Thank you for you purchase!"
                extra={[
                    <Button type="primary" key="buy" onClick={() => setCheckedout(false)} style={{ fontSize: "1.5rem", paddingBottom: "45px" }}>Buy Again</Button>,
                ]}
            />
        )
    }

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
                        icon={<DollarCircleFilled />}
                        onClick={() => setModalVisible(true)}
                        shape={"round"}
                        style={{ fontSize: "1.5rem", paddingBottom: "45px" }}
                        disabled={cartItems.length === 0}
                    >
                        Checkout
                    </Button>
                    <span>
                        Subtotal:{" "}
                        <Statistic
                            value={subtotal}
                            prefix="R$"
                            precision={2}
                            decimalSeparator=","
                            groupSeparator="."
                        />
                    </span>
                </div>

            </Footer>
            <Modal
                visible={modalVisible}
                onOk={handleFinalize}
                onCancel={() => setModalVisible(false)}
                okText={"Checkout"}
                bodyStyle={{ fontSize: "1.25rem" }}
                okButtonProps={{ style: { fontSize: "1.25rem", paddingBottom: "35px" } }}
                cancelButtonProps={{ style: { fontSize: "1.25rem", paddingBottom: "35px" } }}
                cancelText={"Cancel"}
                closable={false}
            >
                Do you want to finish your purchase?
            </Modal>
        </Layout>
    );
}

export default App;
