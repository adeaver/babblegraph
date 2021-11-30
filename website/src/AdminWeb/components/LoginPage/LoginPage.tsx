import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import Page from 'common/components/Page/Page';
import { TypographyColor } from 'common/typography/common';
import { Heading1 } from 'common/typography/Heading';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { setLocation } from 'util/window/Location';

import {
    validateLoginCredentials,
    ValidateLoginCredentialsResponse,
    validateTwoFactorAuthenticationCode,
    ValidateTwoFactorAuthenticationCodeResponse,
} from 'AdminWeb/api/auth/auth';

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
});

type LoginPageProps = {};

const LoginPage = (props: LoginPageProps) => {
    const classes = styleClasses();

    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ password, setPassword ] = useState<string | null>(null);

    const [ twoFactorAuthenticationCode, setTwoFactorAuthenticationCode ] = useState<string>("");

    const [ shouldShowLoginForm, setShouldShowLoginForm ] = useState<boolean>(true);

    const handleSubmit = () => {
        setIsLoading(true);
        validateLoginCredentials({
            emailAddress: emailAddress,
            password: password,
        },
        (resp: ValidateLoginCredentialsResponse) => {
            setIsLoading(false);
            if (resp.success) {
                setShouldShowLoginForm(false);
            }
        },
        (err: Error) => {
            setIsLoading(false);
        });
    }
    const handleSubmitTwoFactorAuthenticationCode = () => {
        setIsLoading(true);
        validateTwoFactorAuthenticationCode({
            emailAddress: emailAddress,
            twoFactorAuthenticationCode: twoFactorAuthenticationCode,
        },
        (resp: ValidateTwoFactorAuthenticationCodeResponse) => {
            setIsLoading(false);
            if (resp.success) {
                setLocation("/ops/dashboard");
            }
        },
        (err: Error) => {
            setIsLoading(false);
        });
    }

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <Heading1 color={TypographyColor.Primary}>
                            babblegraph
                        </Heading1>
                        {
                            shouldShowLoginForm ? (
                                <LoginForm
                                    emailAddress={emailAddress}
                                    setEmailAddress={setEmailAddress}
                                    password={password}
                                    setPassword={setPassword}
                                    handleSubmit={handleSubmit}
                                    isLoading={isLoading} />
                            ) : (
                                <TwoFactorAuthenticationForm
                                    twoFactorAuthenticationCode={twoFactorAuthenticationCode}
                                    setTwoFactorAuthenticationCode={setTwoFactorAuthenticationCode}
                                    handleSubmit={handleSubmitTwoFactorAuthenticationCode}
                                    isLoading={isLoading} />
                            )
                        }
                    </Card>
                </Grid>
            </Grid>
        </Page>
    )
}

type LoginFormProps = {
    emailAddress: string;
    setEmailAddress: (v: string) => void;

    password: string;
    setPassword: (v: string) => void;

    handleSubmit: () => void;
    isLoading: boolean;
}

const LoginForm = (props: LoginFormProps) => {
    const classes = styleClasses();

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setEmailAddress((event.target as HTMLInputElement).value);
    };
    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setPassword((event.target as HTMLInputElement).value);
    };
    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.handleSubmit();
    }
    return  (
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
                    <PrimaryButton type="submit" disabled={!props.emailAddress || !props.password || props.isLoading}>
                        Login
                    </PrimaryButton>
                </Grid>
            </Grid>
        </form>
    );
}

type TwoFactorAuthenticationFormProps = {
    twoFactorAuthenticationCode: string;
    setTwoFactorAuthenticationCode: (v: string) => void;

    handleSubmit: () => void;
    isLoading: boolean;
}

const TwoFactorAuthenticationForm = (props: TwoFactorAuthenticationFormProps) => {
    const classes = styleClasses();

    const handleTwoFactorAuthenticationCodeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setTwoFactorAuthenticationCode((event.target as HTMLInputElement).value);
    };
    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.handleSubmit();
    }
    return (
        <form onSubmit={handleSubmit} noValidate autoComplete="off">
            <Grid container className={classes.formGridContainer}>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6} className={classes.formGridItem}>
                    <PrimaryTextField
                        className={classes.textField}
                        id="two-factor-code"
                        label="Two Factor Authentication Code"
                        variant="outlined"
                        defaultValue={props.twoFactorAuthenticationCode}
                        onChange={handleTwoFactorAuthenticationCodeChange} />
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={3} md={2} className={classes.formGridItem}>
                    <PrimaryButton type="submit" disabled={!props.twoFactorAuthenticationCode || props.isLoading}>
                        Validate
                    </PrimaryButton>
                </Grid>
            </Grid>
        </form>
    );
}

export default LoginPage;
