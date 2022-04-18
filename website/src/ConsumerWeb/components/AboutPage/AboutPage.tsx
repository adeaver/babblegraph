import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import Color from 'common/styles/colors';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PhotoKey } from 'common/data/photos/Photos';
import Page from 'common/components/Page/Page';
import { Heading1, Heading2, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph from 'common/typography/Paragraph';
import Link, { LinkTarget } from 'common/components/Link/Link';

const styleClasses = makeStyles({
    faqRoot: {
        padding: '10px',
    },
});

const AboutPage = () => {
    const classes = styleClasses();
    return (
        <Page withNavbar withBackground={PhotoKey.Cartagena}>
            <CenteredComponent useLargeVersion>
                <DisplayCard>
                    <Heading1 color={TypographyColor.Primary}>
                        You didn’t spend all that time learning Spanish just to forget it
                    </Heading1>
                    <Divider />
                    <Heading3
                        align={Alignment.Left}
                        color={TypographyColor.Primary}>
                        Babblegraph was built to make it effortless to work Spanish into your routine
                    </Heading3>
                    <Paragraph align={Alignment.Left}>
                        It’s not a secret that the best way to maintain your Spanish is to use it. There’s a lot of great ways to use Spanish, like watching TV shows or movies or reading books in Spanish. The problem is that it can be hard to work that into your daily routine. Enter Babblegraph. Babblegraph started as a way for Spanish practice to come to you. This helps you maintain a routine of using Spanish.
                    </Paragraph>
                    <Heading3
                        align={Alignment.Left}
                        color={TypographyColor.Primary}>
                        Learning new words can be difficult as you get more advanced
                    </Heading3>
                    <Paragraph align={Alignment.Left}>
                        When you first started learning Spanish, words like "correr" or "hacer" came up frequently, so it was easy to remember them. As you get more advanced, unfamiliar words get harder to find. How many times a day do you use a word like "respaldar"? Still, it’s good to continue to build up your vocabulary. With Babblegraph’s vocabulary list, you can add words and phrases into your list to ensure you see them again.
                    </Paragraph>
                    <Heading3 color={TypographyColor.Primary}>
                        Try it free for 30 days.
                    </Heading3>
                    <Paragraph>
                        And then $29 for a whole year.
                    </Paragraph>
                    <Link href="/" target={LinkTarget.Self}>
                        Return to home page
                    </Link>
                    <Heading3 align={Alignment.Center} color={TypographyColor.Primary}>
                        Frequently asked questions
                    </Heading3>
                    <Grid container>
                        <FAQItem
                            questionText="I’m learning Spanish from a specific country. Can I only receive articles from that country?"
                            answerText="As of now, you cannot. However, you can choose to select current events from a specific country. All current events articles are from news sources originating in that country. So you’ll only get Mexican news sources if you’re interested in Mexican news." />
                        <FAQItem
                            questionText="Do I need a credit card or an account to sign up?"
                            answerText="Nope! Babblegraph is free to use for 30 days. You won’t need to put in your credit card information or create an account until the end of your free trial." />
                        <FAQItem
                            questionText="Why is there no free option?"
                            answerText="Babblegraph is built a team of one. My name is Andrew, and I’m a software engineer from the United States. Making Babblegraph a paid service means that my primary incentive is to spend my time building out a great product instead of trying to get you to click on ads." />
                        <FAQItem
                            questionText="Why is there no monthly option?"
                            answerText="Basically anytime you use your credit card, the company, website or person who charges you has to pay a fee. These fees really add up for small purchases. I felt the cheapest option for me and for customers is to have a yearly option only." />
                        <FAQItem
                            questionText="I’m skeptical of technology products. What kind of information are you tracking about me? Do you sell my data?"
                            answerText="Babblegraph keeps track of the information that you submit (i.e. your interests, your email, and words that you’re learning, etc). Babblegraph also tracks whether or not emails sent are opened and whether or not links in the newsletter are clicked. Babblegraph doesn’t sell your data." />
                    </Grid>
                </DisplayCard>
            </CenteredComponent>
        </Page>
    );
}

type FAQItemProps = {
    questionText: string;
    answerText: string;
};

const FAQItem = (props: FAQItemProps) => {
    const classes = styleClasses();
    return (
        <Grid className={classes.faqRoot} item xs={12} md={6}>
            <Heading3 color={TypographyColor.Primary} align={Alignment.Left}>{props.questionText}</Heading3>
            <Divider />
            <Paragraph align={Alignment.Left}>{props.answerText}</Paragraph>
        </Grid>
    );
}

export default AboutPage;
