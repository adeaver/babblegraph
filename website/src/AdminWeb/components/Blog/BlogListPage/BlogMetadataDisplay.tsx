import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import ActionCard from 'common/components/ActionCard/ActionCard';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import { setLocation } from 'util/window/Location';
import Paragraph from 'common/typography/Paragraph';

import {
    BlogPostMetadata,
    PostStatus,
} from 'AdminWeb/api/blog/blog';

const styleClasses = makeStyles({
    blogMetadataCard: {
        margin: '10px 0',
    },
});

const BlogMetadataDisplay = (props: BlogPostMetadata) => {
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <ActionCard className={classes.blogMetadataCard} onClick={() => setLocation(`blog-manager/edit/${props.urlPath}`)}>
                    <Heading3 color={TypographyColor.Primary}>
                        {props.title}
                    </Heading3>
                    <Paragraph>
                        {props.description}
                    </Paragraph>
                    <Paragraph color={props.status === PostStatus.Live ? TypographyColor.Confirmation : TypographyColor.Gray}>
                        Status: {props.status}
                    </Paragraph>
                </ActionCard>
            </Grid>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
        </Grid>
    );
}

export default BlogMetadataDisplay;
