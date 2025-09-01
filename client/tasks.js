import { getSwapper, getTemplate } from "./js/domutils.js";
import { setProjectId, getProjectId } from "./projects.js"

const main = document.querySelector('main').shadowRoot;

const taskTemplate = getTemplate(main.querySelector('template'));
function taskView(name) {
  const clone = taskTemplate();
  clone.querySelector('p').textContent = name;
  return clone;
}

const taskList = main.getElementById('task-list');

/**
 * @arg {string} pathname
 * @typedef {{title: string, done: boolean}} Task
 */
export async function replaceTasks(pathname) {
  /** @type {Task[]} */
  setProjectId(pathname);
  const data = await (await fetch("/data" + pathname)).json();
  taskList.replaceChildren(...data.map(t => taskView(t.title)));
}

/**
 * @arg {string} name
 */
function addTask(name) {
  fetch("/addTask", {
    method: "POST",
    body: JSON.stringify({
      name,
      projectId: getProjectId(),
    })
  });

  taskList.append(taskView(name));
}

const addButton = main.getElementById("add-task-button");
const addArea = main.getElementById("add-task-area");

const activate = getSwapper(addButton);
addButton.addEventListener('click', () => {
  activate(addArea);
  addArea.focus();
});

addArea.addEventListener('keydown', event => {
  switch (event.key) {
    case "Enter":
      addTask(addArea.value);
      addArea.value = "";
      break;

    case "Escape":
      activate(addButton);
      addArea.value = "";
      break;
  }
});
