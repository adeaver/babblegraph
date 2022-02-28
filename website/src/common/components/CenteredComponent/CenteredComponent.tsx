import React from 'react';

import Grid from '@material-ui/core/Grid';

type CenteredComponentProps = {
    className?: string;
    children: React.ReactNode;
}

const CenteredComponent = (props: CenteredComponentProps) => {
    return (
        <Grid className={props.className} container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                {props.children}
            </Grid>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
        </Grid>
    );
}

export default CenteredComponent;
