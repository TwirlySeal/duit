import { getTemplate } from "./js/domutils.js";
import { getProjectId } from "./nav.js";

const main = document.querySelector('main');

const taskTempl = getTemplate(main.querySelector('template'));
function taskView(title, id) {
  const clone = taskTempl();
  clone.firstElementChild.dataset.id = id;
  clone.querySelector('p').textContent = title;
  return clone;
}

const heading = main.querySelector('h1');
const taskList = document.getElementById('task-list');

// Insert human-readable text for dates from server
for (const li of taskList.children) {
  /** @type {HTMLTimeElement?} */
  const time = li.querySelector('time');
  if (time === null) continue;

  const date = new Date(time.dateTime);
  time.textContent = date.toLocaleString();
}

/**
 * @arg {number} id
 * @typedef {{title: string, done: boolean}} Task
 */
export async function showProject(id, name) {
  heading.textContent = name;
  document.title = name + " - Duit";
  /** @type {Task[]} */
  const data = await (await fetch("/api/projects/" + id)).json();
  taskList.replaceChildren(...data.map(t => taskView(t.title, t.id)));
}

/**
  @param {import("./js/dates.js").DatetimeExpr} datetime
*/
export async function addTask(title, datetime) {
  const body = {
    title,
    projectId: getProjectId()
  };

  if (datetime !== null) {
    body.date = datetime.date;

    if (datetime.time !== null) {
      body.time = datetime.time;
    }
  }

  const { id } = await (await fetch("/api/tasks", {
    method: "POST",
    body: JSON.stringify(body)
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
  if (li === null) return;

  completeTask(parseInt(li.dataset.id));
  li.remove();
});
