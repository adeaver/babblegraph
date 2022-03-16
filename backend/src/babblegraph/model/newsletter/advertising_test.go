package newsletter

import (
	"babblegraph/model/advertising"
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/util/ptr"
)

type testAdvertisementAccessor struct {
	isUserEligibleForAdvertisement bool

	advertisements []advertising.Advertisement
}

func (t *testAdvertisementAccessor) IsEligibleForAdvertisement() bool {
	return t.isUserEligibleForAdvertisement
}

func (t *testAdvertisementAccessor) LookupAdvertisementForTopic(topic content.TopicID) (*advertising.Advertisement, error) {
	if len(t.advertisements) == 0 {
		return nil, nil
	}
	return &t.advertisements[0], nil
}

func (t *testAdvertisementAccessor) LookupGeneralAdvertisement() (*advertising.Advertisement, error) {
	if len(t.advertisements) == 0 {
		return nil, nil
	}
	return &t.advertisements[0], nil
}

func (t *testAdvertisementAccessor) GetAdvertisementURL(emailRecordID email.ID, ad advertising.Advertisement) (*string, error) {
	return ptr.String("babblegraph.com/advertisement-url"), nil
}
