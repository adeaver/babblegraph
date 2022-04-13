import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import { Alignment, TypographyColor } from 'common/typography/common';
import ActionCard from 'common/components/ActionCard/ActionCard';
import Paragraph from 'common/typography/Paragraph';
import { Heading2 } from 'common/typography/Heading';
import { setLocation } from 'util/window/Location';
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';

const styleClasses = makeStyles({
    navigationCard: {
        padding: '15px',
    },
    navigationCardActionCard: {
        height: '100%',
        boxSizing: 'border-box',
    },
});

type Params = {
    token: string;
}

type SubscriptionManagementHomePageProps = RouteComponentProps<Params>;

const SubscriptionManagementHomePage = withUserProfileInformation<SubscriptionManagementHomePageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [RouteEncryptionKey.WordReinforcement],
    (ownProps: SubscriptionManagementHomePageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.SubscriptionManagement,
    (props: SubscriptionManagementHomePageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;
        const [ reinforcementToken ] = props.userProfile.nextTokens;

        return (
            <div>
                <Grid container spacing={2}>
                    <NavigationCard
                        location={`/manage/${token}/preferences`}
                        title="Schedule and newsletter settings"
                        description="Select which days you receive your newsletter and at what time. Also, toggle other settings, such as how many articles appear per email." />
                    <NavigationCard
                        location={`/manage/${token}/interests`}
                        title="Manage Your Interests"
                        description="Select some topics you’re interested in reading more about or deselect some topics you’d like to read about less. This is a great way to make sure that the content you get is fun and engaging." />
                    <NavigationCard
                        location={`/manage/${reinforcementToken}/vocabulary`}
                        title="Add words to your vocabulary list"
                        description="Learn a new word recently and want to make sure it sticks? You can add it to your vocabulary list, which will send you articles containing these words. Seeing a word frequently is a great way to make sure you remember it." />
                    {
                        props.userProfile.hasAccount && (
                            <NavigationCard
                                location={`/manage/${token}/payment-settings`}
                                title="Subscription and Payment Settings"
                                description="Need to update your preferred payment method or cancel your subscription? Click here!" />
                        )
                    }
                    <NavigationCard
                        location={`/manage/${token}/unsubscribe`}
                        title="Unsubscribe"
                        description="If you’re no longer interested in receiving newsletters, you can unsubscribe here. By unsubscribing, we won’t send you any more emails about anything." />
                </Grid>
                {
                    props.userProfile.isLoggedIn && (
                        <Link href="/logout" target={LinkTarget.Self}>Click here to logout</Link>
                    )
                }
            </div>
        );
    }
);

type NavigationCardProps = {
    location: string;
    title: string;
    description: string;
}

const NavigationCard = (props: NavigationCardProps) => {
    const classes = styleClasses();
    return (
        <Grid className={classes.navigationCard} item xs={12} md={6} xl={4}>
            <ActionCard className={classes.navigationCardActionCard} onClick={() => setLocation(props.location)}>
                <Heading2
                    align={Alignment.Left}
                    color={TypographyColor.Primary}>
                    { props.title }
                </Heading2>
                <Divider />
                <Paragraph align={Alignment.Left}>
                    { props.description }
                </Paragraph>
            </ActionCard>
        </Grid>
    );
}

export default SubscriptionManagementHomePage;
