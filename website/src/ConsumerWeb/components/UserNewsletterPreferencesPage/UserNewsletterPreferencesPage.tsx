import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import RadioGroup from '@material-ui/core/RadioGroup';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading2, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import Paragraph from 'common/typography/Paragraph';
import { PrimaryRadio } from 'common/components/Radio/Radio';

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

enum PodcastDurationPreference {
    LessThanFifteen = 'Less than 15 minutes',
    FifteenToThirty = '15 minutes to 30 minutes',
    ThirtyToOneHour = '30 minutes to 1 hour',
    MoreThanOneHour = 'More than an hour',
}

const getPodcastDurationByMinimumAndMaximium = (minimumDurationSeconds: number | undefined, maximumDurationSeconds: number | undefined) => {
    if (!!minimumDurationSeconds) {
        if (!!maximumDurationSeconds) {
            return null;
        } else if (maximumDurationSeconds / 60 === 15) {
            return PodcastDurationPreference.LessThanFifteen;
        } else {
            return null;
        }
    } else if (minimumDurationSeconds / 60 === 15) {
        if (!!maximumDurationSeconds) {
            return null
        } else if (maximumDurationSeconds / 60 === 30) {
            return PodcastDurationPreference.FifteenToThirty;
        }
        return null
    } else if (minimumDurationSeconds / 60 === 30) {
        if (!!maximumDurationSeconds) {
            return null
        } else if (maximumDurationSeconds / 60 === 60) {
            return PodcastDurationPreference.ThirtyToOneHour;
        }
        return null
    } else if (minimumDurationSeconds / 60 === 60) {
        return PodcastDurationPreference.MoreThanOneHour;
    }
    return null
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

        const [ podcastDuration, setPodcastDuration ] = useState<PodcastDurationPreference>(
            getPodcastDurationByMinimumAndMaximium(
                props.preferences.minimumPodcastDurationSeconds, props.preferences.maximumPodcastDurationSeconds
            )
        );
        const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setPodcastDuration((event.target as HTMLInputElement).value as PodcastDurationPreference);
        };

        const [ isLemmaSpotlightActive, setIsLemmaSpotlightActive ] = useState<boolean>(props.preferences.isLemmaReinforcementSpotlightActive);
        const [ arePodcastsEnabled, setArePodcastsEnabled ] = useState<boolean>(props.preferences.arePodcastsEnabled);
        const [ includeExplicitPodcasts, setIncludeExplicitPodcasts ] = useState<boolean>(props.preferences.includeExplicitPodcasts);
        const [ minimumPodcastDurationSeconds, setMinimumPodcastDurationSeconds ] = useState<number | undefined>();
        const [ maximumPodcastDurationSeconds, setMaximumPodcastDurationSeconds ] = useState<number | undefined>();

        const classes = styleClasses();
        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Manage Preferences"
                        backArrowDestination={`/manage/${props.subscriptionManagementToken}`} />
                    <Grid container>
                        <Grid item xs={10} xl={11}>
                            <Heading4 align={Alignment.Left}>
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
                                    <Grid container>
                                        <Grid item xs={10} xl={11}>
                                            <Heading4 align={Alignment.Left}>
                                                Would you like to include podcasts in your newsletter?
                                            </Heading4>
                                            <Paragraph align={Alignment.Left}>
                                                With your premium subscription, Babblegraph will send you podcasts. But if you don’t want podcasts in your newsletter, you can disable it here.
                                            </Paragraph>
                                        </Grid>
                                        <Grid item
                                            className={classes.toggleContainer}
                                            xs={2}
                                            xl={1}>
                                            <PrimarySwitch
                                                checked={arePodcastsEnabled} onClick={() => {setArePodcastsEnabled(!arePodcastsEnabled)}} />
                                        </Grid>
                                        <Grid item xs={10} xl={11}>
                                            <Heading4 align={Alignment.Left}>
                                                Include potentially explicit podcasts in your newsletter?
                                            </Heading4>
                                            <Paragraph align={Alignment.Left}>
                                                You can disable sending podcasts that deal with explicit subjects or use explicit language in your newsletter.
                                            </Paragraph>
                                        </Grid>
                                        <Grid item
                                            className={classes.toggleContainer}
                                            xs={2}
                                            xl={1}>
                                            <PrimarySwitch
                                                checked={includeExplicitPodcasts && arePodcastsEnabled} onClick={() => {setIncludeExplicitPodcasts(!includeExplicitPodcasts)}} disabled={!arePodcastsEnabled} />
                                        </Grid>
                                    </Grid>
                                    <Grid item xs={12}>
                                        <Heading4 align={Alignment.Left}>
                                            What length of podcasts would you like to receive?
                                        </Heading4>
                                        <Paragraph align={Alignment.Left}>
                                            Maybe you don’t have all day to listen to podcasts, or maybe you have a lot of time to fill.
                                        </Paragraph>
                                    </Grid>
                                    <Grid item xs={12}>
                                        <FormControl component="fieldset">
                                            <RadioGroup aria-label="add-list-type" name="add-list-type1" value={podcastDuration} onChange={handleRadioFormChange}>
                                                <Grid container>
                                                    {
                                                        Object.keys(PodcastDurationPreference).map((p: PodcastDurationPreference) => ((
                                                            <Grid item xs={6}>
                                                                <FormControlLabel value={p} control={<PrimaryRadio />} label={PodcastDurationPreference[p]} />
                                                            </Grid>
                                                        )))
                                                    }
                                                </Grid>
                                            </RadioGroup>
                                        </FormControl>
                                    </Grid>
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