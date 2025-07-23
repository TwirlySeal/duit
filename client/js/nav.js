let firstNavigation = true;
const navigators = new Map(); // could add cleanup callback

/**
 * @template T
 * @arg {T} initialDetail
 * @arg {(detail: T) => void} action
 */
export function getNavigator(id, initialDetail, action) {
  navigators.set(id, action);

  /**
   * @arg {string} url
   * @arg {T} detail
   */
  return (url, detail) => {
    if (firstNavigation) {
      history.replaceState({id, detail: initialDetail}, "");
      firstNavigation = false;
    }

    history.pushState({id, detail}, "", url);
  };
}

// Call stored action
addEventListener("popstate", e => {
  const { id, detail } = e.state;
  navigators.get(id)(detail);
});

/** @arg {Event} event */
export function defaultInteraction(event) {
  return (
    event.button !== 0
    || event.shiftKey
    || event.altKey
    || event.ctrlKey
    || event.metaKey
  );
}


// import { activeClass, getSwapper } from "./domutils.js";

/**
 * Client-side links that swap the "active" class
 *
 * @arg {string} id
 * @arg {Element} container
 * @arg {string} selector - A CSS selector for event delegation
 * @arg {HTMLAnchorElement[]} anchors
 */
// export function navlinks(id, container, selector, ...anchors) {
//   let active, i = 0;
//   for (; i < anchors.length; i++) {
//     active = anchors[i];

//     if (active.classList.contains(activeClass)) {
//       break;
//     }
//   }

//   const activate = getSwapper(active);

//   const navigate = getNavigator(id, activeIndex, (index) => {
//     const anchor = anchors[index];
//     /** @todo call handler for ui update */
//     console.log(index);
//     activate(anchor);
//   });

//   container.addEventListener("click", (event) => {
//     console.log("clicek");
//     if (defaultInteraction(event)) return;

//     const { target } = event;
//     const anchor = target.closest(selector);
//     if (anchor === null) return;

//     event.preventDefault();
//     navigate(anchor.href, anchor[linkIndex]); // change
//   });
// }


//   for (const a of anchors) {
//     a.addEventListener("click", e => {
//       e.preventDefault();
//       navigator(a.href);
//     });
//   }



// const getFragmentURL = (route) => '/html' + new URL(route).pathname;

// function updateUI(fragmentURL) {
//   fetch(fragmentURL)
//     .then(r => r.text())
//     .then(html => handler(html));
// }
