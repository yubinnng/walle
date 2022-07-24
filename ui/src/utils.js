import yaml from "js-yaml";
import mermaid from "mermaid";

const parseSpec = (spec) => {
  try {
    return yaml.load(spec);
  } catch {
    return null;
  }
};

const toMermaid = (workflow) => {
  if (workflow == null) {
    return null;
  }
  let text = `graph\n`;
  try {
    // tasks that no other task depends on
    const tails = new Set();

    workflow.tasks.forEach((task) => {
      tails.add(task.name);
      if (task.depends === undefined || task.depends.length === 0) {
        text += `Start((Start)) --> ${task.name}((${task.name}))\n`;
      } else {
        task.depends.forEach((dep) => {
          text += `${dep}((${dep})) --> ${task.name}((${task.name}))\n`;
          // dep is not a tail
          tails.delete(dep);
        });
      }
    });
    // tails -> Done
    tails.forEach((tail) => {
      text += `${tail}((${tail})) --> Done((Done))\n`;
    });
    return mermaid.render("workflow-graph", text);
  } catch {
    console.log("invalid workflow spec");
    return null;
  }
};

var options = { year: "numeric", month: "long", day: "numeric" };
const formatDatetime = (date) => {};

export { parseSpec, toMermaid };
