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
import Link, { LinkTarget } from 'common/components/Link/Link';
import { withCaptchaToken, loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';
import Form from 'common/components/Form/Form';

import SignupForm from 'ConsumerWeb/components/common/SignupForm/SignupForm';

import {
    SignupUserResponse,
    SignupErrorMessage,
    signupUser,
} from 'ConsumerWeb/api/user/signup';
import {
    SetPageLoadEventResponse,
    setPageLoadEvent,
} from 'ConsumerWeb/api/utm/utm';

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
    verificationButton: {
        width: '100%',
    },
});

const errorMessages = {
    // TODO: think about these
    [SignupErrorMessage.InvalidEmailAddress]: "Hmm, the email address you gave doesn’t appear to be valid. Check to make sure that you spelled everything right.",
    [SignupErrorMessage.RateLimited]: "It looks like we’re having some trouble reaching you. Contact our support so we can get you on the list!",
    [SignupErrorMessage.IncorrectStatus]: "It looks like you’re already signed up for Babblegraph!",
    [SignupErrorMessage.LowScore]: "We’re having some trouble verifying your request. Contact us at hello@babblegraph.com to finish signing up.",
    "default": "Something went wrong. Contact our support so we can get you on the list!"

}

const HomePage = () => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ hadSuccess, setHadSuccess ] = useState<boolean>(false);
    const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);

    useEffect(() => {
        loadCaptchaScript();
        setPageLoadEvent({},
            (resp: SetPageLoadEventResponse) => {},
            (e: Error) => {});
        setHasLoadedCaptcha(true);
    }, []);

    const handleSuccess = (emailAddress: string) => {
        setHadSuccess(true);
    }

    const classes = styleClasses();
    return (
        <Page withBackground={PhotoKey.Seville}>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <Heading1 color={TypographyColor.Primary}>
                            {
                                !hadSuccess ? (
                                    "Don’t lose your Spanish"
                                ) : (
                                    "Almost done! Just one last step."
                                )
                            }
                        </Heading1>
                        <Divider />
                        {
                            !hadSuccess && (
                                <div>
                                    <Paragraph>
                                        Babblegraph picks up where your Spanish class left off by sending you a daily email with real articles from the Spanish-speaking world. You can even customize the content to make keeping up with your Spanish skills more engaging!
                                    </Paragraph>
                                    <Paragraph>
                                        It’s completely free and you can unsubscribe anytime you’d like.
                                    </Paragraph>
                                </div>
                            )
                        }
                        <SignupForm
                            disabled={isLoading || !hasLoadedCaptcha}
                            setIsLoading={setIsLoading}
                            onSuccess={handleSuccess} />
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                !hadSuccess ? (
                                    <PostInitialContent />
                                ) : (
                                    <PostVerificationContent
                                        handleReturnHome={() => setHadSuccess(false)} />
                                )
                            )
                        }
                    </Card>
                    <Card className={classes.displayCard} variant='outlined'>
                        <Paragraph>
                            Have an account?
                        </Paragraph>
                        <Link href="/login" target={LinkTarget.Self}>
                            Click here to login
                        </Link>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

const PostInitialContent = () => {
    const classes = styleClasses();
    return (
        <div>
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
                        Sign up to receive an email every day containing articles from trusted Spanish-language news sources
                    </Paragraph>
                </Grid>
            </Grid>
            <Grid container>
                <Grid className={classes.infoIconContainer} item xs={3} md={2}>
                    <LibraryAddCheckIcon className={classes.infoIcon} />
                </Grid>
                <Grid item xs={9} md={10}>
                    <Paragraph>
                        Select topics that you’re interested in to keep the articles you receive fun and engaging
                    </Paragraph>
                </Grid>
            </Grid>
            <Grid container>
                <Grid className={classes.infoIconContainer} item xs={3} md={2}>
                    <AutorenewIcon className={classes.infoIcon} />
                </Grid>
                <Grid item xs={9} md={10}>
                    <Paragraph>
                        Track words that you’re learning to receive more interesting articles that use those words so that you can reinforce them
                    </Paragraph>
                </Grid>
            </Grid>
            <Link href="/about">
                Learn more about Babblegraph
            </Link>
        </div>
    );
}

type PostVerificationContentProps = {
    handleReturnHome: () => void;
}

const PostVerificationContent = (props: PostVerificationContentProps) => {
    const classes = styleClasses();
    return (
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
    );
}

export default HomePage;
