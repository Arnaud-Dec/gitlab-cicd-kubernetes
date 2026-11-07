describe('Connection status bar', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('shows the backend as connected, with its configured address in the same chip', () => {
    cy.contains('.chip', 'Backend: Connected')
      .should('be.visible')
      .and('contain.text', 'http://localhost:8080/api')
  })

  it('shows the database as connected', () => {
    cy.contains('.chip', 'Database: Connected').should('be.visible')
  })

  it('shows a warning icon and an "Insecure" tooltip for the plain-HTTP backend', () => {
    cy.get('.warning-icon').should('be.visible')
    cy.get('.help-icon').should('have.attr', 'title', 'Insecure (HTTP)')
  })
})
