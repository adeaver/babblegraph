import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import { WordsmithLanguageCode } from 'common/model/language/language';
import { DisplayLanguage } from 'common/model/language/language';
import { ClientError, asReadable } from 'ConsumerWeb/api/clienterror';

import {
    Topic,
} from 'ConsumerWeb/api/content/content';
import {
    UpsertUserContentTopicsResponse,
    upsertUserContentTopics,
} from 'ConsumerWeb/api/user/content';

import {
    TopicSelector,
} from './components';

const styleClasses = makeStyles({
    submitButton: {
        width: '100%',
        padding: '5px',
    },
});

type Params = {
    token: string
}

type InterestSelectionPageProps = RouteComponentProps<Params>;

const InterestSelectionPage = withUserProfileInformation<InterestSelectionPageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [],
    (ownProps: InterestSelectionPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.ContentTopics,
    (props: InterestSelectionPageProps & UserProfileComponentProps) => {
        const subscriptionManagementToken = props.match.params.token;

        const [ emailAddress, setEmailAddress ] = useState<string>(null);
        const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setEmailAddress((event.target as HTMLInputElement).value);
        }

        const [ userTopics, setUserTopics ] = useState<Array<Topic>>([]);

        const [ isLoading, setIsLoading ] = useState<boolean>(false);
        const [ errorMessage, setErrorMessage ] = useState<string>(null);
        const [ wasSuccessful, setWasSuccessful ] = useState<boolean>(false);

        const handleSubmit = () => {
            setIsLoading(true);
            upsertUserContentTopics({
                subscriptionManagementToken: subscriptionManagementToken,
                topicIds: userTopics.map((t: Topic) => t.topicId),
                emailAddress: emailAddress,
                languageCode: WordsmithLanguageCode.Spanish,
            },
            (resp: UpsertUserContentTopicsResponse) => {
                setIsLoading(false);
                setWasSuccessful(resp.success);
                if (resp.error) {
                    setErrorMessage(asReadable(resp.error, DisplayLanguage.English));
                }
            },
            (err: Error) => {
                setIsLoading(false);
                setErrorMessage(asReadable(ClientError.DefaultError, DisplayLanguage.English));
            });
        }

        const classes = styleClasses();
        return (
            <CenteredComponent useLargeVersion>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Manage your interests"
                        backArrowDestination={`/manage/${props.match.params.token}`} />
                    <Form handleSubmit={handleSubmit}>
                        <TopicSelector
                            subscriptionManagementToken={subscriptionManagementToken}
                            isDisabled={isLoading}
                            languageCode={WordsmithLanguageCode.Spanish}
                            handleUserTopicsChange={setUserTopics} />
                        <Grid container>
                            {
                                !props.userProfile.hasAccount ? (
                                    <Grid item xs={8}>
                                        <PrimaryTextField
                                            id="email"
                                            label="Email Address"
                                            variant="outlined"
                                            defaultValue={emailAddress}
                                            disabled={isLoading}
                                            onChange={handleEmailAddressChange} />
                                    </Grid>
                                ) : (
                                    <Grid item xs={false} md={3}>
                                        &nbsp;
                                    </Grid>
                                )
                            }
                            <Grid item xs={!props.userProfile.hasAccount ? 4 : 12} md={!props.userProfile.hasAccount ? 4 : 6}>
                                <PrimaryButton
                                    className={classes.submitButton}
                                    type='submit'
                                    disabled={(!props.userProfile.hasAccount && !emailAddress) || isLoading}>
                                    Update interests
                                </PrimaryButton>
                            </Grid>
                            {
                                isLoading && <LoadingSpinner />
                            }
                        </Grid>
                    </Form>
                    <Snackbar open={!!errorMessage} autoHideDuration={6000} onClose={() => {setErrorMessage(null)}}>
                        <Alert severity="error">{errorMessage}</Alert>
                    </Snackbar>
                    <Snackbar open={wasSuccessful} autoHideDuration={6000} onClose={() => {setWasSuccessful(false)}}>
                        <Alert severity="success">Success! It may take up to 24 hours to see your changes take effect.</Alert>
                    </Snackbar>
                </DisplayCard>
            </CenteredComponent>
        );
    }
);

export default InterestSelectionPage;
