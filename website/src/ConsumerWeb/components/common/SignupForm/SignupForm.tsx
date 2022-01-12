import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';

import Color from 'common/styles/colors';
import Form from 'common/components/Form/Form';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { TypographyColor } from 'common/typography/common';
import { withCaptchaToken, loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import {
    SignupUserResponse,
    SignupErrorMessage,
    signupUser,
} from 'ConsumerWeb/api/user/signup';

const styleClasses = makeStyles({
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    emailField: {
        width: '100%',
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
    iconContainer: {
        alignSelf: 'center',
    },
    confirmationIcon: {
        color: Color.Confirmation,
    },
    warningIcon: {
        color: Color.Warning,
    },
});

type SignupFormProps = {
    disabled: boolean;
    shouldShowVerificationForm: boolean;

    setIsLoading: (isLoading: boolean) => void;
    onSuccess: (emailAddress: string) => void;
}

const SignupForm = (props: SignupFormProps) => {
    const [ emailAddress, setEmailAddress ] = useState<string>(null);

    const handleSuccess = (emailAddress: string) => {
        props.onSuccess(emailAddress);
    }

    const body = props.shouldShowVerificationForm ? (
        <VerificationComponent
            disabled={props.disabled}
            setIsLoading={props.setIsLoading}
            emailAddress={emailAddress} />
    ) : (
        <SignupComponent
            emailAddress={emailAddress}
            disabled={props.disabled}
            setEmailAddress={setEmailAddress}
            setIsLoading={props.setIsLoading}
            onSuccess={handleSuccess} />
    )
    return body;
}

const errorMessages = {
    // TODO: think about these
    [SignupErrorMessage.InvalidEmailAddress]: "Hmm, the email address you gave doesn’t appear to be valid. Check to make sure that you spelled everything right.",
    [SignupErrorMessage.RateLimited]: "It looks like we’re having some trouble reaching you. Contact our support so we can get you on the list!",
    [SignupErrorMessage.IncorrectStatus]: "It looks like you’re already signed up for Babblegraph!",
    [SignupErrorMessage.LowScore]: "We’re having some trouble verifying your request. Contact us at hello@babblegraph.com to finish signing up.",
    "default": "Something went wrong. Contact our support so we can get you on the list!"

}

type SignupComponentProps = {
    emailAddress: string | null;
    disabled: boolean;

    setEmailAddress: (emailAddress: string) => void;
    setIsLoading: (isLoading: boolean) => void;
    onSuccess: (emailAddress: string) => void;
}

const SignupComponent = (props: SignupComponentProps) => {
    const [ errorMessage, setErrorMessage ] = useState<string | null>(null);

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setEmailAddress((event.target as HTMLInputElement).value);
    };

    const handleSubmit = () => {
        props.setIsLoading(true);
        withCaptchaToken("signup", (token: string) => {
            signupUser({
                emailAddress: props.emailAddress,
                captchaToken: token,
            },
            (resp: SignupUserResponse) => {
                props.setIsLoading(false);
                if (!!resp.errorMessage) {
                    setErrorMessage(errorMessages[resp.errorMessage] || errorMessages["default"]);
                } else if (resp.success) {
                    setErrorMessage(null);
                    props.onSuccess(props.emailAddress);
                } else {
                    setErrorMessage(errorMessages["default"]);
                }
            },
            (e: Error) => {
                props.setIsLoading(false);
                setErrorMessage(errorMessages["default"]);
            });
        });
    }

    const classes = styleClasses();
    return (
        <Form
            className={classes.confirmationForm}
            handleSubmit={handleSubmit}>
            <Grid container>
                <Grid item xs={9} md={10}>
                    <PrimaryTextField
                        id="email"
                        className={classes.emailField}
                        label="Email Address"
                        variant="outlined"
                        defaultValue={props.emailAddress}
                        disabled={props.disabled}
                        onChange={handleEmailAddressChange} />
                </Grid>
                <Grid item xs={3} md={2} className={classes.submitButtonContainer}>
                    <PrimaryButton disabled={!props.emailAddress || props.disabled}>
                        Try it!
                    </PrimaryButton>
                </Grid>
            </Grid>
            {
                !!errorMessage && (
                    <Grid container>
                        <Grid className={classes.iconContainer} item xs={1}>
                            <ErrorOutlineIcon className={classes.warningIcon} />
                        </Grid>
                        <Grid item xs={11}>
                            <Paragraph size={Size.Small} color={TypographyColor.Warning}>
                                {errorMessage}
                            </Paragraph>
                        </Grid>
                    </Grid>
                )
            }
        </Form>
    );
}

type VerificationComponentProps = {
    emailAddress: string | null;
    disabled: boolean;

    setIsLoading: (isLoading: boolean) => void;
}

const VerificationComponent = (props: VerificationComponentProps) => {
    const handleReenqueueSignupAttempt = () => {
        props.setIsLoading(true);
        withCaptchaToken("signup", (token: string) => {
            signupUser({
                emailAddress: props.emailAddress,
                captchaToken: token,
            },
            (resp: SignupUserResponse) => {
                props.setIsLoading(false);
            },
            (e: Error) => {
                props.setIsLoading(false);
            });
        });
    }

    const classes = styleClasses();
    return (
        <div>
            <Paragraph>
                Check your email for a verification email from hello@babblegraph.com. We sent it to {props.emailAddress}.
            </Paragraph>
            <Paragraph>
                You’ll need to click the button in the verification email that was just sent to you in order to start receiving emails from Babblegraph.
            </Paragraph>
            <Paragraph>
                It can take up to 5 minutes for the email to make its way to your inbox.
            </Paragraph>
            <Grid container>
                <Grid item xs={3} md={4}>
                    &nbsp;
                </Grid>
                <Grid item xs={6} md={4}>
                    <PrimaryButton className={classes.verificationButton} onClick={handleReenqueueSignupAttempt} disabled={props.disabled}>
                        Resend the verification email
                    </PrimaryButton>
                </Grid>
                <Grid item xs={3} md={4}>
                    &nbsp;
                </Grid>
            </Grid>
        </div>
    );
}

export default SignupForm;
