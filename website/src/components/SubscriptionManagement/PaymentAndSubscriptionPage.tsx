import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import AddPaymentMethodForm from 'common/components/Stripe/AddPaymentMethodForm';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading1, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton, WarningButton } from 'common/components/Button/Button';
import CardIcon from 'common/components/Payment/CardIcon';

import { ContentHeader } from './common';

import {
    getPaymentMethodByID,
    GetPaymentMethodByIDResponse,
    getPaymentMethodsForUser,
    GetPaymentMethodsForUserResponse,
    PaymentMethod,
    setDefaultPaymentMethodForUser,
    SetDefaultPaymentMethodForUserResponse,
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

    const [ isLoadingStripeRequest, setIsLoadingStripeRequest ] = useState<boolean>(false);

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

    const handleSuccessfullyAddedPaymentMethod = (paymentMethodID: string) => {
        getPaymentMethodByID({
            stripePaymentMethodId: paymentMethodID,
        },
        (resp: GetPaymentMethodByIDResponse) => {
            setIsLoadingStripeRequest(false);
            resp.paymentMethod && setPaymentMethods(paymentMethods.concat(resp.paymentMethod));
        },
        (err: Error) => {
            setIsLoadingStripeRequest(false);
            setError(err);
        })
    }
    const setDefaultPaymentMethod = (paymentMethodID: string) => {
        setIsLoadingPaymentMethods(true);
        setDefaultPaymentMethodForUser({
            stripePaymentMethodId: paymentMethodID,
        },
        (resp: SetDefaultPaymentMethodForUserResponse) => {
            setIsLoadingPaymentMethods(false);
            setPaymentMethods(paymentMethods.map((p: PaymentMethod) => ({
                ...p,
                isDefault: p.stripePaymentMethodId === paymentMethodID,
            })));
        },
        (err: Error) => {
            setIsLoadingPaymentMethods(false);
            setError(err);
        });
    }

    const isLoading = isLoadingPaymentMethods;
    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = (
            <Paragraph>
                There was a problem loading your payment information. Try again later!
            </Paragraph>
        );
    } else {
        body = (
            <div>
                <PaymentMethodsDisplay
                    paymentMethods={paymentMethods}
                    isLoadingStripeRequest={isLoadingStripeRequest}
                    setDefaultPaymentMethod={setDefaultPaymentMethod}
                    handleIsLoadingStripeRequest={setIsLoadingStripeRequest}
                    handlePaymentMethodAddedSuccess={handleSuccessfullyAddedPaymentMethod}
                    handleAddPaymentMethodError={setError} />
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
    isLoadingStripeRequest: boolean;

    setDefaultPaymentMethod: (paymentMethodID) => void;

    handleIsLoadingStripeRequest: (isLoading: boolean) => void;
    handlePaymentMethodAddedSuccess: (paymentMethodID: string) => void;
    handleAddPaymentMethodError: (err: Error) => void;
}

const PaymentMethodsDisplay = (props: PaymentMethodsDisplayProps) => {
    const [ showAddPaymentMethodForm, setShowAddPaymentMethodForm ] = useState<boolean>(false);
    const [ addPaymentMethodFailure, setAddPaymentMethodFailure ] = useState<string | null>(null);

    const handlePaymentMethodAddedSuccess = (paymentMethodID: string) => {
        props.handlePaymentMethodAddedSuccess(paymentMethodID);
        setShowAddPaymentMethodForm(false);
    }

    const paymentMethods = props.paymentMethods.map((p: PaymentMethod) => (
        <PaymentMethodDisplay
            key={p.stripePaymentMethodId}
            setDefaultPaymentMethod={props.setDefaultPaymentMethod}
            {...p} />
    ));
    return (
        <div>
            <Heading3 color={TypographyColor.Primary}>
                Your Payment Methods
            </Heading3>
            {
                showAddPaymentMethodForm ? (
                    <AddPaymentMethodForm
                        handleIsStripeRequestLoading={props.handleIsLoadingStripeRequest}
                        handleSuccess={handlePaymentMethodAddedSuccess}
                        handleFailure={setAddPaymentMethodFailure}
                        handleError={props.handleAddPaymentMethodError} />

                ) : (
                    <PrimaryButton
                        disabled={props.isLoadingStripeRequest}
                        onClick={() => setShowAddPaymentMethodForm(true)}>
                        Add new payment method
                    </PrimaryButton>
                )
            }
            {
                props.isLoadingStripeRequest ? (
                    <LoadingSpinner />
                ) : (
                    <Grid container spacing={2}>
                        { paymentMethods }
                    </Grid>
                )
            }
            <Snackbar open={!!addPaymentMethodFailure} close={() => setAddPaymentMethodFailure(null)} autoHideDuration={6000}>
                <Alert severity="error">{addPaymentMethodFailure}</Alert>
            </Snackbar>
        </div>
    );
}

type PaymentMethodDisplayProps = {
    setDefaultPaymentMethod: (paymentMethodID) => void;
} & PaymentMethod;

const PaymentMethodDisplay = (props: PaymentMethodDisplayProps) => {
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
                    {
                        !props.isDefault && (
                            <Grid item xs={12} md={6}>
                                <PrimaryButton onClick={() => props.setDefaultPaymentMethod(props.stripePaymentMethodId)}>
                                    Make default
                                </PrimaryButton>
                            </Grid>
                        )
                    }
                    <Grid item xs={12} md={6}>
                        <WarningButton>
                            Delete
                        </WarningButton>
                    </Grid>
                </Grid>
                {
                    props.isDefault && (
                        <Heading4 color={TypographyColor.Primary}>
                            This is your default payment method
                        </Heading4>
                    )
                }
            </Card>
        </Grid>
    );
}

export default PaymentAndSubscriptionPage;
