import { defaultInteraction, getNavigator } from "./js/nav.js";
import { replaceTasks } from "./tasks.js";
import {activateProject, getProjectElementByPath, initializeProjectNavigation} from "./projects.js";

const nav = document.querySelector('nav').shadowRoot.querySelector('ul');
const {children} = nav;

const navigate = getNavigator("drawer", undefined, () => {
  const {pathname} = location;
  activateProject(getProjectElementByPath(pathname));
  replaceTasks(pathname);
});

initializeProjectNavigation(children, navigate);

nav.addEventListener("click", (event) => {
  if (defaultInteraction(event)) return;
  const link = event.target.closest("a");
  if (link === null) return;

  event.preventDefault();
  activateProject(link);
  replaceTasks(link.pathname);
  navigate(link.href);
});
