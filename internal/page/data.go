// Package page contains the core code and business logic of Make and Register a Lasting Power of Attorney (MRLPA)
//
// Useful links:
//   - [page.Lpa] - details about the LPA being drafted
//   - [actor.Donor] - details about the donor, provided by the applicant
//   - [actor.CertificateProvider] - details about the certificate provider, provided by the applicant
//   - [actor.CertificateProviderProvidedDetails] - details about the certificate provider, provided by the certificate provider
//   - [actor.Attorney] - details about an attorney or replacement attorney, provided by the applicant
//   - [actor.AttorneyDecisions] - details about how an attorney or replacement attorney should act, provided by the applicant
//   - [actor.AttorneyProvidedDetails] - details about an attorney or replacement attorney, provided by the attorney or replacement attorney
//   - [actor.PersonToNotify] - details about a person to notify, provided by the applicant
package page

import (
	"context"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/form"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	"golang.org/x/exp/slices"
)

//go:generate enumerator -type LpaType -linecomment -trimprefix -empty
type LpaType uint8

const (
	LpaTypeHealthWelfare   LpaType = iota + 1 // hw
	LpaTypePropertyFinance                    // pfa
)

func (e LpaType) LegalTermTransKey() string {
	switch e {
	case LpaTypePropertyFinance:
		return "pfaLegalTerm"
	case LpaTypeHealthWelfare:
		return "hwLegalTerm"
	}
	return ""
}

//go:generate enumerator -type CanBeUsedWhen -linecomment -trimprefix -empty
type CanBeUsedWhen uint8

const (
	CanBeUsedWhenCapacityLost CanBeUsedWhen = iota + 1 // when-capacity-lost
	CanBeUsedWhenHasCapacity                           // when-has-capacity
)

//go:generate enumerator -type LifeSustainingTreatment -linecomment -trimprefix -empty
type LifeSustainingTreatment uint8

const (
	LifeSustainingTreatmentOptionA LifeSustainingTreatment = iota + 1 // option-a
	LifeSustainingTreatmentOptionB                                    // option-b
)

//go:generate enumerator -type ReplacementAttorneysStepIn -linecomment -trimprefix -empty
type ReplacementAttorneysStepIn uint8

const (
	ReplacementAttorneysStepInWhenAllCanNoLongerAct ReplacementAttorneysStepIn = iota + 1 // all
	ReplacementAttorneysStepInWhenOneCanNoLongerAct                                       // one
	ReplacementAttorneysStepInAnotherWay                                                  // other
)

//go:generate enumerator -type ApplicationReason -linecomment -empty
type ApplicationReason uint8

const (
	NewApplication             ApplicationReason = iota + 1 // new-application
	RemakeOfInvalidApplication                              // remake
	AdditionalApplication                                   // additional-application
)

//go:generate enumerator -type FeeType
type FeeType uint8

const (
	FullFee FeeType = iota
	HalfFee
	NoFee
	HardshipFee
)

func (i FeeType) Cost() int {
	if i.IsFullFee() {
		return 8200
	}

	if i.IsHalfFee() {
		return 4100
	}

	return 0
}

