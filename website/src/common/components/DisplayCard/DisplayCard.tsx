import React from 'react';

import { withStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';

const DisplayCard = withStyles({
    root: {
        padding: '20px',
        marginTop: '20px',
    },
})((props) => <Card variant='outlined' {...props} />);

export default DisplayCard;
