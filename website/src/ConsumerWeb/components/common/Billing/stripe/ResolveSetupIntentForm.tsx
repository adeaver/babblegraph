import React, { useState } from 'react';

import Grid from '@material-ui/core/Grid';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';

import {
    StripeBeginPaymentMethodSetupResponse,
    stripeBeginPaymentMethodSetup,
} from 'ConsumerWeb/api/billing/stripe';

import { withStripe, WithStripeProps } from './withStripe';
import GenericCardForm from './GenericCardForm';

type ResolveSetupIntentFormOwnProps = {}

const ResolveSetupIntentForm = asBaseComponent<StripeBeginPaymentMethodSetupResponse, ResolveSetupIntentFormOwnProps>(
    withStripe<ResolveSetupIntentFormOwnProps & StripeBeginPaymentMethodSetupResponse & BaseComponentProps>(
        (props: ResolveSetupIntentFormOwnProps & StripeBeginPaymentMethodSetupResponse & BaseComponentProps & WithStripeProps) => {
            const [ cardholderName, setCardholderName ] = useState<string>(null);
            const [ postalCode, setPostalCode ] = useState<string>(null);
            const [ isLoading, setIsLoading ] = useState<boolean>(false);

            const handleSubmit = () => {
                setIsLoading(true);
                // TODO: add endpoint to insert sync request?
            }

            return (
                <Form handleSubmit={handleSubmit}>
                    <Grid container>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid item xs={12} md={6}>
                            <GenericCardForm
                                cardholderName={cardholderName}
                                postalCode={postalCode}
                                isDisabled={isLoading}
                                setCardholderName={setCardholderName}
                                setPostalCode={setPostalCode} />
                        </Grid>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid item xs={12} md={6}>
                            <PrimaryButton
                                type='submit'
                                disabled={isLoading}>
                                Add Payment Method
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                </Form>
            );
        }
    ),
    (
        ownProps: ResolveSetupIntentFormOwnProps,
        onSuccess: (resp: StripeBeginPaymentMethodSetupResponse) => void,
        onError: (err: Error) => void,
    ) => stripeBeginPaymentMethodSetup({}, onSuccess, onError),
    false,
);

export default ResolveSetupIntentForm;
