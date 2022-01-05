import React, { useState, useEffect } from 'react';

import Grid from '@material-ui/core/Grid';

import { TypographyColor } from 'common/typography/common';
import { Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    getBlogPostViewMetrics,
    GetBlogPostViewMetricsResponse,
    BlogPostViewMetrics,
} from 'AdminWeb/api/blog/blog';

type BlogViewMetricsComponentProps = {
    urlPath: string;
}

const BlogViewMetricsComponent = (props: BlogViewMetricsComponentProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ blogPostViewMetrics, setBlogPostViewMetrics ] = useState<BlogPostViewMetrics>(null);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getBlogPostViewMetrics({
            urlPath: props.urlPath,
        },
        (resp: GetBlogPostViewMetricsResponse) => {
            setBlogPostViewMetrics(resp.viewMetrics);
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
                An error has occurred
            </Heading3>
        );
    } else if (!!blogPostViewMetrics) {
        body = (
            <Grid container>
                <Grid item xs={6} md={3}>
                    <Paragraph>
                        {blogPostViewMetrics.totalViews} total views
                    </Paragraph>
                </Grid>
                <Grid item xs={6} md={3}>
                    <Paragraph>
                        {blogPostViewMetrics.uniqueViews} unique views
                    </Paragraph>
                </Grid>
                <Grid item xs={6} md={3}>
                    <Paragraph>
                        {blogPostViewMetrics.lastMonthTotalViews} total views in last month
                    </Paragraph>
                </Grid>
                <Grid item xs={6} md={3}>
                    <Paragraph>
                        {blogPostViewMetrics.lastMonthUniqueViews} unique views in last month
                    </Paragraph>
                </Grid>
            </Grid>
        );
    }

    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <DisplayCard>
                    <Heading3 color={TypographyColor.Primary}>
                        Metrics
                    </Heading3>
                    { body }
                </DisplayCard>
            </Grid>
        </Grid>
    );
}

export default BlogViewMetricsComponent;
