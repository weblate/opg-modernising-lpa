package page

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) Format() string {
	return string(p)
}

type LpaPath string

func (p LpaPath) String() string {
	return string(p)
}

func (p LpaPath) Format(id string) string {
	return "/lpa/" + id + string(p)
}

type AttorneyPath string

func (p AttorneyPath) String() string {
	return string(p)
}

func (p AttorneyPath) Format(id string) string {
	return "/attorney/" + id + string(p)
}

type CertificateProviderPath string

func (p CertificateProviderPath) String() string {
	return string(p)
}

func (p CertificateProviderPath) Format(id string) string {
	return "/certificate-provider/" + id + string(p)
}

type AttorneyPaths struct {
	EnterReferenceNumber Path
	Login                Path
	LoginCallback        Path
	Start                Path

	CheckYourName             AttorneyPath
	CodeOfConduct             AttorneyPath
	MobileNumber              AttorneyPath
	ReadTheLpa                AttorneyPath
	RightsAndResponsibilities AttorneyPath
	Sign                      AttorneyPath
	TaskList                  AttorneyPath
	WhatHappensNext           AttorneyPath
	WhatHappensWhenYouSign    AttorneyPath
}

type CertificateProviderPaths struct {
	Login                Path
	LoginCallback        Path
	EnterReferenceNumber Path
	WhoIsEligible        Path

	CertificateProvided                  CertificateProviderPath
	CheckYourName                        CertificateProviderPath
	EnterDateOfBirth                     CertificateProviderPath
	EnterMobileNumber                    CertificateProviderPath
	IdentityWithBiometricResidencePermit CertificateProviderPath
	IdentityWithDrivingLicencePaper      CertificateProviderPath
	IdentityWithDrivingLicencePhotocard  CertificateProviderPath
	IdentityWithOneLogin                 CertificateProviderPath
	IdentityWithOneLoginCallback         CertificateProviderPath
	IdentityWithOnlineBankAccount        CertificateProviderPath
	IdentityWithPassport                 CertificateProviderPath
	IdentityWithYoti                     CertificateProviderPath
	IdentityWithYotiCallback             CertificateProviderPath
	ProvideCertificate                   CertificateProviderPath
	ReadTheLpa                           CertificateProviderPath
	SelectYourIdentityOptions            CertificateProviderPath
	SelectYourIdentityOptions1           CertificateProviderPath
	SelectYourIdentityOptions2           CertificateProviderPath
	WhatHappensNext                      CertificateProviderPath
	WhatYoullNeedToConfirmYourIdentity   CertificateProviderPath
	YourChosenIdentityOptions            CertificateProviderPath
}

type HealthCheckPaths struct {
	Service    Path
	Dependency Path
}

