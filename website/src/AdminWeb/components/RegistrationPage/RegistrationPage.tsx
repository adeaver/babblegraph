import React, { useState, useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import Page from 'common/components/Page/Page';
import { TypographyColor } from 'common/typography/common';
import { Heading1 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
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

type Params = {
    token: string
}

type RegistrationPageProps = RouteComponentProps<Params>

const RegistrationPage = (props: RegistrationPageProps) => {
    const { token } = props.match.params;

    const [ shouldShowRegistrationForm, setShouldShowRegistrationForm ] = useState<boolean>(true);

    const [ emailAddress, setEmailAddress ] = useState<string>("");
    const [ password, setPassword ] = useState<string>("");
    const [ confirmPassword, setConfirmPassword ] = useState<string>("");

    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const handleSubmit = () => {
    }

    const classes = styleClasses();
    return (
        <Page>
            <Grid container>
                <Grid xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid xs={12} md={6}>
                    <Card className={classes.displayCard}>
                        <Heading1 color={TypographyColor.Primary}>
                            babblegraph
                        </Heading1>
                        {
                            shouldShowRegistrationForm ? (
                                <RegistrationForm
                                    emailAddress={emailAddress}
                                    setEmailAddress={setEmailAddress}
                                    password={password}
                                    setPassword={setPassword}
                                    confirmPassword={confirmPassword}
                                    setConfirmPassword={setConfirmPassword}
                                    handleSubmit={handleSubmit}
                                    isLoading={isLoading} />
                            ) : (
                                <p>Two factor code</p>
                            )
                        }
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

type RegistrationFormProps = {
    emailAddress: string;
    setEmailAddress: (v: string) => void;

    password: string;
    setPassword: (v: string) => void;

    confirmPassword: string;
    setConfirmPassword: (v: string) => void;

    handleSubmit: () => void;
    isLoading: boolean;
}

const RegistrationForm = (props: RegistrationFormProps) => {
    const classes = styleClasses();

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setEmailAddress((event.target as HTMLInputElement).value);
    };
    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setPassword((event.target as HTMLInputElement).value);
    };
    const handleConfirmPasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setConfirmPassword((event.target as HTMLInputElement).value);
    };
    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.handleSubmit();
    }
    return  (
        <form onSubmit={handleSubmit} noValidate autoComplete="off">
            <Paragraph>
                Password must contain at 3 of the following: number, capital letter, lowercase letter, and symbol.
            </Paragraph>
            <Grid container className={classes.formGridContainer}>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6} className={classes.formGridItem}>
                    <PrimaryTextField
                        className={classes.textField}
                        id="email"
                        label="Email Address"
                        variant="outlined"
                        defaultValue={props.emailAddress}
                        onChange={handleEmailAddressChange} />
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6} className={classes.formGridItem}>
                    <PrimaryTextField
                        className={classes.textField}
                        id="password"
                        label="Password"
                        type="password"
                        variant="outlined"
                        defaultValue={props.password}
                        onChange={handlePasswordChange} />
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6} className={classes.formGridItem}>
                    <PrimaryTextField
                        className={classes.textField}
                        id="confirm-password"
                        label="Confirm Password"
                        type="password"
                        variant="outlined"
                        defaultValue={props.confirmPassword}
                        onChange={handleConfirmPasswordChange} />
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={3} md={2} className={classes.formGridItem}>
                    <PrimaryButton type="submit" disabled={!props.emailAddress || !props.password || props.isLoading}>
                       Register
                    </PrimaryButton>
                </Grid>
            </Grid>
        </form>
    );
}

export default RegistrationPage;
