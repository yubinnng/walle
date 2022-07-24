import { Col, Row, Input, Button } from "antd";
import mermaid from "mermaid";
import { useState } from "react";
import { parseSpec, toMermaid } from "./utils";

import axios from "axios";

import { SaveOutlined } from "@ant-design/icons";
import { store } from "./store";

const { TextArea } = Input;

const Edit = () => {
  const [spec, setSpec] = useState();
  const [workflow, setWorkflow] = useState();
  const [graph, setGraph] = useState();

  const handleSpecChange = (val) => {
    let newSpec = val.target.value;
    let wf = parseSpec(newSpec);
    let text = toMermaid(wf);
    if (text) {
      mermaid.render("workflow-graph", text, (graph) => {
        setGraph(graph);
      });
    }
    setSpec(newSpec);
    setWorkflow(wf);
  };

  const handleSubmit = () => {
    if (workflow) {
      axios
        .post("/workflow", spec, {
          headers: {
            "content-type": "text/plain",
          },
        })
        .then((resp) => {
          store.selected = workflow.name;
        })
        .catch((resp) => {
          console.log(resp);
          alert(resp);
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
              valu={spec}
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

export default Edit;
