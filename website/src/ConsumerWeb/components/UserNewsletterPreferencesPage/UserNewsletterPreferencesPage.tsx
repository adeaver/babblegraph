import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading2, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import Paragraph from 'common/typography/Paragraph';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    UserProfileInformation,
} from 'ConsumerWeb/api/useraccounts2/useraccounts';
import {
    getUserNewsletterPreferences,
    GetUserNewsletterPreferencesResponse,

    updateUserNewsletterPreferences,
    UpdateUserNewsletterPreferencesResponse,

    UserNewsletterPreferences,
} from 'ConsumerWeb/api/user/userNewsletterPreferences';

const styleClasses = makeStyles({
    toggleContainer: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
    },
});

type Params = {
    token: string;
}

type UserNewsletterPreferencesPageProps = RouteComponentProps<Params>;

const UserNewsletterPreferencesPage = withUserProfileInformation<UserNewsletterPreferencesPageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [],
    (ownProps: UserNewsletterPreferencesPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.NewsletterPreferences,
    (props: UserNewsletterPreferencesPageProps & UserProfileComponentProps) => {
        return (
            <UserNewsletterPreferencesDisplay
                subscriptionManagementToken={props.match.params.token}
                userProfile={props.userProfile} />
        );
    }
);

type UserNewsletterPreferencesDisplayOwnProps = {
    subscriptionManagementToken: string;
    userProfile: UserProfileInformation;
}

const UserNewsletterPreferencesDisplay = asBaseComponent<GetUserNewsletterPreferencesResponse, UserNewsletterPreferencesDisplayOwnProps>(
    (props: GetUserNewsletterPreferencesResponse & UserNewsletterPreferencesDisplayOwnProps & BaseComponentProps) => {
        if (!!props.error) {
            return (
                <CenteredComponent>
                    <DisplayCard>
                        <DisplayCardHeader
                            title="Manage Preferences"
                            backArrowDestination={`/manage/${props.subscriptionManagementToken}`} />
                        <Heading3 color={TypographyColor.Warning}>
                            There was a problem with your request
                        </Heading3>
                    </DisplayCard>
                </CenteredComponent>
            );
        }

        const [ isLemmaSpotlightActive, setIsLemmaSpotlightActive ] = useState<boolean>(props.preferences.isLemmaReinforcementSpotlightActive);
        const [ arePodcastsEnabled, setArePodcastsEnabled ] = useState<boolean>(props.preferences.arePodcastsEnabled);
        const [ includeExplicitPodcasts, setIncludeExplicitPodcasts ] = useState<boolean>(props.preferences.includeExplicitPodcasts);
        const [ minimumPodcastDurationSeconds, setMinimumPodcastDurationSeconds ] = useState<number | undefined>(props.preferences.minimumPodcastDurationSeconds);
        const [ maximumPodcastDurationSeconds, setMaximumPodcastDurationSeconds ] = useState<number | undefined>(props.preferences.maximumPodcastDurationSeconds);

        const classes = styleClasses();
        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Manage Preferences"
                        backArrowDestination={`/manage/${props.subscriptionManagementToken}`} />
                    <Grid container>
                        <Grid item xs={10} xl={11}>
                            <Heading4 align={Alignment.Left} color={TypographyColor.Primary}>
                                Include word tracking spotlights in your newsletter?
                            </Heading4>
                            <Paragraph align={Alignment.Left}>
                                Word tracking spotlights include a highlighted article that is guaranteed to have a word in your tracking list. It spaces out these spotlights so you can practice new words on your list. If this is disabled, you won't see spotlights in your newsletter.
                            </Paragraph>
                        </Grid>
                        <Grid item
                            className={classes.toggleContainer}
                            xs={2}
                            xl={1}>
                            <PrimarySwitch
                                checked={isLemmaSpotlightActive} onClick={() => {setIsLemmaSpotlightActive(!isLemmaSpotlightActive)}} />
                        </Grid>
                        {
                            /* TODO flip this condition */
                            !props.userProfile.subscriptionLevel && (
                                <Grid item xs={12}>
                                    <Heading2
                                        align={Alignment.Left}
                                        color={TypographyColor.Primary}>
                                        Podcast Settings
                                    </Heading2>
                                </Grid>
                            )
                        }
                    </Grid>
                </DisplayCard>
            </CenteredComponent>
        )
    },
    (
        ownProps: UserNewsletterPreferencesDisplayOwnProps,
        onSuccess: (resp: GetUserNewsletterPreferencesResponse) => void,
        onError: (err: Error) => void,
    ) => {
        getUserNewsletterPreferences({
            languageCode: WordsmithLanguageCode.Spanish,
            subscriptionManagementToken: ownProps.subscriptionManagementToken,
        }, onSuccess, onError);
    },
    false,
)

export default UserNewsletterPreferencesPage;
