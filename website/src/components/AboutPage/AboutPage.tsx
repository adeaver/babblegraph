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

const styleClasses = makeStyles({
    contentCard: {
        padding: '20px',
    },
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
                            Babblegraph is a free daily practice app.<br />
                            Designed for intermediate and advanced learners.
                        </Heading1>
                        <Paragraph>
                            Incorporate Spanish practice into your daily routine effortlessly. Babblegraph sends you articles from trusted, Spanish-language news sources from Spain and Latin America.
                        </Paragraph>
                        <Divider />
                        <Heading2 color={TypographyColor.Primary}>
                            How it works
                        </Heading2>
                        <Paragraph align={Alignment.Left}>
                            Every day, you’ll receive an email from Babblegraph with a handful of articles from Spanish-language news sources, such as elmundo.es (Spain), elsiglo.com (Panama), elespectador.com (Colombia).
                        </Paragraph>
                        <Paragraph align={Alignment.Left}>
                            You can customize the content that you receive by selecting topics that you’re interested in. If you want to read more about science, sports, or even current events from a specific Latin American country or Spain, you can set that up.
                        </Paragraph>
                        <Paragraph align={Alignment.Left}>
                            If you learn a new word, use the word reinforcement feature to make sure you see that word more often. Using a word helps you remember what it means!
                        </Paragraph>
                        <Divider />
                        <Heading2 color={TypographyColor.Primary}>
                            FAQ
                        </Heading2>
                        <FAQItem
                            questionText="Who built this?"
                            answerText="Hello! I’m Andrew. I’m the one-person team who created Babblegraph as a side project. I’m a software engineer from beautiful Boston, Massachusetts. I found that there aren’t a lot of apps that help intermediate or advanced learners expand their vocabulary or keep up with Spanish. Most of the advice seems to be watch movies or read books or listen to podcasts (which is great advice), but I’m lazy. So I decided to build something to make it easier." />
                        <FAQItem
                            questionText="I’m skeptical of technology companies, what information do you track or store about me?"
                            answerText="Babblegraph tracks whether or not an email has been opened (but only a yes or no - if you open it 100 times, I have no way of knowing that) and what topics you’re interested in. I don’t sell that information and Babblegraph doesn’t use any AI to try figure out exactly who you are. Babblegraph also does not track what links you click or what you do once you’re on a news company’s website." />
                        <Divider />
                        <Heading2 color={TypographyColor.Primary}>
                            Contact
                        </Heading2>
                        <Paragraph align={Alignment.Left}>
                            If you have a question, comment, or feedback on Babblegraph; or if you’re affiliated with a blog or news source and you’d like to talk about your content (either adding it to Babblegraph or taking it down). Contact me at hello@babblegraph.com.
                        </Paragraph>
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
