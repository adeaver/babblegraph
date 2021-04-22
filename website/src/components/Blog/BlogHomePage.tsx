import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import { Heading1, Heading2 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import Page from 'common/components/Page/Page';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryButton } from 'common/components/Button/Button';

import {
    GetAllBlogPostsPaginatedResponse,
    BlogPost,
    getAllBlogPostsPaginated,
} from 'api/blog/bloghome';
import {
    getImageURL,
} from 'api/blog/blogpost';

const styleClasses = makeStyles({
    heroImage: {
        width: '100%',
        height: 'auto',
    },
    blogPostCard: {
        padding: '20px',
        boxSizing: 'border-box',
        margin: '15px 0',
    },
});

type BlogHomePageProps = {}

const BlogHomePage = (props: BlogHomePageProps) => {
    const [ pageIndex, setPageIndex ] = useState<number>(0);
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ blogPosts, setBlogPosts ] = useState<Array<BlogPost>>([]);
    const [ hasMoreBlogs, setHasMoreBlogs ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getAllBlogPostsPaginated({
            pageIndex: pageIndex,
        },
        (resp: GetAllBlogPostsPaginatedResponse) => {
            setIsLoading(false);
            setBlogPosts(resp.blogPosts);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }, []);

    const handleUpdatePageIndex = () => {
        setIsLoading(true);
        getAllBlogPostsPaginated({
            pageIndex: pageIndex + 1,
        },
        (resp: GetAllBlogPostsPaginatedResponse) => {
            setIsLoading(false);
            if (!resp.blogPosts) {
                setHasMoreBlogs(false);
            }
            setBlogPosts(blogPosts.concat(resp.blogPosts || []));
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
        setPageIndex(pageIndex + 1);
    }

    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = <Heading1>Something went wrong</Heading1>;
    } else {
        body = (
            <BlogPostsView
                handleUpdatePageIndex={handleUpdatePageIndex}
                shouldShowNextButton={hasMoreBlogs}
                blogPosts={blogPosts} />
        );
    }
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    {body}
                </Grid>
            </Grid>
        </Page>
    );
}

type BlogPostsViewProps = {
    handleUpdatePageIndex: () => void;
    shouldShowNextButton: boolean;
    blogPosts: Array<BlogPost>;
}

const BlogPostsView = (props: BlogPostsViewProps) => {
    return (
        <div>
            <Heading1>Babblegraph Blog</Heading1>
            {
                props.blogPosts.map((blogPost: BlogPost) => {
                    return <BlogPostView key={blogPost.id} {...blogPost} />;
                })
            }
            {
                props.shouldShowNextButton && (
                    <PrimaryButton onClick={props.handleUpdatePageIndex}>
                        More
                    </PrimaryButton>
                )
            }
        </div>
    );
}

const BlogPostView = (props: BlogPost) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <Card onClick={() => { history.push(`/blog/${props.urlPath}`) }} className={classes.blogPostCard}>
            <Grid container>
                <Grid item xs={12} md={4}>
                    <img className={classes.heroImage}
                        src={getImageURL(props.heroImageUrl)}
                        alt={props.heroImageAltText} />
                </Grid>
                <Grid item xs={12} md={8}>
                    <Heading2>
                        {props.title}
                    </Heading2>
                    <Paragraph>
                        {props.description}
                    </Paragraph>
                </Grid>
            </Grid>
        </Card>
    );
}

export default BlogHomePage;
