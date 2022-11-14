package page

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
)

const (
	PayCookieName              = "pay"
	PayCookiePaymentIdValueKey = "paymentId"
	CostOfLpaPence             = 8200
)

type TaskState int

const (
	TaskNotStarted TaskState = iota
	TaskInProgress
	TaskCompleted
)

type Lpa struct {
	ID                                          string
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
	IdentityOptions                             IdentityOptions
	YotiUserData                                identity.UserData
	HowAttorneysMakeDecisions                   string
	HowAttorneysMakeDecisionsDetails            string
	ReplacementAttorneys                        []Attorney
	HowReplacementAttorneysMakeDecisions        string
	HowReplacementAttorneysMakeDecisionsDetails string
}

type PaymentDetails struct {
	PaymentReference string
	PaymentId        string
}

type Tasks struct {
	WhenCanTheLpaBeUsed        TaskState
	Restrictions               TaskState
	CertificateProvider        TaskState
	CheckYourLpa               TaskState
	PayForLpa                  TaskState
	ConfirmYourIdentityAndSign TaskState
	Attorneys                  TaskState
	ReplacementAttorneys       TaskState
}

type Person struct {
	FirstNames  string
	LastName    string
	Email       string
	OtherNames  string
	DateOfBirth time.Time
	Address     place.Address
}

type Attorney struct {
	ID          string
	FirstNames  string
	LastName    string
	Email       string
	DateOfBirth time.Time
	Address     place.Address
}

type CertificateProvider struct {
	FirstNames              string
	LastName                string
	Email                   string
	DateOfBirth             time.Time
	Relationship            []string
	RelationshipDescription string
	RelationshipLength      string
}

type AddressClient interface {
	LookupPostcode(ctx context.Context, postcode string) ([]place.Address, error)
}

type Date struct {
	Day   string
	Month string
	Year  string
}

func readDate(t time.Time) Date {
	return Date{
		Day:   t.Format("2"),
		Month: t.Format("1"),
		Year:  t.Format("2006"),
	}
}

type LpaStore interface {
	Get(context.Context, string) (*Lpa, error)
	Put(context.Context, string, *Lpa) error
}

type lpaStore struct {
	dataStore DataStore
	randomInt func(int) int
}

func (s *lpaStore) Get(ctx context.Context, sessionID string) (*Lpa, error) {
	var lpa Lpa
	if err := s.dataStore.Get(ctx, sessionID, &lpa); err != nil {
		return &lpa, err
	}

	if lpa.ID == "" {
		lpa.ID = "10" + strconv.Itoa(s.randomInt(100000))
	}

	return &lpa, nil
}

func (s *lpaStore) Put(ctx context.Context, sessionID string, lpa *Lpa) error {
	return s.dataStore.Put(ctx, sessionID, lpa)
}

func DecodeAddress(s string) *place.Address {
	var v place.Address
	json.Unmarshal([]byte(s), &v)
	return &v
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
