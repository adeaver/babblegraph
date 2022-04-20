import React from 'react';

import Grid from '@material-ui/core/Grid';
import Box from '@material-ui/core/Box';

type CenteredComponentProps = {
    className?: string;
    children: React.ReactNode;
    useLargeVersion?: boolean;
}

const CenteredComponent = (props: CenteredComponentProps) => {
    return (
        <Grid className={props.className} container>
            <Box
                component={Grid}
                item
                xs={false}
                md={!!props.useLargeVersion ? 2 : 3}
                display={{ xs: "none", sm: "none", md: "block"}}>
                &nbsp;
            </Box>
            <Box
                component={Grid}
                item
                xs={12}
                md={!!props.useLargeVersion ? 8 : 6}
                display={{ xs: "block"}}>
                {props.children}
            </Box>
            <Box
                component={Grid}
                item
                xs={false}
                md={!!props.useLargeVersion ? 2 : 3}
                display={{ xs: "none", sm: "none", md: "block"}}>
                &nbsp;
            </Box>
        </Grid>
    );
}

export default CenteredComponent;
