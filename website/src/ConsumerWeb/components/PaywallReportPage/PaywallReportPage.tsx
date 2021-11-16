import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import VerifiedUserIcon from '@material-ui/icons/VerifiedUser';

import Page from 'common/components/Page/Page';
import { Heading1, Heading2 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import Paragraph from 'common/typography/Paragraph';
import Color from 'common/styles/colors';
import Link, { LinkTarget } from 'common/components/Link/Link';

const styleClasses = makeStyles({
    containingCard: {
        padding: '20px',
    },
    checkIcon: {
        color: Color.Confirmation,
        display: 'block',
        margin: '0 auto',
        fontSize: '48px',
    }
});

type Params = {
    token: string
}

type PaywallReportPageProps = RouteComponentProps<Params>;

const PaywallReportPage = (props: PaywallReportPageProps) => {
    const classes = styleClasses();
    const { token } = props.match.params;
    return (
        <Page>
            <Grid container>
                <Grid item xs={0} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.containingCard}>
                        <VerifiedUserIcon className={classes.checkIcon} />
                        <Heading1 color={TypographyColor.Primary}>
                            Thank you for reporting this exclusive content
                        </Heading1>
                        <Heading2>
                            Gracias por informar de este contenido exclusivo
                        </Heading2>
                        <Paragraph>
                            This content will be investigated so that other articles like it won’t show up in daily emails.
                        </Paragraph>
                        <Paragraph>
                            Este contenido se investigará para que otro contenido similar no parezca en el email diario.
                        </Paragraph>
                        <Link href={`/manage/${token}`} target={LinkTarget.Self}>
                            Go to Subscription Management Page
                        </Link>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

export default PaywallReportPage;
