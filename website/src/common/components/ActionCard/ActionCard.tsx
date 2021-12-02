import React from 'react';

import { withStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';

const ActionCard = withStyles({
    root: {
        padding: '10px',
        "&:hover": {
            boxShadow: "0 0 4px 2px gray",
        },
        cursor: 'pointer',
    },
})((props) => <Card variant='outlined' {...props} />);

export default ActionCard;
