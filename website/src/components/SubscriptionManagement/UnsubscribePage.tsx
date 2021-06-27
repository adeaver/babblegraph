import React, { useEffect, useState } from 'react';
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
import Link from 'common/components/Link/Link';

import { UnsubscribeUser, UnsubscribeResponse } from 'api/user/unsubscribe';
import {
    getUserProfile,
    GetUserProfileResponse
} from 'api/useraccounts/useraccounts';

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

    const [ isUnsubscribeRequestLoading, setIsUnsubscribeRequestLoading ] = useState<boolean>(false);
    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ error, setError ] = useState<Error | null>(null);
    const [ didUpdate, setDidUpdate ] = useState<boolean | null>(null);
    const [ isUserProfileLoading, setIsUserProfileLoading ] = useState<boolean>(true);
    const [ hasUserProfile, setHasUserProfile ] = useState<boolean>(false);

    const handleSubmit = () => {
        setIsUserProfileLoading(true);
        UnsubscribeUser({
            Token: token,
            EmailAddress: emailAddress,
        },
        (resp: UnsubscribeResponse) => {
            setIsUnsubscribeRequestLoading(false);
            setDidUpdate(resp.Success);
        },
        (e: Error) => {
            setIsUnsubscribeRequestLoading(false);
            setError(e);
        });
    }

    useEffect(() => {
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setIsUserProfileLoading(false);
            setHasUserProfile(!!resp.subscriptionLevel);
        },
        (e: Error) => {
            setIsUserProfileLoading(false);
            setError(e);
        });
    }, []);

    const classes = styleClasses();
    const isLoading = isUnsubscribeRequestLoading || isUserProfileLoading;
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
                                    hasUserProfile={hasUserProfile}
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
    hasUserProfile: boolean;
    handleEmailAddressChange: (v: string) => void;
    handleSubmit: () => void;
}

const UnsubscribeForm = (props: UnsubscribeFormProps) => {
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleEmailAddressChange((event.target as HTMLInputElement).value);
    };
    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.handleSubmit();
    }

    const classes = styleClasses();
    return (
        <form onSubmit={handleSubmit} className={classes.formContainer} noValidate autoComplete="off">
            {
                props.hasUserProfile && (
                    <div>
                        <Paragraph>
                            IMPORTANT: unsubscribing does not automatically cancel your subscription. Please make sure that you don’t have any recurring payments setup by going to the BuyMeACoffee Page.
                        </Paragraph>
                        <Link href="https://buymeacoffee.com/babblegraph">
                            Click here for BuyMeACoffee
                        </Link>
                    </div>
                )
            }
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
                    <PrimaryButton type="submit" disabled={!props.emailAddress}>
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
