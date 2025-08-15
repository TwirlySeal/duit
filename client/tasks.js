import { getSwapper, getTemplate } from "./js/domutils.js";

const main = document.querySelector('main').shadowRoot;

const taskTempl = getTemplate(main.getElementById('templ'));
function taskView(name) {
  const clone = taskTempl();
  clone.querySelector('p').textContent = name;
  return clone;
}

const tasks = main.getElementById('tasks');

/**
 * @arg {string} pathname
 * @typedef {{title: string, done: boolean}} Task
 */
export async function replaceTasks(pathname) {
  /** @type {Task[]} */
  const data = await (await fetch("/data" + pathname)).json();
  tasks.replaceChildren(...data.map(t => taskView(t.title)));
}

const addButton = main.getElementById("add-button");
const addArea = main.getElementById("add-area");

const activate = getSwapper(addButton);
addButton.addEventListener('click', () => {
  activate(addArea);
  addArea.focus();
});

addArea.addEventListener('keydown', event => {
  if (event.key === "Enter") {
    tasks.lastElementChild.before(taskView(addArea.value));
    addArea.value = "";
  } else if (event.key === 'Escape') {
    activate(addButton);
    addArea.value = "";
  }
});
