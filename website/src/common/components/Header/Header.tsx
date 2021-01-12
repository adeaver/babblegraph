import './Header.scss';

import React from 'react';

import Grid from '@material-ui/core/Grid';

import { Heading2 } from 'common/typography/Heading';
import { Color, Alignment } from 'common/typography/common';

type HeaderProps = {};

const Header = (props: HeaderProps) => {
    return (
        <div className="Header__root">
            <Grid className="Header__container" container>
                <Grid item xs={11} md={3}>
                    <Heading2 color={Color.Primary} align={Alignment.Left}>babblegraph</Heading2>
                </Grid>
            </Grid>
            <hr className="Header__divider" />
        </div>
    );
}

export default Header;
