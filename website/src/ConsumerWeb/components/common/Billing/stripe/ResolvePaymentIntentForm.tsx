import React, { useState } from 'react';

import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';

import { withStripe, WithStripeProps } from './withStripe';
import GenericCardForm from './GenericCardForm';

type ResolvePaymentIntentFormProps = {
    stripePaymentIntentClientSecret: string;
}

const ResolvePaymentIntentForm = withStripe<ResolvePaymentIntentFormProps>(
    (props: ResolvePaymentIntentFormProps & WithStripeProps) => {
        const [ cardholderName, setCardholderName ] = useState<string>(null);
        const [ postalCode, setPostalCode ] = useState<string>(null);
        const [ isLoading, setIsLoading ] = useState<boolean>(false);

        const handleSubmit = () => {
            setIsLoading(true);
            // TODO: add endpoint to insert sync request
        }

        return (
            <Form handleSubmit={handleSubmit}>
                <GenericCardForm
                    cardholderName={cardholderName}
                    postalCode={postalCode}
                    isDisabled={isLoading}
                    setCardholderName={setCardholderName}
                    setPostalCode={setPostalCode} />
                <PrimaryButton
                    type='submit'
                    disabled={isLoading}>
                    Pay
                </PrimaryButton>
            </Form>
        );
    }
);

export default ResolvePaymentIntentForm;
