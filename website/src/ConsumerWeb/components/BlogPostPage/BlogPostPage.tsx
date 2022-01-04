import React, { useState, useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';

import Paragraph from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import BlogDisplay from 'common/components/BlogDisplay/BlogDisplay';
import Link from 'common/components/Link/Link';

import {
    getBlogMetadata,
    GetBlogMetadataResponse,
    getBlogContent,
    GetBlogContentResponse,
} from 'ConsumerWeb/api/blog/blog';
import {
    ContentNode,
    BlogPostMetadata,
} from 'common/api/blog/content';

type Params = {
    blogPath: string;
}

type BlogPostPageProps = RouteComponentProps<Params>;

const BlogPostPage = (props: BlogPostPageProps) => {
    const { blogPath } = props.match.params;

    const [ error, setError ] = useState<Error>(null);

    const [ blogMetadata, setBlogMetadata ] = useState<BlogPostMetadata>(null);
    const [ isLoadingBlogMetadata, setIsLoadingBlogMetadata ] = useState<boolean>(true);

    const [ blogContent, setBlogContent ] = useState<ContentNode[]>([]);
    const [ isLoadingBlogContent, setIsLoadingBlogContent ] = useState<boolean>(true);

    useEffect(() => {
        getBlogMetadata({
            urlPath: blogPath,
        },
        (resp: GetBlogMetadataResponse) => {
            setBlogMetadata(resp.blogPost);
            setIsLoadingBlogMetadata(false);
        },
        (err: Error) => {
            setError(err);
            setIsLoadingBlogMetadata(false);
        });
        getBlogContent({
            urlPath: blogPath,
        },
        (resp: GetBlogContentResponse) => {
            setBlogContent(resp.content);
            setIsLoadingBlogContent(false);
        },
        (err: Error) => {
            setError(err);
            setIsLoadingBlogContent(false);
        });
    }, []);

    const isLoading = isLoadingBlogContent || isLoadingBlogMetadata;
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={2}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={8}>
                    <DisplayCard>
                        <DisplayCardHeader
                            title="Return to all blog posts"
                            backArrowDestination="/blog" />
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                !!error ? (
                                    <Heading3 color={TypographyColor.Warning}>
                                        Something went wrong! Check back later.
                                    </Heading3>
                                ) : (
                                    <div>
                                        <BlogDisplay
                                            content={blogContent}
                                            metadata={blogMetadata} />
                                        <Heading3
                                            color={TypographyColor.Primary}
                                            align={Alignment.Left}>
                                            About Babblegraph
                                        </Heading3>
                                        <Paragraph
                                            align={Alignment.Left}>
                                            Babblegraph helps intermediate and advanced Spanish students effortlessly get a daily dose of Spanish practice. With Babblegraph you’ll receive a daily newsletter with news articles from trusted, Spanish-language news sources from Spain and Latin America. Sign up for free today!
                                        </Paragraph>
                                        <Link href="/">
                                            Click here to learn more
                                        </Link>
                                    </div>
                                )
                            )
                        }
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default BlogPostPage;
