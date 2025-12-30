import { getTemplate } from "./js/domutils.js";
import { getProjectId } from "./nav.js";
import { parseDate, formatDate } from "./js/dates.js";
import { PlainDate, PlainTime } from "./js/dates.js";

const main = document.querySelector('main');
const taskTempl = getTemplate(main.querySelector('template'));

/**
  @param {{
    id: number;
    title: string;
    date?: PlainDate;
    time?: PlainTime;
  }} task
*/
function taskView(task) {
  const clone = taskTempl();
  clone.firstElementChild.dataset.id = task.id;
  clone.querySelector('p').textContent = task.title;

  /** @type {HTMLTimeElement} */
  const time = clone.querySelector('time');
  if (task.date) {
    time.dateTime = task.date.toString();
    time.textContent = formatDate(task);
  } else {
    time.remove();
  }

  return clone;
}

const heading = main.querySelector('h1');
const taskList = document.getElementById('task-list');

// Insert human-readable text for dates from server
for (const time of taskList.querySelectorAll('time')) {
  time.textContent = formatDate(parseDate(time.dateTime));
}

/**
 @param {number} id
 @param {string} name
*/
export async function showProject(id, name) {
  heading.textContent = name;
  document.title = name + " - Duit";

  /** @type {{
    id: number,
    title: string,
    date: string?,
    time: string?,
  }[]} */
  const data = await (await fetch("/api/projects/" + id)).json();
  taskList.replaceChildren(...data.map(task => {
    if (task.date) {
      Object.assign(task, parseDate(task.date));
    }
    return taskView(task);
  }));
}

/**
  @param {{
    title: string;
    date?: PlainDate;
    time?: PlainTime;
  }} task
*/
export async function addTask(task) {
  task.projectId = getProjectId();

  const { id } = await (await fetch("/api/tasks", {
    method: "POST",
    body: JSON.stringify(task)
  })).json();

  task.id = id;

  taskList.append(taskView(task));
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
