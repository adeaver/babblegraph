import React from 'react';

import Grid from '@material-ui/core/Grid';

import ActionCard from 'common/components/ActionCard/ActionCard';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading1, Heading2 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';

import {
    CardType,
    PaymentMethod,
} from 'ConsumerWeb/api/billing/billing';

type PaymentMethodDisplayProps = {
    paymentMethod: PaymentMethod;

    onClick: (externalID: string) => void;
}

const PaymentMethodDisplay = (props: PaymentMethodDisplayProps) => {
    return (
        <ActionCard onClick={() => props.onClick(props.paymentMethod.externalId)}>
            <Grid container>
                <Grid item xs={2}>
                    <Paragraph align={Alignment.Left}>
                        {props.paymentMethod.cardType}
                    </Paragraph>
                </Grid>
                <Grid item xs={10}>
                    <Paragraph>
                        {`****${props.paymentMethod.displayMask}`}
                    </Paragraph>
                </Grid>
                <Grid item xs={6}>
                    <Paragraph align={Alignment.Left}>
                        Expires {props.paymentMethod.cardExpiration}
                    </Paragraph>
                </Grid>
            </Grid>
        </ActionCard>
    );
}

export default PaymentMethodDisplay;
