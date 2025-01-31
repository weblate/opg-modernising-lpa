import {AddressFormAssertions} from "../support/e2e";

describe('Donor address', () => {
    beforeEach(() => {
        cy.visit('/testing-start?redirect=/your-address');
    });

    it('address can be looked up', () => {
        AddressFormAssertions.assertCanAddAddressFromSelect()
        cy.url().should('contain', '/who-is-the-lpa-for');
    });

    it('address can be entered manually if not found', () => {
        AddressFormAssertions.assertCanAddAddressManually('I can’t find my address in the list')
        cy.url().should('contain', '/who-is-the-lpa-for');
    });

    it('address can be entered manually on invalid postcode', () => {
        AddressFormAssertions.assertCanAddAddressManually('Enter address manually', true)
        cy.url().should('contain', '/who-is-the-lpa-for');
    });

    it('errors when empty postcode', () => {
        AddressFormAssertions.assertErrorsWhenPostcodeEmpty()
    });

    it('errors when invalid postcode', () => {
        AddressFormAssertions.assertErrorsWhenInvalidPostcode()
    });

    it('errors when valid postcode and no addresses', () => {
        AddressFormAssertions.assertErrorsWhenValidPostcodeFormatButNoAddressesFound()
    });

    it('errors when unselected', () => {
        AddressFormAssertions.assertErrorsWhenUnselected()
    });

    it('errors when manual incorrect', () => {
        AddressFormAssertions.assertErrorsWhenManualIncorrect('I can’t find my address in the list')
    });
});
