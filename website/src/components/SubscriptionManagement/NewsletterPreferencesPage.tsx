import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import { ContentHeader } from './common';

import {
    getUserProfile,
    GetUserProfileResponse
} from 'api/useraccounts/useraccounts';

type Params = {
    token: string;
}

type NewsletterPreferencesPageProps = RouteComponentProps<Params>

const NewsletterPreferencesPage = (props: NewsletterPreferencesPageProps) =>  {
    const { token } = props.match.params;

    const [ isWordReinforcementSpotlightActive, setIsWordReinforcementSpotlightActive ] = useState<boolean | null>(null);
    const [ isLoadingUserNewsletterPreferences, setLoadingUserNewsletterPreferences ] = useState<boolean>(false);

    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ subscriptionLevel, setSubscriptionLevel ] = useState<string | undefined>(undefined);

    const [ isLoadingUserProfile, setIsLoadingUserProfile ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setIsLoadingUserProfile(false);
            if (resp.subscriptionLevel) {
                setSubscriptionLevel(resp.subscriptionLevel);
                setEmailAddress(resp.emailAddress);
            }
        },
        (e: Error) => {
            setIsLoadingUserProfile(false);
            setError(e);
        });
    }, []);

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
                            </DisplayCard>
                        )
                    }
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
            <Grid item xs={2} xl={1}>
                <PrimarySwitch checked={props.isActive} onClick={handleToggle} disabled={props.isActive == null} />
            </Grid>
        </Grid>
    );
}

export default NewsletterPreferencesPage;
