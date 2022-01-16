import React, { useState, useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';

import Paragraph from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import BlogDisplay from 'common/components/BlogDisplay/BlogDisplay';
import Link from 'common/components/Link/Link';
import SignupForm from 'ConsumerWeb/components/common/SignupForm/SignupForm';
import { loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';

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

    const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);
    const [ isLoadingSignup, setIsLoadingSignup ] = useState<boolean>(false);
    const [ successfullySignedUp, setSuccessfullySignedUp ] = useState<boolean>(false);

    const handleSignupSuccess = (emailAddress: string) => {
        setSuccessfullySignedUp(true);
    }

    useEffect(() => {
        loadCaptchaScript();
        setHasLoadedCaptcha(true);
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
                            title="Todos los articúlos"
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
                                    <BlogDisplay
                                        content={blogContent}
                                        metadata={blogMetadata} />
                                )
                            )
                        }
                        <div>
                            <Heading3 color={TypographyColor.Primary}>
                                Want to practice more advanced vocabulary without using flash cards?
                            </Heading3>
                            <Paragraph>
                                Babblegraph helps intermediate and advanced Spanish students effortlessly get a daily dose of Spanish practice. With Babblegraph you’ll receive a daily newsletter with news articles from trusted, Spanish-language news sources from Spain and Latin America. Sign up for free today!
                            </Paragraph>
                            <SignupForm
                                disabled={isLoadingSignup || !hasLoadedCaptcha}
                                setIsLoading={setIsLoadingSignup}
                                onSuccess={handleSignupSuccess}
                                shouldShowVerificationForm={successfullySignedUp} />
                            {
                                isLoadingSignup && <LoadingSpinner />
                            }
                            <Link href="/">
                                Click here to learn more
                            </Link>
                        </div>
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default BlogPostPage;
