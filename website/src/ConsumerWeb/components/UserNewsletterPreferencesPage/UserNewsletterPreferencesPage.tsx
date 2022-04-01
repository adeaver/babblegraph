import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import RadioGroup from '@material-ui/core/RadioGroup';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading2, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import Paragraph from 'common/typography/Paragraph';
import { PrimaryRadio } from 'common/components/Radio/Radio';
import Form from 'common/components/Form/Form';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';
import { DisplayLanguage } from 'common/model/language/language';
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
import {
    toTitleCase
} from 'util/string/StringConvert';

import TimeSelector from './TimeSelector';

const minimumNumberOfArticles = 4;
const maximumNumberOfArticles = 12;
const daysOfTheWeekByLanguageCode: { [key: string]: Array<string> } = {
    [DisplayLanguage.Spanish]: [ "Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"],
    [DisplayLanguage.English]: [ "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday" ],
}

const styleClasses = makeStyles({
    toggleContainer: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
    submitButton: {
        display: 'block',
        margin: 'auto',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    emailField: {
        width: '100%',
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
    LessThanFifteen = 'Less than 15 minutes (least number of podcasts)',
    FifteenToThirty = '15 minutes to 30 minutes',
    ThirtyToOneHour = '30 minutes to 1 hour (recommended)',
    MoreThanOneHour = 'More than an hour',
    NoPreference = 'Any podcast length (most number of podcasts)',
}

const getPodcastDurationByMinimumAndMaximium = (minimumDurationSeconds: number | undefined, maximumDurationSeconds: number | undefined) => {
    if (!minimumDurationSeconds) {
        if (!!maximumDurationSeconds && maximumDurationSeconds / 60 === 15) {
            return PodcastDurationPreference.LessThanFifteen;
        }
    } else if (minimumDurationSeconds / 60 === 15) {
        if (!!maximumDurationSeconds && maximumDurationSeconds / 60 === 30) {
            return PodcastDurationPreference.FifteenToThirty;
        }
    } else if (minimumDurationSeconds / 60 === 30) {
        if (!!maximumDurationSeconds && maximumDurationSeconds / 60 === 60) {
            return PodcastDurationPreference.ThirtyToOneHour;
        }
    } else if (minimumDurationSeconds / 60 === 60) {
        return PodcastDurationPreference.MoreThanOneHour;
    }
    return PodcastDurationPreference.NoPreference;
}

const getPodcastDurationBoundsSeconds = (podcastDuration: PodcastDurationPreference | null) => {
    if (!podcastDuration) {
        return [undefined, undefined];
    }
    switch (podcastDuration) {
        case PodcastDurationPreference.LessThanFifteen:
            return [undefined, 15 * 60];
        case PodcastDurationPreference.FifteenToThirty:
            return [15 * 60, 30 * 60];
        case PodcastDurationPreference.ThirtyToOneHour:
            return [30 * 60, 60 * 60];
        case PodcastDurationPreference.MoreThanOneHour:
            return [60 * 60, undefined];
    }
    return [undefined, undefined];
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

        const [ initialIANATimezone, setInitialIANATimezone ] = useState<string>(props.preferences.schedule.ianaTimezone);
        const [ ianaTimezone, setIANATimezone ] = useState<string>(props.preferences.schedule.ianaTimezone);
        const [ hourIndex, setHourIndex ] = useState<number>(props.preferences.schedule.hourIndex);
        const [ quarterHourIndex, setQuarterHourIndex ] = useState<number>(props.preferences.schedule.quarterHourIndex * 15);
        const [ isActiveForDays, setIsActiveForDays ] = useState<Array<boolean>>(props.preferences.schedule.isActiveForDays);
        const [ numberOfArticlesPerEmail, setNumberOfArticlesPerEmail ] = useState<number>(props.preferences.numberOfArticlesPerEmail);
        const handleUpdateNumberOfArticles = (event: React.ChangeEvent<HTMLInputElement>) => {
            const numberOfArticles = parseInt((event.target as HTMLInputElement).value, 10);
            setNumberOfArticlesPerEmail(numberOfArticles);
        }

        const [ isLoading, setIsLoading ] = useState<boolean>(false);
        const [ error, setError ] = useState<ClientError>(null);
        const [ success, setSuccess ] = useState<boolean>(false);

        const [ emailAddress, setEmailAddress ] = useState<string>(undefined);
        const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setEmailAddress((event.target as HTMLInputElement).value);
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

        const handleSubmit = () => {
            const [ minimumPodcastDurationSeconds, maximumPodcastDurationSeconds ] = getPodcastDurationBoundsSeconds(podcastDuration);
            setIsLoading(true);
            updateUserNewsletterPreferences({
                emailAddress: emailAddress,
                subscriptionManagementToken: props.subscriptionManagementToken,
                preferences: {
                    languageCode: WordsmithLanguageCode.Spanish,
                    isLemmaReinforcementSpotlightActive: isLemmaSpotlightActive,
                    arePodcastsEnabled: arePodcastsEnabled,
                    includeExplicitPodcasts: includeExplicitPodcasts,
                    minimumPodcastDurationSeconds: minimumPodcastDurationSeconds,
                    maximumPodcastDurationSeconds: maximumPodcastDurationSeconds,
                    numberOfArticlesPerEmail: numberOfArticlesPerEmail,
                    schedule: {
                        ianaTimezone: ianaTimezone,
                        hourIndex: hourIndex,
                        quarterHourIndex: quarterHourIndex / 15,
                        isActiveForDays: isActiveForDays,
                    },
                },
            },
            (resp: UpdateUserNewsletterPreferencesResponse) => {
                setIsLoading(false);
                if (!!resp.error) {
                    setError(resp.error);
                    return;
                }
                setSuccess(true);
            },
            (err: Error) => {
                setIsLoading(false);
                props.setError(err);
            });
        }

        const classes = styleClasses();
        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Manage Preferences"
                        backArrowDestination={`/manage/${props.subscriptionManagementToken}`} />
                    <Grid container>
                        <Grid item xs={12} md={8}>
                            <Heading4 align={Alignment.Left} color={TypographyColor.Primary}>
                                How many articles would you like to receive per email?
                            </Heading4>
                            <Paragraph align={Alignment.Left}>
                                You can set the number of articles that will appear in each email. I recommend 8 or 12 in case you get articles that you’re less interested in, but it can be overwhelming to receive so many emails. This number must be between {minimumNumberOfArticles} and {maximumNumberOfArticles}.
                            </Paragraph>
                        </Grid>
                        <Grid className={classes.toggleContainer} item xs={12} md={4}>
                            <PrimaryTextField
                                id="number-of-articles"
                                value={numberOfArticlesPerEmail}
                                type="number"
                                label="Number of articles per email"
                                variant="outlined"
                                error={numberOfArticlesPerEmail < minimumNumberOfArticles || numberOfArticlesPerEmail > maximumNumberOfArticles}
                                helperText={`Must select between ${minimumNumberOfArticles} and ${maximumNumberOfArticles}`}
                                onChange={handleUpdateNumberOfArticles} />
                        </Grid>
                        <Grid item xs={12}>
                            <TimeSelector
                                initialIANATimezone={initialIANATimezone}
                                ianaTimezone={ianaTimezone}
                                hourIndex={hourIndex}
                                quarterHourIndex={quarterHourIndex}
                                handleUpdateIANATimezone={setIANATimezone}
                                handleUpdateHourIndex={setHourIndex}
                                handleUpdateQuarterHourIndex={setQuarterHourIndex} />
                        </Grid>
                        <Grid item xs={12}>
                            <Heading4 align={Alignment.Left} color={TypographyColor.Primary}>
                                Which days would you like to receive your newsletter?
                            </Heading4>
                        </Grid>
                        {
                            (daysOfTheWeekByLanguageCode[DisplayLanguage.English] || []).map((day: string, idx: number) => (
                                <Grid item xs={12}>
                                    <Grid container>
                                        <Grid item xs={10} xl={11}>
                                            <Paragraph align={Alignment.Left}>
                                                {toTitleCase(day)}
                                            </Paragraph>
                                        </Grid>
                                        <Grid item
                                            className={classes.toggleContainer}
                                            xs={2}
                                            xl={1}>
                                            <PrimarySwitch
                                                checked={isActiveForDays[idx]}
                                                onClick={() => {
                                                    setIsActiveForDays(
                                                        isActiveForDays.map((val: boolean, i: number) => i === idx ? !val : val)
                                                    )
                                                }}
                                                disabled={isLoading} />
                                        </Grid>
                                    </Grid>
                                </Grid>
                            ))
                        }
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
                                checked={isLemmaSpotlightActive}
                                onClick={() => {setIsLemmaSpotlightActive(!isLemmaSpotlightActive)}}
                                disabled={isLoading} />
                        </Grid>
                        {
                            !!props.userProfile.subscriptionLevel && (
                                <Grid item xs={12}>
                                    <Heading2
                                        align={Alignment.Left}>
                                        Podcast Settings
                                    </Heading2>
                                    <Grid container>
                                        <Grid item xs={10} xl={11}>
                                            <Heading4 align={Alignment.Left} color={TypographyColor.Primary}>
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
                                                checked={arePodcastsEnabled}
                                                onClick={() => {setArePodcastsEnabled(!arePodcastsEnabled)}}
                                                disabled={isLoading} />
                                        </Grid>
                                        <Grid item xs={10} xl={11}>
                                            <Heading4 align={Alignment.Left} color={TypographyColor.Primary}>
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
                                                checked={includeExplicitPodcasts && arePodcastsEnabled}
                                                onClick={() => {setIncludeExplicitPodcasts(!includeExplicitPodcasts)}}
                                                disabled={!arePodcastsEnabled || isLoading} />

                                        </Grid>
                                    </Grid>
                                    <Grid item xs={12}>
                                        <Heading4 align={Alignment.Left} color={TypographyColor.Primary}>
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
                                                            <Grid key={`podcast-option-${p}`} item xs={12}>
                                                                <FormControlLabel value={PodcastDurationPreference[p]} control={<PrimaryRadio disabled={isLoading} />} label={PodcastDurationPreference[p]} />
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
                        <Divider />
                        <Form className={classes.confirmationForm} handleSubmit={handleSubmit}>
                            <Grid container>
                                {
                                    props.userProfile.hasAccount ? (
                                        <Grid item xs={4} md={5}>
                                            &nbsp;
                                        </Grid>
                                    ) : (
                                        <Grid item xs={8} md={10}>
                                            <PrimaryTextField
                                                id="email"
                                                disabled={isLoading}
                                                className={classes.emailField}
                                                value={emailAddress}
                                                label="Email Address"
                                                variant="outlined"
                                                onChange={handleEmailAddressChange} />
                                        </Grid>
                                    )
                                }
                                <Grid item xs={4} md={2} className={classes.submitButtonContainer}>
                                    <PrimaryButton
                                        className={classes.submitButton}
                                        disabled={(!emailAddress && !props.userProfile.hasAccount) || isLoading}
                                        type="submit">
                                        Submit
                                    </PrimaryButton>
                                </Grid>
                            </Grid>
                        </Form>
                    </Grid>
                    {
                        isLoading && <LoadingSpinner />
                    }
                </DisplayCard>
                <Snackbar open={!!error} autoHideDuration={6000} onClose={() => setError(null)}>
                    <Alert severity="error">Something went wrong processing your request.</Alert>
                </Snackbar>
                <Snackbar open={success} autoHideDuration={6000} onClose={() => setSuccess(false)}>
                    <Alert severity="success">Successfully updated your preferences. Changes may take up to 24 hours to take effect!</Alert>
                </Snackbar>
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