// Lpa contains all the data related to the LPA application
type Lpa struct {
	PK, SK string
	// Identifies the LPA being drafted
	ID string
	// A unique identifier created after sending basic LPA details to the UID service
	UID string `dynamodbav:",omitempty"`
	// Tracking when the LPA is updated
	UpdatedAt time.Time
	// The donor the LPA relates to
	Donor actor.Donor
	// Attorneys named in the LPA
	Attorneys actor.Attorneys
	// Information on how the applicant wishes their attorneys to act
	AttorneyDecisions actor.AttorneyDecisions
	// The certificate provider named in the LPA
	CertificateProvider actor.CertificateProvider
	// Who the LPA is being drafted for (set, but not used)
	WhoFor string
	// Type of LPA being drafted
	Type LpaType
	// ApplicationReason is why the application is being made
	ApplicationReason ApplicationReason
	// PreviousApplicationNumber if the application is related to an existing application
	PreviousApplicationNumber string
	// Whether the applicant wants to add replacement attorneys
	WantReplacementAttorneys form.YesNo
	// When the LPA can be used
	WhenCanTheLpaBeUsed CanBeUsedWhen
	// Preferences on life sustaining treatment (applicable to personal welfare LPAs only)
	LifeSustainingTreatmentOption LifeSustainingTreatment
	// Restrictions on attorneys actions
	Restrictions string
	// Used to show the task list
	Tasks Tasks
	// Whether the applicant has checked the LPA and is happy to share the LPA with the certificate provider
	CheckedAndHappy bool
	// Used as part of GOV.UK Pay
	PaymentDetails PaymentDetails
	// Which option has been used to complete applicant identity checks
	DonorIdentityOption identity.Option
	// Information returned by the identity service related to the applicant
	DonorIdentityUserData identity.UserData
	// Replacement attorneys named in the LPA
	ReplacementAttorneys actor.Attorneys
	// Information on how the applicant wishes their replacement attorneys to act
	ReplacementAttorneyDecisions actor.AttorneyDecisions
	// How to bring in replacement attorneys, if set
	HowShouldReplacementAttorneysStepIn ReplacementAttorneysStepIn
	// Details on how replacement attorneys must step in if HowShouldReplacementAttorneysStepIn is set to "other"
	HowShouldReplacementAttorneysStepInDetails string
	// Whether the applicant wants to notify people about the application
	DoYouWantToNotifyPeople form.YesNo
	// People to notify about the application
	PeopleToNotify actor.PeopleToNotify
	// Codes used for the certificate provider to witness signing
	WitnessCodes WitnessCodes
	// Confirmation that the applicant wants to apply to register the LPA
	WantToApplyForLpa bool
	// Confirmation that the applicant wants to sign the LPA
	WantToSignLpa bool
	// When the Lpa was signed
	Submitted time.Time
	// Whether the signing was witnessed by the certificate provider
	CPWitnessCodeValidated bool
	// Used to rate limit witnessing requests
	WitnessCodeLimiter *Limiter
	// FeeType is the type of fee the user is applying for
	FeeType FeeType
	// EvidenceFormAddress is where the form to provide evidence for a fee reduction will be sent
	EvidenceFormAddress place.Address
	// EvidenceKey is the S3 key for uploaded evidence
	EvidenceKey string

	HasSentPreviousApplicationLinkedEvent bool
	HasSentEvidenceFormRequiredEvent      bool
	HasSentReducedFeeRequestedEvent       bool
}

type PaymentDetails struct {
	// Reference generated for the payment
	PaymentReference string
	// ID returned from GOV.UK Pay
	PaymentId string
	// Amount is the amount paid in pence
	Amount int
}

type Tasks struct {
	YourDetails                actor.TaskState
	ChooseAttorneys            actor.TaskState
	ChooseReplacementAttorneys actor.TaskState
	WhenCanTheLpaBeUsed        actor.TaskState // pfa only
	LifeSustainingTreatment    actor.TaskState // hw only
	Restrictions               actor.TaskState
	CertificateProvider        actor.TaskState
	CheckYourLpa               actor.TaskState
	PayForLpa                  actor.PaymentTask
	ConfirmYourIdentityAndSign actor.TaskState
	PeopleToNotify             actor.TaskState
}

type Progress struct {
	LpaSigned                   actor.TaskState
	CertificateProviderDeclared actor.TaskState
	AttorneysDeclared           actor.TaskState
	LpaSubmitted                actor.TaskState
	StatutoryWaitingPeriod      actor.TaskState
	LpaRegistered               actor.TaskState
}

type SessionData struct {
	SessionID string
	LpaID     string
}

type SessionMissingError struct{}

func (s SessionMissingError) Error() string {
	return "Session data not set"
}

func SessionDataFromContext(ctx context.Context) (*SessionData, error) {
	data, ok := ctx.Value((*SessionData)(nil)).(*SessionData)

	if !ok {
		return nil, SessionMissingError{}
	}

	return data, nil
}

func ContextWithSessionData(ctx context.Context, data *SessionData) context.Context {
	return context.WithValue(ctx, (*SessionData)(nil), data)
}

