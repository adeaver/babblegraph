import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { Alignment, TypographyColor } from 'common/typography/common';
import Form from 'common/components/Form/Form';

import {
    UnsubscribeError,
    UnsubscribeResponse,
    unsubscribeUser,
} from 'ConsumerWeb/api/user/unsubscribe';

import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';

const styleClasses = makeStyles({
    submitButtonContainer: {
        display: 'flex',
        justifyContent: 'center',
        alignSelf: 'center',
        padding: '5px',
    },
    textField: {
        width: '100%',
        margin: '10px 0',
    },
    formContainer: {
        padding: '10px 0',
    },
});

const errorMessages = {
    [UnsubscribeError.MissingEmail]: "You need to include your email address to unsubscribe.",
    [UnsubscribeError.IncorrectEmail]: "The email address you entered is incorrect",
    [UnsubscribeError.NoAuth]: "The email address you entered is incorrect",
    [UnsubscribeError.IncorrectKey]: "The email address you entered is incorrect",
    [UnsubscribeError.InvalidToken]: "The email address you entered is incorrect",
    "default": "There was a problem processing your request. Try again later or contact hello@babblegraph.com to unsubscribe",
}

type Params = {
    token: string;
}

type UnsubscribePageOwnProps = RouteComponentProps<Params>;

const UnsubscribePage = withUserProfileInformation<UnsubscribePageOwnProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [],
    (ownProps: UnsubscribePageOwnProps) => {
        return ownProps.match.params.token;
    },
    undefined,
    (props:  UnsubscribePageOwnProps & UserProfileComponentProps) => {
        const { token } = props.match.params;

        const [ isLoading, setIsLoading ] = useState<boolean>(false);

        const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
        const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setEmailAddress((event.target as HTMLInputElement).value);
        };

        const [ unsubscribeReason, setUnsubscribeReason ] = useState<string | null>(null);
        const [ validationError, setValidationError ] = useState<string | null>(null);
        const handleUnsubscribeReasonChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            const unsubscribeReason = (event.target as HTMLInputElement).value;
            if (unsubscribeReason.length > 500) {
                setValidationError("Unsubscribe reason needs to be less than 500 characters");
            } else {
                setValidationError(null);
            }
            setUnsubscribeReason(unsubscribeReason);
        };

        const [ didUpdate, setDidUpdate ] = useState<boolean>(false);
        const [ errorMessage, setErrorMessage ] = useState<string>(null);
        const handleSubmit = () => {
            setIsLoading(true);
            unsubscribeUser({
                token: token,
                unsubscribeReason: unsubscribeReason,
                emailAddress: props.userProfile.hasAccount ? undefined : emailAddress,
            },
            (resp: UnsubscribeResponse) => {
                setIsLoading(false);
                if (!!resp.error) {
                    setErrorMessage(errorMessages[resp.error] || errorMessage["default"]);
                } else {
                    setDidUpdate(resp.success);
                }
            },
            (err: Error) => {
                setIsLoading(false);
                setErrorMessage(errorMessage["default"]);
            });
        }

        const classes = styleClasses();
        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Unsubscribe"
                        backArrowDestination={`/manage/${props.match.params.token}`} />
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                <Form handleSubmit={handleSubmit} className={classes.formContainer}>
                                    <Paragraph size={Size.Medium} align={Alignment.Left}>
                                        We’re sorry to see you go! When you unsubscribe, you won’t receive any more daily emails or any other communication from us. But you can always come back by resubscribing on the homepage.
                                    </Paragraph>
                                    <Paragraph size={Size.Medium} align={Alignment.Left}>
                                    {
                                        props.userProfile.hasAccount ? (
                                            "To unsubscribe, just click the unsubscribe button."
                                        ) : (
                                            "To unsubscribe, just enter your email and click the unsubscribe button."
                                        )
                                    }
                                    </Paragraph>
                                    <Grid container>
                                        <Grid item xs={12}>
                                            <PrimaryTextField
                                                id="unsubscribe-reason"
                                                className={classes.textField}
                                                label="Reason for unsubscribing (optional)"
                                                variant="outlined"
                                                rows={3}
                                                error={validationError}
                                                onChange={handleUnsubscribeReasonChange}
                                                multiline />
                                        </Grid>
                                        {
                                            validationError && (
                                                <Grid item xs={12}>
                                                    <Paragraph color={TypographyColor.Warning}>
                                                        { validationError }
                                                    </Paragraph>
                                                </Grid>
                                            )
                                        }
                                        {
                                            !props.userProfile.hasAccount && (
                                                <Grid item xs={9} md={10}>
                                                    <PrimaryTextField
                                                        id="email"
                                                        className={classes.textField}
                                                        label="Email Address"
                                                        variant="outlined"
                                                        onChange={handleEmailAddressChange} />
                                                </Grid>
                                            )
                                        }
                                        <Grid item xs={props.userProfile.hasAccount ? 12 : 3} md={props.userProfile.hasAccount ? 12 : 2} className={classes.submitButtonContainer}>
                                            <PrimaryButton type="submit" disabled={(!emailAddress && !props.userProfile.hasAccount) || !!validationError}>
                                                Unsubscribe
                                            </PrimaryButton>
                                        </Grid>
                                    </Grid>
                                    {
                                        !!props.userProfile.subscriptionLevel && (
                                            <div>
                                                <Paragraph>
                                                    This will also cancel your subscription to Babblegraph Premium.
                                                </Paragraph>
                                            </div>
                                        )
                                    }
                            </Form>
                        )
                    }
                </DisplayCard>
                <Snackbar open={!!errorMessage} autoHideDuration={6000} onClose={() => {setErrorMessage(null)}}>
                    <Alert severity="error">{errorMessage}</Alert>
                </Snackbar>
                <Snackbar open={didUpdate} autoHideDuration={6000} onClose={() => {setDidUpdate(false)}}>
                    <Alert severity="success">You’ve been successfully unsubscribed from Babblegraph.</Alert>
                </Snackbar>
            </CenteredComponent>
        );
    }
);

export default UnsubscribePage;
