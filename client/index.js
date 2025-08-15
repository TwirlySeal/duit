import { defaultInteraction, getNavigator } from "./js/nav.js";
import { getSwapper } from "./js/domutils.js";
import { replaceTasks } from "./tasks.js";

const nav = document.body.firstElementChild.shadowRoot;
const {children} = nav;

const map = new Map(); // associate anchor elements to URL paths
for (let i = 1; i < children.length; i++) { // Skips <style> element
  const link = children[i];
  map.set(link.pathname, link);
}

const activate = getSwapper(map.get(location.pathname));

const navigate = getNavigator("drawer", undefined, () => {
  const {pathname} = location;
  activate(map.get(pathname));
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
