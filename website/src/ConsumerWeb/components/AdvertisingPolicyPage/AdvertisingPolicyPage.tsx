import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import Page from 'common/components/Page/Page';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { Heading1, Heading2, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph from 'common/typography/Paragraph';
import Link, { LinkTarget } from 'common/components/Link/Link';

const AdvertisingPolicyPage = () => {
    return (
        <Page>
            <CenteredComponent>
                <DisplayCard>
                    <Heading1
                        align={Alignment.Left}
                        color={TypographyColor.Primary}>
                        Babblegraph Advertising Policy
                    </Heading1>
                    <Divider />
                    <Paragraph align={Alignment.Left}>
                        Hello! This is in part a policy page and part an open letter from the creator of (and only person working on) Babblegraph!
                    </Paragraph>
                    <Paragraph align={Alignment.Left}>
                        I completely understand that advertisements are most people’s least favorite part of using the internet. They are a necessary evil since it costs money to run and grow Babblegraph.
                    </Paragraph>
                    <Paragraph align={Alignment.Left}>
                        These are the main points of Babblegraph’s advertising policy.
                    </Paragraph>
                    <ul>
                        <li>
                            <Paragraph align={Alignment.Left} isBold>
                                Babblegraph primarily uses affiliate-based advertisements. This means that Babblegraph may earn commission if you purchase items through links in your newsletter. It’s a great way to help support Babblegraph at no additional cost to you.
                            </Paragraph>
                        </li>
                        <Paragraph align={Alignment.Left}>
                            Unlike other advertising strategies, affiliate-based marketing means that I’m incentivized to show products that you’re actually interested in purchasing because Babblegraph only compensated if you purchase an item. As a general rule, the website from which you purchase a product is who is paying the commission that we earn. All advertisements will contain a disclaimer that it is an adveritsement.
                        </Paragraph>
                        <li>
                            <Paragraph align={Alignment.Left} isBold>
                                All products are manually approved by a human before they appear in your newsletter. Babblegraph does not allow advertisers to run their own ads.
                            </Paragraph>
                        </li>
                        <Paragraph align={Alignment.Left}>
                            I rely on research, testimony from experts and others that I trust, as well as personal experience in choosing which products and brands are featured in your newsletter.
                        </Paragraph>
                        <li>
                            <Paragraph align={Alignment.Left} isBold>
                                My primary hope is that the frequency and the contents of advertisements on Babblegraph strike a balance between not impacting your experience with Babblegraph and helping to support Babblegraph’s continued operation.
                            </Paragraph>
                        </li>
                        <Paragraph align={Alignment.Left}>
                            If you have any feedback about advertisements, please reach out to hello@babblegraph.com. Likewise, if you want to support Babblegraph without ads, you can subscribe to Babblegraph Premium (and get some new exclusive features in the process.)
                        </Paragraph>
                        <li>
                            <Paragraph align={Alignment.Left} isBold>
                                Lastly, and perhaps most importantly, Babblegraph does not sell any information about you. While we try to partner with brands and vendors that are respectful about privacy, it is worth noting that we are not responsible for their policies around data privacy.
                            </Paragraph>
                        </li>
                    </ul>
                </DisplayCard>
            </CenteredComponent>
        </Page>
    )
}

export default AdvertisingPolicyPage;
