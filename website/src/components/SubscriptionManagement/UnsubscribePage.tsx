import React, { useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import Card from '@material-ui/core/Card';
import CircularProgress from '@material-ui/core/CircularProgress';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import MuiAlert from '@material-ui/lab/Alert';
import Snackbar from '@material-ui/core/Snackbar';

import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import { UnsubscribeUser, UnsubscribeResponse } from 'api/user/unsubscribe';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
    },
    contentHeaderBackArrow: {
        alignSelf: 'center',
        cursor: 'pointer',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    emailField: {
        width: '100%',
    },
    formContainer: {
        padding: '10px 0',
    },
    loadingSpinner: {
        color: Color.Primary,
        display: 'block',
        margin: 'auto',
    },
});

type Params = {
    token: string
}

type UnsubscribePageProps = RouteComponentProps<Params>

const UnsubscribePage = (props: UnsubscribePageProps) => {
    const { token } = props.match.params;

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ error, setError ] = useState<Error | null>(null);
    const [ didUpdate, setDidUpdate ] = useState<boolean | null>(null);

    const handleSubmit = () => {
        UnsubscribeUser({
            Token: token,
            EmailAddress: emailAddress,
        },
        (resp: UnsubscribeResponse) => {
            setIsLoading(false);
            setDidUpdate(resp.Success);
        },
        (e: Error) => {
            setIsLoading(false);
            setError(e);
        });
    }

    const classes = styleClasses();
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <ContentHeader token={token} />
                        <Divider />
                        <Paragraph size={Size.Medium} align={Alignment.Left}>
                            We’re sorry to see you go! When you unsubscribe, you won’t receive any more daily emails or any other communication from us. But you can always come back by resubscribing on the homepage. To unsubscribe, just enter your email and click the unsubscribe button.
                        </Paragraph>
                        {
                            isLoading ? (
                                <LoadingScreen />
                            ) : (
                                <UnsubscribeForm
                                    emailAddress={emailAddress}
                                    handleEmailAddressChange={setEmailAddress}
                                    handleSubmit={handleSubmit} />
                            )
                        }
                    </Card>
                    <Snackbar open={!!error} autoHideDuration={6000}>
                        <Alert severity="error">Something went wrong processing your request.</Alert>
                    </Snackbar>
                    <Snackbar open={didUpdate} autoHideDuration={6000}>
                        <Alert severity="success">Successfully unsubscribed from Babblegraph.</Alert>
                    </Snackbar>
                </Grid>
            </Grid>
        </Page>
    );
}

const Alert = (props) => {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

type ContentHeaderProps = {
    token: string;
}

const ContentHeader = (props: ContentHeaderProps) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <Grid container>
            <Grid className={classes.contentHeaderBackArrow} onClick={() => history.push(`/manage/${props.token}`)} item xs={1}>
                <ArrowBackIcon color='action' />
            </Grid>
            <Grid item xs={11}>
                <Paragraph size={Size.Large} color={TypographyColor.Primary} align={Alignment.Left}>
                    Unsubscribe
                </Paragraph>
            </Grid>
        </Grid>
    );
}

type UnsubscribeFormProps = {
    emailAddress: string;
    handleEmailAddressChange: (v: string) => void;
    handleSubmit: () => void;
}

const UnsubscribeForm = (props: UnsubscribeFormProps) => {
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleEmailAddressChange((event.target as HTMLInputElement).value);
    };

    const classes = styleClasses();
    return (
        <form className={classes.formContainer} noValidate autoComplete="off">
            <Grid container>
                <Grid item xs={9} md={10}>
                    <PrimaryTextField
                        id="email"
                        className={classes.emailField}
                        label="Email Address"
                        variant="outlined"
                        onChange={handleEmailAddressChange} />
                </Grid>
                <Grid item xs={3} md={2} className={classes.submitButtonContainer}>
                    <PrimaryButton onClick={props.handleSubmit} disabled={!props.emailAddress}>
                        Submit
                    </PrimaryButton>
                </Grid>
            </Grid>
        </form>
    );
}

const LoadingScreen = () => {
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <CircularProgress className={classes.loadingSpinner} />
                <Paragraph size={Size.Medium} align={Alignment.Center}>
                    Loading, please wait.
                </Paragraph>
            </Grid>
        </Grid>
    )
}

export default UnsubscribePage;
