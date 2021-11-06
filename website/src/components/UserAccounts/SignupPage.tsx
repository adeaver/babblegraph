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

import {
    createUser,
    CreateUserError,
    CreateUserResponse,
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

const createUserErrorMessages = {
   [CreateUserError.AlreadyExists]: "There’s already an existing account for that email address",
   [CreateUserError.InvalidToken]: "The email submitted didn’t match the email address this unique link is for. Make sure you entered the same email address that you received the signup link with.",
   [CreateUserError.PasswordRequirements]: "The password entered did not match the minimum password requirements",
   [CreateUserError.NoSubscription]: "At this time, you need to be subscribed to create an account. Please see the BuyMeACoffee Link to subscribe",
   [CreateUserError.PasswordsNoMatch]: "The passwords entered did not match.",
   "default": "Something went wrong processing your request. Try again, or email hello@babblegraph.com for help.",
}

type Params = {
    token: string
}

type SignupPageProps = RouteComponentProps<Params>

const SignupPage = (props: SignupPageProps) => {
    const classes = styleClasses();
    const history = useHistory();
    const { token } = props.match.params;

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ password, setPassword ] = useState<string | null>(null);
    const [ confirmPassword, setConfirmPassword ] = useState<string | null>(null);
    const [ errorMessage, setErrorMessage ] = useState<string | null>(null);

    const handleSubmit = () => {
        setIsLoading(true);
        createUser({
            createUserToken: token,
            emailAddress: emailAddress,
            password: password,
            confirmPassword: confirmPassword,
        },
        (resp: CreateUserResponse) => {
            setIsLoading(false);
            if (!!resp.checkoutToken) {
                history.push(`/checkout/${resp.checkoutToken}`);
            } else if (!!resp.createUserError) {
                setErrorMessage(createUserErrorMessages[resp.createUserError] || createUserErrorMessages["default"]);
            } else {
                setErrorMessage(createUserErrorMessages["default"]);
            }
        },
        (e: Error) => {
            setIsLoading(false);
            setErrorMessage(createUserErrorMessages["default"]);
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
                                <SignupForm
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
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );

}

type SignupFormProps = {
    emailAddress: string | null;
    password: string | null;
    confirmPassword: string | null;
    errorMessage: string | null;

    updateEmailAddress: (emailAddress: string) => void;
    updatePassword: (password: string) => void;
    updateConfirmPassword: (password: string) => void;
    handleSubmit: () => void;
}

const SignupForm = (props: SignupFormProps) => {
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
                First step, sign up for a Babblegraph account
            </Heading1>
            <Paragraph>
                Why do you need to sign up for an account to access Babblegraph Premium? Security is a big concern when dealing with payment information. Accounts are more secure than managing your Babblegraph subscription.
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
                        label="Confirm Your Email Address"
                        variant="outlined"
                        defaultValue={props.emailAddress}
                        onChange={handleEmailAddressChange} />
                    <Paragraph align={Alignment.Left}>
                        Password Requirements:
                        <ul>
                            <PasswordConstraint isConstraintMet={props.password && props.password.length > 8}>
                                At least 8 characters
                            </PasswordConstraint>
                            <PasswordConstraint isConstraintMet={props.password && props.password.length < 32} >
                                No more than 32 characters
                            </PasswordConstraint>
                            <PasswordConstraint isConstraintMet={false}>
                                At least three of the following:
                            </PasswordConstraint>
                                <ul>
                                    <PasswordConstraint isConstraintMet={props.password && !!props.password.match(/[a-z]/)}>
                                        Lower Case Latin Letter (a-z)
                                    </PasswordConstraint>
                                    <PasswordConstraint isConstraintMet={props.password && !!props.password.match(/[A-Z]/)}>
                                        Upper Case Latin Letter (A-Z)
                                    </PasswordConstraint>
                                    <PasswordConstraint isConstraintMet={props.password && !!props.password.match(/[0-9]/)}>
                                        Number (0-9)
                                    </PasswordConstraint>
                                    <PasswordConstraint isConstraintMet={props.password && !!props.password.match(/[^0-9a-zA-Z]/)}>
                                        Special Character (such as !@#$%^&*)
                                    </PasswordConstraint>
                                </ul>
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
                        Sign Up
                    </PrimaryButton>
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
            </Grid>
        </div>
    );
}

type PasswordConstraintProps = {
    isConstraintMet: boolean;
    children: string | JSX.Element;
}

const PasswordConstraint = (props: PasswordConstraintProps) => {
    return (
        <li>
            <Paragraph
                color={props.isConstraintMet ? TypographyColor.Confirmation : TypographyColor.Gray}
                align={Alignment.Left}>
                {props.children}
            </Paragraph>
        </li>
    )
}

export default SignupPage;
