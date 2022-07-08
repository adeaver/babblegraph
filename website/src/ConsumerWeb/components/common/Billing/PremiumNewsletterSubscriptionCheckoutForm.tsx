import React, { useState } from 'react';

import { Heading3 } from 'common/typography/Heading';
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    PaymentState,
    PremiumNewsletterSubscription
} from 'common/api/billing/billing';

import ResolvePaymentIntentForm from './stripe/ResolvePaymentIntentForm';
import ResolveSetupIntentForm from './stripe/ResolveSetupIntentForm';

type PremiumNewsletterSubscriptionCardFormProps = {
    premiumNewsletterSusbcription: PremiumNewsletterSubscription;

    subscriptionManagementToken?: string;
}

const PremiumNewsletterSubscriptionCardForm = (props: PremiumNewsletterSubscriptionCardFormProps) => {
    return (
        <div>
            <Heading3>
                This payment form is no longer accessible
            </Heading3>
        </div>
    );
}

export default PremiumNewsletterSubscriptionCardForm;
