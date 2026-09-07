// Fills in and submits the add-todo form.
Cypress.Commands.add('addTodo', (title, description = '') => {
  cy.get('input[aria-label="Todo title"]').clear().type(title)
  if (description) {
    cy.get('input[aria-label="Todo description"]').clear().type(description)
  }
  cy.contains('button', 'Add').click()
})

// Deletes the todo with the given title if it's currently in the list.
// Safe to call even if it isn't (e.g. a test already deleted it itself) —
// used from afterEach hooks to clean up test data unconditionally.
Cypress.Commands.add('deleteTodoByTitle', (title) => {
  cy.get('body').then(($body) => {
    const exists = [...$body.find('.todo-item')].some((el) => el.innerText.includes(title))
    if (exists) {
      cy.contains('.todo-item', title).find('button[aria-label="Delete todo"]').click()
    }
  })
})