type AppPaths struct {
	Attorney            AttorneyPaths
	CertificateProvider CertificateProviderPaths
	HealthCheck         HealthCheckPaths

	AuthRedirect                       Path
	Login                              Path
	LoginCallback                      Path
	Root                               Path
	SignOut                            Path
	Fixtures                           Path
	YourLegalRightsAndResponsibilities Path
	CertificateProviderStart           Path
	Start                              Path
	TestingStart                       Path
	Dashboard                          Path
	YotiRedirect                       Path
	CookiesConsent                     Path

	LpaYourLegalRightsAndResponsibilities                      LpaPath
	AboutPayment                                               LpaPath
	AreYouHappyIfOneAttorneyCantActNoneCan                     LpaPath
	AreYouHappyIfOneReplacementAttorneyCantActNoneCan          LpaPath
	AreYouHappyIfRemainingAttorneysCanContinueToAct            LpaPath
	AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct LpaPath
	CertificateProviderAddress                                 LpaPath
	CertificateProviderDetails                                 LpaPath
	CertificateProviderOptOut                                  LpaPath
	CheckYourLpa                                               LpaPath
	ChooseAttorneys                                            LpaPath
	ChooseAttorneysAddress                                     LpaPath
	ChooseAttorneysSummary                                     LpaPath
	ChoosePeopleToNotify                                       LpaPath
	ChoosePeopleToNotifyAddress                                LpaPath
	ChoosePeopleToNotifySummary                                LpaPath
	ChooseReplacementAttorneys                                 LpaPath
	ChooseReplacementAttorneysAddress                          LpaPath
	ChooseReplacementAttorneysSummary                          LpaPath
	DoYouWantReplacementAttorneys                              LpaPath
	DoYouWantToNotifyPeople                                    LpaPath
	HowDoYouKnowYourCertificateProvider                        LpaPath
	HowLongHaveYouKnownCertificateProvider                     LpaPath
	HowShouldAttorneysMakeDecisions                            LpaPath
	HowShouldReplacementAttorneysMakeDecisions                 LpaPath
	HowShouldReplacementAttorneysStepIn                        LpaPath
	HowToConfirmYourIdentityAndSign                            LpaPath
	HowWouldCertificateProviderPreferToCarryOutTheirRole       LpaPath
	IdentityConfirmed                                          LpaPath
	IdentityWithBiometricResidencePermit                       LpaPath
	IdentityWithDrivingLicencePaper                            LpaPath
	IdentityWithDrivingLicencePhotocard                        LpaPath
	IdentityWithOneLogin                                       LpaPath
	IdentityWithOneLoginCallback                               LpaPath
	IdentityWithOnlineBankAccount                              LpaPath
	IdentityWithPassport                                       LpaPath
	IdentityWithYoti                                           LpaPath
	IdentityWithYotiCallback                                   LpaPath
	LifeSustainingTreatment                                    LpaPath
	LpaType                                                    LpaPath
	PaymentConfirmation                                        LpaPath
	Progress                                                   LpaPath
	ReadYourLpa                                                LpaPath
	RemoveAttorney                                             LpaPath
	RemovePersonToNotify                                       LpaPath
	RemoveReplacementAttorney                                  LpaPath
	ResendWitnessCode                                          LpaPath
	Restrictions                                               LpaPath
	SelectYourIdentityOptions                                  LpaPath
	SelectYourIdentityOptions1                                 LpaPath
	SelectYourIdentityOptions2                                 LpaPath
	SignYourLpa                                                LpaPath
	TaskList                                                   LpaPath
	UseExistingAddress                                         LpaPath
	WhatYoullNeedToConfirmYourIdentity                         LpaPath
	WhenCanTheLpaBeUsed                                        LpaPath
	WhoDoYouWantToBeCertificateProviderGuidance                LpaPath
	WhoIsTheLpaFor                                             LpaPath
	WitnessingAsCertificateProvider                            LpaPath
	WitnessingYourSignature                                    LpaPath
	YouHaveSubmittedYourLpa                                    LpaPath
	YourAddress                                                LpaPath
	YourChosenIdentityOptions                                  LpaPath
	YourDetails                                                LpaPath
}

