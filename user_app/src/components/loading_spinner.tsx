import { LoadingOutlined } from "@ant-design/icons";
import { Spin } from "antd";

export const LoadingSpinner = () => {
  return <Spin indicator={<LoadingOutlined style={{ fontSize: 24 }} spin />} />;
};
