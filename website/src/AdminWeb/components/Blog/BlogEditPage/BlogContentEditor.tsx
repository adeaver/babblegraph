import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';

const styleClasses = makeStyles({
    blogContentEditorContainer: {
        padding: '10px',
    },
});


type BlogContentEditorProps = {
    urlPath: string;
}

const BlogContentEditor = (props: BlogContentEditorProps) => {
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid className={classes.blogContentEditorContainer} item xs={6}>
                <DisplayCard>
                    Hello
                </DisplayCard>
            </Grid>
            <Grid className={classes.blogContentEditorContainer} item xs={6}>
                <DisplayCard>
                    Hello
                </DisplayCard>
            </Grid>
        </Grid>
    );
}

export default BlogContentEditor;
