import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import DateRangeIcon from '@material-ui/icons/DateRange';
import FindInPageIcon from '@material-ui/icons/FindInPage';
import BallotIcon from '@material-ui/icons/Ballot';

import Color from 'common/styles/colors';
import { Heading1, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';

const styleClasses = makeStyles({
    featureIcon: {
        color: Color.Primary,
    },
});

type PremiumInformationPanelProps = {};

const PremiumInformationPanel = (props: PremiumInformationPanelProps) => {
    const classes = styleClasses();
    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                Babblegraph Premium
            </Heading1>
            <Heading3>
                 Get access to new features that improve your learning and give you more control over how Babblegraph fits into your journey!
            </Heading3>
            <Divider />
            <List>
                <PremiumFeatureListItem
                    title="Newsletter Scheduling"
                    description="Choose which days you receive newsletters from Babblegraph and which days you don’t."
                    icon={<DateRangeIcon className={classes.featureIcon} />} />
                <PremiumFeatureListItem
                    title="Word Tracking Spotlights"
                    description="With word tracking spotlights, every newsletter will include an article that highlights a word on your tracking list. Spotlights rotates through all your words, so you’ll even see less common words regularly!"
                    icon={<FindInPageIcon className={classes.featureIcon} />} />
                <PremiumFeatureListItem
                    title="Newsletter Customization"
                    description="Want fewer articles in your newsletter? Want guarantee that you’ll see a certain topic in every email? Premium Subscribers can do just that with the newsletter customization tool!"
                    icon={<BallotIcon className={classes.featureIcon} />} />
            </List>
            <Heading3 color={TypographyColor.Primary}>
                $3/month or $34/year
            </Heading3>
            <Heading3>
                The best part is that you can try it for free for 14 days!
            </Heading3>
        </div>
    )
}

type PremiumFeatureListItemProps = {
    title: string;
    description: string;
    icon: JSX.Element;
}

const PremiumFeatureListItem = (props: PremiumFeatureListItemProps) => {
    const classes = styleClasses();
    return (
        <ListItem>
            <ListItemIcon>
                {props.icon}
            </ListItemIcon>
            <ListItemText>
                <Heading4
                    align={Alignment.Left}
                    color={TypographyColor.Primary}>
                    {props.title}
                </Heading4>
                <Paragraph align={Alignment.Left} size={Size.Small}>
                    {props.description}
                </Paragraph>
            </ListItemText>
        </ListItem>
    );
}

export default PremiumInformationPanel;
