describe('Task list', () => {
    beforeEach(() => {
        cy.visit('/testing-start?redirect=/task-list&lpa.complete=1&attorneyProvided=1&loginAs=attorney');
    });

    it('shows tasks', () => {
        cy.checkA11yApp();

        cy.contains('li', 'Confirm your details').should('contain', 'Not started');
        cy.contains('li', 'Read the LPA').should('contain', 'Not started');
        cy.contains('li', 'Sign the LPA').should('contain', 'Cannot start yet');
    });
});
