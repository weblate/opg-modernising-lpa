describe('Date of birth', () => {
    beforeEach(() => {
        cy.visit('/testing-start?withAttorney=1&asAttorney=1&redirect=/attorney-date-of-birth&loginAs=attorney');
    });

    it('can be completed', () => {
        cy.checkA11yApp();

        cy.get('#f-date-of-birth').type('1');
        cy.get('#f-date-of-birth-month').type('2');
        cy.get('#f-date-of-birth-year').type('1990');

        cy.contains('button', 'Continue').click();
        cy.url().should('contain', '/attorney-mobile-number');
    });

    it('can be over 100', () => {
        cy.checkA11yApp();

        cy.get('#f-date-of-birth').type('1');
        cy.get('#f-date-of-birth-month').type('2');
        cy.get('#f-date-of-birth-year').type('1900');

        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/attorney-date-of-birth');
        cy.contains('By saving this section, you confirm that the person is more than 100 years old')

        cy.checkA11yApp();

        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/attorney-mobile-number');
    });

    it('errors when empty', () => {
        cy.contains('button', 'Continue').click();

        cy.get('.govuk-error-summary').within(() => {
            cy.contains('Enter your date of birth');
        });

        cy.contains('#date-of-birth-hint + .govuk-error-message', 'Enter your date of birth');
    });

    it('errors when invalid dates of birth', () => {
        cy.get('#f-date-of-birth').type('1');
        cy.contains('button', 'Continue').click();
        cy.contains('#date-of-birth-hint + .govuk-error-message', 'Your date of birth must include a month and year');

        cy.get('#f-date-of-birth-month').type('2');
        cy.get('#f-date-of-birth-year').type('2222');
        cy.contains('button', 'Continue').click();
        cy.contains('#date-of-birth-hint + .govuk-error-message', 'Your date of birth must be in the past');

        cy.get('#f-date-of-birth').type('not');
        cy.get('#f-date-of-birth-month').type('valid');
        cy.get('#f-date-of-birth-year').clear().type('values');
        cy.contains('button', 'Continue').click();
        cy.contains('#date-of-birth-hint + .govuk-error-message', 'Your date of birth must be a real date');
    });

    it('errors when not over 18', () => {
        const lastYear = (new Date(new Date().setFullYear(new Date().getFullYear() - 1))).getFullYear();

        cy.get('#f-date-of-birth').type('1');
        cy.get('#f-date-of-birth-month').type('2');
        cy.get('#f-date-of-birth-year').type(lastYear.toString());
        cy.contains('button', 'Continue').click();

        cy.get('.govuk-error-summary').within(() => {
            cy.contains('You’ve entered a date of birth that means you are under 18, you must be 18 to be an attorney');
        });

        cy.contains('#date-of-birth-hint + .govuk-error-message', 'You’ve entered a date of birth that means you are under 18, you must be 18 to be an attorney');
    });
});