var Paths = AppPaths{
	AboutPayment:                                               "/about-payment",
	AreYouHappyIfOneAttorneyCantActNoneCan:                     "/are-you-happy-if-one-attorney-cant-act-none-can",
	AreYouHappyIfOneReplacementAttorneyCantActNoneCan:          "/are-you-happy-if-one-replacement-attorney-cant-act-none-can",
	AreYouHappyIfRemainingAttorneysCanContinueToAct:            "/are-you-happy-if-remaining-attorneys-can-continue-to-act",
	AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct: "/are-you-happy-if-remaining-replacement-attorneys-can-continue-to-act",
	Attorney: AttorneyPaths{
		CheckYourName:             "/attorney-check-your-name",
		CodeOfConduct:             "/attorney-code-of-conduct",
		EnterReferenceNumber:      "/attorney-enter-reference-number",
		Login:                     "/attorney-login",
		LoginCallback:             "/attorney-login-callback",
		MobileNumber:              "/attorney-mobile-number",
		ReadTheLpa:                "/attorney-read-the-lpa",
		RightsAndResponsibilities: "/attorney-legal-rights-and-responsibilities",
		Sign:                      "/attorney-sign",
		Start:                     "/attorney-start",
		TaskList:                  "/attorney-task-list",
		WhatHappensNext:           "/attorney-what-happens-next",
		WhatHappensWhenYouSign:    "/attorney-what-happens-when-you-sign-the-lpa",
	},
	AuthRedirect:               "/auth/redirect",
	CertificateProviderDetails: "/certificate-provider-details",
	CertificateProvider: CertificateProviderPaths{
		CertificateProvided:                  "/certificate-provided",
		CheckYourName:                        "/certificate-provider-check-your-name",
		EnterDateOfBirth:                     "/certificate-provider-enter-date-of-birth",
		EnterMobileNumber:                    "/certificate-provider-enter-mobile-number",
		EnterReferenceNumber:                 "/certificate-provider-enter-reference-number",
		IdentityWithBiometricResidencePermit: "/certificate-provider/id/brp",
		IdentityWithDrivingLicencePaper:      "/certificate-provider/id/dlpaper",
		IdentityWithDrivingLicencePhotocard:  "/certificate-provider/id/dlphoto",
		IdentityWithOneLogin:                 "/certificate-provider-identity-with-one-login",
		IdentityWithOneLoginCallback:         "/certificate-provider-identity-with-one-login-callback",
		IdentityWithOnlineBankAccount:        "/certificate-provider/id/bank",
		IdentityWithPassport:                 "/certificate-provider/id/passport",
		IdentityWithYoti:                     "/certificate-provider-identity-with-yoti",
		IdentityWithYotiCallback:             "/certificate-provider-identity-with-yoti-callback",
		Login:                                "/certificate-provider-login",
		LoginCallback:                        "/certificate-provider-login-callback",
		ProvideCertificate:                   "/provide-certificate",
		ReadTheLpa:                           "/certificate-provider-read-the-lpa",
		SelectYourIdentityOptions1:           "/certificate-provider-select-identity-document",
		SelectYourIdentityOptions2:           "/certificate-provider-select-identity-document-2",
		SelectYourIdentityOptions:            "/certificate-provider-select-your-identity-options",
		WhatHappensNext:                      "/certificate-provider-what-happens-next",
		WhatYoullNeedToConfirmYourIdentity:   "/certificate-provider-what-youll-need-to-confirm-your-identity",
		WhoIsEligible:                        "/certificate-provider-who-is-eligible",
		YourChosenIdentityOptions:            "/certificate-provider-your-chosen-identity-options",
	},
	CertificateProviderOptOut:         "/certificate-provider-opt-out",
	CertificateProviderAddress:        "/certificate-provider-address",
	CertificateProviderStart:          "/certificate-provider-start",
	CheckYourLpa:                      "/check-your-lpa",
	ChooseAttorneys:                   "/choose-attorneys",
	ChooseAttorneysAddress:            "/choose-attorneys-address",
	ChooseAttorneysSummary:            "/choose-attorneys-summary",
	ChoosePeopleToNotify:              "/choose-people-to-notify",
	ChoosePeopleToNotifyAddress:       "/choose-people-to-notify-address",
	ChoosePeopleToNotifySummary:       "/choose-people-to-notify-summary",
	ChooseReplacementAttorneys:        "/choose-replacement-attorneys",
	ChooseReplacementAttorneysAddress: "/choose-replacement-attorneys-address",
	ChooseReplacementAttorneysSummary: "/choose-replacement-attorneys-summary",
	CookiesConsent:                    "/cookies-consent",
	Dashboard:                         "/dashboard",
	DoYouWantReplacementAttorneys:     "/do-you-want-replacement-attorneys",
	DoYouWantToNotifyPeople:           "/do-you-want-to-notify-people",
	Fixtures:                          "/fixtures",
	HealthCheck: HealthCheckPaths{
		Service:    "/health-check/service",
		Dependency: "/health-check/dependency",
	},
	HowDoYouKnowYourCertificateProvider:                  "/how-do-you-know-your-certificate-provider",
	HowLongHaveYouKnownCertificateProvider:               "/how-long-have-you-known-certificate-provider",
	HowShouldAttorneysMakeDecisions:                      "/how-should-attorneys-make-decisions",
	HowShouldReplacementAttorneysMakeDecisions:           "/how-should-replacement-attorneys-make-decisions",
	HowShouldReplacementAttorneysStepIn:                  "/how-should-replacement-attorneys-step-in",
	HowToConfirmYourIdentityAndSign:                      "/how-to-confirm-your-identity-and-sign",
	HowWouldCertificateProviderPreferToCarryOutTheirRole: "/how-would-certificate-provider-prefer-to-carry-out-their-role",
	IdentityConfirmed:                                    "/identity-confirmed",
	IdentityWithBiometricResidencePermit:                 "/id/biometric-residence-permit",
	IdentityWithDrivingLicencePaper:                      "/id/driving-licence-paper",
	IdentityWithDrivingLicencePhotocard:                  "/id/driving-licence-photocard",
	IdentityWithOneLogin:                                 "/id/one-login",
	IdentityWithOneLoginCallback:                         "/id/one-login/callback",
	IdentityWithOnlineBankAccount:                        "/id/online-bank-account",
	IdentityWithPassport:                                 "/id/passport",
	IdentityWithYoti:                                     "/id/yoti",
	IdentityWithYotiCallback:                             "/id/yoti/callback",
	LifeSustainingTreatment:                              "/life-sustaining-treatment",
	Login:                                                "/login",
	LoginCallback:                                        "/login-callback",
	LpaType:                                              "/lpa-type",
	PaymentConfirmation:                                  "/payment-confirmation",
	Progress:                                             "/progress",
	ReadYourLpa:                                          "/read-your-lpa",
	RemoveAttorney:                                       "/remove-attorney",
	RemovePersonToNotify:                                 "/remove-person-to-notify",
	RemoveReplacementAttorney:                            "/remove-replacement-attorney",
	ResendWitnessCode:                                    "/resend-witness-code",
	Restrictions:                                         "/restrictions",
	Root:                                                 "/",
	SelectYourIdentityOptions1:                           "/select-identity-document",
	SelectYourIdentityOptions2:                           "/select-identity-document-2",
	SelectYourIdentityOptions:                            "/select-your-identity-options",
	SignOut:                                              "/sign-out",
	SignYourLpa:                                          "/sign-your-lpa",
	Start:                                                "/start",
	TaskList:                                             "/task-list",
	TestingStart:                                         "/testing-start",
	UseExistingAddress:                                   "/use-existing-address",
	WhatYoullNeedToConfirmYourIdentity:                   "/what-youll-need-to-confirm-your-identity",
	WhenCanTheLpaBeUsed:                                  "/when-can-the-lpa-be-used",
	WhoDoYouWantToBeCertificateProviderGuidance:          "/who-do-you-want-to-be-certificate-provider-guidance",
	WhoIsTheLpaFor:                                       "/who-is-the-lpa-for",
	WitnessingAsCertificateProvider:                      "/witnessing-as-certificate-provider",
	WitnessingYourSignature:                              "/witnessing-your-signature",
	YotiRedirect:                                         "/yoti/redirect",
	YouHaveSubmittedYourLpa:                              "/you-have-submitted-your-lpa",
	YourAddress:                                          "/your-address",
	YourChosenIdentityOptions:                            "/your-chosen-identity-options",
	YourDetails:                                          "/your-details",
	YourLegalRightsAndResponsibilities:                   "/your-legal-rights-and-responsibilities",
	LpaYourLegalRightsAndResponsibilities:                "/lpa-your-legal-rights-and-responsibilities",
}
