import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import { Heading1, Heading2, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph from 'common/typography/Paragraph';
import Link, { LinkTarget } from 'common/components/Link/Link';

const styleClasses = makeStyles({
    contentCard: {
        padding: '20px',
        margin: '10px 0',
    },
    topicsExampleImage: {
        minWidth: '100%',
        minHeight: '350px',
        backgroundImage: 'url("https://static.babblegraph.com/assets/topics-example.png")',
        backgroundSize: 'contain',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
    },
    reinforcementExampleImage: {
        width: '100%',
        minHeight: '350px',
        backgroundImage: 'url("https://static.babblegraph.com/assets/reinforcement-example.png")',
        backgroundSize: 'contain',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
    }
});

const AboutPage = () => {
    const classes = styleClasses();
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={2}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={8}>
                    <Card className={classes.contentCard}>
                        <Heading1 color={TypographyColor.Primary}>
                            Free daily Spanish practice for intermediate and advanced learners.
                        </Heading1>
                        <Paragraph>
                            Incorporate Spanish practice into your daily routine effortlessly.  Babblegraph sends you news articles from trusted, Spanish-language news sources from Spain and Latin America.
                        </Paragraph>
                    </Card>
                    <Card className={classes.contentCard}>
                        <Grid container>
                            <Grid item xs={12} md={6}>
                                <Heading2 color={TypographyColor.Primary}>
                                    Keeping up with your practice routine is easier when it's fun and engaging
                                </Heading2>
                                <Paragraph>
                                    Babblegraph helps with this by allowing you to select categories for interesting topics. Read anything from articles about film to current events from individual Spanish-speaking countries.
                                </Paragraph>
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <div className={classes.topicsExampleImage} />
                            </Grid>
                        </Grid>
                    </Card>
                    <Card className={classes.contentCard}>
                        <Grid container>
                            <Grid item xs={12} md={6}>
                                <Heading2 color={TypographyColor.Primary}>
                                    Reinforcing new vocabulary depends on seeing words you’re learning multiple times.
                                </Heading2>
                                <Paragraph>
                                    Babblegraph helps with this by sending you articles that contain words on your tracking list. Use Babblegraph’s tracking tool to lookup new words, get their definitions, and add them to your tracking list.
                                </Paragraph>
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <div className={classes.reinforcementExampleImage} />
                            </Grid>
                        </Grid>
                    </Card>
                    <Card className={classes.contentCard}>
                        <Heading2 color={TypographyColor.Primary}>
                            Frequently Asked Questions
                        </Heading2>
                        <FAQItem
                            questionText="I’m learning Spanish from a specific country. Can I only receive articles from that country?"
                            answerText="Currently, no. However, you can choose to select current events from a specific country. All current events articles are from news sources originating in that country. So you’ll only get Mexican news sources if you’re interested in Mexican news." />
                        <FAQItem
                            questionText="I’m skeptical of technology products. What kind of information are you tracking about me?"
                            answerText="Babblegraph keeps track of the information that you submit (i.e. your interests, your email, and words that you’re learning, etc). Babblegraph also tracks whether or not emails sent are opened and whether or not links in the daily email are clicked." />
                        <FAQItem
                            questionText="Why does Babblegraph need to keep track of whether or not emails are opened?"
                            answerText="I keep track of whether or not emails are opened to know whether or not changes to Babblegraph drastically affect open rate. This can be a feature that people either really like or really dislike that causes the open rate to go up or go down. It can also be a bug that causes emails to be undeliverable or something." />
                        <FAQItem
                            questionText="Why does Babblegraph need to keep track of whether or not I’ve clicked links in the email?"
                            answerText="Many of the websites that Babblegraph links to have a monthly limit of the number of articles that you can read for free (this is called a metered paywall). I keep track of the links that are clicked on so that the daily email doesn’t include articles from a news source with a metered paywall once the monthly quota is reached. There is not currently an option to indicate if you have a digital subscription to a specific news source that will cause Babblegraph to ignore the monthly limit of articles - but you should consider supporting news organizations that you find interesting!" />
                        <FAQItem
                            questionText="Does Babblegraph use the links I’ve clicked on to send more personalized articles to me?"
                            answerText="Absolutely not - and it will not as long as I am running it. Hyper-personalized content on the internet can lead to polarization. For instance, if you bias towards either conservative or liberal news, personalization will mean that the other side will be hidden from you - I tend to feel that this isn’t good." />
                        <FAQItem
                            questionText="You keep saying “I”... Who built this?"
                            answerText="Hello! My name is Andrew. I’m the sole creator behind Babblegraph and the only person working on it. I built Babblegraph because I found that it’s hard to use new words that you learn as an intermediate or advanced speaker. I read a lot of content in Spanish, but there’s never any guarantee that I’ll see a new word again." />
                        <FAQItem
                            questionText="Do you sell my data?"
                            answerText="Nope, I don’t sell anything as of now. Babblegraph will never sell your data." />
                    </Card>
                    <Card className={classes.contentCard}>
                        <Heading2 color={TypographyColor.Primary}>
                            Contact Information
                        </Heading2>
                        <Paragraph>
                            Feel free to reach out to me by email at hello@babblegraph.com. I love receiving feedback or compliments. I can also help remove or add your content to Babblegraph.
                        </Paragraph>
                        <Link href="/" target={LinkTarget.Self}>
                            Return to Main Page
                        </Link>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

type FAQItemProps = {
    questionText: string;
    answerText: string;
};

const FAQItem = (props: FAQItemProps) => {
    return (
        <div>
            <Heading3 align={Alignment.Left}>{props.questionText}</Heading3>
            <Paragraph align={Alignment.Left}>{props.answerText}</Paragraph>
        </div>
    );
}

export default AboutPage;
