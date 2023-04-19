package page

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/notify"
)

var useTestCode = false

type ShareCodeData struct {
	SessionID  string
	LpaID      string
	Identity   bool
	AttorneyID string
}

type ShareCodeSender struct {
	dataStore    DataStore
	notifyClient NotifyClient
	appPublicURL string
	randomString func(int) string
}

func NewShareCodeSender(dataStore DataStore, notifyClient NotifyClient, appPublicURL string, randomString func(int) string) *ShareCodeSender {
	return &ShareCodeSender{
		dataStore:    dataStore,
		notifyClient: notifyClient,
		appPublicURL: appPublicURL,
		randomString: randomString,
	}
}

func (s *ShareCodeSender) UseTestCode() {
	useTestCode = true
}

func (s *ShareCodeSender) SendCertificateProvider(ctx context.Context, template notify.TemplateId, appData AppData, identity bool, lpa *Lpa) error {
	var shareCode string

	if useTestCode {
		shareCode = "abcdef123456"
		useTestCode = false
	} else {
		shareCode = s.randomString(12)
	}

	if err := s.dataStore.Put(ctx, "CERTIFICATEPROVIDERSHARE#"+shareCode, "#METADATA#"+shareCode, ShareCodeData{
		SessionID: appData.SessionID,
		LpaID:     appData.LpaID,
		Identity:  identity,
	}); err != nil {
		return fmt.Errorf("creating sharecode failed: %w", err)
	}

	if _, err := s.notifyClient.Email(ctx, notify.Email{
		TemplateID:   s.notifyClient.TemplateID(template),
		EmailAddress: lpa.CertificateProvider.Email,
		Personalisation: map[string]string{
			"shareCode":         shareCode,
			"cpFullName":        lpa.CertificateProvider.FullName(),
			"donorFirstNames":   lpa.Donor.FirstNames,
			"donorFullName":     lpa.Donor.FullName(),
			"lpaLegalTerm":      appData.Localizer.T(lpa.TypeLegalTermTransKey()),
			"cpLandingPageLink": fmt.Sprintf("%s%s", s.appPublicURL, Paths.CertificateProviderStart),
			"optOutLink":        fmt.Sprintf("%s%s?share-code=%s", s.appPublicURL, Paths.CertificateProviderOptOut, shareCode),
		},
	}); err != nil {
		return fmt.Errorf("email failed: %w", err)
	}

	return nil
}

func (s *ShareCodeSender) SendAttorneys(ctx context.Context, template notify.TemplateId, appData AppData, lpa *Lpa) error {
	for _, attorney := range lpa.Attorneys {
		if attorney.Email == "" {
			continue
		}

		var shareCode string

		if useTestCode {
			shareCode = "abcdef123456"
			useTestCode = false
		} else {
			shareCode = s.randomString(12)
		}

		if err := s.dataStore.Put(ctx, "ATTORNEYSHARE#"+shareCode, "#METADATA#"+shareCode, ShareCodeData{
			SessionID:  appData.SessionID,
			LpaID:      appData.LpaID,
			AttorneyID: attorney.ID,
		}); err != nil {
			return fmt.Errorf("creating attorney share failed: %w", err)
		}

		if _, err := s.notifyClient.Email(ctx, notify.Email{
			TemplateID:   s.notifyClient.TemplateID(template),
			EmailAddress: attorney.Email,
			Personalisation: map[string]string{
				"shareCode":        shareCode,
				"attorneyFullName": attorney.FullName(),
				"donorFirstNames":  lpa.Donor.FirstNames,
				"donorFullName":    lpa.Donor.FullName(),
				"lpaLegalTerm":     appData.Localizer.T(lpa.TypeLegalTermTransKey()),
				"landingPageLink":  fmt.Sprintf("%s%s", s.appPublicURL, Paths.Attorney.Start),
			},
		}); err != nil {
			return fmt.Errorf("email failed: %w", err)
		}
	}

	return nil
}
