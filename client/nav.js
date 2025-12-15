import { defaultInteraction, getNavigator } from "./js/nav.js";
import { getSwapper, activeClass } from "./js/domutils.js";
import { showProject } from "./tasks.js";

const nav = document.getElementById("project-list");

const links = new Map();
let activeLi, projectId;
for (const li of nav.children) {
  const id = parseInt(li.dataset.id);
  const link = li.firstElementChild;
  links.set(id, link);

  if (link.classList.contains(activeClass)) {
    activeLi = li;
    projectId = id;
  }
}

export const getProjectId = () => projectId;

const activate = getSwapper(activeLi.firstElementChild);

/** @arg {number} id */
function swapProject(id, link) {
  projectId = id;
  activate(link);
  showProject(id, link.textContent);
}

const navigate = getNavigator("drawer", projectId, (id) => {
  swapProject(id, links.get(id));
});

nav.addEventListener("click", (event) => {
  if (defaultInteraction(event)) return;
  const li = event.target.closest("li");
  if (li === null) return;
  event.preventDefault();
  
  const id = parseInt(li.dataset.id);
  const link = li.firstElementChild;
  swapProject(id, link);
  navigate(link.href, id);
});
