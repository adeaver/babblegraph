import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Link, { LinkTarget } from 'common/components/Link/Link';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Page from 'common/components/Page/Page';
import Paragraph from 'common/typography/Paragraph';
import { Heading1 } from 'common/typography/Heading';
import { PrimaryButton } from 'common/components/Button/Button';
import { TypographyColor } from 'common/typography/common';

import PremiumInformationPanel from 'components/PremiumInformation/PremiumInformationPanel';

import {
    getCreateUserToken,
    GetCreateUserTokenResponse
} from 'api/token/createUserToken';
import {
    getPremiumCheckoutToken,
    GetPremiumCheckoutTokenResponse
} from 'api/token/premiumCheckoutToken';
import {
    getUserProfile,
    GetUserProfileResponse
} from 'api/useraccounts/useraccounts';

const styleClasses = makeStyles({
    callToActionButton: {
        margin: '15px 0',
        width: '100%',
    },
});

type Params = {
    token: string;
}

type PremiumInformationPageProps = RouteComponentProps<Params>

const PremiumInformationPage = (props: PremiumInformationPageProps) => {
    const { token } = props.match.params;

    const [ isLoadingUserProfile, setIsLoadingUserProfile ] = useState<boolean>(true);
    const [ hasAccount, setHasAccount ] = useState<boolean>(false);
    const [ hasSubscription, setHasSubscription ] = useState<boolean>(false);

    const [ isLoadingCreateUserToken, setIsLoadingCreateUserToken ] = useState<boolean>(true);
    const [ createUserToken, setCreateUserToken ] = useState<string | null>(null);

    const [ isLoadingPremiumCheckoutToken, setIsLoadingPremiumCheckoutToken ] = useState<boolean>(true);
    const [ premiumCheckoutToken, setPremiumCheckoutToken ] = useState<string | null>(null);

    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setIsLoadingUserProfile(false);
            const hasAccount = !!resp.emailAddress;
            const hasSubscription = !!resp.subscriptionLevel;
            setHasAccount(hasAccount);
            setHasSubscription(hasSubscription);
            if (hasAccount && hasSubscription) {
                // If a user is already subscribed, then don't load
                // the create user token or the premium checkout token
                setIsLoadingCreateUserToken(false);
                setIsLoadingPremiumCheckoutToken(false);
            } else if (hasAccount && !hasSubscription) {
                // If a user has an account, but is not subscribed, then they
                // were previously subscribed. They need a separate call to action.
                // We should not load the create user token, but we need to load the premium checkout token
                setIsLoadingCreateUserToken(false);
                // There are two options here:
                // 1) User has previously had a subscription. In which case, we need to reactivate.
                // 2) User has made an account but not paid.
                getPremiumCheckoutToken({
                    token: token,
                },
                (resp: GetPremiumCheckoutTokenResponse) => {
                    setIsLoadingPremiumCheckoutToken(false);
                    setPremiumCheckoutToken(resp.token);
                },
                (err: Error) => {
                    setIsLoadingPremiumCheckoutToken(false);
                    setError(err);
                });
            } else if (!hasAccount && hasSubscription) {
                // This should be impossible
                setIsLoadingCreateUserToken(false);
                setIsLoadingPremiumCheckoutToken(false);
                setError(new Error("invalid state"));
            } else {
                // Load the create user token
                setIsLoadingPremiumCheckoutToken(false);
                getCreateUserToken({
                    token: token,
                },
                (resp: GetCreateUserTokenResponse) => {
                    setIsLoadingCreateUserToken(false);
                    setCreateUserToken(resp.token);
                },
                (err: Error) => {
                    setIsLoadingCreateUserToken(false);
                    setError(err);
                });
            }
        },
        (err: Error) => {
            setIsLoadingUserProfile(false);
            setError(err);
        });
    }, []);

    const classes = styleClasses();
    const history = useHistory();
    const isLoading = isLoadingUserProfile || isLoadingCreateUserToken || isLoadingPremiumCheckoutToken;
    let body;
    if (isLoading) {
        body = <LoadingSpinner />
    } else if (!!error) {
        body = (
            <Heading1>
                Something went wrong getting your information. Try again later or reach out to hello@babblegraph.com!
            </Heading1>
        );
    } else {
        let callToAction = null;
        if (hasAccount && hasSubscription) {
            callToAction = (
                <Paragraph color={TypographyColor.Confirmation}>
                    You already have a Babblegraph Premium Subscription.
                </Paragraph>
            );
        } else if (hasAccount && !hasSubscription && premiumCheckoutToken) {
            // TODO: This needs to go somewhere
            callToAction = (
                <PrimaryButton
                    onClick={() => history.push(`/checkout/${premiumCheckoutToken}`)}
                    className={classes.callToActionButton}
                    size="large">
                    Renew your Babblegraph Premium Subscription
                </PrimaryButton>
            );
        } else if (!hasAccount && !hasSubscription && createUserToken) {
            callToAction = (
                <PrimaryButton
                    onClick={() => history.push(`/signup/${createUserToken}`)}
                    className={classes.callToActionButton}
                    size="large">
                    Try Babblegraph Premium
                </PrimaryButton>
            );
        }
        body = (
            <PremiumInformationWithCallToAction
                subscriptionManagementToken={token}
                callToAction={callToAction} />
        );
    }
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    {body}
                </Grid>
            </Grid>
        </Page>
    );
}

type PremiumInformationWithCallToActionProps = {
    subscriptionManagementToken: string;
    callToAction: JSX.Element | string;
};

const PremiumInformationWithCallToAction = (props: PremiumInformationWithCallToActionProps) => {
    const classes = styleClasses();
    return (
        <DisplayCard>
            <PremiumInformationPanel />
            <Divider />
            <Grid container>
                <Grid item xs={2}>
                    &nbsp;
                </Grid>
                <Grid item xs={8}>
                    {props.callToAction}
                </Grid>
            </Grid>
            <Link href={`/manage/${props.subscriptionManagementToken}`} target={LinkTarget.Self}>
                Go back to subscription management
            </Link>
        </DisplayCard>
    );
}

export default PremiumInformationPage;
