import {TestEmail} from "../../support/e2e";

describe.skip('Choose replacement attorneys task', () => {
    it('is not started when no replacement attorneys are set', () => {
        cy.visit('/testing-start?redirect=/task-list&donorDetails=1&cookiesAccepted=1');

        cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Not started');
    });

    it('is completed if I do not want replacement attorneys', () => {
        cy.visit('/testing-start?redirect=/task-list&donorDetails=1&cookiesAccepted=1');
        cy.contains('a', 'Choose your replacement attorneys (optional)').click();

        cy.contains('label', 'No').click();
        cy.contains('button', 'Continue').click();

        cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed');
    });

    it('is in progress if I do want replacement attorneys', () => {
        cy.visit('/testing-start?redirect=/task-list&donorDetails=1&cookiesAccepted=1');
        cy.contains('a', 'Choose your replacement attorneys (optional)').click();

        cy.contains('label', 'Yes').click();
        cy.contains('button', 'Continue').click();

        cy.visitLpa('/task-list');
        cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress');
    });

    it('is completed if enter a replacement attorneys details', () => {
        cy.visit('/testing-start?redirect=/task-list&donorDetails=1&cookiesAccepted=1');
        cy.contains('a', 'Choose your replacement attorneys (optional)').click();

        cy.contains('label', 'Yes').click();
        cy.contains('button', 'Continue').click();

        cy.get('#f-first-names').type('John');
        cy.get('#f-last-name').type('Doe');
        cy.get('#f-email').type(TestEmail);
        cy.get('#f-date-of-birth').type('1');
        cy.get('#f-date-of-birth-month').type('2');
        cy.get('#f-date-of-birth-year').type('1990');
        cy.contains('button', 'Continue').click();

        cy.visitLpa('/task-list');
        cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (1)');
    });

    it('is in progress if enter a replacement attorneys details then add attorneys', () => {
        cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorney=1&cookiesAccepted=1');
        cy.contains('a', 'Choose your replacement attorneys (optional)').click();

        cy.contains('label', 'Yes').click();
        cy.contains('button', 'Continue').click();

        cy.get('#f-first-names').type('John');
        cy.get('#f-last-name').type('Doe');
        cy.get('#f-email').type(TestEmail);
        cy.get('#f-date-of-birth').type('1');
        cy.get('#f-date-of-birth-month').type('2');
        cy.get('#f-date-of-birth-year').type('1990');
        cy.contains('button', 'Continue').click();

        cy.visitLpa('/task-list');

        cy.contains('a', 'Choose your attorneys').click();

        cy.contains('label', 'Yes').click();
        cy.contains('button', 'Continue').click();

        cy.get('#f-first-names').type('Janet');
        cy.get('#f-last-name').type('Doe');
        cy.get('#f-email').type(TestEmail);
        cy.get('#f-date-of-birth').type('1');
        cy.get('#f-date-of-birth-month').type('2');
        cy.get('#f-date-of-birth-year').type('1990');
        cy.contains('button', 'Continue').click();
        cy.contains('button', 'Skip').click();

        cy.contains('label', 'No').click();
        cy.contains('button', 'Continue').click();

        cy.get('input[value=jointly-and-severally]').click();
        cy.contains('button', 'Continue').click();

        cy.visitLpa('/task-list');
        cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (1)');
    });

    describe('having a single attorney and a single replacement attorney', () => {
        it('is completed', () => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorney=1&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (1)');
        });
    });

    describe('having a single attorney and multiple replacement attorneys', () => {
        beforeEach(() => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorney=1&withReplacementAttorney=1&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();
        });

        it('is in progress', () => {
            cy.visitLpa('/task-list');

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is completed if replacements act jointly and severally', () => {
            cy.get('input[value=jointly-and-severally]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is in progress if replacement act jointly', () => {
            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is in progress if replacement act mixed', () => {
            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is completed if replacement act jointly happily', () => {
            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is completed if replacement act mixed happily', () => {
            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });
    });

    describe('having jointly and severally attorneys and a single replacement attorney', () => {
        beforeEach(() => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorneys=1&howAttorneysAct=jointly-and-severally&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();
        });

        it('is completed if step in as soon as one', () => {
            cy.contains('label', 'As soon as one').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (1)');
        });

        it('is completed if step in when none', () => {
            cy.contains('label', 'When none').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (1)');
        });

        it('is completed if step in some other way', () => {
            cy.contains('label', 'In some other way').click();
            cy.get('textarea').type('Details');
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (1)');
        });
    });

    describe('having jointly attorneys and a single replacement attorney', () => {
        it('is completed', () => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorneys=1&howAttorneysAct=jointly&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (1)');
        });
    });

    describe('having jointly for some attorneys and a single replacement attorney', () => {
        it('is completed', () => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorneys=1&howAttorneysAct=mixed&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (1)');
        });
    });

    describe('having jointly and severally attorneys and multiple replacement attorneys', () => {
        beforeEach(() => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorneys=1&howAttorneysAct=jointly-and-severally&withReplacementAttorney=1&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();
        });

        it('is completed if step in as soon as one', () => {
            cy.contains('label', 'As soon as one').click();
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is in progress if step in when none', () => {
            cy.contains('label', 'When none').click();
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is completed if step in when none and jointly and severally', () => {
            cy.contains('label', 'When none').click();
            cy.contains('button', 'Continue').click();

            cy.get('input[value=jointly-and-severally]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is in progress if step in when none and jointly', () => {
            cy.contains('label', 'When none').click();
            cy.contains('button', 'Continue').click();

            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is in progress if step in when none and mixed', () => {
            cy.contains('label', 'When none').click();
            cy.contains('button', 'Continue').click();

            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is complete if step in when none and jointly and happy', () => {
            cy.contains('label', 'When none').click();
            cy.contains('button', 'Continue').click();

            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is complete if step in when none and mixed and happy', () => {
            cy.contains('label', 'When none').click();
            cy.contains('button', 'Continue').click();

            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('does not allow in some other way', () => {
            cy.contains('label', 'In some other way').should('not.exist');
        });
    });

    describe('having jointly attorneys and multiple replacement attorneys', () => {
        beforeEach(() => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorneys=1&howAttorneysAct=jointly&withReplacementAttorney=1&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();
        });

        it('is completed if jointly and severally', () => {
            cy.get('input[value=jointly-and-severally]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is in progress if jointly', () => {
            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is completed if jointly and happy', () => {
            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is in progress if mixed', () => {
            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is completed if mixed and happy', () => {
            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });
    });

    describe('having jointly for some attorneys and multiple replacement attorneys', () => {
        beforeEach(() => {
            cy.visit('/testing-start?redirect=/task-list&donorDetails=1&withAttorneys=1&howAttorneysAct=mixed&withReplacementAttorney=1&cookiesAccepted=1');
            cy.contains('a', 'Choose your replacement attorneys (optional)').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.get('#f-first-names').type('John');
            cy.get('#f-last-name').type('Doe');
            cy.get('#f-email').type(TestEmail);
            cy.get('#f-date-of-birth').type('1');
            cy.get('#f-date-of-birth-month').type('2');
            cy.get('#f-date-of-birth-year').type('1990');
            cy.contains('button', 'Continue').click();
            cy.contains('button', 'Skip').click();

            cy.contains('label', 'No').click();
            cy.contains('button', 'Continue').click();
        });

        it('is completed if replacements act jointly and severally', () => {
            cy.get('input[value=jointly-and-severally]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is in progress if replacements act jointly', () => {
            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is in progress if replacements act mixed', () => {
            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.visitLpa('/task-list');
            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('In progress (2)');
        });

        it('is completed if replacements act jointly happily', () => {
            cy.get('input[value=jointly]').click();
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });

        it('is completed if replacements act mixed happily', () => {
            cy.get('input[value=mixed]').click();
            cy.get('textarea').type('Some details');
            cy.contains('button', 'Continue').click();

            cy.contains('label', 'Yes').click();
            cy.contains('button', 'Continue').click();

            cy.contains('a', 'Choose your replacement attorneys (optional)').parent().parent().contains('Completed (2)');
        });
    });
});