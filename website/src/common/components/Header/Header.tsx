import React from 'react';
import { useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Color from 'common/styles/colors';
import { Heading2 } from 'common/typography/Heading';
import {
    TypographyColor,
    Alignment
} from 'common/typography/common';

const styleClasses = makeStyles({
    headerRoot: {
        height: '80px',
        background: Color.White,
    },
    headerContainer: {
        padding: '0 10px',
    },
    headerDivider: {
        color: Color.Primary,
    },
    headerAnchor: {
        textDecoration: 'none',
        cursor: 'pointer',
        color: 'inherit',
    },
});

type HeaderProps = {};

const Header = (props: HeaderProps) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <div className={classes.headerRoot}>
            <Grid className={classes.headerContainer} container>
                <Grid item xs={11} md={3}>
                    <Heading2 color={TypographyColor.Primary} align={Alignment.Left}>
                        <a className={classes.headerAnchor} href="/">
                            babblegraph
                        </a>
                    </Heading2>
                </Grid>
            </Grid>
            <hr className={classes.headerDivider} />
        </div>
    );
}

export default Header;
