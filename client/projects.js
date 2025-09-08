import { getSwapper, getTemplate } from "./js/domutils.js";

const map = new Map();

export function getProjectElementByPath(path) {
    return map.get(path);
}

export function initializeProjectNavigation(elements, navigate) {
    for (const {firstElementChild: link} of elements) {
        map.set(link.pathname, link);
    }
    activateProject = getSwapper(getProjectElementByPath(location.pathname))
    navigateFn = navigate;
}

export let activateProject;
let navigateFn;

const main = document.querySelector('nav').shadowRoot;

const projectTemplate = getTemplate(main.querySelector('template'));
function projectView(name, id, activate) {
    const clone = projectTemplate();
    let a = clone.querySelector('a');
    a.textContent = name;
    a.href = id;
    map.set(id, a);
    if (activate) {
        activateProject(a);
        navigateFn(id);
    }
    return clone;
}

const projectList = main.getElementById('project-list');

/**
 * @arg {string} name
 */
function addProject(name) {
    fetch("/addProject", {
        method: "POST",
        body: JSON.stringify({
            name,
        })
    }).then(response => {
        let id = response.headers.get("Location")
        projectList.append(projectView(name, id, true));
    });
}

let projectId;
export function setProjectId(pathname) {
    projectId = pathname.split('/')[1];
}
export function getProjectId() {
    return projectId;
}
setProjectId(location.pathname);

const addButton = main.getElementById("add-project-button");
const addArea = main.getElementById("add-project-area");

const activate = getSwapper(addButton);
addButton.addEventListener('click', () => {
    activate(addArea);
    addArea.focus();
});

addArea.addEventListener('keydown', event => {
    switch (event.key) {
        case "Enter":
            addProject(addArea.value);
            addArea.value = "";
            break;

        case "Escape":
            activate(addButton);
            addArea.value = "";
            break;
    }
});