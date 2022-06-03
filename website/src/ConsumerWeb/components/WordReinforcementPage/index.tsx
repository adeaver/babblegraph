import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';

import WordSearchForm from "./WordSearchForm";
import UserVocabularyDisplay from "./UserVocabularyDisplay"
import {
    withUserVocabulary,
    InjectedUserVocabularyComponentProps,
} from './withUserVocabulary';

type Params = {
    token: string;
}

type WordReinforcementPageProps = RouteComponentProps<Params>;

const WordReinforcementPage = withUserProfileInformation<WordReinforcementPageProps>(
    RouteEncryptionKey.WordReinforcement,
    [ RouteEncryptionKey.SubscriptionManagement ],
    (ownProps: WordReinforcementPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.Vocabulary,
    (props: WordReinforcementPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;
        const [ subscriptionManagementToken ] = props.userProfile.nextTokens;

        const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);

        useEffect(() => {
            loadCaptchaScript();
            setHasLoadedCaptcha(true);
        }, []);

        if (!hasLoadedCaptcha) {
            return <LoadingSpinner />;
        }

        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Add vocabulary words"
                        backArrowDestination={`/manage/${subscriptionManagementToken}`} />
                    <WordReinforcementBody
                        wordReinforcementToken={token}
                        subscriptionManagementToken={subscriptionManagementToken} />
                </DisplayCard>
            </CenteredComponent>
        );
    }
);

type WordReinforcementBodyProps = {
    wordReinforcementToken: string;
    subscriptionManagementToken: string;
}

const WordReinforcementBody = withUserVocabulary(
    (props: WordReinforcementBodyProps & InjectedUserVocabularyComponentProps) => (
        <Grid container>
            <Grid item xs={12}>
                <WordSearchForm
                    wordReinforcementToken={props.wordReinforcementToken}
                    subscriptionManagementToken={props.subscriptionManagementToken}
                    userVocabularyEntries={props.userVocabularyEntries}
                    handleAddNewUserVocabularyEntry={props.handleAddNewVocabularyEntry} />
            </Grid>
            <Grid item xs={12}>
                <UserVocabularyDisplay
                    subscriptionManagementToken={props.subscriptionManagementToken}
                    userVocabularyEntries={props.userVocabularyEntries}
                    handleRemoveVocabularyEntry={props.handleRemoveVocabularyEntry} />
            </Grid>
        </Grid>
    )
)

export default WordReinforcementPage;
