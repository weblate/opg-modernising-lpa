package page

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/date"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	"golang.org/x/exp/slices"
)

const (
	AllCanNoLongerAct                = "all"
	CostOfLpaPence                   = 8200
	Jointly                          = "jointly"
	JointlyAndSeverally              = "jointly-and-severally"
	JointlyForSomeSeverallyForOthers = "mixed"
	LpaTypeCombined                  = "both"
	LpaTypeHealthWelfare             = "hw"
	LpaTypePropertyFinance           = "pfa"
	PayCookieName                    = "pay"
	PayCookiePaymentIdValueKey       = "paymentId"
	OneCanNoLongerAct                = "one"
	SomeOtherWay                     = "other"
	UsedWhenCapacityLost             = "when-capacity-lost"
	UsedWhenRegistered               = "when-registered"
)

type TaskState int

const (
	TaskNotStarted TaskState = iota
	TaskInProgress
	TaskCompleted
)

func (t TaskState) InProgress() bool { return t == TaskInProgress }
func (t TaskState) Completed() bool  { return t == TaskCompleted }

func (t TaskState) String() string {
	switch t {
	case TaskNotStarted:
		return "notStarted"
	case TaskInProgress:
		return "inProgress"
	case TaskCompleted:
		return "completed"
	}
	return ""
}

type Lpa struct {
	ID                                          string
	UpdatedAt                                   time.Time
	You                                         Person
	Attorneys                                   []Attorney
	CertificateProvider                         CertificateProvider
	WhoFor                                      string
	Contact                                     []string
	Type                                        string
	WantReplacementAttorneys                    string
	WhenCanTheLpaBeUsed                         string
	Restrictions                                string
	Tasks                                       Tasks
	Checked                                     bool
	HappyToShare                                bool
	PaymentDetails                              PaymentDetails
	CheckedAgain                                bool
	ConfirmFreeWill                             bool
	SignatureCode                               string
	EnteredSignatureCode                        string
	SignatureEmailID                            string
	SignatureSmsID                              string
	IdentityOption                              IdentityOption
	YotiUserData                                identity.UserData
	OneLoginUserData                            identity.UserData
	HowAttorneysMakeDecisions                   string
	HowAttorneysMakeDecisionsDetails            string
	ReplacementAttorneys                        []Attorney
	HowReplacementAttorneysMakeDecisions        string
	HowReplacementAttorneysMakeDecisionsDetails string
	HowShouldReplacementAttorneysStepIn         string
	HowShouldReplacementAttorneysStepInDetails  string
	DoYouWantToNotifyPeople                     string
	PeopleToNotify                              []PersonToNotify
	WitnessCode                                 WitnessCode
	CPWitnessedDonorSign                        bool
	WantToApplyForLpa                           bool
	WantToSignLpa                               bool
	Submitted                                   time.Time
	CPWitnessCodeValidated                      bool
}

type PaymentDetails struct {
	PaymentReference string
	PaymentId        string
}

type Tasks struct {
	YourDetails                TaskState
	ChooseAttorneys            TaskState
	ChooseReplacementAttorneys TaskState
	WhenCanTheLpaBeUsed        TaskState
	Restrictions               TaskState
	CertificateProvider        TaskState
	CheckYourLpa               TaskState
	PayForLpa                  TaskState
	ConfirmYourIdentityAndSign TaskState
	PeopleToNotify             TaskState
}

type Person struct {
	FirstNames  string
	LastName    string
	Email       string
	OtherNames  string
	DateOfBirth date.Date
	Address     place.Address
}

type PersonToNotify struct {
	FirstNames string
	LastName   string
	Email      string
	Address    place.Address
	ID         string
}

type Attorney struct {
	ID          string
	FirstNames  string
	LastName    string
	Email       string
	DateOfBirth date.Date
	Address     place.Address
}

type CertificateProvider struct {
	FirstNames              string
	LastName                string
	Email                   string
	Address                 place.Address
	Mobile                  string
	DateOfBirth             date.Date
	CarryOutBy              string
	Relationship            string
	RelationshipDescription string
	RelationshipLength      string
}

type Progress struct {
	LpaSigned                   TaskState
	CertificateProviderDeclared TaskState
	AttorneysDeclared           TaskState
	LpaSubmitted                TaskState
	StatutoryWaitingPeriod      TaskState
	LpaRegistered               TaskState
}

type AddressClient interface {
	LookupPostcode(ctx context.Context, postcode string) ([]place.Address, error)
}

type WitnessCode struct {
	Code    string
	Created time.Time
}

func (w *WitnessCode) HasExpired() bool {
	return w.Created.Before(time.Now().Add(-30 * time.Minute))
}

