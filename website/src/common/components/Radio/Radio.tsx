import React from 'react';

import { withStyles } from '@material-ui/core/styles';
import Radio from '@material-ui/core/Radio';

import Color from 'common/styles/colors';

export const PrimaryRadio = withStyles({
    root: {
        color: Color.Primary,
        '&$checked': {
            color: Color.Primary,
        },
    },
    checked: {},
})((props) => <Radio color="default" {...props} />);
