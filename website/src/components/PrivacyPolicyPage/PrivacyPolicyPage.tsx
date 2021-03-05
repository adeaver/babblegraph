import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph from 'common/typography/Paragraph';
import { PhotoKey } from 'common/data/photos/Photos';
import Link, { LinkTarget } from 'common/components/Link/Link';

const styleClasses = makeStyles({
    privacyPolicyContainingCard: {
        padding: '20px',
    },
});

const PrivacyPolicyPage = () => {
    const classes = styleClasses();
    return (
        <Page withBackground={PhotoKey.Seville}>
            <Grid container>
                <Grid item xs={0} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.privacyPolicyContainingCard}>
                        <Heading1 align={Alignment.Left}>
                            Babblegraph Privacy Policy
                        </Heading1>
                        <Paragraph align={Alignment.Left}>
                            This Privacy Policy describes how your personal information is collected, used, and shared when you visit or make a purchase from babblegraph.com (the “Site”).
                        </Paragraph>
                        <Heading3 align={Alignment.Left}>
                            PERSONAL INFORMATION WE COLLECT
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            When you visit the Site, we automatically collect certain information about your device, including information about your web browser, IP address, time zone, and some of the cookies that are installed on your device. Additionally, when you open emails sent by Babblegraph, we collect information when an email is first opened, including timestamps. At this time, we do not collect any information regarding which links were followed in the email. We also collect information about how you interact with the Site. We refer to this automatically-collected information as “Device Information.”
                        </Paragraph>
                        <Paragraph align={Alignment.Left}>
                            We collect Device Information using the following technologies:
                        </Paragraph>
                        <ul>
                            <li>
                                <Paragraph align={Alignment.Left}>
                                    “Log files” track actions occurring on the Site, and collect data including your IP address, browser type, Internet service provider, referring/exit pages, and date/time stamps.
                                </Paragraph>
                            </li>
                            <li>
                                <Paragraph align={Alignment.Left}>
                                    “Web beacons,” “tags,” and “pixels” are electronic files used to record information about how you browse the Site, including when emails sent by Babblegraph are opened.
                                </Paragraph>
                            </li>
                        </ul>
                        <Heading3 align={Alignment.Left}>
                            HOW DO WE USE YOUR PERSONAL INFORMATION?
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            We use the Device Information that we collect to help us screen for potential risk and fraud (in particular, your IP address), and more generally to improve and optimize our Site (for example, by generating analytics about how our users browse and interact with the Site, and to assess the success of our email product).
                        </Paragraph>
                        <Paragraph align={Alignment.Left}>
                            We do not advertise or share any of the information that you enter in order to advertise as of this time.
                        </Paragraph>
                        <Heading3 align={Alignment.Left}>
                           SHARING YOUR PERSONAL INFORMATION
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            We share your Personal Information with third parties to help us use your Personal Information, as described above.  For example, we use Amazon Web Services Simple Email Service to power our online email platform--you can read more about how Amazon Web Services uses your Personal Information here:  https://aws.amazon.com/compliance/data-privacy-faq.
                        </Paragraph>
                        <Paragraph align={Alignment.Left}>
                            Finally, we may also share your Personal Information to comply with applicable laws and regulations, to respond to a subpoena, search warrant or other lawful request for information we receive, or to otherwise protect our rights.
                        </Paragraph>
                        <Heading3 align={Alignment.Left}>
                           DO NOT TRACK
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            Please note that we do not alter our Site’s data collection and use practices when we see a Do Not Track signal from your browser.
                        </Paragraph>
                        <Heading3 align={Alignment.Left}>
                           YOUR RIGHTS
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            If you are a European resident, you have the right to access personal information we hold about you and to ask that your personal information be corrected, updated, or deleted. If you would like to exercise this right, please contact us through the contact information below.
                        </Paragraph>
                        <Paragraph align={Alignment.Left}>
                            Additionally, if you are a European resident we note that we are processing your information to pursue our legitimate business interests listed above. Additionally, please note that your information will be transferred outside of Europe, including to the United States.
                        </Paragraph>
                        <Heading3 align={Alignment.Left}>
                           MINORS
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            The Site is not intended for individuals under the age of 13.
                        </Paragraph>
                        <Heading3 align={Alignment.Left}>
                           CHANGES
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            We may update this privacy policy from time to time in order to reflect, for example, changes to our practices or for other operational, legal or regulatory reasons.
                        </Paragraph>
                        <Heading3 align={Alignment.Left}>
                           CONTACT US
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            For more information about our privacy practices, if you have questions, or if you would like to make a complaint, please contact us by e-mail at hello@babblegraph.com.
                        </Paragraph>
                        <Link href="/" target={LinkTarget.Self}>
                            Home
                        </Link>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

export default PrivacyPolicyPage;