func (l *Lpa) DonorIdentityConfirmed() bool {
	return l.DonorIdentityUserData.OK && l.DonorIdentityUserData.Provider != identity.UnknownOption &&
		l.DonorIdentityUserData.MatchName(l.Donor.FirstNames, l.Donor.LastName) &&
		l.DonorIdentityUserData.DateOfBirth.Equals(l.Donor.DateOfBirth)
}

func (l *Lpa) AttorneysAndCpSigningDeadline() time.Time {
	return l.Submitted.Add((24 * time.Hour) * 28)
}

func (l *Lpa) CanGoTo(url string) bool {
	path, _, _ := strings.Cut(url, "?")
	if path == "" {
		return false
	}

	if strings.HasPrefix(path, "/lpa/") {
		_, lpaPath, _ := strings.Cut(strings.TrimPrefix(path, "/lpa/"), "/")
		return l.canGoToLpaPath("/" + lpaPath)
	}

	return true
}

func (l *Lpa) canGoToLpaPath(path string) bool {
	section1Completed := l.Tasks.YourDetails.Completed() &&
		l.Tasks.ChooseAttorneys.Completed() &&
		l.Tasks.ChooseReplacementAttorneys.Completed() &&
		(l.Type == LpaTypeHealthWelfare && l.Tasks.LifeSustainingTreatment.Completed() || l.Type == LpaTypePropertyFinance && l.Tasks.WhenCanTheLpaBeUsed.Completed()) &&
		l.Tasks.Restrictions.Completed() &&
		l.Tasks.CertificateProvider.Completed() &&
		l.Tasks.PeopleToNotify.Completed() &&
		l.Tasks.CheckYourLpa.Completed()

	switch path {
	case Paths.ReadYourLpa.String(), Paths.SignYourLpa.String(), Paths.WitnessingYourSignature.String(), Paths.WitnessingAsCertificateProvider.String(), Paths.YouHaveSubmittedYourLpa.String():
		return l.DonorIdentityConfirmed()
	case Paths.WhenCanTheLpaBeUsed.String(), Paths.LifeSustainingTreatment.String(), Paths.Restrictions.String(), Paths.WhatACertificateProviderDoes.String(), Paths.DoYouWantToNotifyPeople.String(), Paths.DoYouWantReplacementAttorneys.String():
		return l.Tasks.YourDetails.Completed() &&
			l.Tasks.ChooseAttorneys.Completed()
	case Paths.CheckYourLpa.String():
		return l.Tasks.YourDetails.Completed() &&
			l.Tasks.ChooseAttorneys.Completed() &&
			l.Tasks.ChooseReplacementAttorneys.Completed() &&
			(l.Type == LpaTypeHealthWelfare && l.Tasks.LifeSustainingTreatment.Completed() || l.Tasks.WhenCanTheLpaBeUsed.Completed()) &&
			l.Tasks.Restrictions.Completed() &&
			l.Tasks.CertificateProvider.Completed() &&
			l.Tasks.PeopleToNotify.Completed()
	case Paths.AboutPayment.String():
		return section1Completed
	case Paths.SelectYourIdentityOptions.String(), Paths.HowToConfirmYourIdentityAndSign.String():
		return section1Completed && l.Tasks.PayForLpa.IsCompleted()
	case "":
		return false
	default:
		return true
	}
}

func (l *Lpa) Progress(certificateProvider *actor.CertificateProviderProvidedDetails) Progress {
	p := Progress{
		LpaSigned:                   actor.TaskInProgress,
		CertificateProviderDeclared: actor.TaskNotStarted,
		AttorneysDeclared:           actor.TaskNotStarted,
		LpaSubmitted:                actor.TaskNotStarted,
		StatutoryWaitingPeriod:      actor.TaskNotStarted,
		LpaRegistered:               actor.TaskNotStarted,
	}

	if !l.Submitted.IsZero() {
		p.LpaSigned = actor.TaskCompleted
		p.CertificateProviderDeclared = actor.TaskInProgress
	}

	if !certificateProvider.Certificate.Agreed.IsZero() {
		p.CertificateProviderDeclared = actor.TaskCompleted
		p.AttorneysDeclared = actor.TaskInProgress
	}

	return p
}

