import { getSwapper, getTemplate } from "./js/domutils.js";
import { getProjectId } from "./index.js"

const main = document.querySelector('main').shadowRoot;

const taskTempl = getTemplate(main.querySelector('template'));
function taskView(title, id) {
  const clone = taskTempl();
  clone.firstElementChild.dataset.id = id;
  clone.querySelector('p').textContent = title;
  return clone;
}

const taskList = main.getElementById('task-list');

/**
 * @arg {number} projectId
 * @typedef {{title: string, done: boolean}} Task
 */
export async function replaceTasks(projectId) {
  /** @type {Task[]} */
  const data = await (await fetch("/api/projects/" + projectId)).json();
  taskList.replaceChildren(...data.map(t => taskView(t.title)));
}

async function addTask(title) {
  const { id } = await (await fetch("/api/tasks", {
    method: "POST",
    body: JSON.stringify({
      title,
      projectId: getProjectId()
    })
  })).json();

  taskList.append(taskView(title, id));
}

/** @arg {number} id */
function completeTask(id){
  fetch("/api/tasks", {
    method: "PATCH",
    body: JSON.stringify({
      id,
      done: true
    })
  });
}

taskList.addEventListener('click', (event) => {
  if (event.target.tagName != 'BUTTON') {
    return
  }

  const li = event.target.closest("li");
  if (li == null) return;

  completeTask(parseInt(li.dataset.id));
  li.remove();
});

const addButton = main.getElementById("add-button");
const addArea = main.getElementById("add-area");
const checkButton = main.querySelector(".check");

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
