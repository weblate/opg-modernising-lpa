import { TestEmail } from "../../support/e2e";

describe('Choose attorneys summary', () => {
    beforeEach(() => {
        cy.visit('/testing-start?redirect=/choose-attorneys-summary&lpa.attorneys=-1&cookiesAccepted=1');
    });

    it('multiple attorneys details are listed', () => {
        cy.checkA11yApp();

        cy.contains('You have added 2 attorneys');

        cy.contains('Jessie Jones');
        cy.contains('2 January 2000');
        cy.contains('2 RICHMOND PLACE');
        cy.contains('B14 7ED');

        cy.contains('Robin Redcar');
        cy.contains('2 January 2000');

        cy.visitLpa('/task-list')
        cy.contains('a', 'Choose your attorneys').parent().parent().contains('In progress (2)')
    });

    it('can amend attorney details', () => {
        cy.checkA11yApp();

        cy.get('#attorney-name-1').contains('a', 'Change').click();

        cy.url().should('contain', '/choose-attorneys?from=');
        cy.url().should('contain', '/choose-attorneys-summary');
        cy.url().should('match', /id=\w*/);

        cy.get('#f-first-names').clear().type('Mark');

        cy.contains('button', 'Save and continue').click();

        cy.url().should('contain', '/choose-attorneys-summary');

        cy.contains('Mark Jones');
    });

    it('can amend attorney address', () => {
        cy.checkA11yApp();

        cy.get('#attorney-address-2').contains('a', 'Change').click();

        cy.url().should('contain', '/choose-attorneys-address?from=');
        cy.url().should('contain', '/choose-attorneys-summary');
        cy.url().should('match', /id=\w*/);

        cy.contains('label', 'Enter a new address').click();
        cy.contains('button', 'Continue').click();

        cy.get('#f-lookup-postcode').type('B14 7ED');
        cy.contains('button', 'Find address').click();

        cy.get('#f-select-address').select('1 RICHMOND PLACE, BIRMINGHAM, B14 7ED');
        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys-address');
        cy.get('#f-address-line-1').should('have.value', '1 RICHMOND PLACE');

        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys-summary');

        cy.contains('dd', '1 RICHMOND PLACE');

        cy.visitLpa('/task-list')
        cy.contains('a', 'Choose your attorneys').parent().parent().contains('Completed (2)')
    });

    it('can add another attorney from summary page', () => {
        cy.checkA11yApp();

        cy.get('input[name="yes-no"]').check('yes');
        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys');

        cy.get('#f-first-names').clear().type('Bob Arnold');
        cy.get('#f-last-name').clear().type('Jones');
        cy.get('#f-email').clear().type(TestEmail);
        cy.get('input[name="date-of-birth-day"]').clear().type('31');
        cy.get('input[name="date-of-birth-month"]').clear().type('12');
        cy.get('input[name="date-of-birth-year"]').clear().type('1995');
        cy.contains('button', 'Save and continue').click();

        cy.contains('label', 'Enter a new address').click();
        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys-address');

        cy.get('#f-lookup-postcode').type('B14 7ED');
        cy.contains('button', 'Find address').click();

        cy.get('#f-select-address').select('5 RICHMOND PLACE, BIRMINGHAM, B14 7ED');
        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys-address');
        cy.get('#f-address-line-1').should('have.value', '5 RICHMOND PLACE');
        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys-summary');

        cy.contains('Bob Arnold Jones');
        cy.contains('31 December 1995');
        cy.contains('5 RICHMOND PLACE');
        cy.contains('B14 7ED');
    });

    it('can remove an attorney', () => {
        cy.checkA11yApp();

        cy.get('#remove-attorney-1').contains('a', 'Remove').click();

        cy.url().should('contain', '/remove-attorney');
        cy.url().should('match', /id=\w*/);

        cy.checkA11yApp();

        cy.contains('Are you sure you want to remove Jessie Jones?');

        cy.get('input[name="yes-no"]').check('yes');
        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys-summary');

        cy.checkA11yApp();

        cy.get('main').should('not.contain', 'Jessie Jones');

        cy.get('#remove-attorney-1').contains('a', 'Remove').click();
        cy.get('input[name="yes-no"]').check('yes');
        cy.contains('button', 'Continue').click();

        cy.url().should('contain', '/choose-attorneys');
    });

    it('errors when remove not selected', () => {
        cy.get('#remove-attorney-1').contains('a', 'Remove').click();

        cy.contains('button', 'Continue').click();

        cy.get('.govuk-error-summary').within(() => {
            cy.contains('Select yes to remove the attorney');
        });

        cy.contains('.govuk-fieldset .govuk-error-message', 'Select yes to remove the attorney');
    });
});
