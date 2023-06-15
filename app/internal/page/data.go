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
	"encoding/json"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/place"
	"golang.org/x/exp/slices"
)

const (
	AllCanNoLongerAct      = "all"
	CostOfLpaPence         = 8200
	LpaTypeHealthWelfare   = "hw"
	LpaTypePropertyFinance = "pfa"
	OneCanNoLongerAct      = "one"
	SomeOtherWay           = "other"
	UsedWhenCapacityLost   = "when-capacity-lost"
	UsedWhenRegistered     = "when-registered"
	OptionA                = "option-a"
	OptionB                = "option-b"
)

// Lpa contains all the data related to the LPA application
type Lpa struct {
	// Identifies the LPA being drafted
	ID string
	// A unique identifier created after sending basic LPA details to the UID service
	UID string
	// Tracking when the LPA is updated
	UpdatedAt time.Time
	// The donor the LPA relates to
	Donor actor.Donor
	// Attorney/s named in the LPA
	Attorneys actor.Attorneys
	// Information on how the applicant wishes their attorneys to act
	AttorneyDecisions actor.AttorneyDecisions
	// The certificate provider named in the LPA
	CertificateProvider actor.CertificateProvider
	// Who the LPA is being drafted for (set, but not used)
	WhoFor string
	// Type of LPA being drafted
	Type string
	// Whether the applicant wants to add replacement attorneys
	WantReplacementAttorneys string
	// When the LPA can be used
	WhenCanTheLpaBeUsed string
	// Preferences on life sustaining treatment (applicable to personal welfare LPAs only)
	LifeSustainingTreatmentOption string
	// Restrictions on attorneys actions
	Restrictions string
	// Used to show the task list
	Tasks Tasks
	// Whether the applicant has checked the LPA
	Checked bool
	// Whether the applicant is happy to share the LPA with the certificate provider
	HappyToShare bool
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
	HowShouldReplacementAttorneysStepIn string
	// Details on how replacement attorneys must step in if HowShouldReplacementAttorneysStepIn is set to "other"
	HowShouldReplacementAttorneysStepInDetails string
	// Whether the applicant wants to notify people about the application
	DoYouWantToNotifyPeople string
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
}

