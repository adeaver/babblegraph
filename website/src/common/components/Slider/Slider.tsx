import React from 'react';

import { withStyles } from '@material-ui/core/styles';
import Slider from '@material-ui/core/Slider';

import Color from 'common/styles/colors';

export const PrimarySlider = withStyles({
    root: {
        color: Color.Primary,
            height: 8,
    },
    thumb: {
        height: 24,
        width: 24,
        backgroundColor: Color.White,
        border: '2px solid currentColor',
        marginTop: -8,
        marginLeft: -12,
        '&:focus, &:hover, &$active': {
            boxShadow: 'inherit',
        },
    },
    active: {},
    valueLabel: {
        left: 'calc(-50% + 4px)',
    },
    track: {
        height: 8,
        borderRadius: 4,
    },
    rail: {
        height: 8,
        borderRadius: 4,
    },
})((props) => <Slider {...props} />);
