describe('Check the LPA', () => {
    beforeEach(() => {
        cy.visit('/testing-start?redirect=/check-your-lpa&lpa.yourDetails=1&lpa.certificateProvider=1&lpa.attorneys=1&lpa.replacementAttorneys=2&lpa.chooseWhenCanBeUsed=1&lpa.restrictions=1&lpa.peopleToNotify=1');
    });

    it("can submit the completed LPA", () => {
        cy.contains('h1', "Check your LPA")

        cy.checkA11yApp();

        cy.contains('h2', "LPA decisions")

        cy.contains('dt', "When your attorneys can use your LPA")
        cy.contains('dt', "Who is your attorney")
        cy.contains('dt', "Who are your replacement attorneys")

        cy.contains('h2', "People named on the LPA")
        cy.contains('h3', "Donor")
        cy.contains('h3', "Certificate provider")
        cy.contains('h3', "Attorneys")

        cy.get('#f-checked-and-happy').check()

        cy.contains('button', 'Confirm').click();

        cy.url().should('contain', '/lpa-details-saved');
    });

    describe('CP acting on paper', () => {
        describe('on first check', () => {
            it('content is tailored for paper CPs and a details component is shown', () => {
                cy.visit('/testing-start?redirect=/check-your-lpa&lpa.yourDetails=1&lpa.certificateProvider=1&lpa.attorneys=1&lpa.replacementAttorneys=2&lpa.chooseWhenCanBeUsed=1&lpa.restrictions=1&lpa.peopleToNotify=1');

                cy.get('label[for=f-checked-and-happy]').contains('I’ve checked this LPA and I’m happy to show it to my certificate provider, Jessie Jones')
                cy.get('details[data-module=govuk-details]').contains('What happens if I need to make changes later?')

                cy.get('#f-checked-and-happy').check()
                cy.contains('button', 'Confirm').click();

                cy.url().should('contain', '/lpa-details-saved');

                cy.get('div[data-module=govuk-notification-banner]').contains('You should show your LPA to your certificate provider, Jessie Jones.')

                cy.contains('a', 'Continue').click();

                cy.url().should('contain', '/about-payment');
            })
        })

        describe('on subsequent check after completing LPA', () => {
            it('content is tailored for paper CPs and a warning component is shown', () => {
                cy.visit('/testing-start?redirect=/check-your-lpa&lpa.complete=1');

                cy.get('label[for=f-checked-and-happy]').contains('I’ve checked this LPA and I’m happy to show it to my certificate provider, Jessie Jones')
                cy.get('.govuk-warning-text').contains('Once you select the confirm button, your certificate provider will be sent a text telling them you have changed your LPA.')

                cy.get('#f-checked-and-happy').check()
                cy.contains('button', 'Confirm').click();

                cy.url().should('contain', '/lpa-details-saved');

                cy.get('div[data-module=govuk-notification-banner]').contains('We’ve saved your changes and sent a text to your certificate provider, Jessie Jones, to tell them that your LPA is ready for review. You should show them your LPA.')

                cy.contains('a', 'Return to dashboard').click();

                cy.url().should('contain', '/dashboard');
            })
        })
    })

    describe('CP acting online', () => {
        describe('on first check', () => {
            it('content is tailored for online CPs and a details component is shown', () => {
                cy.visit('/testing-start?redirect=/check-your-lpa&lpa.yourDetails=1&lpa.certificateProvider=1&lpa.attorneys=1&lpa.replacementAttorneys=2&lpa.chooseWhenCanBeUsed=1&lpa.restrictions=1&lpa.peopleToNotify=1&lpa.certificateProviderActOnline=1');

                cy.get('label[for=f-checked-and-happy]').contains('I’ve checked this LPA and I’m happy for OPG to share it with my certificate provider, Jessie Jones')
                cy.get('details[data-module=govuk-details]').contains('What happens if I need to make changes later?')

                cy.get('#f-checked-and-happy').check()
                cy.contains('button', 'Confirm').click();

                cy.url().should('contain', '/lpa-details-saved');

                cy.get('div[data-module=govuk-notification-banner]').contains('We’ve sent an email to your certificate provider, Jessie Jones, to tell them what they need to do next. You should tell them to expect an email from us.')

                cy.contains('a', 'Continue').click();

                cy.url().should('contain', '/about-payment');
            })
        })

        describe('on subsequent check after completing LPA', () => {
            it('content is tailored for online CPs and a warning component is shown', () => {
                cy.visit('/testing-start?redirect=/check-your-lpa&lpa.complete=1&lpa.certificateProviderActOnline=1');

                cy.get('label[for=f-checked-and-happy]').contains('I’ve checked this LPA and I’m happy for OPG to share it with my certificate provider, Jessie Jones')
                cy.get('.govuk-warning-text').contains('Once you select the confirm button, your certificate provider will be sent a text telling them you have changed your LPA.')

                cy.get('#f-checked-and-happy').check()
                cy.contains('button', 'Confirm').click();

                cy.url().should('contain', '/lpa-details-saved');

                cy.get('div[data-module=govuk-notification-banner]').contains('We’ve saved your changes and sent a text to your certificate provider, Jessie Jones, to tell them that they should review your LPA online.')

                cy.contains('a', 'Return to dashboard').click();

                cy.url().should('contain', '/dashboard');
            })
        })
    })

    it("errors when not selected", () => {
        cy.contains('button', 'Confirm').click();

        cy.get('.govuk-error-summary').within(() => {
            cy.contains('Select the box if you have checked your LPA and are happy to share it with your certificate provider');
        });

        cy.contains('.govuk-form-group .govuk-error-message', 'Select the box if you have checked your LPA and are happy to share it with your certificate provider');
    })
});
