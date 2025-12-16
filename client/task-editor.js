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

/** @typedef {{
  start: number;
  end: number;
}} OffsetRange */

/** @param {OffsetRange} range */
function selectRange({start, end}) {
  const range = new Range();

  let current = 0;
  let found = 0;
  for (const node of nameInput.childNodes) {    
    const textNode = node.firstChild ?? node;
    const len = textNode.textContent.length;
    const currentEnd = current + len;

    if (start >= current && start <= currentEnd) {
      range.setStart(textNode, start - current);
      found++;
    }

    if (end >= current && end <= currentEnd) {
      range.setEnd(textNode, end - current);
      found++;
    }

    if (found === 2) {
      const sel = getSelection();
      sel.removeAllRanges();
      sel.addRange(range);

      return;
    }

    current = currentEnd;
  }
}

/**
  start will be greater than end for backwards selections
  @returns {OffsetRange | null}
*/
function selectionToRange() {
  const sel = getSelection();

  if (sel.type === "None" || !nameInput.contains(sel.anchorNode)) {
    return null;
  }

  const range = new Range();
  range.selectNodeContents(nameInput);

  range.setEnd(sel.anchorNode, sel.anchorOffset);
  const start = range.toString().length;

  range.setEnd(sel.focusNode, sel.focusOffset);
  const end = range.toString().length;

  return { start, end };
}

nameInput.addEventListener('input', () => {
  if (nameInput.textContent.trim() === "") {
    // Clear leftover elements that hide the placeholder
    clearNameInput();
    return;
  }

  const highlights = parseTaskName(nameInput.textContent);
  const range = new Range();

  const sel = selectionToRange();

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

  selectRange(sel);
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
