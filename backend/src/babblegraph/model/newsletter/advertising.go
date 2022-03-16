package newsletter

import (
	"babblegraph/model/advertising"
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

func lookupAdvertisement(c ctx.LogContext, emailRecordID email.ID, userAccessor userPreferencesAccessor, advertisementAccessor advertisementAccessor) (*NewsletterAdvertisement, error) {
	if !advertisementAccessor.IsEligibleForAdvertisement() {
		return nil, nil
	}
	for _, t := range userAccessor.getUserTopics() {
		advertisement, err := advertisementAccessor.LookupAdvertisementForTopic(t)
		switch {
		case err != nil:
			c.Errorf("Error getting advertisement for user %s: %s", userAccessor.getUserID(), err.Error())
			return nil, err
		case advertisement != nil:
			advertisementURL, err := advertisementAccessor.GetAdvertisementURL(emailRecordID, *advertisement)
			if err != nil {
				c.Errorf("Error inserting advertisement for user %s: %s", userAccessor.getUserID(), err.Error())
				return nil, err
			}
			return &NewsletterAdvertisement{
				Title:       advertisement.Title,
				Description: advertisement.Description,
				ImageURL:    advertisement.ImageURL,
				URL:         *advertisementURL,
			}, nil
		}
	}
	advertisement, err := advertisementAccessor.LookupGeneralAdvertisement()
	switch {
	case err != nil:
		c.Errorf("Error getting advertisement for user %s: %s", userAccessor.getUserID(), err.Error())
		return nil, err
	case advertisement != nil:
		advertisementURL, err := advertisementAccessor.GetAdvertisementURL(emailRecordID, *advertisement)
		if err != nil {
			c.Errorf("Error inserting advertisement for user %s: %s", userAccessor.getUserID(), err.Error())
			return nil, err
		}
		return &NewsletterAdvertisement{
			Title:       advertisement.Title,
			Description: advertisement.Description,
			ImageURL:    advertisement.ImageURL,
			URL:         *advertisementURL,
		}, nil
	}
	return nil, nil
}

type advertisementAccessor interface {
	IsEligibleForAdvertisement() bool

	LookupAdvertisementForTopic(t content.TopicID) (*advertising.Advertisement, error)
	LookupGeneralAdvertisement() (*advertising.Advertisement, error)
	GetAdvertisementURL(emailRecordID email.ID, ad advertising.Advertisement) (*string, error)
}

type DefaultAdvertisementAccessor struct {
	tx *sqlx.Tx

	userID       users.UserID
	languageCode wordsmith.LanguageCode

	userAdvertisementEligibility advertising.UserAdvertisementEligibility
}

func GetDefaultAdvertisementAccessor(tx *sqlx.Tx, userID users.UserID, languageCode wordsmith.LanguageCode) (*DefaultAdvertisementAccessor, error) {
	userEligibility, err := advertising.GetUserAdvertisementEligibility(tx, userID)
	if err != nil {
		return nil, err
	}
	return &DefaultAdvertisementAccessor{
		tx:                           tx,
		userID:                       userID,
		languageCode:                 languageCode,
		userAdvertisementEligibility: *userEligibility,
	}, nil
}

func (d *DefaultAdvertisementAccessor) IsEligibleForAdvertisement() bool {
	return d.userAdvertisementEligibility.IsUserEligibleForAdvertisement
}

func (d *DefaultAdvertisementAccessor) LookupAdvertisementForTopic(t content.TopicID) (*advertising.Advertisement, error) {
	return advertising.QueryAdvertisementsForUser(d.tx, d.userID, t.Ptr(), d.languageCode, d.userAdvertisementEligibility.IneligibleCampaignIDs)
}

func (d *DefaultAdvertisementAccessor) LookupGeneralAdvertisement() (*advertising.Advertisement, error) {
	return advertising.QueryAdvertisementsForUser(d.tx, d.userID, nil, d.languageCode, d.userAdvertisementEligibility.IneligibleCampaignIDs)
}

func (d *DefaultAdvertisementAccessor) GetAdvertisementURL(emailRecordID email.ID, ad advertising.Advertisement) (*string, error) {
	id, err := advertising.InsertUserAdvertisementAndGetID(d.tx, d.userID, ad, emailRecordID)
	if err != nil {
		return nil, err
	}
	return ptr.String(id.GetAdvertisementURL()), nil
}