type LpaStore interface {
	Create(context.Context) (*Lpa, error)
	GetAll(context.Context) ([]*Lpa, error)
	Get(context.Context) (*Lpa, error)
	Put(context.Context, *Lpa) error
}

type sessionData struct {
	SessionID string
	LpaID     string
}

func sessionDataFromContext(ctx context.Context) *sessionData {
	data, _ := ctx.Value((*sessionData)(nil)).(*sessionData)

	return data
}

func contextWithSessionData(ctx context.Context, data *sessionData) context.Context {
	return context.WithValue(ctx, (*sessionData)(nil), data)
}

type lpaStore struct {
	dataStore DataStore
	randomInt func(int) int
}

func (s *lpaStore) Create(ctx context.Context) (*Lpa, error) {
	lpa := &Lpa{
		ID: "10" + strconv.Itoa(s.randomInt(100000)),
	}

	err := s.Put(ctx, lpa)

	return lpa, err
}

func (s *lpaStore) GetAll(ctx context.Context) ([]*Lpa, error) {
	var lpas []*Lpa
	err := s.dataStore.GetAll(ctx, sessionDataFromContext(ctx).SessionID, &lpas)

	slices.SortFunc(lpas, func(a, b *Lpa) bool {
		return a.UpdatedAt.After(b.UpdatedAt)
	})

	return lpas, err
}

func (s *lpaStore) Get(ctx context.Context) (*Lpa, error) {
	data := sessionDataFromContext(ctx)
	if data.LpaID == "" {
		return nil, errors.New("lpaStore.Get requires LpaID to retrieve")
	}

	var lpa Lpa
	if err := s.dataStore.Get(ctx, data.SessionID, data.LpaID, &lpa); err != nil {
		return nil, err
	}

	return &lpa, nil
}

func (s *lpaStore) Put(ctx context.Context, lpa *Lpa) error {
	lpa.UpdatedAt = time.Now()

	return s.dataStore.Put(ctx, sessionDataFromContext(ctx).SessionID, lpa.ID, lpa)
}

func DecodeAddress(s string) *place.Address {
	var v place.Address
	json.Unmarshal([]byte(s), &v)
	return &v
}

func (l *Lpa) IdentityConfirmed() bool {
	return l.YotiUserData.OK || l.OneLoginUserData.OK
}

func (l *Lpa) GetAttorney(id string) (Attorney, bool) {
	idx := slices.IndexFunc(l.Attorneys, func(a Attorney) bool { return a.ID == id })

	if idx == -1 {
		return Attorney{}, false
	}

	return l.Attorneys[idx], true
}

func (l *Lpa) PutAttorney(attorney Attorney) bool {
	idx := slices.IndexFunc(l.Attorneys, func(a Attorney) bool { return a.ID == attorney.ID })

	if idx == -1 {
		return false
	}

	l.Attorneys[idx] = attorney

	return true
}

func (l *Lpa) DeleteAttorney(attorney Attorney) bool {
	idx := slices.IndexFunc(l.Attorneys, func(a Attorney) bool { return a.ID == attorney.ID })

	if idx == -1 {
		return false
	}

	l.Attorneys = slices.Delete(l.Attorneys, idx, idx+1)

	return true
}

func (l *Lpa) GetReplacementAttorney(id string) (Attorney, bool) {
	idx := slices.IndexFunc(l.ReplacementAttorneys, func(a Attorney) bool { return a.ID == id })

	if idx == -1 {
		return Attorney{}, false
	}

	return l.ReplacementAttorneys[idx], true
}

func (l *Lpa) PutReplacementAttorney(attorney Attorney) bool {
	idx := slices.IndexFunc(l.ReplacementAttorneys, func(a Attorney) bool { return a.ID == attorney.ID })

	if idx == -1 {
		return false
	}

	l.ReplacementAttorneys[idx] = attorney

	return true
}

func (l *Lpa) DeleteReplacementAttorney(attorney Attorney) bool {
	idx := slices.IndexFunc(l.ReplacementAttorneys, func(a Attorney) bool { return a.ID == attorney.ID })

	if idx == -1 {
		return false
	}

	l.ReplacementAttorneys = slices.Delete(l.ReplacementAttorneys, idx, idx+1)

	return true
}

func (l *Lpa) GetPersonToNotify(id string) (PersonToNotify, bool) {
	idx := slices.IndexFunc(l.PeopleToNotify, func(p PersonToNotify) bool { return p.ID == id })

	if idx == -1 {
		return PersonToNotify{}, false
	}

	return l.PeopleToNotify[idx], true
}

func (l *Lpa) PutPersonToNotify(person PersonToNotify) bool {
	idx := slices.IndexFunc(l.PeopleToNotify, func(p PersonToNotify) bool { return p.ID == person.ID })

	if idx == -1 {
		return false
	}

	l.PeopleToNotify[idx] = person

	return true
}

