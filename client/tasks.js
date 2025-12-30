import { getTemplate } from "./js/domutils.js";
import { getProjectId } from "./nav.js";
import { parseDate, formatDate } from "./js/dates.js";
import { PlainDate, PlainTime } from "./js/dates.js";

const main = document.querySelector('main');
const taskTempl = getTemplate(main.querySelector('template'));

/** @typedef {{
  id: number;
  title: string;
  done: boolean;
  date: PlainDate?;
  time: PlainTime?;
}} Task */

/**
  @param {Task} task
*/
function taskView(task) {
  const clone = taskTempl();
  clone.firstElementChild.dataset.id = task.id;
  clone.querySelector('p').textContent = task.title;

  /** @type {HTMLTimeElement} */
  const time = clone.querySelector('time');
  if (task.date !== null) {
    time.dateTime = task.date.toString();
    time.textContent = formatDate(task.date);
  } else {
    time.remove();
  }

  return clone;
}

const heading = main.querySelector('h1');
const taskList = document.getElementById('task-list');

// Insert human-readable text for dates from server
for (const li of taskList.children) {
  /** @type {HTMLTimeElement?} */
  const time = li.querySelector('time');
  if (time === null) continue;

  const floatingDate = parseDate(time.dateTime);
  time.textContent = formatDate(floatingDate.date);
}

/**
 @param {number} id
 @param {string} name
*/
export async function showProject(id, name) {
  heading.textContent = name;
  document.title = name + " - Duit";

  /** @type {Task[]} */
  const data = await (await fetch("/api/projects/" + id)).json();
  taskList.replaceChildren(...data.map(task => {
    if (task.date) {
      Object.assign(task, parseDate(task.date));
    }
    return taskView(task);
  }));
}

/**
  @param {import("./js/dates.js").DatetimeExpr} datetime
*/
export async function addTask(title, datetime) {
  const task = {
    title,
    projectId: getProjectId()
  };

  if (datetime !== null) {
    task.date = datetime.date;

    if (datetime.time !== null) {
      task.time = datetime.time;
    }
  }

  const { id } = await (await fetch("/api/tasks", {
    method: "POST",
    body: JSON.stringify(task)
  })).json();

  taskList.append(taskView({
    ...task,
    id,
    date: datetime.date,
    time: datetime.time,
  }));
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
