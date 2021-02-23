import React, { useEffect, useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CheckCircleOutlineIcon from '@material-ui/icons/CheckCircleOutline';
import Divider from '@material-ui/core/Divider';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import Grid from '@material-ui/core/Grid';

import MailOutlineIcon from '@material-ui/icons/MailOutline';
import LibraryAddCheckIcon from '@material-ui/icons/LibraryAddCheck';
import AutorenewIcon from '@material-ui/icons/Autorenew';

import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { TypographyColor } from 'common/typography/common';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PhotoKey } from 'common/data/photos/Photos';
import Link from 'common/components/Link/Link';
import { withCaptchaToken, loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';

import {
    SignupUserResponse,
    SignupErrorMessage,
    signupUser,
} from 'api/user/signup';
import {
    SetPageLoadEventResponse,
    setPageLoadEvent,
} from 'api/utm/utm';

const styleClasses = makeStyles({
    displayCard: {
        padding: '20px',
        marginTop: '20px',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    emailField: {
        width: '100%',
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
    iconContainer: {
        alignSelf: 'center',
    },
    confirmationIcon: {
        color: Color.Confirmation,
    },
    warningIcon: {
        color: Color.Warning,
    },
    infoIcon: {
        color: Color.Primary,
        display: 'block',
        margin: '0 auto',
        fontSize: '48px',
    },
    infoIconContainer: {
        alignSelf: 'center',
    },
    mainPageLinkContainer: {
        cursor: 'pointer',
    },
});

const errorMessages = {
    // TODO: think about these
    [SignupErrorMessage.InvalidEmailAddress]: "Hmm, the email address you gave doesn’t appear to be valid. Check to make sure that you spelled everything right.",
    [SignupErrorMessage.RateLimited]: "It looks like we’re having some trouble reaching you. Contact our support so we can get you on the list!",
    [SignupErrorMessage.IncorrectStatus]: "It looks like you’re already signed up for Babblegraph!",
    "default": "Something went wrong. Contact our support so we can get you on the list!"

}

const HomePage = () => {
    const [ emailAddress, setEmailAddress ] = useState<string>('');
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ errorMessage, setErrorMessage ] = useState<string | null>(null);
    const [ hadSuccess, setHadSuccess ] = useState<boolean>(false);
    const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);

    const handleSubmit = () => {
        setIsLoading(true);
        withCaptchaToken("signup", (token: string) => {
            signupUser({
                emailAddress: emailAddress,
                captchaToken: token,
            },
            (resp: SignupUserResponse) => {
                setIsLoading(false);
                if (!!resp.errorMessage) {
                    setErrorMessage(errorMessages[resp.errorMessage] || errorMessages["default"]);
                    setHadSuccess(false);
                } else if (resp.success) {
                    setErrorMessage(null);
                    setHadSuccess(true);
                } else {
                    setErrorMessage(errorMessages["default"]);
                    setHadSuccess(false);
                }
            },
            (e: Error) => {
                setIsLoading(false);
                setErrorMessage(errorMessages["default"]);
                setHadSuccess(false);
            });
        });
    }
    useEffect(() => {
        loadCaptchaScript();
        setPageLoadEvent({},
            (resp: SetPageLoadEventResponse) => {},
            (e: Error) => {});
        setHasLoadedCaptcha(true);
    }, []);

    const classes = styleClasses();
    let body;
    if (isLoading) {
        body = (<LoadingSpinner />);
    } else if (hadSuccess) {
        body = (
            <SuccessConfirmation
                emailAddress={emailAddress}
                handleResendVerificationEmail={handleSubmit}
                handleReturnHome={() => { setHadSuccess(false)}} />
        );
    } else {
        body = (
            <SignupForm
                emailAddress={emailAddress}
                errorMessage={errorMessage}
                canSubmit={hasLoadedCaptcha}
                handleSubmit={handleSubmit}
                handleEmailAddressChange={setEmailAddress} />
        );
    }
    return (
        <Page withBackground={PhotoKey.Seville}>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        {body}
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

type SignupFormProps = {
    emailAddress: string;
    canSubmit: boolean;
    errorMessage: string | null;

    handleSubmit: () => void;
    handleEmailAddressChange: (emailAddress: string) => void;
}

const SignupForm = (props: SignupFormProps) => {
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleEmailAddressChange((event.target as HTMLInputElement).value);
    };
    const preventDefault = (event: React.SyntheticEvent) => event.preventDefault();

    const classes = styleClasses();
    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                Don’t lose your Spanish
            </Heading1>
            <Divider />
            <Paragraph>
                Babblegraph picks up where your Spanish class left off by sending you a daily email with real articles from the Spanish-speaking world. You can even customize the content to make keeping up with your Spanish skills more engaging!
            </Paragraph>
            <Paragraph>
                It’s completely free and you can unsubscribe anytime you’d like.
            </Paragraph>
            <form className={classes.confirmationForm} noValidate autoComplete="off">
                <Grid container>
                    <Grid item xs={9} md={10}>
                        <PrimaryTextField
                            id="email"
                            className={classes.emailField}
                            label="Email Address"
                            variant="outlined"
                            defaultValue={props.emailAddress}
                            onChange={handleEmailAddressChange} />
                    </Grid>
                    <Grid item xs={3} md={2} className={classes.submitButtonContainer}>
                        <PrimaryButton onClick={props.handleSubmit} disabled={!props.emailAddress && props.canSubmit}>
                            Try it!
                        </PrimaryButton>
                    </Grid>
                </Grid>
            </form>
            {
                !!props.errorMessage && (
                    <Grid container>
                        <Grid className={classes.iconContainer} item xs={1}>
                            <ErrorOutlineIcon className={classes.warningIcon} />
                        </Grid>
                        <Grid item xs={11}>
                            <Paragraph size={Size.Small} color={TypographyColor.Warning}>
                                {props.errorMessage}
                            </Paragraph>
                        </Grid>
                    </Grid>
                )
            }
            <Link href="/privacy-policy">
                View our Privacy Policy
            </Link>
            <Divider />
            <Heading3 color={TypographyColor.Primary}>
                How it works
            </Heading3>
            <Grid container>
                <Grid className={classes.infoIconContainer} item xs={3} md={2}>
                    <MailOutlineIcon className={classes.infoIcon} />
                </Grid>
                <Grid item xs={9} md={10}>
                    <Paragraph>
                        Sign up to receive an email every day from a trusted Spanish-language news source
                    </Paragraph>
                </Grid>
            </Grid>
            <Grid container>
                <Grid className={classes.infoIconContainer} item xs={3} md={2}>
                    <LibraryAddCheckIcon className={classes.infoIcon} />
                </Grid>
                <Grid item xs={9} md={10}>
                    <Paragraph>
                        Select topics that you’re interested in to keep your articles fun and engaging.
                    </Paragraph>
                </Grid>
            </Grid>
            <Grid container>
                <Grid className={classes.infoIconContainer} item xs={3} md={2}>
                    <AutorenewIcon className={classes.infoIcon} />
                </Grid>
                <Grid item xs={9} md={10}>
                    <Paragraph>
                        Track words that you’re learning to receive more interesting articles that use those words in order to reinforce them.
                    </Paragraph>
                </Grid>
            </Grid>
        </div>
    )
}

type SuccessConfirmationProps = {
    emailAddress: string;

    handleResendVerificationEmail: () => void;
    handleReturnHome: () => void;
};

const SuccessConfirmation = (props: SuccessConfirmationProps) => {
    const classes = styleClasses();
    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>Almost done! Just one last step.</Heading1>
            <Paragraph>
                Check your email for a verification email from babblegraph@gmail.com. We sent it to {props.emailAddress}.
            </Paragraph>
            <Paragraph>
                You’ll need to click the button in the verification email that was just sent to you in order to start receiving emails from Babblegraph.
            </Paragraph>
            <Paragraph>
                It can take up to 5 minutes for the email to make its way to your inbox.
            </Paragraph>
            <Grid container>
                <Grid item xs={3} md={4}>
                    &nbsp;
                </Grid>
                <Grid item xs={6} md={4}>
                    <PrimaryButton onClick={props.handleResendVerificationEmail}>
                        Resend the verification email
                    </PrimaryButton>
                </Grid>
            </Grid>
            <Grid container>
                <Grid item xs={3} md={4}>
                    &nbsp;
                </Grid>
                <Grid className={classes.mainPageLinkContainer} item xs={6} md={4} onClick={props.handleReturnHome}>
                    <Paragraph color={TypographyColor.LinkBlue}>
                        Return to main page
                    </Paragraph>
                </Grid>
            </Grid>
        </div>
    )
}

export default HomePage;
