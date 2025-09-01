import { defaultInteraction, getNavigator } from "./js/nav.js";
import { getSwapper } from "./js/domutils.js";
import { replaceTasks } from "./tasks.js";
import {getProjectElementByPath, initializeProjectNavigation} from "./projects.js";

const nav = document.querySelector('nav').shadowRoot.querySelector('ul');
const {children} = nav;
initializeProjectNavigation(children);

const activate = getSwapper(getProjectElementByPath(location.pathname));

const navigate = getNavigator("drawer", undefined, () => {
  const {pathname} = location;
  activate(getProjectElementByPath(pathname));
  replaceTasks(pathname);
});

nav.addEventListener("click", (event) => {
  if (defaultInteraction(event)) return;
  const link = event.target.closest("a");
  if (link === null) return;

  event.preventDefault();
  activate(link);
  replaceTasks(link.pathname);
  navigate(link.href);
});
