import React, { useState, useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Page from 'common/components/Page/Page';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import BlogMetadataEditForm from './BlogMetadataEditForm';
import BlogContentEditor from './BlogContentEditor';
import BlogViewMetricsComponent from './BlogViewMetricsComponent';

import {
    BlogPostMetadata,
} from 'common/api/blog/content.ts';
import {
    getBlogPostMetadataByURLPath,
    GetBlogPostMetadataByURLPathResponse,
    getBlogPostViewMetrics,
    GetBlogPostViewMetricsResponse,
    BlogPostViewMetrics,
} from 'AdminWeb/api/blog/blog';

type Params = {
    blogPath: string;
}

type BlogEditPageProps = RouteComponentProps<Params>;

const BlogEditPage = (props: BlogEditPageProps) => {
    const { blogPath } = props.match.params;

    const [ blogPostMetadata, setBlogPostMetadata ] = useState<BlogPostMetadata | null>(null);
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getBlogPostMetadataByURLPath({
            urlPath: blogPath,
        },
        (resp: GetBlogPostMetadataByURLPathResponse) => {
            setIsLoading(false);
            setBlogPostMetadata(resp.blogPost);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }, []);


    return (
        <Page>
            <BlogViewMetricsComponent urlPath={blogPath} />
            {
                isLoading ? (
                    <LoadingSpinner />
                ) : (
                    !!blogPostMetadata ? (
                        <BlogMetadataEditForm
                            setIsLoading={setIsLoading}
                            updateBlogPostMetadata={setBlogPostMetadata}
                            blogPostMetadata={blogPostMetadata}
                            urlPath={blogPath} />
                    ) : (
                        <Heading3 color={TypographyColor.Warning}>
                            An error has occurred
                        </Heading3>
                    )
                )
            }
            <BlogContentEditor
                urlPath={blogPath}
                blogPostMetadata={blogPostMetadata} />
        </Page>
    );
}

export default BlogEditPage;
