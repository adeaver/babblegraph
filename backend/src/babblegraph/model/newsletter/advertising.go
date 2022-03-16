package newsletter

import (
	"babblegraph/model/advertising"
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"time"

	"github.com/jmoiron/sqlx"
)

func lookupAdvertisement(c ctx.LogContext, emailRecordID email.ID, userAccessor userPreferencesAccessor, advertisementAccessor advertisementAccessor) (*NewsletterAdvertisement, error) {
	switch {
	case !advertisementAccessor.IsEligibleForAdvertisement():
		c.Debugf("User is not eligible for ad")
		return nil, nil
	case userAccessor.getUserCreatedDate().Add(advertising.MinimumUserAccountAge).After(time.Now()):
		c.Debugf("User account is too new for ad")
	default:
		c.Debugf("Finding ad for user")
	}
	premiumInformationLink, err := routes.MakePremiumInformationLink(userAccessor.getUserID())
	if err != nil {
		return nil, err
	}
	advertisingPolicyLink := routes.GetAdvertisingPolicyURL()
	for _, t := range userAccessor.getUserTopics() {
		advertisement, err := advertisementAccessor.LookupAdvertisementForTopic(c, t)
		switch {
		case err != nil:
			c.Errorf("Error getting advertisement for user %s: %s", userAccessor.getUserID(), err.Error())
			return nil, err
		case advertisement != nil:
			c.Debugf("ad found for topic %s", t)
			advertisementURL, err := advertisementAccessor.GetAdvertisementURL(emailRecordID, *advertisement)
			if err != nil {
				c.Errorf("Error inserting advertisement for user %s: %s", userAccessor.getUserID(), err.Error())
				return nil, err
			}
			return &NewsletterAdvertisement{
				Title:                   advertisement.Title,
				Description:             advertisement.Description,
				ImageURL:                advertisement.ImageURL,
				URL:                     *advertisementURL,
				PremiumLink:             *premiumInformationLink,
				AdvertisementPolicyLink: advertisingPolicyLink,
			}, nil
		}
	}
	advertisement, err := advertisementAccessor.LookupGeneralAdvertisement(c)
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
			Title:                   advertisement.Title,
			Description:             advertisement.Description,
			ImageURL:                advertisement.ImageURL,
			URL:                     *advertisementURL,
			PremiumLink:             *premiumInformationLink,
			AdvertisementPolicyLink: advertisingPolicyLink,
		}, nil
	}
	return nil, nil
}

type advertisementAccessor interface {
	IsEligibleForAdvertisement() bool

	LookupAdvertisementForTopic(c ctx.LogContext, t content.TopicID) (*advertising.Advertisement, error)
	LookupGeneralAdvertisement(c ctx.LogContext) (*advertising.Advertisement, error)
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

func (d *DefaultAdvertisementAccessor) LookupAdvertisementForTopic(c ctx.LogContext, t content.TopicID) (*advertising.Advertisement, error) {
	return advertising.QueryAdvertisementsForUser(c, d.tx, d.userID, t.Ptr(), d.languageCode, d.userAdvertisementEligibility.IneligibleCampaignIDs)
}

func (d *DefaultAdvertisementAccessor) LookupGeneralAdvertisement(c ctx.LogContext) (*advertising.Advertisement, error) {
	return advertising.QueryAdvertisementsForUser(c, d.tx, d.userID, nil, d.languageCode, d.userAdvertisementEligibility.IneligibleCampaignIDs)
}

func (d *DefaultAdvertisementAccessor) GetAdvertisementURL(emailRecordID email.ID, ad advertising.Advertisement) (*string, error) {
	id, err := advertising.InsertUserAdvertisementAndGetID(d.tx, d.userID, ad, emailRecordID)
	if err != nil {
		return nil, err
	}
	return ptr.String(id.GetAdvertisementURL()), nil
}
