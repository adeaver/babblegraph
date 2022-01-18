import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import ActionCard from 'common/components/ActionCard/ActionCard';
import Link, { LinkTarget } from 'common/components/Link/Link';
import { TypographyColor } from 'common/typography/common';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { getStaticContentURLForPath } from 'util/static/static';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { setLocation } from 'util/window/Location';

import {
    GetAllBlogPostsResponse, getAllBlogPosts,
} from 'ConsumerWeb/api/blog/blog';
import {
    BlogPostMetadata,
} from 'common/api/blog/content';

const styleClasses = makeStyles({
    image: {
        borderRadius: '5px',
        width: '100%',
        height: 'auto',
    },
    displayCard: {
        margin: '10px 0',
    },
    blogPostRow: {
        display: 'flex',
        alignItems: 'center',
    },
});

type BlogListPageProps = {};

const BlogListPage = (props: BlogListPageProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ blogPosts, setBlogPosts ] = useState<BlogPostMetadata[]>([]);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getAllBlogPosts({},
            (resp: GetAllBlogPostsResponse) => {
                setBlogPosts(resp.blogPosts);
                setIsLoading(false);
            },
            (err: Error) => {
                setError(err);
                setIsLoading(false);
            });
    }, []);

    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Warning}>
                An error occurred. Check back later!
            </Heading3>
        );
    } else {
        body = blogPosts.map((post: BlogPostMetadata, idx: number) => (
            <BlogPostDisplay key={`blog-post-${idx}`} {...post} />
        ));
    }

    return (
        <Page>
            <Grid container>
                <Grid xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid xs={12} md={6}>
                    <Heading1 color={TypographyColor.Primary}>
                        Babblegraph Blog
                    </Heading1>
                    <Link href="/" target={LinkTarget.Self}>
                        Return to Home Page
                    </Link>
                    { body }
                </Grid>
            </Grid>
        </Page>
    );
}

const BlogPostDisplay = (props: BlogPostMetadata) => {
    const classes = styleClasses();
    return (
        <ActionCard
            className={classes.displayCard}
            onClick={() => setLocation(`/blog/${props.urlPath}`)}>
            <Grid className={classes.blogPostRow} container>
                <Grid item xs={12} md={4}>
                    <img
                        className={classes.image}
                        src={getStaticContentURLForPath(props.heroImage.path)}
                        alt={props.heroImage.altText} />
                </Grid>
                <Grid item xs={12} md={8}>
                    <Heading3 color={TypographyColor.Primary}>
                        { props.title}
                    </Heading3>
                    <Paragraph>
                        { props.description }
                    </Paragraph>
                    <Link href={`/blog/${props.urlPath}`} target={LinkTarget.Self}>
                        Click here to read
                    </Link>
                </Grid>
            </Grid>
        </ActionCard>
    );
}

export default BlogListPage;
