import { getSwapper, getTemplate } from "./js/domutils.js";

const main = document.querySelector('main').shadowRoot;

const taskTempl = getTemplate(main.querySelector('template'));
function taskView(name) {
  const clone = taskTempl();
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

let projectId;
function setProjectId(pathname) {
  projectId = pathname.split('/')[1];
}
setProjectId(location.pathname);

/**
 * @arg {string} name
 */
function addTask(name) {
  fetch("/addTask", {
    method: "POST",
    body: JSON.stringify({
      name,
      projectId
    })
  });

  taskList.append(taskView(name));
}

const addButton = main.getElementById("add-button");
const addArea = main.getElementById("add-area");

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
