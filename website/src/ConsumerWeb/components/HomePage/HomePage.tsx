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

import {
    LookupPromotionCodeResponse,
    lookupPromotionCode
} from 'ConsumerWeb/api/billing/billing';
import { PromotionCode } from 'common/api/billing/billing';

import { asRoundedFixedDecimal } from 'util/string/NumberString';

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
    promotionContainer: {
        width: '100%',
        backgroundColor: Color.Primary,
        padding: '2px 0',
        borderRadius: '5px',
    },
});

const HomePage = () => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ hadSuccess, setHadSuccess ] = useState<boolean>(false);
    const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);

    const [ promotionCode, setPromotionCode ] = useState<PromotionCode>(null);

    useEffect(() => {
        lookupPromotionCode({},
        (resp: LookupPromotionCodeResponse) => {
            setPromotionCode(resp.promotionCode);
        },
        (err: Error) => {});
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
                        Thank you for your interest!
                    </Heading1>
                    <Divider />
                    <Paragraph align={Alignment.Left}>
                        Babblegraph started as a side project of mine when I wanted a better way to incorporate Spanish into my daily routine. About a year and a half ago, I decided to see how much I could grow it! Since then it’s been a whirlwind adventure that has seen thousands of people sign up and give it a try. I’m truly humbled by how much interest it has gotten. I truly could have never foreseen it.
                    </Paragraph>
                    <Paragraph align={Alignment.Left}>
                        I really loved building Babblegraph, and for months, I really strived to turn Babblegraph from a side project into a full-fledged company. Along the way, I learned a lot about what makes a great product and what makes a great company.
                    </Paragraph>
                    <Paragraph align={Alignment.Left}>
                        Unfortunately, I also couldn't forsee what the challenges of that transition were. As a solo developer, it really put a lot of stress on me. In the end, it really ended up draining the love I had for building Babblegraph when it first got started.
                    </Paragraph>
                    <Paragraph align={Alignment.Left}>
                         As with any product or service that you use, you should really expect that the person providing it to you has their heart in it, and for Babblegraph, that just isn't as true as it was a year and a half ago. Therefore, I’ve made the tough decision to let it go and stop working on it.
                    </Paragraph>
                    <Paragraph align={Alignment.Left}>
                         At the end of the month, we are shutting down the site, so we’re not currently accepting new signups. But once again, thank you so much for your interest!
                    </Paragraph>
                    <Paragraph align={Alignment.Right}>
                        - Andrew
                    </Paragraph>
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
