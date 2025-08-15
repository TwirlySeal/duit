export const activeClass = "active";

/**
 * Get a function that swaps the "active" class
 *
 * @arg {Element} active
 * @returns {(el: Element) => void}
 */
export function getSwapper(active) {
  return (el) => {
    active.classList.remove(activeClass);
    el.classList.add(activeClass);
    active = el;
  };
}

export function getTemplate(el) {
  return () => el.content.cloneNode(true);
}

// export function render(element, data) {
//   for (const [key, value] of Object.entries(data)) {
//     element.querySelector(`[data-bind="${key}"]`).textContent = value;
//   }
//   return element;
// }

// export function getTemplate(template) {
//   return (data) => render(template.content.cloneNode(true), data);
// }
