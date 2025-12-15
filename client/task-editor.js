import { getSwapper } from "./js/domutils.js";
import { parseTaskName } from "./js/dates.js";
import { addTask } from "./tasks.js";

const showButton = document.getElementById("show-task-editor");
const taskEditor = document.getElementById("task-editor");
const nameInput = document.getElementById("name-input");

const activate = getSwapper(showButton);

showButton.addEventListener('click', () => {
  activate(taskEditor);
  nameInput.focus();
});

function clearNameInput() {
  nameInput.replaceChildren();
}
clearNameInput();

nameInput.addEventListener('input', () => {
  if (nameInput.textContent.trim() === "") {
    // Clear leftover elements that hide the placeholder
    clearNameInput();
    return;
  }

  const highlights = parseTaskName(nameInput.textContent);
  const range = new Range();

  // Clear spans and normalise text nodes
  nameInput.textContent = nameInput.textContent;
  const textNode = nameInput.firstChild;

  // Iterate in reverse so that offsets remain valid as spans are inserted
  for (let i = highlights.length - 1; i >= 0; i--) {
    const {start, end} = highlights[i];
    range.setStart(textNode, start);
    range.setEnd(textNode, end);
    range.surroundContents(document.createElement("span"));
  }
});

function submitTask() {
  addTask(nameInput.textContent);
  clearNameInput();
}

function cancelTask() {
  activate(showButton);
  clearNameInput();
}

taskEditor.addEventListener('submit', event => {
  event.preventDefault();
  submitTask();
});

const cancelButton = document.getElementById("cancel-button");
cancelButton.addEventListener('click', cancelTask);

nameInput.addEventListener('keydown', event => {
  switch (event.key) {
    case "Enter":
      submitTask();
      break;

    case "Escape":
      cancelTask();
      break;
  }
});
