import yaml from "js-yaml";

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

var options = { year: "numeric", month: "long", day: "numeric" };
const formatDatetime = (date) => {};

export { parseSpec, toMermaid };
