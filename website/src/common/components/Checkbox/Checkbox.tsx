import React from 'react';

import { withStyles } from '@material-ui/core/styles';
import Checkbox from '@material-ui/core/Checkbox';

import Color from 'common/styles/colors';

export const PrimaryCheckbox = withStyles({
    root: {
        color: Color.Primary,
        '&$checked': {
            color: Color.Primary,
        },
    },
    checked: {},
})((props) => <Checkbox color="default" {...props} />);
