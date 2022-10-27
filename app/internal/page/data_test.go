package page

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadDate(t *testing.T) {
	date := readDate(time.Date(2020, time.March, 12, 0, 0, 0, 0, time.Local))

	assert.Equal(t, Date{Day: "12", Month: "3", Year: "2020"}, date)
}

type mockDataStore struct {
	data interface{}
	mock.Mock
}

func (m *mockDataStore) Get(ctx context.Context, id string, v interface{}) error {
	data, _ := json.Marshal(m.data)
	json.Unmarshal(data, v)
	return m.Called(ctx, id).Error(0)
}

func (m *mockDataStore) Put(ctx context.Context, id string, v interface{}) error {
	return m.Called(ctx, id, v).Error(0)
}

func TestLpaStoreGet(t *testing.T) {
	ctx := context.Background()

	dataStore := &mockDataStore{}
	dataStore.On("Get", ctx, "an-id").Return(nil)

	lpaStore := &lpaStore{dataStore: dataStore, randomInt: func(x int) int { return x }}

	lpa, err := lpaStore.Get(ctx, "an-id")
	assert.Nil(t, err)
	assert.Equal(t, &Lpa{ID: "10100000"}, lpa)
}

func TestLpaStoreGetWhenExists(t *testing.T) {
	existingLpa := &Lpa{ID: "an-id"}
	ctx := context.Background()

	dataStore := &mockDataStore{data: existingLpa}
	dataStore.On("Get", ctx, "an-id").Return(nil)

	lpaStore := &lpaStore{dataStore: dataStore, randomInt: func(x int) int { return x }}

	lpa, err := lpaStore.Get(ctx, "an-id")
	assert.Nil(t, err)
	assert.Equal(t, existingLpa, lpa)
}

func TestLpaStoreGetWhenDataStoreError(t *testing.T) {
	ctx := context.Background()

	dataStore := &mockDataStore{}
	dataStore.On("Get", ctx, "an-id").Return(expectedError)

	lpaStore := &lpaStore{dataStore: dataStore, randomInt: func(x int) int { return x }}

	_, err := lpaStore.Get(ctx, "an-id")
	assert.Equal(t, expectedError, err)
}

func TestLpaStorePut(t *testing.T) {
	ctx := context.Background()
	lpa := &Lpa{ID: "5"}

	dataStore := &mockDataStore{}
	dataStore.On("Put", ctx, "an-id", lpa).Return(expectedError)

	lpaStore := &lpaStore{dataStore: dataStore}

	err := lpaStore.Put(ctx, "an-id", lpa)
	assert.Equal(t, expectedError, err)
}

func TestGetAttorney(t *testing.T) {
	want := Attorney{ID: "1"}
	otherAttorney := Attorney{ID: "2"}

	lpa := &Lpa{
		Attorneys: []Attorney{
			want,
			otherAttorney,
		},
	}

	got, found := lpa.GetAttorney("1")

	assert.True(t, found)
	assert.Equal(t, want, got)
}

func TestGetAttorneyIdDoesNotMatch(t *testing.T) {
	attorney := Attorney{ID: "1"}
	lpa := &Lpa{
		Attorneys: []Attorney{
			attorney,
		},
	}

	_, found := lpa.GetAttorney("2")

	assert.False(t, found)
}

func TestPutAttorney(t *testing.T) {
	attorney := Attorney{ID: "1"}

	lpa := &Lpa{
		Attorneys: []Attorney{
			attorney,
		},
	}

	updatedAttorney := Attorney{ID: "1", FirstNames: "Bob"}

	updated := lpa.PutAttorney(updatedAttorney)

	assert.True(t, updated)
}

func TestPutAttorneyIdDoesNotMatch(t *testing.T) {
	attorney := Attorney{ID: "2"}

	lpa := &Lpa{
		Attorneys: []Attorney{
			attorney,
		},
	}

	updatedAttorney := Attorney{ID: "1", FirstNames: "Bob"}

	updated := lpa.PutAttorney(updatedAttorney)

	assert.False(t, updated)
}

func TestAttorneysFullNames(t *testing.T) {
	l := &Lpa{
		Attorneys: []Attorney{
			{
				FirstNames: "Bob Alan George",
				LastName:   "Jones",
			},
			{
				FirstNames: "Samantha",
				LastName:   "Smith",
			},
			{
				FirstNames: "Abby Helen",
				LastName:   "Burns-Simpson",
			},
		},
	}

	assert.Equal(t, "Bob Alan George Jones, Samantha Smith and Abby Helen Burns-Simpson", l.AttorneysFullNames())
}

func TestAttorneysFirstNames(t *testing.T) {
	l := &Lpa{
		Attorneys: []Attorney{
			{
				FirstNames: "Bob Alan George",
				LastName:   "Jones",
			},
			{
				FirstNames: "Samantha",
				LastName:   "Smith",
			},
			{
				FirstNames: "Abby Helen",
				LastName:   "Burns-Simpson",
			},
		},
	}

	assert.Equal(t, "Bob Alan George, Samantha and Abby Helen", l.AttorneysFirstNames())
}

func TestConcatSentence(t *testing.T) {
	assert.Equal(t, "Bob Smith, Alice Jones, John Doe and Paul Compton", concatSentence([]string{"Bob Smith", "Alice Jones", "John Doe", "Paul Compton"}))
	assert.Equal(t, "Bob Smith, Alice Jones and John Doe", concatSentence([]string{"Bob Smith", "Alice Jones", "John Doe"}))
	assert.Equal(t, "Bob Smith and John Doe", concatSentence([]string{"Bob Smith", "John Doe"}))
	assert.Equal(t, "Bob Smith", concatSentence([]string{"Bob Smith"}))
}
