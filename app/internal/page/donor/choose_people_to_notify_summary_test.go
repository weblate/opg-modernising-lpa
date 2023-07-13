package donor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/form"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetChoosePeopleToNotifySummary(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpa := &page.Lpa{PeopleToNotify: actor.PeopleToNotify{{}}}

	template := newMockTemplate(t)
	template.
		On("Execute", w, &choosePeopleToNotifySummaryData{
			App:     testAppData,
			Lpa:     lpa,
			Form:    &form.YesNoForm{},
			Options: form.YesNoValues,
		}).
		Return(nil)

	err := ChoosePeopleToNotifySummary(template.Execute)(testAppData, w, r, lpa)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetChoosePeopleToNotifySummaryWhenNoPeopleToNotify(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	err := ChoosePeopleToNotifySummary(nil)(testAppData, w, r, &page.Lpa{
		ID: "lpa-id",
		Tasks: page.Tasks{
			YourDetails:                actor.TaskCompleted,
			ChooseAttorneys:            actor.TaskCompleted,
			ChooseReplacementAttorneys: actor.TaskCompleted,
			WhenCanTheLpaBeUsed:        actor.TaskCompleted,
			Restrictions:               actor.TaskCompleted,
			CertificateProvider:        actor.TaskCompleted,
		},
	})

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.DoYouWantToNotifyPeople.Format("lpa-id"), resp.Header.Get("Location"))
}

func TestPostChoosePeopleToNotifySummaryAddPersonToNotify(t *testing.T) {
	form := url.Values{
		"yes-no": {form.Yes.String()},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	err := ChoosePeopleToNotifySummary(nil)(testAppData, w, r, &page.Lpa{ID: "lpa-id", PeopleToNotify: actor.PeopleToNotify{{ID: "123"}}})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.ChoosePeopleToNotify.Format("lpa-id")+"?addAnother=1", resp.Header.Get("Location"))
}

func TestPostChoosePeopleToNotifySummaryNoFurtherPeopleToNotify(t *testing.T) {
	form := url.Values{
		"yes-no": {form.No.String()},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	err := ChoosePeopleToNotifySummary(nil)(testAppData, w, r, &page.Lpa{
		ID:             "lpa-id",
		PeopleToNotify: actor.PeopleToNotify{{ID: "123"}},
		Tasks: page.Tasks{
			YourDetails:                actor.TaskCompleted,
			ChooseAttorneys:            actor.TaskCompleted,
			ChooseReplacementAttorneys: actor.TaskCompleted,
			WhenCanTheLpaBeUsed:        actor.TaskCompleted,
			Restrictions:               actor.TaskCompleted,
			CertificateProvider:        actor.TaskCompleted,
			PeopleToNotify:             actor.TaskCompleted,
		},
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.TaskList.Format("lpa-id"), resp.Header.Get("Location"))
}

func TestPostChoosePeopleToNotifySummaryFormValidation(t *testing.T) {
	form := url.Values{
		"yes-no": {""},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	validationError := validation.With("yes-no", validation.SelectError{Label: "yesToAddAnotherPersonToNotify"})

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *choosePeopleToNotifySummaryData) bool {
			return assert.Equal(t, validationError, data.Errors)
		})).
		Return(nil)

	err := ChoosePeopleToNotifySummary(template.Execute)(testAppData, w, r, &page.Lpa{PeopleToNotify: actor.PeopleToNotify{{}}})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
