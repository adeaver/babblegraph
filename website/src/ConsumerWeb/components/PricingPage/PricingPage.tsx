import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';

import LibraryAddCheckIcon from '@material-ui/icons/LibraryAddCheck';
import FavoriteIcon from '@material-ui/icons/Favorite';
import DynamicFeedIcon from '@material-ui/icons/DynamicFeed';

import Color from 'common/styles/colors';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Page from 'common/components/Page/Page';
import { PhotoKey } from 'common/data/photos/Photos';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph from 'common/typography/Paragraph';
import { Heading1, Heading3, Heading4 } from 'common/typography/Heading';
import Link, { LinkTarget } from 'common/components/Link/Link';

const styleClasses = makeStyles({
    featureIcon: {
        color: Color.Primary,
    },
});

const PricingPage = () => {
    const classes = styleClasses();
    return (
        <Page withNavbar withBackground={PhotoKey.MachuPicchu}>
            <CenteredComponent useLargeVersion>
                <DisplayCard>
                    <Heading1 color={TypographyColor.Primary}>
                        Babblegraph costs $29/year.
                    </Heading1>
                    <Paragraph>
                        But you can try it before you buy it with a 30-day, no credit card required free trial.
                    </Paragraph>
                    <Heading3 color={TypographyColor.Primary}>
                        Why is there no free option?
                    </Heading3>
                    <Paragraph>
                        Babblegraph is built a team of one. My name is Andrew, and I’m a software engineer from the United States. Making Babblegraph a paid service means that my primary incentive is to spend my time building out a great product instead of trying to get you to click on ads.
                    </Paragraph>
                    <Heading3 color={TypographyColor.Primary}>
                        Why is there no monthly option?
                    </Heading3>
                    <Paragraph>
                        Basically anytime you use your credit card, the company, website or person who charges you has to pay a fee. These fees really add up for small purchases. Taking into account these fees, the cost to you would be much higher if it were monthly. Babblegraph may offer a monthly plan in the future, though.
                    </Paragraph>
                    <Heading3 color={TypographyColor.Primary}>
                        What does $29/year get you with your subscription?
                    </Heading3>
                    <List>
                        <FeatureListItem
                            title="Receive a newsletter with articles and podcasts from real news sources and creators around Spain and Latin America"
                            description="Babblegraph curates content so that you don’t have to. It makes it really easy to discover new podcasts or newspapers from Spain and Latin America that make maintaining your Spanish fun and enjoyable."
                            icon={<FavoriteIcon className={classes.featureIcon} />} />
                        <FeatureListItem
                            title="Learn new vocabulary words and phrases with spaced repetition"
                            description="You can add vocabulary words and phrases that you learn while reading and listening to a vocabulary list. Babblegraph will send you more content that uses your new words so that you can expand your vocabulary."
                            icon={<DynamicFeedIcon className={classes.featureIcon} />} />
                        <FeatureListItem
                            title="Tailor your newsletter to your interests"
                            description="Babblegraph allows you to personalize your newsletter by allowing you to select topics that you find interesting. You’ll receive news articles and podcasts about these topics."
                            icon={<LibraryAddCheckIcon className={classes.featureIcon} />} />
                    </List>
                    <Link href="/" target={LinkTarget.Self}>
                        Return to home page
                    </Link>
                </DisplayCard>
            </CenteredComponent>
        </Page>
    )
}

type FeatureListItemProps = {
    title: string;
    description: string;
    icon: JSX.Element;
}

const FeatureListItem = (props: FeatureListItemProps) => {
    const classes = styleClasses();
    return (
        <ListItem alignItems="center">
            <ListItemIcon>
                {props.icon}
            </ListItemIcon>
            <ListItemText>
                <Heading4
                    align={Alignment.Left}
                    color={TypographyColor.Primary}>
                    {props.title}
                </Heading4>
                <Paragraph align={Alignment.Left}>
                    {props.description}
                </Paragraph>
            </ListItemText>
        </ListItem>
    );
}

export default PricingPage;
