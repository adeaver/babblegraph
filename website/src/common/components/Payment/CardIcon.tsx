import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import CreditCardIcon from '@material-ui/icons/CreditCard';

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
    "amex": true,
    "mastercard": true,
    "visa": true,
}

type CardIconProps = {
    cardType: string;
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

export default CardIcon;