type PaymentDetails struct {
	// Reference generated for the payment
	PaymentReference string
	// ID returned from GOV.UK Pay
	PaymentId string
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
	PayForLpa                  actor.TaskState
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

func DecodeAddress(s string) *place.Address {
	var v place.Address
	json.Unmarshal([]byte(s), &v)
	return &v
}

func (l *Lpa) DonorIdentityConfirmed() bool {
	return l.DonorIdentityUserData.OK && l.DonorIdentityUserData.Provider != identity.UnknownOption &&
		l.DonorIdentityUserData.MatchName(l.Donor.FirstNames, l.Donor.LastName) &&
		l.DonorIdentityUserData.DateOfBirth.Equals(l.Donor.DateOfBirth)
}

func (l *Lpa) TypeLegalTermTransKey() string {
	switch l.Type {
	case LpaTypePropertyFinance:
		return "pfaLegalTerm"
	case LpaTypeHealthWelfare:
		return "hwLegalTerm"
	}
	return ""
}

func (l *Lpa) AttorneysAndCpSigningDeadline() time.Time {
	return l.Submitted.Add((24 * time.Hour) * 28)
}

func (l *Lpa) CanGoTo(url string) bool {
	path, _, _ := strings.Cut(url, "?")

	section1Completed := l.Tasks.YourDetails.Completed() &&
		l.Tasks.ChooseAttorneys.Completed() &&
		l.Tasks.ChooseReplacementAttorneys.Completed() &&
		(l.Type == LpaTypeHealthWelfare && l.Tasks.LifeSustainingTreatment.Completed() || l.Type == LpaTypePropertyFinance && l.Tasks.WhenCanTheLpaBeUsed.Completed()) &&
		l.Tasks.Restrictions.Completed() &&
		l.Tasks.CertificateProvider.Completed() &&
		l.Tasks.PeopleToNotify.Completed() &&
		l.Tasks.CheckYourLpa.Completed()

	switch path {
	case Paths.ReadYourLpa, Paths.SignYourLpa, Paths.WitnessingYourSignature, Paths.WitnessingAsCertificateProvider, Paths.YouHaveSubmittedYourLpa:
		return l.DonorIdentityConfirmed()
	case Paths.WhenCanTheLpaBeUsed, Paths.LifeSustainingTreatment, Paths.Restrictions, Paths.WhoDoYouWantToBeCertificateProviderGuidance, Paths.DoYouWantToNotifyPeople, Paths.DoYouWantReplacementAttorneys:
		return l.Tasks.YourDetails.Completed() &&
			l.Tasks.ChooseAttorneys.Completed()
	case Paths.CheckYourLpa:
		return l.Tasks.YourDetails.Completed() &&
			l.Tasks.ChooseAttorneys.Completed() &&
			l.Tasks.ChooseReplacementAttorneys.Completed() &&
			(l.Type == LpaTypeHealthWelfare && l.Tasks.LifeSustainingTreatment.Completed() || l.Tasks.WhenCanTheLpaBeUsed.Completed()) &&
			l.Tasks.Restrictions.Completed() &&
			l.Tasks.CertificateProvider.Completed() &&
			l.Tasks.PeopleToNotify.Completed()
	case Paths.AboutPayment:
		return section1Completed
	case Paths.SelectYourIdentityOptions, Paths.HowToConfirmYourIdentityAndSign:
		return section1Completed && l.Tasks.PayForLpa.Completed()
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

	for _, attorney := range l.Attorneys {
		if attorney.Address.String() != "" && !slices.Contains(addresses, attorney.Address) {
			addresses = append(addresses, attorney.Address)
		}
	}

	for _, replacementAttorney := range l.ReplacementAttorneys {
		if replacementAttorney.Address.String() != "" && !slices.Contains(addresses, replacementAttorney.Address) {
			addresses = append(addresses, replacementAttorney.Address)
		}
	}

	return addresses
}

func ChooseAttorneysState(attorneys actor.Attorneys, decisions actor.AttorneyDecisions) actor.TaskState {
	if len(attorneys) == 0 {
		return actor.TaskNotStarted
	}

	for _, a := range attorneys {
		if a.FirstNames == "" || (a.Address.Line1 == "" && a.Email == "") {
			return actor.TaskInProgress
		}
	}

	if len(attorneys) > 1 && !decisions.IsComplete(len(attorneys)) {
		return actor.TaskInProgress
	}

	return actor.TaskCompleted
}

func ChooseReplacementAttorneysState(lpa *Lpa) actor.TaskState {
	if lpa.WantReplacementAttorneys == "no" {
		return actor.TaskCompleted
	}

	if len(lpa.ReplacementAttorneys) == 0 {
		if lpa.WantReplacementAttorneys == "" {
			return actor.TaskNotStarted
		}

		return actor.TaskInProgress
	}

	for _, a := range lpa.ReplacementAttorneys {
		if a.FirstNames == "" || (a.Address.Line1 == "" && a.Email == "") {
			return actor.TaskInProgress
		}
	}

	if len(lpa.ReplacementAttorneys) > 1 &&
		lpa.HowShouldReplacementAttorneysStepIn != OneCanNoLongerAct &&
		!lpa.ReplacementAttorneyDecisions.IsComplete(len(lpa.ReplacementAttorneys)) {
		return actor.TaskInProgress
	}

	if lpa.AttorneyDecisions.How == actor.Jointly &&
		len(lpa.ReplacementAttorneys) > 1 &&
		!lpa.ReplacementAttorneyDecisions.IsComplete(len(lpa.ReplacementAttorneys)) {
		return actor.TaskInProgress
	}

	if lpa.AttorneyDecisions.How == actor.JointlyAndSeverally {
		if lpa.HowShouldReplacementAttorneysStepIn == "" {
			return actor.TaskInProgress
		}

		if len(lpa.ReplacementAttorneys) > 1 &&
			lpa.HowShouldReplacementAttorneysStepIn == AllCanNoLongerAct &&
			!lpa.ReplacementAttorneyDecisions.IsComplete(len(lpa.ReplacementAttorneys)) {
			return actor.TaskInProgress
		}
	}

	return actor.TaskCompleted
}

func (l *Lpa) IsHealthAndWelfareLpa() bool {
	return l.Type == LpaTypeHealthWelfare
}
