import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import { Heading1, Heading2, Heading4 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import Link from 'common/components/Link/Link';
import Color from 'common/styles/colors';
import { getStaticContentURLForPath } from 'util/static/static';

import {
    BlogPostMetadata,
    ContentNodeType,
    ContentNode,
    ContentNodeBody,
    Heading as HeadingContent,
    Paragraph as ParagraphContent,
    Image as ImageContent,
    Link as LinkContent,
    List as ListContent,
    ListType,
    getDefaultContentNodeForType,
} from 'common/api/blog/content';

const styleClasses = makeStyles({
    image: {
        borderRadius: '5px',
        width: '100%',
        height: 'auto',
    },
    orderedList: {
        fontFamily: "'Roboto', sans-serif",
        lineHeight: 1.5,
        fontSize: '16px',
        color: Color.TextGray,
    },
});

type BlogDisplayProps = {
    metadata: BlogPostMetadata;
    content: ContentNode[];
}

const BlogDisplay = (props: BlogDisplayProps) => {
    const classes = styleClasses();
    if (!props.metadata) {
        return (
            <div />
        );
    }
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
            <Divider />
            {
                props.content.map((node: ContentNode, idx: number) => {
                    if (node.type === ContentNodeType.Heading) {
                        return (
                            <HeadingDisplay key={`${props.metadata.id}-display-${idx}`} {...node.body as HeadingContent} />
                        );
                    } else if (node.type === ContentNodeType.Paragraph) {
                        return (
                            <ParagraphDisplay key={`${props.metadata.id}-display-${idx}`} {...node.body as ParagraphContent} />
                        );
                    } else if (node.type === ContentNodeType.Image) {
                        return (
                            <ImageDisplay key={`${props.metadata.id}-display-${idx}`} {...node.body as ImageContent} />
                        );
                    } else if (node.type === ContentNodeType.Link) {
                        return (
                            <LinkDisplay key={`${props.metadata.id}-display-${idx}`} {...node.body as LinkContent} />
                        );
                    } else if (node.type === ContentNodeType.List) {
                        return (
                            <ListDisplay key={`${props.metadata.id}-display-${idx}`} {...node.body as ListContent} />
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
    if (!props.path || !props.path.length) {
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

const LinkDisplay = (props: LinkContent) => {
    return (
        <Link href={props.destinationUrl}>
            { props.text }
        </Link>
    );
}

const ListDisplay = (props: ListContent) => {
    const classes = styleClasses();
    const body = props.items.map((item: string, idx: number) => (
        <li key={`list-item-${idx}`}>
            <Paragraph align={Alignment.Left}>
                {item}
            </Paragraph>
        </li>
    ))
    if (props.type === ListType.Unordered) {
        return (
            <ul>
                {body}
            </ul>
        );
    } else if (props.type === ListType.Ordered) {
        return (
            <ol className={classes.orderedList}>
                {body}
            </ol>
        );
    }
    throw new Error(`unrecognized list type ${props.type}`);
}

export default BlogDisplay;
