import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import { Heading1, Heading2, Heading4 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { getStaticContentURLForPath } from 'util/static/static';

import {
    BlogPostMetadata,
    ContentNodeType,
    ContentNode,
    ContentNodeBody,
    Heading as HeadingContent,
    Paragraph as ParagraphContent,
    Image as ImageContent,
    getDefaultContentNodeForType,
} from 'common/api/blog/content.ts';

const styleClasses = makeStyles({
    image: {
        borderRadius: '5px',
        width: '100%',
        height: 'auto',
    },
});

type BlogDisplayProps = {
    metadata: BlogPostMetadata;
    content: ContentNode[];
}

const BlogDisplay = (props: BlogDisplayProps) => {
    const classes = styleClasses();
    return (
        <div>
            {
                !!props.metadata.heroImage && (
                    <img
                        className={classes.image}
                        src={getStaticContentURLForPath(props.metadata.heroImage.path)}
                        alt={props.metadata.heroImage.altText} />
                )
            }
            <Heading1
                color={TypographyColor.Primary}
                align={Alignment.Center}>
                { props.metadata.title }
            </Heading1>
            <Paragraph
                align={Alignment.Center}>
                { props.metadata.description }
            </Paragraph>
            {
                props.content.map((node: ContentNode, idx: number) => {
                    if (node.type === ContentNodeType.Heading) {
                        return (
                            <HeadingDisplay {...node.body as HeadingContent} />
                        );
                    } else if (node.type === ContentNodeType.Paragraph) {
                        return (
                            <ParagraphDisplay {...node.body as ParagraphContent} />
                        );
                    } else if (node.type === ContentNodeType.Image) {
                        return (
                            <ImageDisplay {...node.body as ImageContent} />
                        );
                    } else {
                        throw new Error("Unsupported node type");
                    }
                })
            }
        </div>
    );
}

const HeadingDisplay = (props: HeadingContent) => {
    return (
        <Heading2
            align={Alignment.Center}
            color={TypographyColor.Primary}>
            { props.text}
        </Heading2>
    );
}

const ParagraphDisplay = (props: ParagraphContent) => {
    return (
        <Paragraph align={Alignment.Left}>
            { props.text }
        </Paragraph>
    );
}

const ImageDisplay = (props: ImageContent) => {
    if (!props.path.length) {
        return null;
    }
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={2}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={8}>
                <img
                    className={classes.image}
                    src={getStaticContentURLForPath(props.path)}
                    alt={props.altText} />
                {
                    !!props.caption && (
                        <Paragraph align={Alignment.Center} size={Size.Small}>
                            {props.caption}
                        </Paragraph>
                    )
                }
            </Grid>
        </Grid>
    )
}

export default BlogDisplay;
