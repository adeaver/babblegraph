import React from 'react';

import { withStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';

import Color from 'common/styles/colors';

export const PrimarySwitch = withStyles({
    switchBase: {
        color: Color.BorderGray,
        '&$checked': {
              color: Color.Primary,
        },
        '&$checked + $track': {
            backgroundColor: Color.Primary,
        },
    },
    checked: {},
    track: {
        display: 'inline-block',
    },
})((props) => <Switch color="default" {...props} />);
