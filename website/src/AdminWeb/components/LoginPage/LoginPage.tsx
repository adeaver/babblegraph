import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import Page from 'common/components/Page/Page';
import { TypographyColor } from 'common/typography/common';
import { Heading1 } from 'common/typography/Heading';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

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

    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ password, setPassword ] = useState<string | null>(null);

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmailAddress((event.target as HTMLInputElement).value);
    };
    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPassword((event.target as HTMLInputElement).value);
    };
    const handleSubmit = () => {
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
                        <form onSubmit={handleSubmit} noValidate autoComplete="off">
                            <Grid container className={classes.formGridContainer}>
                                <Grid item xs={12} md={5} className={classes.formGridItem}>
                                    <PrimaryTextField
                                        className={classes.textField}
                                        id="email"
                                        label="Email Address"
                                        variant="outlined"
                                        defaultValue={emailAddress}
                                        onChange={handleEmailAddressChange} />
                                </Grid>
                                <Grid item xs={12} md={5} className={classes.formGridItem}>
                                    <PrimaryTextField
                                        className={classes.textField}
                                        id="password"
                                        label="Password"
                                        type="password"
                                        variant="outlined"
                                        defaultValue={password}
                                        onChange={handlePasswordChange} />
                                </Grid>
                                <Grid item xs={3} md={2} className={classes.formGridItem}>
                                    <PrimaryButton type="submit" disabled={!emailAddress || !password}>
                                        Login
                                    </PrimaryButton>
                                </Grid>
                            </Grid>
                        </form>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    )
}

export default LoginPage;
