import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading1, Heading3, Heading4 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';

import {
    PremiumNewsletterSubscription,
    PaymentState,
    PromotionType,
} from 'common/api/billing/billing';
import { SubscriptionLevel } from 'common/api/useraccounts/useraccounts';
import {
    UserBillingInformation,

    GetBillingInformationForEmailAddressResponse,
    getBillingInformationForEmailAddress,

    ForceSyncForUserResponse,
    forceSyncForUser,

    CreatePromotionCodeResponse,
    createPromotionCode,
} from 'AdminWeb/api/billing/billing';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
    },
});

const BillingManagementPage = () => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const [ userBillingInformation, setUserBillingInformation ] = useState<UserBillingInformation>(null);
    const [ userAccountStatus, setUserAccountStatus ] = useState<SubscriptionLevel>(null);

    const [ emailAddress, setEmailAddress ] = useState<string>(null);
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmailAddress((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        setIsLoading(true);
        getBillingInformationForEmailAddress({
            emailAddress: emailAddress,
        },
        (resp: GetBillingInformationForEmailAddressResponse) => {
            setIsLoading(false);
            setUserBillingInformation(resp.billingInformation);
            setUserAccountStatus(resp.userAccountStatus);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }

    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Warning}>
                An error occurred
            </Heading3>
        );
    } else if (userBillingInformation !== null) {
        body = (
            <UserBillingInformationDisplay
                userAccountStatus={userAccountStatus}
                userBillingInformation={userBillingInformation} />
        );
    }

    const classes = styleClasses();
    return (
        <Page>
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Billing Manager"
                        backArrowDestination="/ops/dashboard" />
                    <PromotionCodeForm />
                    <Divider />
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                            <Grid item xs={12}>
                                <PrimaryTextField
                                    id="email-address"
                                    className={classes.formComponent}
                                    label="Email Address"
                                    variant="outlined"
                                    defaultValue={emailAddress}
                                    onChange={handleEmailAddressChange} />
                            </Grid>
                            <Grid item xs={6}>
                                <PrimaryButton
                                    className={classes.formComponent}
                                    disabled={!emailAddress}
                                    type="submit">
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                </DisplayCard>
            </CenteredComponent>
            {body}
        </Page>
    );
}

type UserBillingInformationDisplayProps = {
    userBillingInformation: UserBillingInformation | undefined;
    userAccountStatus: SubscriptionLevel | undefined;
}

