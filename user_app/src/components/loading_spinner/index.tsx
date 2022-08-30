import { LoadingOutlined } from "@ant-design/icons";
import { Spin } from "antd";

export interface LoadingSpinnerProps {
    fontSize?: number
}

export const LoadingSpinner = (props: LoadingSpinnerProps) => {
    const { fontSize } = props;
    return <Spin indicator={<LoadingOutlined style={{ fontSize: fontSize ?? 24 }} spin />
    } />;
};
