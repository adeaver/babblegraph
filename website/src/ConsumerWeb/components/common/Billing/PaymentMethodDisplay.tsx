import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import CreditCardIcon from '@material-ui/icons/CreditCard';

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
    isHighlighted?: boolean;

    onClick: (externalID: string) => void;
}

const PaymentMethodDisplay = (props: PaymentMethodDisplayProps) => {
    return (
        <ActionCard onClick={() => props.onClick(props.paymentMethod.externalId)}>
            <Grid container>
                <Grid item xs={2}>
                    <CardIcon cardType={props.paymentMethod.cardType} />
                </Grid>
                <Grid item xs={10}>
                    <Paragraph color={props.isHighlighted ? TypographyColor.Primary : TypographyColor.Gray}>
                        {`****${props.paymentMethod.displayMask}`}
                    </Paragraph>
                </Grid>
                <Grid item xs={12} md={6}>
                    <Paragraph
                        align={Alignment.Left}
                        color={props.isHighlighted ? TypographyColor.Primary : TypographyColor.Gray}>
                        Expires {props.paymentMethod.cardExpiration}
                    </Paragraph>
                </Grid>
                {
                    !!props.isHighlighted && (
                        <Grid item xs={12} md={6}>
                            <Paragraph
                                align={Alignment.Right}
                                color={props.isHighlighted ? TypographyColor.Primary : TypographyColor.Gray}>
                                selected
                            </Paragraph>
                        </Grid>
                    )
                }
            </Grid>
        </ActionCard>
    );
}

const styleClasses = makeStyles({
    cardIconRoot: (props: CardIconProps) => {
        const baseProperties = {
            height: '42px',
            width: 'auto'
        };
        if (!!validCardTypes[props.cardType]) {
            return {
                ...baseProperties,
                backgroundImage: `url("https://static.babblegraph.com/assets/payment/${props.cardType}.png")`,
                backgroundSize: 'contain',
                backgroundRepeat: 'no-repeat',
                backgroundPosition: 'left',
            };
        }
        return baseProperties;
    }
});

const validCardTypes = {
    [CardType.Amex]: true,
    [CardType.Visa]: true,
    [CardType.Mastercard]: true,
}

type CardIconProps = {
    cardType: CardType;
}

const CardIcon = (props: CardIconProps) => {
    const classes = styleClasses(props);
    if (!!validCardTypes[props.cardType]) {
        return (
            <div className={classes.cardIconRoot} />
        );
    } else {
        return <CreditCardIcon className={classes.cardIconRoot} />;
    }
}

export default PaymentMethodDisplay;
