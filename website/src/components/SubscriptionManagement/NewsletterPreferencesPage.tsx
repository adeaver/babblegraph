import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryButton } from 'common/components/Button/Button';

import { ContentHeader } from './common';

import {
    getUserProfile,
    GetUserProfileResponse
} from 'api/useraccounts/useraccounts';
import {
    getUserNewsletterPreferences,
    GetUserNewsletterPreferencesResponse,

    updateUserNewsletterPreferences,
    UpdateUserNewsletterPreferencesResponse,

    UserNewsletterPreferences,
} from 'api/user/userNewsletterPreferences';

const styleClasses = makeStyles({
    buttonContainer: {
        padding: '20px',
        minWidth: '100%',
        maxWidth: '100%',
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
    },
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

type NewsletterPreferencesPageProps = RouteComponentProps<Params>

const NewsletterPreferencesPage = (props: NewsletterPreferencesPageProps) =>  {
    const { token } = props.match.params;

    const [ isWordReinforcementSpotlightActive, setIsWordReinforcementSpotlightActive ] = useState<boolean | null>(null);
    const [ isLoadingUserNewsletterPreferences, setIsLoadingUserNewsletterPreferences ] = useState<boolean>(true);

    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);

    const [ didUpdate, setDidUpdate ] = useState<boolean | null>(null);

    const [ isLoadingUserProfile, setIsLoadingUserProfile ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setIsLoadingUserProfile(false);
            if (resp.subscriptionLevel) {
                setEmailAddress(resp.emailAddress);
                getUserNewsletterPreferences({
                    emailAddress: resp.emailAddress,
                    // TODO(multiple-languages): don't hardcode this
                    languageCode: "es",
                    subscriptionManagementToken: token,
                },
                (resp: GetUserNewsletterPreferencesResponse) => {
                    setIsLoadingUserNewsletterPreferences(false);
                    setIsWordReinforcementSpotlightActive(resp.preferences.isLemmaReinforcementSpotlightActive);
                },
                (e: Error) => {
                    setIsLoadingUserNewsletterPreferences(false);
                    setError(e);
                });
            } else {
                setIsLoadingUserNewsletterPreferences(false);
            }
        },
        (e: Error) => {
            setIsLoadingUserProfile(false);
            setError(e);
        });
    }, []);

    const handleSubmit = () => {
        setIsLoadingUserNewsletterPreferences(true);
        updateUserNewsletterPreferences({
            // TODO(multiple-languages): don't hardcode this
           languageCode: "es",
           emailAddress: emailAddress,
           subscriptionManagementToken: token,
           preferences: {
                isLemmaReinforcementSpotlightActive: isWordReinforcementSpotlightActive,
           },
        },
        (resp: UpdateUserNewsletterPreferencesResponse) => {
            setIsLoadingUserNewsletterPreferences(false);
            setDidUpdate(resp.success);
        },
        (e: Error) => {
            setIsLoadingUserNewsletterPreferences(false);
            setError(e);
        });
    }
    const closeSnackbar = () => {
        setError(null);
        setDidUpdate(null);
    }

    const classes = styleClasses();
    const isLoading = isLoadingUserProfile || isLoadingUserNewsletterPreferences;
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    {
                        isLoading ? (
                            <LoadingSpinner />
                        ) : (
                            <DisplayCard>
                                <ContentHeader
                                    title="Newsletter General Settings"
                                    token={token} />
                                <Paragraph align={Alignment.Left}>
                                    You can adjust some general settings for your newsletter here.
                                </Paragraph>
                                <Divider />
                                {
                                    !subscriptionLevel && (
                                        <Paragraph align={Alignment.Left} color={TypographyColor.Warning}>
                                            Since you’re not subscribed, there are no settings that you will be able to adjust.
                                        </Paragraph>
                                    )
                                }
                                <LemmaReinforcementHighlightToggle
                                    isActive={isWordReinforcementSpotlightActive}
                                    toggleIsActive={setIsWordReinforcementSpotlightActive} />
                                <Divider />
                                <div className={classes.buttonContainer}>
                                    <PrimaryButton onClick={handleSubmit}>
                                        Save your preferences
                                    </PrimaryButton>
                                </div>
                            </DisplayCard>
                        )
                    }
                    <Snackbar open={(!didUpdate && didUpdate != null) || !!error} autoHideDuration={6000} onClose={closeSnackbar}>
                        <Alert severity="error">Something went wrong processing your request.</Alert>
                    </Snackbar>
                    <Snackbar open={didUpdate} autoHideDuration={6000} onClose={closeSnackbar}>
                        <Alert severity="success">Successfully updated your email preferences.</Alert>
                    </Snackbar>
                </Grid>
            </Grid>
        </Page>
    );
}

type LemmaReinforcementHighlightToggleProps = {
    isActive: boolean | null;
    toggleIsActive: (boolean) => void;
}

const LemmaReinforcementHighlightToggle = (props: LemmaReinforcementHighlightToggleProps) => {
    const handleToggle = () => {
        props.toggleIsActive(!props.isActive);
    }

    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={10} xl={11}>
                <Heading3 align={Alignment.Left} color={TypographyColor.Primary}>
                    Include word tracking spotlights in your newsletter?
                </Heading3>
                <Paragraph align={Alignment.Left}>
                    Word tracking spotlights include a highlighted article that is guaranteed to have a word in your tracking list. It spaces out these spotlights so you can practice new words on your list. If this is disabled, you won't see spotlights in your newsletter.
                </Paragraph>
            </Grid>
            <Grid item
                className={classes.toggleContainer}
                xs={2}
                xl={1}>
                <PrimarySwitch checked={props.isActive} onClick={handleToggle} disabled={props.isActive == null} />
            </Grid>
        </Grid>
    );
}

export default NewsletterPreferencesPage;
