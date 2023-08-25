package app

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/dynamo"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
)

// An lpaLink is used to join an actor to an LPA.
type lpaLink struct {
	// PK is the same as the PK for the LPA
	PK string
	// SK is the subKey for the current user
	SK string
	// DonorKey is the donorKey for the donor
	DonorKey string
	// ActorType is the type for the current user
	ActorType actor.Type
}

type dashboardStore struct {
	dynamoClient DynamoClient
}

type keys struct {
	PK, SK string
}

func (k keys) isLpa() bool {
	return strings.HasPrefix(k.SK, donorKey(""))
}

func (k keys) isCertificateProviderDetails() bool {
	return strings.HasPrefix(k.SK, certificateProviderKey(""))
}

func (k keys) isAttorneyDetails() bool {
	return strings.HasPrefix(k.SK, attorneyKey(""))
}

func (s *dashboardStore) GetAll(ctx context.Context) (donor, attorney, certificateProvider []page.LpaAndActorTasks, err error) {
	data, err := page.SessionDataFromContext(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	if data.SessionID == "" {
		return nil, nil, nil, errors.New("donorStore.GetAll requires SessionID")
	}

	var links []lpaLink
	if err := s.dynamoClient.GetAllByGsi(ctx, "ActorIndex", subKey(data.SessionID), &links); err != nil {
		return nil, nil, nil, err
	}

	var searchKeys []dynamo.Key
	keyMap := map[string]actor.Type{}
	for _, key := range links {
		searchKeys = append(searchKeys, dynamo.Key{PK: key.PK, SK: key.DonorKey})

		if key.ActorType == actor.TypeAttorney {
			searchKeys = append(searchKeys, dynamo.Key{PK: key.PK, SK: attorneyKey(data.SessionID)})
		}

		if key.ActorType == actor.TypeCertificateProvider {
			searchKeys = append(searchKeys, dynamo.Key{PK: key.PK, SK: certificateProviderKey(data.SessionID)})
		}

		_, id, _ := strings.Cut(key.PK, "#")
		keyMap[id] = key.ActorType
	}

	if len(searchKeys) == 0 {
		return nil, nil, nil, nil
	}

	lpasOrProvidedDetails, err := s.dynamoClient.GetAllByKeys(ctx, searchKeys)
	if err != nil {
		return nil, nil, nil, err
	}

	certificateProviderMap := make(map[string]page.LpaAndActorTasks)
	attorneyMap := make(map[string]page.LpaAndActorTasks)

	for _, item := range lpasOrProvidedDetails {
		var ks keys
		if err = attributevalue.UnmarshalMap(item, &ks); err != nil {
			return nil, nil, nil, err
		}

		if ks.isLpa() {
			lpa := &page.Lpa{}
			if err := attributevalue.UnmarshalMap(item, lpa); err != nil {
				return nil, nil, nil, err
			}

			switch keyMap[lpa.ID] {
			case actor.TypeDonor:
				donor = append(donor, page.LpaAndActorTasks{Lpa: lpa})
			case actor.TypeAttorney:
				attorneyMap[lpa.ID] = page.LpaAndActorTasks{Lpa: lpa}
			case actor.TypeCertificateProvider:
				certificateProviderMap[lpa.ID] = page.LpaAndActorTasks{Lpa: lpa}
			}
		}
	}

	for _, item := range lpasOrProvidedDetails {
		var ks keys
		if err = attributevalue.UnmarshalMap(item, &ks); err != nil {
			return nil, nil, nil, err
		}

		if ks.isAttorneyDetails() {
			attorneyProvidedDetails := &actor.AttorneyProvidedDetails{}
			if err := attributevalue.UnmarshalMap(item, attorneyProvidedDetails); err != nil {
				return nil, nil, nil, err
			}

			if entry, ok := attorneyMap[attorneyProvidedDetails.LpaID]; ok {
				entry.AttorneyTasks = attorneyProvidedDetails.Tasks
				attorneyMap[attorneyProvidedDetails.LpaID] = entry
				continue
			}
		}

		if ks.isCertificateProviderDetails() {
			certificateProviderProvidedDetails := &actor.CertificateProviderProvidedDetails{}
			if err := attributevalue.UnmarshalMap(item, certificateProviderProvidedDetails); err != nil {
				return nil, nil, nil, err
			}

			if entry, ok := certificateProviderMap[certificateProviderProvidedDetails.LpaID]; ok {
				entry.CertificateProviderTasks = certificateProviderProvidedDetails.Tasks
				certificateProviderMap[certificateProviderProvidedDetails.LpaID] = entry
			}
		}
	}

	for _, value := range certificateProviderMap {
		certificateProvider = append(certificateProvider, value)
	}

	for _, value := range attorneyMap {
		attorney = append(attorney, value)
	}

	byUpdatedAt := func(a, b page.LpaAndActorTasks) int {
		if a.Lpa.UpdatedAt.After(b.Lpa.UpdatedAt) {
			return -1
		}
		return 1
	}

	slices.SortFunc(donor, byUpdatedAt)
	slices.SortFunc(attorney, byUpdatedAt)
	slices.SortFunc(certificateProvider, byUpdatedAt)

	return donor, attorney, certificateProvider, nil
}