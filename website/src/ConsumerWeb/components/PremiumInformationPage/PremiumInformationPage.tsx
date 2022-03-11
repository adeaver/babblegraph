import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import Grid from '@material-ui/core/Grid';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import DateRangeIcon from '@material-ui/icons/DateRange';
import FindInPageIcon from '@material-ui/icons/FindInPage';
import BallotIcon from '@material-ui/icons/Ballot';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';

import Color from 'common/styles/colors';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading1, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { setLocation } from 'util/window/Location';
import { PrimaryButton } from 'common/components/Button/Button';

import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';

import { DisplayLanguage } from 'common/model/language/language';
import { TextBlock, getTextBlocksForLanguage } from './translations';

const styleClasses = makeStyles({
    featureIcon: {
        color: Color.Primary,
    },
    callToActionButton: {
        margin: '15px 0',
        width: '100%',
    },
});

type Params = {
    token: string;
}

type PremiumInformationPageOwnProps = RouteComponentProps<Params>;

const PremiumInformationPage = withUserProfileInformation<PremiumInformationPageOwnProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [RouteEncryptionKey.CreateUser, RouteEncryptionKey.PremiumSubscriptionCheckout],
    (ownProps: PremiumInformationPageOwnProps) => {
        return ownProps.match.params.token;
    },
    // TODO: Create Login Redirect for Premium Information
    undefined,
    (props: PremiumInformationPageOwnProps & UserProfileComponentProps) => {
        const { token } = props.match.params;
        const [ createUserToken, premiumSubscriptionCheckoutToken ] = props.userProfile.nextTokens;

        const classes = styleClasses();

        // TODO(i18n): convert this to Spanish
        const translations = getTextBlocksForLanguage(DisplayLanguage.English);

        let callToAction;
        // Invariants here
        // If the user has an account, they are logged in because the higher level component will redirect them
        // If the user has a subscription level, they have an account and are logged in
        // If the user doesn't have an account, they do not have a subscription
        if (!!props.userProfile.subscriptionLevel) {
            // The user already has a subscription
            callToAction = (
                <Heading4 color={TypographyColor.Primary}>
                    {translations[TextBlock.AccountDisclaimer]}
                </Heading4>
            );
        } else {
            callToAction = (
                <PrimaryButton
                    onClick={() => setLocation(
                        props.userProfile.hasAccount ? (
                            `/checkout/${premiumSubscriptionCheckoutToken}`
                        ) : (
                            `/signup/${createUserToken}`
                        ))}
                    className={classes.callToActionButton}
                    size="large">
                    {translations[TextBlock.CallToAction]}
                </PrimaryButton>
            );
        }

        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <DisplayCardHeader
                            title="Babblegraph Premium"
                            backArrowDestination={`/manage/${token}`} />
                        <Heading3>
                            {translations[TextBlock.Subheading]}
                        </Heading3>
                        <Divider />
                        <List>
                            <PremiumFeatureListItem
                                title={translations[TextBlock.PodcastsTitle]}
                                description={translations[TextBlock.PodcastsDescription]}
                                icon={<GraphicEqIcon className={classes.featureIcon} />} />
                            <PremiumFeatureListItem
                                title={translations[TextBlock.SchedulingTitle]}
                                description={translations[TextBlock.SchedulingDescription]}
                                icon={<DateRangeIcon className={classes.featureIcon} />} />
                            <PremiumFeatureListItem
                                title={translations[TextBlock.CustomizationTitle]}
                                description={translations[TextBlock.CustomizationDescription]}
                                icon={<BallotIcon className={classes.featureIcon} />} />
                        </List>
                        <Heading3 color={TypographyColor.Primary}>
                            {translations[TextBlock.Price]}
                        </Heading3>
                        {
                            /* TODO(i18n): fix this */
                            !!props.userProfile.trialEligibilityDays && (
                                <Heading3>
                                    The best part is that you can try it for free for {props.userProfile.trialEligibilityDays} days!
                                </Heading3>
                            )
                        }
                        <Grid container>
                            <Grid item xs={2}>
                                &nbsp;
                            </Grid>
                            <Grid item xs={8}>
                                { callToAction }
                            </Grid>
                        </Grid>
                    </DisplayCard>
                </Grid>
            </Grid>
        )
    }
);

type PremiumFeatureListItemProps = {
    title: string;
    description: string;
    icon: JSX.Element;
}

const PremiumFeatureListItem = (props: PremiumFeatureListItemProps) => {
    const classes = styleClasses();
    return (
        <ListItem>
            <ListItemIcon>
                {props.icon}
            </ListItemIcon>
            <ListItemText>
                <Heading4
                    align={Alignment.Left}
                    color={TypographyColor.Primary}>
                    {props.title}
                </Heading4>
                <Paragraph align={Alignment.Left}>
                    {props.description}
                </Paragraph>
            </ListItemText>
        </ListItem>
    );
}

export default PremiumInformationPage;
