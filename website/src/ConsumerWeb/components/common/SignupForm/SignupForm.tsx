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
import { DisplayLanguage } from 'common/model/language/language';

import { TextBlock, getTextBlocksForLanguage } from './translations';

import {
    SignupUserResponse,
    SignupErrorMessage,
    signupUser,
} from 'ConsumerWeb/api/user/signup';

const styleClasses = makeStyles({
    submitButtonContainer: {
        display: 'flex',
        justifyContent: 'end',
        padding: '5px 0',
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

declare const window: any;

type SignupFormProps = {
    disabled: boolean;
    shouldShowVerificationForm: boolean;
    displayLanguage?: DisplayLanguage;

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
            displayLanguage={props.displayLanguage}
            disabled={props.disabled}
            setIsLoading={props.setIsLoading}
            emailAddress={emailAddress} />
    ) : (
        <SignupComponent
            displayLanguage={props.displayLanguage}
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
    [SignupErrorMessage.InvalidEmailAddress]: TextBlock.ErrorMessageInvalidEmailAddress,
    [SignupErrorMessage.RateLimited]: TextBlock.ErrorMessageRateLimited,
    [SignupErrorMessage.IncorrectStatus]: TextBlock.ErrorMessageIncorrectStatus,
    [SignupErrorMessage.LowScore]: TextBlock.ErrorMessageLowScore,
    "default": TextBlock.ErrorMessageDefault,

}

type SignupComponentProps = {
    emailAddress: string | null;
    disabled: boolean;
    displayLanguage?: DisplayLanguage;

    setEmailAddress: (emailAddress: string) => void;
    setIsLoading: (isLoading: boolean) => void;
    onSuccess: (emailAddress: string) => void;
}

const SignupComponent = (props: SignupComponentProps) => {
    const [ errorMessage, setErrorMessage ] = useState<string | null>(null);

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setEmailAddress((event.target as HTMLInputElement).value);
    };

    const translations = getTextBlocksForLanguage(props.displayLanguage);

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
                    setErrorMessage(translations[errorMessages[resp.errorMessage] || errorMessages["default"]]);
                } else if (resp.success) {
                    window.gtag_report_conversion(undefined);
                    setErrorMessage(null);
                    props.onSuccess(props.emailAddress);
                } else {
                    setErrorMessage(translations[errorMessages["default"]]);
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
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="email"
                        className={classes.emailField}
                        label={translations[TextBlock.EmailAddressInputLabel]}
                        variant="outlined"
                        defaultValue={props.emailAddress}
                        disabled={props.disabled}
                        onChange={handleEmailAddressChange} />
                </Grid>
                <Grid item xs={12} className={classes.submitButtonContainer}>
                    <PrimaryButton type="submit" disabled={!props.emailAddress || props.disabled}>
                        { translations[TextBlock.SignupFormButtonText] }
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
    displayLanguage?: DisplayLanguage;

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

    const translations = getTextBlocksForLanguage(props.displayLanguage);

    const classes = styleClasses();
    return (
        <div>
            <Paragraph>
                {`${translations[TextBlock.VerificationLocationConfirmation]} ${props.emailAddress}`}.
            </Paragraph>
            <Paragraph>
                { translations[TextBlock.VerificationInstructions] }
            </Paragraph>
            <Paragraph>
                { translations[TextBlock.VerificationWarningDisclaimer] }
            </Paragraph>
            <Grid container>
                <Grid item xs={3} md={4}>
                    &nbsp;
                </Grid>
                <Grid item xs={6} md={4}>
                    <PrimaryButton className={classes.verificationButton} onClick={handleReenqueueSignupAttempt} disabled={props.disabled}>
                        { translations[TextBlock.VerificationResendButtonText] }
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
