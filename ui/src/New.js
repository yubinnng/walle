import { Col, Row, Input, Button } from "antd";
import { useState, useEffect } from "react";
import { parseSpec, toMermaid } from "./utils";
import axios from "axios";
import { SaveOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";

const { TextArea } = Input;

const New = () => {
  const [spec, setSpec] = useState();
  const [graph, setGraph] = useState();
  const navigate = useNavigate();

  const handleSpecChange = (val) => {
    setSpec(val.target.value);
  };

  useEffect(() => {
    setGraph(toMermaid(parseSpec(spec)));
  }, [spec]);

  const handleSubmit = () => {
    if (spec) {
      axios
        .post("/api/workflow", spec, {
          headers: {
            "content-type": "text/plain",
          },
        })
        .then((resp) => {
          navigate("/" + resp.data.name);
        })
        .catch((error) => {
          console.log(error);
          if (error.response) {
            alert(error.response.data.error);
          }
        });
    } else {
      alert("invalid specification");
    }
  };

  return (
    <>
      <p style={{ fontSize: "1.5rem" }}>Create a New Workflow</p>
      <Row gutter={16}>
        <Col span={16}>
          <Row gutter={[0, 10]}>
            <TextArea
              autoSize={{ minRows: 30 }}
              value={spec}
              onChange={handleSpecChange}
            />
            <Button
              type="primary"
              icon={<SaveOutlined />}
              size="large"
              onClick={handleSubmit}
            >
              Save
            </Button>
          </Row>
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

export default New;
