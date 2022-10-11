package page

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/ordnance_survey"

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
	You                      Person
	Attorney                 Attorney
	CertificateProvider      CertificateProvider
	WhoFor                   string
	Contact                  []string
	Type                     string
	WantReplacementAttorneys string
	WhenCanTheLpaBeUsed      string
	Restrictions             string
	Tasks                    Tasks
	Checked                  bool
	HappyToShare             bool
	PaymentDetails           PaymentDetails
	CheckedAgain             bool
	ConfirmFreeWill          bool
	SignatureCode            string
	IdentityOptions          IdentityOptions
	YotiUserData             identity.UserData
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
}

type Person struct {
	FirstNames  string
	LastName    string
	Email       string
	OtherNames  string
	DateOfBirth time.Time
	Address     Address
}

type Attorney struct {
	FirstNames  string
	LastName    string
	Email       string
	DateOfBirth time.Time
	Address     Address
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

type Address struct {
	Line1      string
	Line2      string
	Line3      string
	TownOrCity string
	Postcode   string
}

type AddressClient interface {
	LookupPostcode(string) ([]Address, error)
}

func (a Address) Encode() string {
	x, _ := json.Marshal(a)
	return string(x)
}

func DecodeAddress(s string) *Address {
	var v Address
	json.Unmarshal([]byte(s), &v)
	return &v
}

func (a Address) String() string {
	var parts []string

	if a.Line1 != "" {
		parts = append(parts, a.Line1)
	}
	if a.Line2 != "" {
		parts = append(parts, a.Line2)
	}
	if a.TownOrCity != "" {
		parts = append(parts, a.TownOrCity)
	}
	if a.Postcode != "" {
		parts = append(parts, a.Postcode)
	}

	return strings.Join(parts, ", ")
}

func TransformAddressDetailsToAddress(ad ordnance_survey.AddressDetails) Address {
	a := Address{}

	if len(ad.BuildingName) > 0 {
		a.Line1 = ad.BuildingName

		if len(ad.BuildingNumber) > 0 {
			a.Line2 = fmt.Sprintf("%s %s", ad.BuildingNumber, ad.ThoroughFareName)
		} else {
			a.Line2 = ad.ThoroughFareName
		}

		a.Line3 = ad.DependentLocality
	} else {
		a.Line1 = fmt.Sprintf("%s %s", ad.BuildingNumber, ad.ThoroughFareName)
		a.Line2 = ad.DependentLocality
	}

	a.TownOrCity = ad.Town
	a.Postcode = ad.Postcode

	return a
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
