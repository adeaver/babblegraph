import React from 'react';

import { withStripe, WithStripeProps } from './withStripe';

type ResolvePaymentIntentFormProps = {}

const ResolvePaymentIntentForm = withStripe(
    (props: ResolvePaymentIntentFormProps & WithStripeProps) => {
        return <div />;
    }
);

export default ResolvePaymentIntentForm;
