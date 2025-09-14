import { defaultInteraction, getNavigator } from "./js/nav.js";
import { getSwapper, activeClass } from "./js/domutils.js";
import { replaceTasks } from "./tasks.js";

const nav = document.querySelector('nav').shadowRoot.querySelector('ul');

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
function swapProject(id) {
  projectId = id;
  replaceTasks(id);
}

const navigate = getNavigator("drawer", projectId, (id) => {
  swapProject(id);
  activate(links.get(id));
});

nav.addEventListener("click", (event) => {
  if (defaultInteraction(event)) return;
  const li = event.target.closest("li");
  if (li === null) return;
  event.preventDefault();
  
  const id = parseInt(li.dataset.id);
  const link = li.firstElementChild;
  swapProject(id);
  activate(link);
  navigate(link.href, id);
});
