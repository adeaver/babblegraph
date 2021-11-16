import React, { useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';

import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import { Heading1 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    loginUser,
    LoginUserResponse,
    LoginError,
} from 'ConsumerWeb/api/useraccounts/useraccounts';

const styleClasses = makeStyles({
    displayCard: {
        padding: '20px',
        marginTop: '20px',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    textField: {
        width: '100%',
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

const loginErrorMessages = {
    [LoginError.InvalidCredentials]: "The username or password entered is not correct",
    "default": "Something went wrong, try again later.",
}

type LoginPageProps = {};

const LoginPage = (props: LoginPageProps) => {
    const classes = styleClasses();
    const history = useHistory();

    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ password, setPassword ] = useState<string | null>(null);
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ errorMessage, setErrorMessage ] = useState<string | null>(null);


    const urlSearchParams = new URLSearchParams(window.location.search);

    const handleSubmit = () => {
        setIsLoading(true);
        loginUser({
            emailAddress: emailAddress,
            password: password,
            redirectKey: urlSearchParams.get("d") || "",
        },
        (resp: LoginUserResponse) => {
            setIsLoading(false);
            if (!!resp.location) {
                history.push(resp.location);
            } else if (!!resp.loginError) {
                setErrorMessage(loginErrorMessages[resp.loginError] || loginErrorMessages["default"]);
            } else {
                setErrorMessage(loginErrorMessages["default"]);
            }
        },
        (e: Error) => {
            setIsLoading(false);
            setErrorMessage(loginErrorMessages["default"]);
        })
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
                                <LoginForm
                                    emailAddress={emailAddress}
                                    password={password}
                                    errorMessage={errorMessage}
                                    updateEmailAddress={setEmailAddress}
                                    updatePassword={setPassword}
                                    handleSubmit={handleSubmit} />
                            )
                        }
                    </Card>
                </Grid>
            </Grid>

        </Page>
    );
}

type LoginFormProps = {
    emailAddress: string | null;
    password: string | null;
    errorMessage: string | null;

    updateEmailAddress: (emailAddress: string) => void;
    updatePassword: (password: string) => void;
    handleSubmit: () => void;
}

const LoginForm = (props: LoginFormProps) => {
    const classes = styleClasses();

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.handleSubmit();
    }
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updateEmailAddress((event.target as HTMLInputElement).value);
    };
    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updatePassword((event.target as HTMLInputElement).value);
    };

    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                Login to Babblegraph Premium
            </Heading1>
            <Paragraph>
                If you have a Babblegraph Premium Account, you can sign in here.
            </Paragraph>
            <Paragraph>
                If you don’t have a premium account and you’re trying to add new words or change your interest settings, you can do that by clicking the “Manage your subscription” link at the bottom of your most recent daily email.
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
            <form onSubmit={handleSubmit} noValidate autoComplete="off">
                <Grid container className={classes.formGridContainer}>
                    <Grid item xs={12} md={5} className={classes.formGridItem}>
                        <PrimaryTextField
                            className={classes.textField}
                            id="email"
                            label="Email Address"
                            variant="outlined"
                            defaultValue={props.emailAddress}
                            onChange={handleEmailAddressChange} />
                    </Grid>
                    <Grid item xs={12} md={5} className={classes.formGridItem}>
                        <PrimaryTextField
                            className={classes.textField}
                            id="password"
                            label="Password"
                            type="password"
                            variant="outlined"
                            defaultValue={props.password}
                            onChange={handlePasswordChange} />
                    </Grid>
                    <Grid item xs={3} md={2} className={classes.formGridItem}>
                        <PrimaryButton type="submit" disabled={!props.emailAddress || !props.password}>
                            Login
                        </PrimaryButton>
                    </Grid>
                </Grid>
            </form>
            <Link href="/forgot-password" target={LinkTarget.Self}>
                Forgot your password?
            </Link>
            <Paragraph>
                If you don’t have a premium account and would like one, check out the manage your subscription link on the bottom of your daily newsletter.
            </Paragraph>
        </div>
    );
}

export default LoginPage;
