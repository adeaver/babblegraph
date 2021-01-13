import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';

import Color from 'common/styles/colors';

export const PrimaryTextField = withStyles({
    root: {
        '& label.Mui-focused': {
            color: Color.Primary,
        },
        '& .MuiInput-underline:after': {
            borderBottomColor: Color.Primary,
        },
        '& .MuiOutlinedInput-root': {
            '&:hover fieldset': {
                borderColor: Color.Primary,
            },
            '&.Mui-focused fieldset': {
                borderColor: Color.Primary,
            },
        },
    },
})(TextField);
