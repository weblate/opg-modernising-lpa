package certificateprovider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/onelogin"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/random"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/sesh"
)

//go:generate mockery --testonly --inpackage --name Logger --structname mockLogger
type Logger interface {
	Print(v ...interface{})
}

//go:generate mockery --testonly --inpackage --name LpaStore --structname mockLpaStore
type LpaStore interface {
	Create(context.Context) (*page.Lpa, error)
	GetAll(context.Context) ([]*page.Lpa, error)
	Get(context.Context) (*page.Lpa, error)
	Put(context.Context, *page.Lpa) error
}

//go:generate mockery --testonly --inpackage --name OneLoginClient --structname mockOneLoginClient
type OneLoginClient interface {
	AuthCodeURL(state, nonce, locale string, identity bool) string
	Exchange(ctx context.Context, code, nonce string) (string, error)
	UserInfo(ctx context.Context, accessToken string) (onelogin.UserInfo, error)
	ParseIdentityClaim(ctx context.Context, userInfo onelogin.UserInfo) (identity.UserData, error)
}

type DataStore interface {
	GetAll(context.Context, string, interface{}) error
	Get(context.Context, string, string, interface{}) error
	Put(context.Context, string, string, interface{}) error
}

//go:generate mockery --testonly --inpackage --name AddressClient --structname mockAddressClient
type AddressClient interface {
	LookupPostcode(ctx context.Context, postcode string) ([]place.Address, error)
}

func Register(
	rootMux *http.ServeMux,
	logger Logger,
	tmpls template.Templates,
	sessionStore sesh.Store,
	lpaStore LpaStore,
	oneLoginClient OneLoginClient,
	dataStore DataStore,
	addressClient AddressClient,
) {
	handleRoot := makeHandle(rootMux, logger, sessionStore, None)

	handleRoot(page.Paths.CertificateProviderStart, None,
		Start(tmpls.Get("certificate_provider_start.gohtml"), lpaStore, dataStore))
	handleRoot(page.Paths.CertificateProviderLogin, None,
		Login(logger, oneLoginClient, sessionStore, random.String))
	handleRoot(page.Paths.CertificateProviderLoginCallback, None,
		LoginCallback(tmpls.Get("identity_with_one_login_callback.gohtml"), oneLoginClient, sessionStore, lpaStore))
	handleRoot(page.Paths.CertificateProviderYourDetails, RequireSession,
		YourDetails(tmpls.Get("certificate_provider_your_details.gohtml"), lpaStore))
	handleRoot(page.Paths.CertificateProviderYourAddress, RequireSession,
		YourAddress(logger, tmpls.Get("your_address.gohtml"), addressClient, lpaStore))
	handleRoot(page.Paths.CertificateProviderReadTheLpa, RequireSession,
		page.Guidance(tmpls.Get("certificate_provider_read_the_lpa.gohtml"), lpaStore))
	handleRoot(page.Paths.CertificateProviderGuidance, RequireSession,
		page.Guidance(tmpls.Get("certificate_provider_guidance.gohtml"), lpaStore))
	handleRoot(page.Paths.CertificateProviderConfirmation, RequireSession,
		page.Guidance(tmpls.Get("certificate_provider_confirmation.gohtml"), lpaStore))
	handleRoot(page.Paths.ProvideCertificate, RequireSession,
		ProvideCertificate(tmpls.Get("provide_certificate.gohtml"), lpaStore, time.Now))
	handleRoot(page.Paths.CertificateProvided, RequireSession,
		page.Guidance(tmpls.Get("certificate_provided.gohtml"), lpaStore))
}

type handleOpt byte

const (
	None handleOpt = 1 << iota
	RequireSession
	CanGoBack
)

func makeHandle(mux *http.ServeMux, logger Logger, store sesh.Store, defaultOptions handleOpt) func(string, handleOpt, page.Handler) {
	return func(path string, opt handleOpt, h page.Handler) {
		opt = opt | defaultOptions

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			appData := page.AppDataFromContext(ctx)
			appData.Page = path
			appData.CanGoBack = opt&CanGoBack != 0

			if opt&RequireSession != 0 {
				session, err := sesh.CertificateProvider(store, r)
				if err != nil {
					logger.Print(err)
					http.Redirect(w, r, page.Paths.CertificateProviderStart, http.StatusFound)
					return
				}

				appData.SessionID = session.DonorSessionID
				appData.LpaID = session.LpaID

				ctx = page.ContextWithSessionData(ctx, &page.SessionData{SessionID: appData.SessionID, LpaID: appData.LpaID})
			}

			if err := h(appData, w, r.WithContext(page.ContextWithAppData(ctx, appData))); err != nil {
				str := fmt.Sprintf("Error rendering page for path '%s': %s", path, err.Error())

				logger.Print(str)
				http.Error(w, "Encountered an error", http.StatusInternalServerError)
			}
		})
	}
}
