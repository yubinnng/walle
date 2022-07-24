import { Routes, Route, useNavigate } from "react-router-dom";
import { Layout, Menu } from "antd";
import { PlusSquareOutlined } from "@ant-design/icons";
import Edit from "./Edit";
import Workflow from "./Workflow";
import { useEffect, useState } from "react";
import axios from "axios";
import { store } from "./store";
import Execution from "./Execution";

const { Content, Footer, Sider } = Layout;
const NEW = "new";

const LeftSider = () => {
  const navigate = useNavigate();
  const [menuItems, setMenuItems] = useState([]);
  const { selected } = store;

  const handleMenuClick = ({ key }) => {
    store.selected = key;
  };

  useEffect(() => {
    navigate("/" + selected);
  }, [selected]);

  useEffect(() => {
    axios.get("/workflow/list").then((resp) => {
      let names = [NEW].concat(resp.data);
      let items = names.map((name) => ({
        key: name,
        icon: name === NEW ? <PlusSquareOutlined /> : null,
        label: name === NEW ? "New Workflow" : name,
      }));
      setMenuItems(items);
    });
  }, [selected]);

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
        selectedKeys={[selected]}
        items={menuItems}
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
            <Route index element={<Edit />} />
            <Route path="/new" element={<Edit />} />
            <Route path="/exec/:id" element={<Execution />} />
            <Route path="/:name" element={<Workflow />} />
          </Routes>
        </Content>
        <Footer
          style={{
            textAlign: "center",
          }}
        >
          Walle Serverless Workflow @2022
        </Footer>
      </Layout>
    </Layout>
  );
};

export default Home;
