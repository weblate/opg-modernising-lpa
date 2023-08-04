describe('Previous application number', () => {
    it('can be submitted', () => {
        cy.visit('/testing-start?redirect=/previous-application-number&lpa.yourDetails=1');
        cy.checkA11yApp();

        cy.get('#f-previous-application-number').type('ABC');

        cy.contains('button', 'Continue').click();
        cy.url().should('contain', '/task-list');
    });

    it('errors when unselected', () => {
        cy.visit('/testing-start?redirect=/previous-application-number');

        cy.contains('button', 'Continue').click();

        cy.get('.govuk-error-summary').within(() => {
            cy.contains('Enter previousApplicationNumber');
        });

        cy.contains('.govuk-fieldset .govuk-error-message', 'Enter previousApplicationNumber');
    });
});