import React from 'react';

import Grid from '@material-ui/core/Grid';

import { Heading1 } from 'common/typography/Heading';

type HeaderProps = {};

const Header = (props: HeaderProps) => {
    return (
        <div className="Header__root">
            <Grid container>
                <Grid item xs={11} md={3}>
                    <Heading1>babblegraph</Heading1>
                </Grid>
            </Grid>
            <hr />
        </div>
    );
}

export default Header;
