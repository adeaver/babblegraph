import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Divider from '@material-ui/core/Divider';
import ArrowForwardIcon from '@material-ui/icons/ArrowForward';

import Link, { LinkTarget } from 'common/components/Link/Link';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { getReinforcementToken, GetReinforcementTokenResponse } from 'api/token/reinforcementToken';

import {
    getUserProfile,
    GetUserProfileResponse
} from 'api/useraccounts/useraccounts';

type Params = {
    token: string
}

const styleClasses = makeStyles({
    loginInfoContainer: {
        marginTop: '45px',
        display: 'inline-block',
        width: '100%',
    },
    actionCard: {
        padding: '10px',
        height: '100%',
        "&:hover": {
            boxShadow: "0 0 4px 2px gray",
        },
        cursor: 'pointer',
    },
    headerArrow: {
        alignSelf: 'center',
    },
    gridComponent: {
        marginTop: '15px',
    },
});

type ActionCardProps = {
    title: string;
    redirectURL: string;
    children: string;
}

const ActionCard = (props: ActionCardProps) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <Grid className={classes.gridComponent} item xs={12} md={6}>
            <Card onClick={() => { history.push(props.redirectURL) }} className={classes.actionCard} variant='outlined'>
                <Grid container>
                    <Grid item xs={11}>
                        <Paragraph size={Size.Large} color={TypographyColor.Primary} align={Alignment.Left}>{props.title}</Paragraph>
                    </Grid>
                    <Grid className={classes.headerArrow} item xs={1}>
                        <ArrowForwardIcon color="action" />
                    </Grid>
                </Grid>
                <Divider />
                <Paragraph size={Size.Medium} color={TypographyColor.Gray} align={Alignment.Left}>{props.children}</Paragraph>
            </Card>
        </Grid>
    );
}

type SubscriptionManagementDashboardPageProps = RouteComponentProps<Params>

const SubscriptionManagementDashboardPage = (props: SubscriptionManagementDashboardPageProps) => {
    const classes = styleClasses();
    const { token } = props.match.params;

    const [ isReinforcementTokenLoading, setIsReinforcementTokenLoading ] = useState<boolean>(true);
    const [ isUserProfileLoading, setIsUserProfileLoading ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);
    const [ reinforcementToken, setReinforcementToken ] = useState<string | null>(null);
    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ subscriptionLevel, setSubscriptionLevel ] = useState<string | null>(null);

    useEffect(() => {
        getReinforcementToken({
            token: token,
        },
        (resp: GetReinforcementTokenResponse) => {
            setIsReinforcementTokenLoading(false);
            setReinforcementToken(resp.token);
        },
        (e: Error) => {
            setIsReinforcementTokenLoading(false);
            setError(e);
        });
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setIsUserProfileLoading(false);
            setEmailAddress(resp.emailAddress);
            resp.subscriptionLevel && setSubscriptionLevel(resp.subscriptionLevel)
        },
        (e: Error) => {
            setIsUserProfileLoading(false);
            setError(e);
        });
    }, []);

    const isLoading = isReinforcementTokenLoading || isUserProfileLoading;
    return (
        <Page>
            {
                isLoading ? (
                    <LoadingSpinner />
                ) : (
                    <div>
                        <Grid container spacing={2}>
                            <ActionCard redirectURL={`/manage/${token}/interests`} title='Manage Your Interests'>
                                Select some topics you’re interested in reading more about or deselect some topics you’d like to read about less. This is a great way to make sure that the content you get is fun and engaging.
                            </ActionCard>
                            <ActionCard redirectURL={`/manage/${reinforcementToken}/vocabulary`} title='Track words to reinforce'>
                                Learn a new word recently and want to make sure it sticks? You can track it, which will send you articles containing these words. Seeing a word frequently is a great way to make sure you remember it.
                            </ActionCard>
                            {
                                subscriptionLevel && (
                                    <ActionCard redirectURL={`/manage/${token}/schedule`} title='Newsletter schedule and preferences'>
                                        Select which days you receive newsletter emails from Babblegraph. You can also configure how many articles you receive in each email and what topics are in each email.
                                    </ActionCard>
                                )
                            }
                            <ActionCard redirectURL={`/manage/${token}/level`} title='Set your difficulty level'>
                                If your daily email is too hard or too easy, you can change the difficulty level here.
                            </ActionCard>
                            <ActionCard redirectURL={`/manage/${token}/unsubscribe`} title='Unsubscribe'>
                                If you’re no longer interested in receiving daily emails, you can unsubscribe here. By unsubscribing, we won’t send you any more emails about anything.
                            </ActionCard>
                        </Grid>
                        {
                            !!emailAddress && (
                                <div className={classes.loginInfoContainer}>
                                    <Heading4 color={TypographyColor.Primary}>Logged in as {emailAddress}</Heading4>
                                    <Link href="/logout" target={LinkTarget.Self}>Click here to logout</Link>
                                </div>
                            )
                        }
                    </div>
                )
            }
        </Page>
    );
}

export default SubscriptionManagementDashboardPage;