type AddressDetail struct {
	Name    string
	Role    actor.Type
	Address place.Address
	ID      string
}

func (l *Lpa) ActorAddresses() []place.Address {
	var addresses []place.Address

	if l.Donor.Address.String() != "" {
		addresses = append(addresses, l.Donor.Address)
	}

	if l.CertificateProvider.Address.String() != "" && !slices.Contains(addresses, l.CertificateProvider.Address) {
		addresses = append(addresses, l.CertificateProvider.Address)
	}

	for _, address := range l.Attorneys.Addresses() {
		if address.String() != "" && !slices.Contains(addresses, address) {
			addresses = append(addresses, address)
		}
	}

	for _, address := range l.ReplacementAttorneys.Addresses() {
		if address.String() != "" && !slices.Contains(addresses, address) {
			addresses = append(addresses, address)
		}
	}

	return addresses
}

func (l *Lpa) AllLayAttorneysFirstNames() []string {
	var names []string

	for _, a := range l.Attorneys.Attorneys {
		names = append(names, a.FirstNames)
	}

	for _, a := range l.ReplacementAttorneys.Attorneys {
		names = append(names, a.FirstNames)
	}

	return names
}

func (l *Lpa) AllLayAttorneysFullNames() []string {
	var names []string

	for _, a := range l.Attorneys.Attorneys {
		names = append(names, a.FullName())
	}

	for _, a := range l.ReplacementAttorneys.Attorneys {
		names = append(names, a.FullName())
	}

	return names
}

func (l *Lpa) TrustCorporationsNames() []string {
	var names []string

	if l.Attorneys.TrustCorporation.Name != "" {
		names = append(names, l.Attorneys.TrustCorporation.Name)
	}

	if l.ReplacementAttorneys.TrustCorporation.Name != "" {
		names = append(names, l.ReplacementAttorneys.TrustCorporation.Name)
	}

	return names
}

func ChooseAttorneysState(attorneys actor.Attorneys, decisions actor.AttorneyDecisions) actor.TaskState {
	if attorneys.Len() == 0 {
		return actor.TaskNotStarted
	}

	if !attorneys.Complete() {
		return actor.TaskInProgress
	}

	if attorneys.Len() > 1 && !decisions.IsComplete() {
		return actor.TaskInProgress
	}

	return actor.TaskCompleted
}

func ChooseReplacementAttorneysState(lpa *Lpa) actor.TaskState {
	if lpa.WantReplacementAttorneys == form.No {
		return actor.TaskCompleted
	}

	if lpa.ReplacementAttorneys.Len() == 0 {
		if lpa.WantReplacementAttorneys != form.Yes && lpa.WantReplacementAttorneys != form.No {
			return actor.TaskNotStarted
		}

		return actor.TaskInProgress
	}

	if !lpa.ReplacementAttorneys.Complete() {
		return actor.TaskInProgress
	}

	if lpa.ReplacementAttorneys.Len() > 1 &&
		lpa.HowShouldReplacementAttorneysStepIn != ReplacementAttorneysStepInWhenOneCanNoLongerAct &&
		!lpa.ReplacementAttorneyDecisions.IsComplete() {
		return actor.TaskInProgress
	}

	if lpa.AttorneyDecisions.How.IsJointly() &&
		lpa.ReplacementAttorneys.Len() > 1 &&
		!lpa.ReplacementAttorneyDecisions.IsComplete() {
		return actor.TaskInProgress
	}

	if lpa.AttorneyDecisions.How.IsJointlyAndSeverally() {
		if lpa.HowShouldReplacementAttorneysStepIn.Empty() {
			return actor.TaskInProgress
		}

		if lpa.ReplacementAttorneys.Len() > 1 &&
			lpa.HowShouldReplacementAttorneysStepIn == ReplacementAttorneysStepInWhenAllCanNoLongerAct &&
			!lpa.ReplacementAttorneyDecisions.IsComplete() {
			return actor.TaskInProgress
		}
	}

	return actor.TaskCompleted
}