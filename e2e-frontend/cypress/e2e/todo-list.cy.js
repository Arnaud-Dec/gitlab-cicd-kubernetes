function uniqueTitle() {
  return `Cypress todo ${Date.now()}-${Math.floor(Math.random() * 10000)}`
}

describe('Todo list', () => {
  let title

  beforeEach(() => {
    title = uniqueTitle()
    cy.visit('/')
  })

  afterEach(() => {
    cy.deleteTodoByTitle(title)
  })

  it('shows the app title', () => {
    cy.contains('h1', 'Todo List').should('be.visible')
  })

  it('disables the Add button until a title is entered', () => {
    cy.contains('button', 'Add').should('be.disabled')

    cy.get('input[aria-label="Todo title"]').type('Temporary')
    cy.contains('button', 'Add').should('not.be.disabled')

    cy.get('input[aria-label="Todo title"]').clear()
    cy.contains('button', 'Add').should('be.disabled')
  })

  it('creates a todo and shows it in the list', () => {
    cy.addTodo(title, 'created by Cypress')

    cy.contains('.todo-item', title)
      .should('be.visible')
      .find('.description')
      .should('have.text', 'created by Cypress')
  })

  it('clears the form after adding a todo', () => {
    cy.addTodo(title, 'some notes')

    cy.get('input[aria-label="Todo title"]').should('have.value', '')
    cy.get('input[aria-label="Todo description"]').should('have.value', '')
  })

  it('toggles a todo as completed and back', () => {
    cy.addTodo(title)

    cy.contains('.todo-item', title).should('not.have.class', 'completed')

    cy.contains('.todo-item', title).find('input[type="checkbox"]').click()
    cy.contains('.todo-item', title).should('have.class', 'completed')

    cy.contains('.todo-item', title).find('input[type="checkbox"]').click()
    cy.contains('.todo-item', title).should('not.have.class', 'completed')
  })

  it('deletes a todo', () => {
    cy.addTodo(title)
    cy.contains('.todo-item', title).should('exist')

    cy.deleteTodoByTitle(title)
    cy.contains('.todo-item', title).should('not.exist')
  })
})
