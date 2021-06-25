import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';

import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import { Heading1 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    ResetPasswordResponse,
    resetPassword,
    ResetPasswordError,
} from 'api/useraccounts/useraccounts';

const styleClasses = makeStyles({
    displayCard: {
        padding: '20px',
        marginTop: '20px',
    },
    submitButton: {
        margin: '10px 0',
    },
    textField: {
        margin: '10px 0',
        width: '100%',
    },
    textFieldContainer: {
        margin: '10px 0',
    },
    formGridContainer: {
        alignItems: 'center',
    },
    formGridItem: {
       padding: '5px',
    },
    iconContainer: {
        alignSelf: 'center',
    },
    warningIcon: {
        color: Color.Warning,
    },
});

const errorMessagesByType = {
    [ResetPasswordError.InvalidToken]: "The request could not be validated. Make sure that the email address that you entered is correct.",
    [ResetPasswordError.TokenExpired]: "The link to reset your password is expired. If you still need to change your password, request a new reset link from the login page.",
    [ResetPasswordError.PasswordRequirements]: "The password entered did not match the minimum password requirements",
    [ResetPasswordError.PasswordsNoMatch]: "The passwords entered did not match.",
    [ResetPasswordError.NoAccount]: "There is no account associated with the email address entered.",
   "default": "Something went wrong processing your request. Try again, or email hello@babblegraph.com for help.",
}

type Params = {
    token: string
}

type ResetPasswordPageProps = RouteComponentProps<Params>

const ResetPasswordPage = (props: ResetPasswordPageProps) => {
    const classes = styleClasses();
    const { token } = props.match.params;

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ password, setPassword ] = useState<string | null>(null);
    const [ confirmPassword, setConfirmPassword ] = useState<string | null>(null);
    const [ errorMessage, setErrorMessage ] = useState<string | null>(null);
    const [ managementToken, setManagementToken ] = useState<string | null>(null);

    const handleSubmit = () => {
        setIsLoading(true);
        resetPassword({
            resetPasswordToken: token,
            password: password,
            confirmPassword: confirmPassword,
            emailAddress: emailAddress,
        },
        (resp: ResetPasswordResponse) => {
            setIsLoading(false);
            if (!!resp.resetPasswordError) {
                setErrorMessage(errorMessagesByType[resp.resetPasswordError] || errorMessagesByType["default"]);
                setManagementToken(null);
            } else {
                setErrorMessage(null);
                setManagementToken(resp.managementToken);
            }
        },
        (e: Error) => {
            setIsLoading(false);
            setErrorMessage(errorMessagesByType["default"]);
            setManagementToken(null);
        });
    }

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard}>
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                <ResetPasswordForm
                                    emailAddress={emailAddress}
                                    password={password}
                                    confirmPassword={confirmPassword}
                                    errorMessage={errorMessage}
                                    updateEmailAddress={setEmailAddress}
                                    updatePassword={setPassword}
                                    updateConfirmPassword={setConfirmPassword}
                                    handleSubmit={handleSubmit} />
                            )
                        }
                        {
                            !!managementToken && (
                                <div>
                                    <Paragraph color={TypographyColor.Confirmation}>
                                        Your password was successfully updated!
                                    </Paragraph>
                                    <Link href={`/manage/${managementToken}`} target={LinkTarget.Self}>
                                        Go to subscription management page
                                    </Link>
                                </div>
                            )
                        }
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );

}

type ResetPasswordFormProps = {
    emailAddress: string | null;
    password: string | null;
    confirmPassword: string | null;
    errorMessage: string | null;

    updateEmailAddress: (emailAddress: string) => void;
    updatePassword: (password: string) => void;
    updateConfirmPassword: (password: string) => void;
    handleSubmit: () => void;
}

const ResetPasswordForm = (props: ResetPasswordFormProps) => {
    const classes = styleClasses();

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updateEmailAddress((event.target as HTMLInputElement).value);
    };
    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updatePassword((event.target as HTMLInputElement).value);
    };
    const handleConfirmPasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updateConfirmPassword((event.target as HTMLInputElement).value);
    };

    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                Reset your Password
            </Heading1>
            <Paragraph>
                Reset the password for your Babblegraph account
            </Paragraph>
            {
                !!props.errorMessage && (
                    <Grid container>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid className={classes.iconContainer} item xs={1} md={1}>
                            <ErrorOutlineIcon className={classes.warningIcon} />
                        </Grid>
                        <Grid item xs={10} md={5}>
                            <Paragraph size={Size.Small} color={TypographyColor.Warning}>
                                {props.errorMessage}
                            </Paragraph>
                        </Grid>
                    </Grid>
                )
            }
            <Grid container className={classes.textFieldContainer}>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <PrimaryTextField
                        className={classes.textField}
                        id="email"
                        label="Email Address"
                        variant="outlined"
                        defaultValue={props.emailAddress}
                        onChange={handleEmailAddressChange} />
                    <Paragraph align={Alignment.Left}>
                        Password Requirements:
                        <ul>
                            <li>At least 8 characters</li>
                            <li>No more than 32 characters</li>
                            <li>Must include 3 of the following:
                                <ul>
                                    <li>Upper Case Latin Letter (A-Z)</li>
                                    <li>Lower Case Latin Letter (a-z)</li>
                                    <li>Number (0-9)</li>
                                    <li>Special Character (such as !@#$%^&*)</li>
                                </ul>
                            </li>
                        </ul>
                    </Paragraph>
                    <PrimaryTextField
                        className={classes.textField}
                        id="password"
                        label="Password"
                        type="password"
                        variant="outlined"
                        defaultValue={props.password}
                        onChange={handlePasswordChange} />
                    <PrimaryTextField
                        className={classes.textField}
                        id="confirm-password"
                        label="Confirm Password"
                        type="password"
                        variant="outlined"
                        defaultValue={props.confirmPassword}
                        onChange={handleConfirmPasswordChange} />
                    <PrimaryButton
                        className={classes.submitButton}
                        onClick={props.handleSubmit}
                        disabled={!props.emailAddress && !props.password && !props.confirmPassword}>
                        Reset Password
                    </PrimaryButton>
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
            </Grid>
        </div>
    );
}

export default ResetPasswordPage;
