document.addEventListener('DOMContentLoaded', () => {
  initializeListeners(); // Initialize all listeners on page load
  setupEscapeKeyListener(); // Set up listener for escape key to close modals
  setupHtmxSwapListener(); // Reinitialize listeners after HTMX swaps content
});

function initializeListeners() {
  initializeModalListeners();
  initializeNotificationListeners();
  initializeNavbarBurgerListeners();
}

function openModal($el) {
  $el.classList.add('is-active');
}

function closeModal($el) {
  $el.classList.remove('is-active');
  initializeNavbarBurgerListeners();
}

function closeAllModals() {
  (document.querySelectorAll('.modal') || []).forEach(($modal) => {
    closeModal($modal);
  });
}

function initializeModalListeners() {
  // Modal trigger buttons
  (document.querySelectorAll('.js-modal-trigger') || []).forEach(($trigger) => {
    const modal = $trigger.dataset.target;
    const $target = document.getElementById(modal);
    $trigger.addEventListener('click', () => openModal($target));
  });

  // Modal close buttons and other elements
  (document.querySelectorAll('.modal-background, .modal-close, .modal-card-head .delete, .modal-card-foot .button') || []).forEach(($close) => {
    const $target = $close.closest('.modal');
    $close.addEventListener('click', () => closeModal($target));
  });
}

// Functions to handle notifications
function initializeNotificationListeners() {
  (document.querySelectorAll('.notification .delete') || []).forEach(($delete) => {
    const $notification = $delete.parentNode;
    $delete.addEventListener('click', () => $notification.parentNode.removeChild($notification));
  });
}

// Functions to handle navbar burger toggle
function initializeNavbarBurgerListeners() {
  // Get all "navbar-burger" elements
  const $navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);

  // Add a click event on each of them
  $navbarBurgers.forEach( el => {
    el.addEventListener('click', () => {

      // Get the target from the "data-target" attribute
      const target = el.dataset.target;
      const $target = document.getElementById(target);

      // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
      el.classList.toggle('is-active');
      $target.classList.toggle('is-active');

    });
  });
}

// Event listener for escape key to close all modals
function setupEscapeKeyListener() {
  document.addEventListener('keydown', (event) => {
    if (event.key === "Escape") {
      closeAllModals();
    }
  });
}

// Reinitialize listeners after HTMX swaps content
function setupHtmxSwapListener() {
  document.body.addEventListener('htmx:afterSwap', () => {
    initializeListeners();
    initializeNavbarBurgerListeners()
  });
}
