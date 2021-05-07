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
import Link from 'common/components/Link/Link';

import {
    loginUser,
    LoginUserResponse,
    LoginError,
} from 'api/useraccounts/useraccounts';

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


    // TODO: Add useEffect here to see if
    // user is already logged in.

    const handleSubmit = () => {
        setIsLoading(true);
        loginUser({
            emailAddress: emailAddress,
            password: password,
        },
        (resp: LoginUserResponse) => {
            setIsLoading(false);
            if (!!resp.managementToken) {
                history.push(`/manage/${resp.managementToken}`);
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

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updateEmailAddress((event.target as HTMLInputElement).value);
    };
    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updatePassword((event.target as HTMLInputElement).value);
    };

    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                Login to Babblegraph
            </Heading1>
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
                <Grid item xs={12} md={2} className={classes.formGridItem}>
                    <PrimaryButton onClick={props.handleSubmit} disabled={!props.emailAddress && !props.password}>
                        Login
                    </PrimaryButton>
                </Grid>
            </Grid>
            <Link href="/forgot-password">
                Forgot your password?
            </Link>
        </div>
    );
}

export default LoginPage;
