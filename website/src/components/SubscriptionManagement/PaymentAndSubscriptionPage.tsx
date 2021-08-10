import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Snackbar from '@material-ui/core/Snackbar';
import Divider from '@material-ui/core/Divider';

import Alert from 'common/components/Alert/Alert';
import AddPaymentMethodForm from 'common/components/Stripe/AddPaymentMethodForm';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading1, Heading3, Heading5 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton, ConfirmationButton, WarningButton } from 'common/components/Button/Button';
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
    deletePaymentMethodForUser,
    DeletePaymentMethodForUserResponse,
    DeletePaymentMethodError
} from 'api/stripe/payment_method';
import {
    DeleteStripeSubscriptionForUserResponse,
    deleteStripeSubscriptionForUser,
} from 'api/stripe/subscription';

const styleClasses = makeStyles({
    paymentMethodDisplayCard: {
        padding: '10px',
        boxSizing: 'border-box',
        height: '100%',
    },
    buttonContainer: {
        width: '100%',
        margin: '10px 0',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        flexDirection: 'column',
    },
    divider: {
        margin: '15px 0',
    },
});

const deletePaymentMethodErrorMessage = {
    [DeletePaymentMethodError.DeleteDefault]: "You cannot delete your default card",
    [DeletePaymentMethodError.OnlyCard]: "You cannot delete your last card!",
}

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

    const [ deletePaymentError, setDeletePaymentError ] = useState<string | null>(null);

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
    const deletePaymentMethod = (paymentMethodID: string) => {
        setIsLoadingPaymentMethods(true);
        deletePaymentMethodForUser({
            stripePaymentMethodId: paymentMethodID,
        },
        (resp: DeletePaymentMethodForUserResponse) => {
            setIsLoadingPaymentMethods(false);
            if (!resp.error) {
                setPaymentMethods(paymentMethods.filter((p: PaymentMethod) => p.stripePaymentMethodId !== paymentMethodID));
            } else {
                setDeletePaymentError(deletePaymentMethodErrorMessage[resp.error]);
            }
        },
        (err: Error) => {
            setIsLoadingPaymentMethods(false);
            setError(err);
        });
    }

    const classes = styleClasses();
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
                    deletePaymentMethod={deletePaymentMethod}
                    handleIsLoadingStripeRequest={setIsLoadingStripeRequest}
                    handlePaymentMethodAddedSuccess={handleSuccessfullyAddedPaymentMethod}
                    handleAddPaymentMethodError={setError} />
                <Divider className={classes.divider} />
                <CancelSubscriptionButton />
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

    deletePaymentMethod: (paymentMethodID) => void;
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
            deletePaymentMethod={props.deletePaymentMethod}
            {...p} />
    ));
    const classes = styleClasses();
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
                    <div className={classes.buttonContainer}>
                        <PrimaryButton
                            disabled={props.isLoadingStripeRequest}
                            onClick={() => setShowAddPaymentMethodForm(true)}>
                            Add new payment method
                        </PrimaryButton>
                    </div>
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
    deletePaymentMethod: (paymentMethodID) => void;
    setDefaultPaymentMethod: (paymentMethodID) => void;
} & PaymentMethod;

const PaymentMethodDisplay = (props: PaymentMethodDisplayProps) => {
    const [ showDeleteConfirmation, setShowDeleteConfirmation ] = useState<boolean>(false);

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
                {
                    !props.isDefault && (
                        <Grid container>
                            <Grid item xs={12} md={6}>
                                <PrimaryButton onClick={() => props.setDefaultPaymentMethod(props.stripePaymentMethodId)}>
                                    Make default
                                </PrimaryButton>
                            </Grid>
                            <Grid item xs={12} md={6}>
                            {
                                showDeleteConfirmation ? (
                                    <ConfirmationButton onClick={() => props.deletePaymentMethod(props.stripePaymentMethodId)}>
                                        Confirm
                                    </ConfirmationButton>
                                ) : (
                                    <WarningButton onClick={() => setShowDeleteConfirmation(true)}>
                                        Delete
                                    </WarningButton>
                                )
                            }
                            </Grid>
                        </Grid>
                    )
                }
                {
                    props.isDefault && (
                        <Heading5 color={TypographyColor.Primary}>
                            This is your default payment method
                        </Heading5>
                    )
                }
            </Card>
        </Grid>
    );
}

type CancelSubscriptionButtonProps = {};

const CancelSubscriptionButton = (props:  CancelSubscriptionButtonProps) => {
    const [ showConfirmation, setShowConfirmation ] = useState<boolean>(false);

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ didDelete, setDidDelete ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const handleCancelSubscription = () => {
        setIsLoading(true);
        deleteStripeSubscriptionForUser({},
        (resp: DeleteStripeSubscriptionForUserResponse) => {
            setIsLoading(false);
            setDidDelete(true);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }
    const classes = styleClasses();
    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (didDelete) {
        body = (
            <Paragraph>
                Your subscription will be canceled at the end of the current period! If you want to cancel now and receive a refund, reach out to hello@babblegraph.com via email!
            </Paragraph>
        )
    } else {
        body = (
            <div className={classes.buttonContainer}>
                <Paragraph>
                    You can cancel your subscription here. You will retain your premium benefits until the end of the current period, after which you won’t be billed. You will still receive the daily newsletter. To unsubscribe, go to the unsubscribe tab on the previous page.
                </Paragraph>
                {
                    showConfirmation ? (
                        <WarningButton onClick={handleCancelSubscription}>
                            Confirm Subscription Cancellation.
                        </WarningButton>
                    ) : (
                        <PrimaryButton onClick={() => setShowConfirmation(true)}>
                            Cancel Subscription
                        </PrimaryButton>
                    )
                }
            </div>
        );
    }
    return (
        <div>
            <Heading3 color={TypographyColor.Primary}>
                Cancel your subscription
            </Heading3>
            { body }
        </div>
    )
}

export default PaymentAndSubscriptionPage;