func (l *Lpa) DeletePersonToNotify(personToNotify PersonToNotify) bool {
	idx := slices.IndexFunc(l.PeopleToNotify, func(p PersonToNotify) bool { return p.ID == personToNotify.ID })

	if idx == -1 {
		return false
	}

	l.PeopleToNotify = slices.Delete(l.PeopleToNotify, idx, idx+1)

	return true
}

func (l *Lpa) AttorneysFullNames() string {
	var names []string

	for _, a := range l.Attorneys {
		names = append(names, fmt.Sprintf("%s %s", a.FirstNames, a.LastName))
	}

	return concatSentence(names)
}

func (l *Lpa) AttorneysFirstNames() string {
	var names []string

	for _, a := range l.Attorneys {
		names = append(names, a.FirstNames)
	}

	return concatSentence(names)
}

func (l *Lpa) ReplacementAttorneysFullNames() string {
	var names []string

	for _, a := range l.ReplacementAttorneys {
		names = append(names, fmt.Sprintf("%s %s", a.FirstNames, a.LastName))
	}

	return concatSentence(names)
}

func (l *Lpa) ReplacementAttorneysFirstNames() string {
	var names []string

	for _, a := range l.ReplacementAttorneys {
		names = append(names, a.FirstNames)
	}

	return concatSentence(names)
}

func concatSentence(list []string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		last := len(list) - 1
		return fmt.Sprintf("%s and %s", strings.Join(list[:last], ", "), list[last])
	}
}

func (l *Lpa) DonorFullName() string {
	return fmt.Sprintf("%s %s", l.You.FirstNames, l.You.LastName)
}

func (l *Lpa) CertificateProviderFullName() string {
	return fmt.Sprintf("%s %s", l.CertificateProvider.FirstNames, l.CertificateProvider.LastName)
}

func (l *Lpa) TypeLegalTermTransKey() string {
	switch l.Type {
	case LpaTypePropertyFinance:
		return "pfaLegalTerm"
	case LpaTypeHealthWelfare:
		return "hwLegalTerm"
	case LpaTypeCombined:
		return "combinedLegalTerm"
	}
	return ""
}

func (l *Lpa) AttorneysAndCpSigningDeadline() time.Time {
	return l.Submitted.Add((24 * time.Hour) * 28)
}

func (l *Lpa) CanGoTo(url string) bool {
	path, _, _ := strings.Cut(url, "?")

	switch path {
	case Paths.WhenCanTheLpaBeUsed, Paths.Restrictions, Paths.WhoDoYouWantToBeCertificateProviderGuidance, Paths.DoYouWantToNotifyPeople:
		return l.Tasks.YourDetails.Completed() &&
			l.Tasks.ChooseAttorneys.Completed()
	case Paths.CheckYourLpa:
		return l.Tasks.YourDetails.Completed() &&
			l.Tasks.ChooseAttorneys.Completed() &&
			l.Tasks.ChooseReplacementAttorneys.Completed() &&
			l.Tasks.WhenCanTheLpaBeUsed.Completed() &&
			l.Tasks.Restrictions.Completed() &&
			l.Tasks.CertificateProvider.Completed() &&
			l.Tasks.PeopleToNotify.Completed()
	case Paths.AboutPayment:
		return l.Tasks.YourDetails.Completed() &&
			l.Tasks.ChooseAttorneys.Completed() &&
			l.Tasks.ChooseReplacementAttorneys.Completed() &&
			l.Tasks.WhenCanTheLpaBeUsed.Completed() &&
			l.Tasks.Restrictions.Completed() &&
			l.Tasks.CertificateProvider.Completed() &&
			l.Tasks.PeopleToNotify.Completed() &&
			l.Tasks.CheckYourLpa.Completed()
	case Paths.SelectYourIdentityOptions, Paths.HowToConfirmYourIdentityAndSign:
		return l.Tasks.PayForLpa.Completed()
	case "":
		return false
	default:
		return true
	}
}

func (l *Lpa) Progress() Progress {
	p := Progress{
		LpaSigned:                   TaskInProgress,
		CertificateProviderDeclared: TaskNotStarted,
		AttorneysDeclared:           TaskNotStarted,
		LpaSubmitted:                TaskNotStarted,
		StatutoryWaitingPeriod:      TaskNotStarted,
		LpaRegistered:               TaskNotStarted,
	}

	if !l.Submitted.IsZero() {
		p.LpaSigned = TaskCompleted
	}

	if p.LpaSigned.Completed() {
		p.CertificateProviderDeclared = TaskInProgress
	}

	// Further logic to be added as we build the rest of the flow

	return p
}
