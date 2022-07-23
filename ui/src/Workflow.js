import { Card, Space } from "antd";
import { Routes, Route, useParams } from 'react-router-dom';

const Workflow = () => {
  let { name } = useParams();

  return (
    <Space
      direction="vertical"
      size="middle"
      style={{
        display: "flex",
      }}
    >
      <div style={{ fontSize: "1.5rem" }}>{name}</div>
      <Card title="Info" size="small">
        <div>
          <p>URL: http://localhost:8080/example-workflow</p>
          <p>Created: 28 July 2022</p>
          <p>Updated: 28 July 2022</p>
        </div>
      </Card>
      {/* <Card title="Triggers" size="small">
        <p>Card content</p>
        <p>Card content</p>
      </Card> */}
      <Card title="Executions (2)" size="small">
        <p>Card content</p>
        <p>Card content</p>
      </Card>
    </Space>
  );
};

export default Workflow;
