import React, { useEffect, useState } from 'react';

import Grid from '@material-ui/core/Grid';

import { Heading1 } from 'common/typography/Heading';
import Page from 'common/components/Page/Page';

type BlogHomePageProps = {}

const BlogHomePage = (props: BlogHomePageProps) => {
   return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Heading1>Babblegraph Blog</Heading1>
                </Grid>
            </Grid>
        </Page>
   );
}

export default BlogHomePage;
