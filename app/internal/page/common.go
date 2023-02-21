package page

import (
	"context"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/notify"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/onelogin"
)

const FormUrlEncoded = "application/x-www-form-urlencoded"

//go:generate mockery --testonly --inpackage --name Logger --structname mockLogger
type Logger interface {
	Print(v ...interface{})
}

type DataStore interface {
	GetAll(context.Context, string, interface{}) error
	Get(context.Context, string, string, interface{}) error
	Put(context.Context, string, string, interface{}) error
}

//go:generate mockery --testonly --inpackage --name NotifyClient --structname mockNotifyClient
type NotifyClient interface {
	Email(ctx context.Context, email notify.Email) (string, error)
	Sms(ctx context.Context, sms notify.Sms) (string, error)
	TemplateID(id notify.TemplateId) string
}

//go:generate mockery --testonly --inpackage --name OneLoginClient --structname mockOneLoginClient
type OneLoginClient interface {
	AuthCodeURL(state, nonce, locale string, identity bool) string
	Exchange(ctx context.Context, code, nonce string) (string, error)
	UserInfo(ctx context.Context, accessToken string) (onelogin.UserInfo, error)
	ParseIdentityClaim(ctx context.Context, userInfo onelogin.UserInfo) (identity.UserData, error)
}

//go:generate mockery --testonly --inpackage --name LpaStore --structname mockLpaStore
type LpaStore interface {
	Create(context.Context) (*Lpa, error)
	GetAll(context.Context) ([]*Lpa, error)
	Get(context.Context) (*Lpa, error)
	Put(context.Context, *Lpa) error
}

func PostFormString(r *http.Request, name string) string {
	return strings.TrimSpace(r.PostFormValue(name))
}

type Handler func(data AppData, w http.ResponseWriter, r *http.Request) error
