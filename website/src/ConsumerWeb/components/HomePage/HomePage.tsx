import React, { useEffect, useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';

import MailOutlineIcon from '@material-ui/icons/MailOutline';
import LibraryAddCheckIcon from '@material-ui/icons/LibraryAddCheck';
import AutorenewIcon from '@material-ui/icons/Autorenew';

import Color from 'common/styles/colors';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { PhotoKey } from 'common/data/photos/Photos';
import Link from 'common/components/Link/Link';
import { withCaptchaToken, loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';

import SignupForm from 'ConsumerWeb/components/common/SignupForm/SignupForm';

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
        <Page withNavbar withBackground={PhotoKey.Seville}>
            <CenteredComponent useLargeVersion>
                <DisplayCard>
                    <Heading1 color={TypographyColor.Primary}>
                        {
                            hadSuccess ? (
                                'Success! You’re one step closer to receiving interesting Spanish language content'
                            ) : (
                                'Don’t lose your Spanish'
                            )
                        }
                    </Heading1>
                    <Divider />
                    {
                        !hadSuccess && (
                            <div>
                                <Paragraph>
                                    Babblegraph picks up where your Spanish class left off by sending you an email newsletter with real articles and podcasts from the Spanish-speaking world. You can even customize the content to make keeping up with your Spanish skills more engaging!
                                </Paragraph>
                                <Paragraph>
                                    Try it risk-free for 30 days. No account or credit card required!
                                </Paragraph>
                            </div>
                        )
                    }
                    <SignupForm
                        disabled={isLoading}
                        shouldShowVerificationForm={hadSuccess}
                        setIsLoading={setIsLoading}
                        onSuccess={handleSuccess} /> {
                        hadSuccess ? (
                            <PostVerificationContent handleReturnHome={() => setHadSuccess(false)} />
                        ) : (
                            <PostInitialContent />
                        )
                    }
                </DisplayCard>
            </CenteredComponent>
        </Page>
    );
}

const PostInitialContent = () => {
    const classes = styleClasses();
    return (
        <div>
            <Paragraph size={Size.Small}>
                By signing up, you’re acknowledging that you’ve read and agreed to our privacy policy.
            </Paragraph>
            <Link href="/privacy-policy">
                Read our privacy policy
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
                    <Paragraph align={Alignment.Left}>
                        Sign up to receive a newsletter containing articles from trusted Spanish-language news sources and podcasts on your schedule, as frequently as every day or as a little as every week.
                    </Paragraph>
                </Grid>
            </Grid>
            <Grid container>
                <Grid className={classes.infoIconContainer} item xs={3} md={2}>
                    <LibraryAddCheckIcon className={classes.infoIcon} />
                </Grid>
                <Grid item xs={9} md={10}>
                    <Paragraph align={Alignment.Left}>
                        Select topics that you’re interested in to keep the articles you receive fun and engaging
                    </Paragraph>
                </Grid>
            </Grid>
            <Grid container>
                <Grid className={classes.infoIconContainer} item xs={3} md={2}>
                    <AutorenewIcon className={classes.infoIcon} />
                </Grid>
                <Grid item xs={9} md={10}>
                    <Paragraph align={Alignment.Left}>
                        Add new vocabulary words to receive more interesting articles that use those words so that you can reinforce them
                    </Paragraph>
                </Grid>
            </Grid>
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
