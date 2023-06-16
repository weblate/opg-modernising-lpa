describe('Sign', () => {
    describe('as an attorney', () => {
        beforeEach(() => {
            cy.visit('/testing-start?cookiesAccepted=1&redirect=/attorney-sign&completeLpa=1&asAttorney=1&asCertificateProvider=certified&loginAs=attorney');
        });

        it('can be signed', () => {
            cy.checkA11yApp();

            cy.contains('Sign as an attorney on this LPA');

            cy.contains('label', 'I, John Smith, confirm').click();
            cy.contains('button', 'Submit signature').click();

            cy.url().should('contain', '/attorney-what-happens-next');
            cy.checkA11yApp();

            cy.contains('h1', 'You’ve formally agreed to be an attorney');
        });

        it('shows an error when not selected', () => {
            cy.contains('button', 'Submit signature').click();

            cy.get('.govuk-error-summary').within(() => {
                cy.contains('You must select the box to sign as an attorney');
            });

            cy.contains('.govuk-form-group .govuk-error-message', 'You must select the box to sign as an attorney');
        });
    });

    describe('as a replacement attorney', () => {
        beforeEach(() => {
            cy.visit('/testing-start?cookiesAccepted=1&redirect=/attorney-sign&completeLpa=1&withReplacementAttorney=1&signedByDonor=1&asCertificateProvider=certified&asReplacementAttorney=1');
        });

        it('can be signed', () => {
            cy.checkA11yApp();

            cy.contains('Sign as a replacement attorney on this LPA');

            cy.contains('label', 'I, Jane Smith, confirm').click();
            cy.contains('button', 'Submit signature').click();

            cy.url().should('contain', '/attorney-what-happens-next');
            cy.checkA11yApp();

            cy.contains('h1', 'You’ve formally agreed to be a replacement attorney');
        });

        it('shows an error when not selected', () => {
            cy.contains('button', 'Submit signature').click();

            cy.get('.govuk-error-summary').within(() => {
                cy.contains('You must select the box to sign as a replacement attorney');
            });

            cy.contains('.govuk-form-group .govuk-error-message', 'You must select the box to sign as a replacement attorney');
        });
    });
});
