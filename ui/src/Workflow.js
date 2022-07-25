import { Card, Space, Button, Row, Col, Modal, message } from "antd";
import {
  DeleteOutlined,
  CaretRightOutlined,
  EditOutlined,
} from "@ant-design/icons";
import axios from "axios";
import { useEffect, useState } from "react";
import { Link, useParams, useNavigate } from "react-router-dom";
import mermaid from "mermaid";
import { formatDatetime, parseSpec, toMermaid } from "./utils";

const Workflow = () => {
  const navigate = useNavigate();
  let { name } = useParams();

  const [workflow, setWorkflow] = useState({
    name: "",
    spec: "",
    createdAt: "",
    updatedAt: "",
  });

  const [graph, setGraph] = useState();

  useEffect(() => {
    axios.get("/workflow/" + name).then((resp) => {
      setWorkflow(resp.data);
      setGraph(toMermaid(parseSpec(resp.data.spec)));
    });
  }, [name]);

  const [isModalVisible, setIsDeleteModalVisible] = useState(false);
  const showDeleteModal = () => {
    setIsDeleteModalVisible(true);
  };
  const handleDeleteCancel = () => {
    setIsDeleteModalVisible(false);
  };
  const handleDelete = () => {
    axios.delete("/workflow/" + name).then((resp) => {
      navigate("/");
      setIsDeleteModalVisible(false);
    });
  };

  const handleEdit = () => {
    navigate("/edit/" + workflow.name);
  };

  const Info = () => {
    return (
      <Card
        title={
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <div>Metadata</div>
            <Button
              type="primary"
              icon={<EditOutlined />}
              style={{ marginLeft: "10px" }}
              onClick={handleEdit}
            >
              Edit
            </Button>
          </div>
        }
        size="small"
      >
        <div>
          <p>URL: http://localhost:8080/{workflow.name}</p>
          <p>Created: {workflow.createdAt}</p>
          <p>Updated: {workflow.updatedAt}</p>
        </div>
      </Card>
    );
  };

  const Executions = () => {
    const [executions, setExecutions] = useState([]);
    const fetchExecutions = () => {
      axios.get("/execution/list?workflow_name=" + name).then((resp) => {
        setExecutions(resp.data);
      });
    };
    useEffect(() => {
      fetchExecutions();
      const inter = setInterval(fetchExecutions, 5000);
      return () => {
        clearInterval(inter);
      };
    }, [name]);

    const handleExec = () => {
      axios.post("/workflow/" + workflow.name + "/exec");
      message.success("Successfully executed");
    };

    return (
      <Card
        title={
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <div>Executions ({executions.length})</div>
            <Button
              type="primary"
              icon={<CaretRightOutlined />}
              style={{ marginLeft: "10px" }}
              onClick={handleExec}
            >
              Execute
            </Button>
          </div>
        }
        size="small"
      >
        <Col>
          <Space
            direction="vertical"
            size="middle"
            style={{
              display: "flex",
            }}
          >
            <Row>
              <Col span={8}>ID</Col>
              <Col span={4}>Status</Col>
              <Col span={6}>Started At</Col>
              <Col span={6}>Ended At</Col>
            </Row>
            {executions.map((exec, key) => {
              return (
                <Row key={key}>
                  <Col span={8}>
                    <Link to={"/exec/" + exec.id}>{exec.id}</Link>
                  </Col>
                  <Col span={4}>{exec.status}</Col>
                  <Col span={6}>{formatDatetime(exec.startAt)}</Col>
                  <Col span={6}>{formatDatetime(exec.endAt)}</Col>
                </Row>
              );
            })}
          </Space>
        </Col>
      </Card>
    );
  };

  mermaid.initialize({ startOnLoad: true });

  return (
    <>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          marginBottom: "20px",
        }}
      >
        <div style={{ fontSize: "1.5rem" }}>{workflow.name}</div>
        <Button
          type="primary"
          danger
          icon={<DeleteOutlined />}
          style={{ marginLeft: "25px" }}
          onClick={showDeleteModal}
          shape="round"
        />
        <Modal
          title="Delete Confirmation"
          visible={isModalVisible}
          onOk={handleDelete}
          onCancel={handleDeleteCancel}
        >
          <p>Are you sure to delete workflow {workflow.name}?</p>
        </Modal>
      </div>
      <Row gutter={[16]}>
        <Col span={16}>
          <Space direction="vertical" style={{ display: "flex" }}>
            <Info />
            <Executions />
          </Space>
        </Col>
        <Col span={8}>
          <Row justify="center">
            <div dangerouslySetInnerHTML={{ __html: graph }} />
          </Row>
        </Col>
      </Row>
    </>
  );
};

export default Workflow;