const UserBillingInformationDisplay = (props: UserBillingInformationDisplayProps) => {

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const handleSubmit = () => {
        setIsLoading(true);
        forceSyncForUser({
            userId: props.userBillingInformation.userId,
        },
        (resp: ForceSyncForUserResponse) => {
            setIsLoading(false);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }

    const classes = styleClasses();
    return (
        <CenteredComponent>
            <DisplayCard>
                <Heading3 color={!!props.userAccountStatus ? TypographyColor.Primary : TypographyColor.Gray}>
                    Subscription Status: {props.userAccountStatus || "inactive"}
                </Heading3>
                <Paragraph>
                    Account Type: {props.userBillingInformation.externalIdType}
                </Paragraph>
                {
                    isLoading ? (
                        <LoadingSpinner />
                    ) : (
                        <CenteredComponent>
                            <PrimaryButton className={classes.formComponent} onClick={handleSubmit}>
                                Force Sync
                            </PrimaryButton>
                        </CenteredComponent>
                    )
                }
                {
                    props.userBillingInformation.subscriptions.map((s: PremiumNewsletterSubscription) => (
                        <div>
                            <PremiumNewsletterSubscriptionView key={`subscription-${s.id}`} subscription={s} />
                            <Divider />
                        </div>
                    ))
                }
            </DisplayCard>
        </CenteredComponent>
    );
}

type PremiumNewsletterSubscriptionViewProps = {
    subscription: PremiumNewsletterSubscription;
}

const PremiumNewsletterSubscriptionView = (props: PremiumNewsletterSubscriptionViewProps) => {
    let stateDisplay;
    let titleColor = TypographyColor.Gray;
    const { paymentState } = props.subscription;
    switch (paymentState) {
        case PaymentState.CreatedUnpaid:
            stateDisplay = "Created, Unpaid";
            break;
        case PaymentState.TrialNoPaymentMethod:
            titleColor = TypographyColor.Primary;
            stateDisplay = "Trial, No Payment Method";
            break;
        case PaymentState.TrialPaymentMethodAdded:
            titleColor = TypographyColor.Primary;
            stateDisplay = "Trial, Added Payment Method";
            break;
        case PaymentState.Active:
            titleColor = TypographyColor.Primary;
            stateDisplay = "Active";
            break;
        case PaymentState.Errored:
            titleColor = TypographyColor.Warning;
            stateDisplay = "Errored";
            break;
        case PaymentState.Terminated:
            stateDisplay = "Terminated";
            break;
        default:
            throw new Error(`Unrecognized payment state ${paymentState}`);
    }
    return (
        <div>
            <Heading4 align={Alignment.Left} color={titleColor}>
                {stateDisplay} subscription, ending period {new Date(props.subscription.currentPeriodEnd).toLocaleDateString()}
            </Heading4>
            <Paragraph align={Alignment.Left} color={props.subscription.isAutoRenewEnabled ? TypographyColor.Confirmation : TypographyColor.Gray}>
                {
                    props.subscription.isAutoRenewEnabled ? "Auto-Renew is enabled" : "Auto-Renew is disabled"
                }
            </Paragraph>
        </div>
    );
}

const PromotionCodeForm = () => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ promotionCode, setPromotionCode ] = useState<string>(null);
    const handlePromotionCodeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPromotionCode((event.target as HTMLInputElement).value);
    }

    // TODO: make this dynamic later
    const [ promotionType, setPromotionType ] = useState<PromotionType>(PromotionType.URL);

    const [ amountOffCents, setAmountOffCents ] = useState<number>(null);
    const handleAmountOffCentsChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const amount = parseFloat((event.target as HTMLInputElement).value.replace(/\$/g, ""));
        setAmountOffCents(amount * 100);
        setPercentOffBps(null);
    }

    const [ percentOffBps, setPercentOffBps ] = useState<number>(null);
    const handlePercentOffBpsChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const percent = parseFloat((event.target as HTMLInputElement).value.replace(/%/g, ""));
        setPercentOffBps(percent * 100);
        setAmountOffCents(null);
    }

   const [ maxRedemptions, setMaxRedemptions ] = useState<number>(null);
   const handleMaxRedemptionsChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const amount = parseInt((event.target as HTMLInputElement).value, 10);
        setMaxRedemptions(amount);
    }

    const handleSubmit = () => {
        setIsLoading(true);
        createPromotionCode({
            code: promotionCode,
            discount: {
                amountOffCents: amountOffCents,
                percentOffBps: percentOffBps,
            },
            promotionType: promotionType,
            maxRedemptions: maxRedemptions,
        },
        (resp: CreatePromotionCodeResponse) => {
            setIsLoading(false);
        },
        (err: Error) => {
            setIsLoading(false);
        });
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid container>
                <Grid item xs={3}>
                    <PrimaryTextField
                        id="promotion-code"
                        className={classes.formComponent}
                        label="Promotion Code"
                        variant="outlined"
                        defaultValue={promotionCode}
                        onChange={handlePromotionCodeChange} />
                </Grid>
                <Grid item xs={3}>
                    <PrimaryTextField
                        id="amount-off"
                        className={classes.formComponent}
                        label="Amount Off"
                        variant="outlined"
                        defaultValue={amountOffCents != null ? `$${amountOffCents/100}` : "--"}
                        onChange={handleAmountOffCentsChange} />
                </Grid>
                <Grid item xs={3}>
                    <PrimaryTextField
                        id="percent-off"
                        className={classes.formComponent}
                        label="Percent Off"
                        variant="outlined"
                        defaultValue={percentOffBps != null ? `${percentOffBps/100}%` : "--"}
                        onChange={handlePercentOffBpsChange} />
                </Grid>
                <Grid item xs={3}>
                    <PrimaryTextField
                        id="max-redemptions"
                        className={classes.formComponent}
                        label="Max Redemptions"
                        variant="outlined"
                        defaultValue={maxRedemptions}
                        onChange={handleMaxRedemptionsChange} />
                </Grid>
                <Grid item xs={6}>
                    <PrimaryButton
                        className={classes.formComponent}
                        disabled={!promotionCode && !(amountOffCents && percentOffBps)}
                        type="submit">
                        Submit
                    </PrimaryButton>
                </Grid>
            </Grid>
        </Form>
    );
}

export default BillingManagementPage;
