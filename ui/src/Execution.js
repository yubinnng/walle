import { Card, Space, Row, Col, Button, Popover } from "antd";
import { LeftOutlined } from "@ant-design/icons";
import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { formatDatetime, parseSpec, toMermaid } from "./utils";

const Execution = () => {
  let { id } = useParams();
  const navigate = useNavigate();

  const [execution, setExecution] = useState({
    id: "",
    workflowName: "",
    status: "",
    tasks: [],
    startAt: "",
    endAt: "",
  });
  const [graph, setGraph] = useState();

  useEffect(() => {
    axios.get("/api/execution/" + id).then((resp) => {
      let execution = resp.data;
      setExecution(execution);
      axios.get("/api/workflow/" + execution.workflowName).then((resp) => {
        setGraph(toMermaid(parseSpec(resp.data.spec)));
      });
    });
  }, []);

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
          </div>
        }
        size="small"
      >
        <Space
          direction="vertical"
          size="middle"
          style={{
            display: "flex",
          }}
        >
          <Row>
            <Col span={12}>Status</Col>
            <Col span={6}>Started At</Col>
            <Col span={6}>Ended At</Col>
          </Row>
          <Row>
            <Col span={12}>{execution.status}</Col>
            <Col span={6}>{formatDatetime(execution.startAt)}</Col>
            <Col span={6}>{formatDatetime(execution.endAt)}</Col>
          </Row>
        </Space>
      </Card>
    );
  };

  const Execution = () => {
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
            <div>Tasks</div>
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
              <Col span={4}>Name</Col>
              <Col span={4}>Status</Col>
              <Col span={8}>Started At</Col>
              <Col span={8}>Updated At</Col>
            </Row>
            {execution.tasks.map((task, key) => {
              return (
                <Row key={key}>
                  <Col span={4}>
                    <Popover content={task.log} title="Task Log">
                      <Button type="primary">{task.name}</Button>
                    </Popover>
                  </Col>
                  <Col span={4}>{task.status}</Col>
                  <Col span={8}>{formatDatetime(task.startedAt)}</Col>
                  <Col span={8}>{formatDatetime(task.updatedAt)}</Col>
                </Row>
              );
            })}
          </Space>
        </Col>
      </Card>
    );
  };

  const handleBack = () => {
    navigate(-1);
  };

  return (
    <>
      <div
        style={{ display: "flex", alignItems: "center", marginBottom: "20px" }}
      >
        <Button icon={<LeftOutlined />} size="large" onClick={handleBack} />
        <div style={{ fontSize: "1.5rem", marginLeft: "25px" }}>
          {execution.id}
        </div>
      </div>
      <Row gutter={[16]}>
        <Col span={16}>
          <Space direction="vertical" style={{ display: "flex" }}>
            <Info />
            <Execution />
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

export default Execution;
