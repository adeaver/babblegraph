import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';

import Color from 'common/styles/colors';

export const PrimaryButton = withStyles({
    root: {
        color: Color.White,
        backgroundColor: Color.Primary,
        '&:hover': {
            backgroundColor: Color.Primary,
        },
    },
})(Button);

export const WarningButton = withStyles({
    root: {
        color: Color.White,
        backgroundColor: Color.Warning,
        '&:hover': {
            backgroundColor: Color.Warning,
        },
    },
})(Button);

export const ConfirmationButton = withStyles({
    root: {
        color: Color.White,
        backgroundColor: Color.Confirmation,
        '&:hover': {
            backgroundColor: Color.Confirmation,
        },
    },
})(Button);
