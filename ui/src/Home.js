import { Routes, Route, useNavigate } from "react-router-dom";
import { Layout, Menu } from "antd";
import { PlusSquareOutlined } from "@ant-design/icons";
import Edit from "./Edit";
import Workflow from "./Workflow";
import { useState } from "react";

const { Content, Footer, Sider } = Layout;
const NEW = "new";

const items = [NEW, "workflow-1", "workflow-2"].map((name) => ({
  key: name,
  icon: name === NEW ? <PlusSquareOutlined /> : null,
  label: name === NEW ? "New Workflow" : name,
}));

const LeftSider = () => {
  const navigate = useNavigate();
  const [selected, setSelected] = useState("new");

  const handleMenuClick = ({ key }) => {
    setSelected(key);
    navigate("/workflow/" + key);
  };

  return (
    <Sider
      theme="light"
      style={{
        overflow: "auto",
        height: "100vh",
        position: "fixed",
        left: 0,
        top: 0,
        bottom: 0,
      }}
    >
      <div
        style={{
          textAlign: "center",
          height: "32px",
          margin: "16px",
          // color: "#fff",
          fontSize: "1.2rem",
        }}
      >
        WALLE
      </div>
      <Menu
        mode="inline"
        defaultSelectedKeys={[NEW]}
        selectedKeys={[selected]}
        items={items}
        onClick={handleMenuClick}
      />
    </Sider>
  );
};

const Home = () => {
  return (
    <Layout hasSider style={{ minHeight: "100vh" }}>
      <LeftSider />
      <Layout
        style={{
          marginLeft: 200,
        }}
      >
        <Content
          style={{
            margin: "24px 16px 0",
            overflow: "initial",
          }}
        >
          <Routes>
            <Route path="/workflow/new" element={<Edit />} />
            <Route path="/workflow/:name" element={<Workflow />} />
          </Routes>
        </Content>
        <Footer
          style={{
            textAlign: "center",
          }}
        >
          Ant Design ©2018 Created by Ant UED
        </Footer>
      </Layout>
    </Layout>
  );
};

export default Home;
