import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Paragraph from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    GetOrCreateBillingInformationResponse,
    getOrCreateBillingInformation,
} from 'ConsumerWeb/api/billing/billing';

type Params = {
    token: string;
}

type PremiumNewsletterSubscriptionCheckoutPageProps = RouteComponentProps<Params>;

const PremiumNewsletterSubscriptionCheckoutPage = withUserProfileInformation<PremiumNewsletterSubscriptionCheckoutPageProps>(
    RouteEncryptionKey.PremiumSubscriptionCheckout,
    [RouteEncryptionKey.SubscriptionManagement],
    (ownProps: PremiumNewsletterSubscriptionCheckoutPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.CheckoutPage,
    (props: PremiumNewsletterSubscriptionCheckoutPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;
        const [ subscriptionManagementToken ] = props.userProfile.nextTokens;

        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <DisplayCardHeader
                            title="Babblegraph Premium Checkout"
                            backArrowDestination={`/manage/${subscriptionManagementToken}`} />
                        <OrderDetailsSection premiumSubscriptionCheckoutToken={token} />
                    </DisplayCard>
                </Grid>
            </Grid>
        );
    }
);

type OrderDetailsSectionProps = {
    premiumSubscriptionCheckoutToken: string;
}

const OrderDetailsSection = asBaseComponent<GetOrCreateBillingInformationResponse, OrderDetailsSectionProps>(
    (props: GetOrCreateBillingInformationResponse & OrderDetailsSectionProps & BaseComponentProps) => {
        return (
            <Grid container>
                <Grid item xs={12}>
                    <Heading3 align={Alignment.Left}>
                        Your Order
                    </Heading3>
                </Grid>
                <Grid item xs={10}>
                    <Paragraph align={Alignment.Left}>
                        1-year Babblegraph Premium Subscription
                    </Paragraph>
                </Grid>
                <Grid item xs={2}>
                    <Paragraph align={Alignment.Right}>
                        US$29.00
                    </Paragraph>
                </Grid>
                <Divider />
                <Grid item xs={10}>
                    <Paragraph align={Alignment.Left}>
                        Total Due Now
                    </Paragraph>
                </Grid>
                <Grid item xs={2}>
                    <Paragraph align={Alignment.Right}>
                        US$29.00
                    </Paragraph>
                </Grid>
            </Grid>
        );
    },
    (
        ownProps: OrderDetailsSectionProps,
        onSuccess: (GetOrCreateBillingInformationResponse) => void,
        onError: (err: Error) => void,
    ) => {
        getOrCreateBillingInformation({
            premiumSubscriptionCheckoutToken: ownProps.premiumSubscriptionCheckoutToken,
        },
        onSuccess,
        onError);
    },
    false,
);
export default PremiumNewsletterSubscriptionCheckoutPage;
