import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton, WarningButton } from 'common/components/Button/Button';
import CardIcon from 'common/components/Payment/CardIcon';

import { ContentHeader } from './common';

import {
    getPaymentMethodsForUser,
    GetPaymentMethodsForUserResponse,
    PaymentMethod
} from 'api/stripe/payment_method';

const styleClasses = makeStyles({
    paymentMethodDisplayCard: {
        padding: '10px',
    },
});

type Params = {
    token: string;
}

type PaymentAndSubscriptionPageProps = RouteComponentProps<Params>

const PaymentAndSubscriptionPage = (props: PaymentAndSubscriptionPageProps) => {
    const { token } = props.match.params;

    const [ isLoadingPaymentMethods, setIsLoadingPaymentMethods ] = useState<boolean>(true);
    const [ paymentMethods, setPaymentMethods ] = useState<Array<PaymentMethod>>([]);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getPaymentMethodsForUser({},
            (resp: GetPaymentMethodsForUserResponse) => {
                setIsLoadingPaymentMethods(false);
                setPaymentMethods(resp.paymentMethods);
            },
            (err: Error) => {
                setIsLoadingPaymentMethods(false);
                setError(err);
            });
    }, []);

    const isLoading = isLoadingPaymentMethods;
    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Primary}>
                There was a problem loading your payment information. Try again later!
            </Heading3>
        );
    } else {
        body = (
            <div>
                <PaymentMethodsDisplay
                    paymentMethods={paymentMethods} />
            </div>
        )
    }
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <ContentHeader
                            title="Subscription and Payment Settings"
                            token={token} />
                            { body }
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

type PaymentMethodsDisplayProps = {
    paymentMethods: Array<PaymentMethod>;
}

const PaymentMethodsDisplay = (props: PaymentMethodsDisplayProps) => {
    const paymentMethods = props.paymentMethods.map((p: PaymentMethod) => (
        <PaymentMethodDisplay key={p.stripePaymentMethodId} {...p} />
    ));
    return (
        <div>
            <Heading3 color={TypographyColor.Primary}>
                Your Payment Methods
            </Heading3>
            <Grid container>
                { paymentMethods }
            </Grid>
        </div>
    );
}

const PaymentMethodDisplay = (props: PaymentMethod) => {
    const classes = styleClasses();
    return (
        <Grid item xs={12} md={6}>
            <Card className={classes.paymentMethodDisplayCard}>
                <CardIcon cardType={props.cardType} />
                <Paragraph align={Alignment.Left}>
                    Ending in { props.lastFourDigits }
                </Paragraph>
                <Paragraph align={Alignment.Left}>
                    Expires { props.expirationMonth }/{ props.expirationYear }
                </Paragraph>
                <Grid container>
                    <Grid item xs={12} md={6}>
                        <PrimaryButton>
                            Set as default
                        </PrimaryButton>
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <WarningButton>
                            Delete Payment Method
                        </WarningButton>
                    </Grid>
                </Grid>
            </Card>
        </Grid>
    );
}

export default PaymentAndSubscriptionPage;
