import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import reportWebVitals from "./reportWebVitals";
import { CartServiceCartProvider } from "src/service/cart_provider/cart_service_cart_provider";
import { StubCartProvider } from "src/service/cart_provider/stub_cart_provider";

const root = ReactDOM.createRoot(
    document.getElementById("root") as HTMLElement
);

const provider = process.env.REACT_APP_CART_SERVICE_URL
    ? new CartServiceCartProvider(process.env.REACT_APP_CART_SERVICE_URL)
    : new StubCartProvider();

root.render(
    <React.StrictMode>
        <App cartProvider={provider} />
    </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
