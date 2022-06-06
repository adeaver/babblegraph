import React, { useState } from 'react'

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import Alert from 'common/components/Alert/Alert';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import { LoginRedirectKey } from 'ConsumerWeb/api/routes/consts';
import {
    loginUser,
    LoginUserResponse,
    LoginError,
} from 'ConsumerWeb/api/useraccounts/useraccounts';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
    },
    formContainer: {
        padding: '5px',
    },
});

const loginErrorMessages = {
    [LoginError.InvalidCredentials]: "The username or password entered is not correct",
    "default": "Something went wrong, try again later.",
}

type LoginFormProps = {
    redirectKey?: LoginRedirectKey;
    onLoginSuccess: (location: string) => void;
}

const LoginForm = (props: LoginFormProps) => {
    const [ emailAddress, setEmailAddress ] = useState<string>(null);
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmailAddress((event.target as HTMLInputElement).value);
    };

    const [ password, setPassword ] = useState<string>(null);
    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPassword((event.target as HTMLInputElement).value);
    }

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ errorMessage, setErrorMessage ] = useState<string | null>(null);

    const handleSubmit = () => {
        setIsLoading(true);
        loginUser({
            emailAddress: emailAddress,
            password: password,
            redirectKey: props.redirectKey || undefined,
        },
        (resp: LoginUserResponse) => {
            setIsLoading(false);
            if (!!resp.loginError) {
                setErrorMessage(loginErrorMessages[resp.loginError] || loginErrorMessages["default"]);
            } else if (!!resp.location) {
                props.onLoginSuccess(resp.location);
            } else {
                setErrorMessage(loginErrorMessages["default"]);
            }
        },
        (e: Error) => {
            setIsLoading(false);
            setErrorMessage(loginErrorMessages["default"]);
        });
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid container>
                <Grid className={classes.formContainer} item xs={12} md={6}>
                    <PrimaryTextField
                        className={classes.formComponent}
                        id="email"
                        label="Email Address"
                        variant="outlined"
                        defaultValue={emailAddress}
                        onChange={handleEmailAddressChange} />
                </Grid>
                <Grid className={classes.formContainer} item xs={12} md={6}>
                    <PrimaryTextField
                        className={classes.formComponent}
                        id="password"
                        label="Password"
                        type="password"
                        variant="outlined"
                        defaultValue={password}
                        onChange={handlePasswordChange} />

                </Grid>
                <Grid className={classes.formContainer} item xs={12}>
                    <CenteredComponent>
                        <PrimaryButton
                            className={classes.formComponent}
                            type="submit"
                            disabled={!emailAddress || !password}>
                            Login
                        </PrimaryButton>
                    </CenteredComponent>
                </Grid>
            </Grid>
            <Snackbar open={!!errorMessage} autoHideDuration={6000} onClose={() => setErrorMessage(null)}>
                <Alert severity="error">{errorMessage}</Alert>
            </Snackbar>
        </Form>
    )
}

export default LoginForm;
