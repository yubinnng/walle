import { Col, Row, Input } from "antd";
import mermaid from "mermaid";
import { useState } from "react";
import yaml from "js-yaml";

const { TextArea } = Input;

const parseSpec = (spec) => {
  const workflow = yaml.load(spec);
  if (
    workflow === undefined ||
    workflow.tasks === undefined ||
    workflow.tasks.length === 0
  ) {
    return null;
  }
  return workflow;
};

const toMermaid = (workflow) => {
  if (workflow == null) {
    return null;
  }
  let graph = `graph\n`;
  try {
    // tasks that no other task depends on
    const tails = new Set();

    workflow.tasks.forEach((task) => {
      tails.add(task.name);
      if (task.depends === undefined || task.depends.length === 0) {
        graph += `Start((Start)) --> ${task.name}((${task.name}))\n`;
      } else {
        task.depends.forEach((dep) => {
          graph += `${dep}((${dep})) --> ${task.name}((${task.name}))\n`;
          // dep is not a tail
          tails.delete(dep);
        });
      }
    });
    // tails -> Done
    tails.forEach((tail) => {
      graph += `${tail}((${tail})) --> Done((Done))\n`;
    });
    return graph;
  } catch {
    console.log("failed to parse workflow spec");
    return null;
  }
};

const NewWorkflow = () => {
  const [spec, SetSpec] = useState();
  const [graph, setGraph] = useState();

  return (
    <>
      <p style={{ fontSize: "1.5rem" }}>Create a New Workflow</p>
      <Row gutter={16}>
        <Col span={14}>
          <TextArea
            autoSize={{ minRows: 30 }}
            valu={spec}
            onChange={(val) => {
              let newSpec = val.target.value;
              SetSpec(newSpec);
              let text = toMermaid(parseSpec(newSpec));
              if (text) {
                mermaid.render("workflow-graph", text, (graph) => {
                  setGraph(graph);
                });
              }
            }}
          />
        </Col>
        <Col span={10}>
          <Row justify="center">
            <div dangerouslySetInnerHTML={{ __html: graph }} />
          </Row>
        </Col>
      </Row>
    </>
  );
};

export default NewWorkflow;
