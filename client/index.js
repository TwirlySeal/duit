import { defaultInteraction, getNavigator } from "./js/nav.js";
import { getSwapper } from "./js/domutils.js";

const nav = document.body.firstElementChild.shadowRoot;
const {children} = nav;

const map = new Map();
for (let i = 1; i < children.length; i++) { // Skips <style> element
  const link = children[i];
  map.set(link.pathname, link);
}

const main = document.body.children[1];

/**
 * @arg {string} pathname
 * @typedef {{title: string, done: boolean}} Task
 */
async function getTasks(pathname) {
  /** @type {Task[]} */
  const tasks = await (await fetch("/data" + pathname)).json();
  main.replaceChildren(...tasks.map(t => {
    const p = document.createElement('p');
    p.textContent = t.title;
    return p;
  }));
}

// function renderTask(task) {
//   const update = {
//     title(v) {
//       p.textContent = v;
//     },

//     done(v) {}
//   };
// }

const activate = getSwapper(map.get(location.pathname));

const navigate = getNavigator("drawer", undefined, () => {
  const {pathname} = location;
  activate(map.get(pathname));
  getTasks(pathname);
});

nav.addEventListener("click", (event) => {
  if (defaultInteraction(event)) return;
  const link = event.target.closest("a");
  if (link === null) return;

  event.preventDefault();
  activate(link);
  getTasks(link.pathname);
  navigate(link.href);
});
