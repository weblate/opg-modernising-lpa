const { TestMobile } = require("../../support/e2e");

describe('As a trust corporation', () => {
    beforeEach(() => {
        cy.visit('/testing-start?redirect=/attorney-start&lpa.complete=1&lpa.trustCorporation=complete&useTestShareCode=1&sendAttorneyShare=1&lpa.signedByDonor=1&asCertificateProvider=certified');

        // start
        cy.contains('a', 'Start').click();

        // enter reference number
        cy.get('#f-reference-number').type('abcdef123456');
        cy.contains('button', 'Continue').click();

        // acting as an attorney
        cy.contains('We have identified the trust corporation’s attorney reference number');
        cy.contains('a', 'Continue').click();

        // task list
        cy.contains('a', 'Confirm your details').click();

        // mobile number
        cy.get('#f-mobile').type(TestMobile);
        cy.contains('button', 'Continue').click();

        // confirm your company details
        cy.contains(TestMobile);
        cy.contains('Confirm your company details');
        cy.contains('My company');
        cy.contains('555555555');
        cy.contains('simulate-delivered@notifications.service.gov.uk');
        cy.contains('123 Fake Street');
        cy.contains('FF1 1FF');
        cy.contains('button', 'Continue').click();

        // read the lpa
        cy.contains('Read Jamie Smith’s LPA');
        cy.contains('button', 'Continue').click();

        // legal rights and responsibilities
        cy.contains('Before signing, you must read the trust corporation’s legal rights and responsibilities as an attorney.');
        cy.contains('a', 'Continue').click();

        // what happens when you sign the lpa
        cy.contains('What happens when you sign the LPA');
        cy.contains('a', 'Continue to signing page').click();
    });

    it('allows a single signatory', () => {
        // sign
        cy.contains('Sign the LPA on behalf of the trust corporation');
        cy.get('#f-first-names').type('Sign');
        cy.get('#f-last-name').type('Signson');
        cy.get('#f-professional-title').type('Pro signer');
        cy.get('#f-confirm').check();
        cy.contains('button', 'Submit signature').click();

        // would like a 2nd signatory
        cy.contains('label', 'No').click();
        cy.contains('button', 'Continue').click();

        // what happens next
        cy.contains('My company has formally agreed to be an attorney');
        cy.contains('a', 'Go to your dashboard');
    });

    it('allows a second signatory', () => {
        // sign
        cy.contains('Sign the LPA on behalf of the trust corporation');
        cy.get('#f-first-names').type('Sign');
        cy.get('#f-last-name').type('Signson');
        cy.get('#f-professional-title').type('Pro signer');
        cy.get('#f-confirm').check();
        cy.contains('button', 'Submit signature').click();

        // would like a 2nd signatory
        cy.contains('label', 'Yes').click();
        cy.contains('button', 'Continue').click();

        // task list
        cy.contains('a', 'Return to task list').click();
        cy.contains('Sign the LPA (signatory 1)');
        cy.contains('Sign the LPA (signatory 2)').click();

        // sign
        cy.get('#f-first-names').type('Sign2');
        cy.get('#f-last-name').type('Signson2');
        cy.get('#f-professional-title').type('Pro signer2');
        cy.get('#f-confirm').check();
        cy.contains('button', 'Submit signature').click();

        // what happens next
        cy.contains('My company has formally agreed to be an attorney');
        cy.contains('a', 'Go to your dashboard');
    });

    it('can remove second signatory', () => {
        // sign
        cy.contains('Sign the LPA on behalf of the trust corporation');
        cy.get('#f-first-names').type('Sign');
        cy.get('#f-last-name').type('Signson');
        cy.get('#f-professional-title').type('Pro signer');
        cy.get('#f-confirm').check();
        cy.contains('button', 'Submit signature').click();

        // would like a 2nd signatory
        cy.contains('label', 'Yes').click();
        cy.contains('button', 'Continue').click();

        // task list
        cy.contains('a', 'Return to task list').click();
        cy.contains('Sign the LPA (signatory 1)');
        cy.contains('Sign the LPA (signatory 2)').click();

        // sign
        cy.contains('a', 'The trust corporation no longer requires a second signatory').click();

        // would like a 2nd signatory
        cy.contains('label', 'No').click();
        cy.contains('button', 'Continue').click();

        // what happens next
        cy.contains('My company has formally agreed to be an attorney');
        cy.contains('a', 'Go to your dashboard');
    });
});
