import React, { useState, useEffect, useRef, useCallback } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Color from 'common/styles/colors';
import Paragraph from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import {
    PrimaryButton,
    WarningButton,
    ConfirmationButton,
} from 'common/components/Button/Button';

import {
    getArticleMetadata,
    GetArticleMetadataResponse,
} from 'ConsumerWeb/api/article';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    content: {
        border: 'none',
        borderTop: `2px solid ${Color.Primary}`,
        position: 'absolute',
        left: '0',
        top: '140px',
        width: '100%',
        minHeight: 'calc(100vh - 145px)',
    },
    navbar: {
        position: 'absolute',
        left: '0',
        top: '80px',
        width: '100%',
        height: '60px',
    },
    navbarContainer: {
        height: '100%',
    },
    navbarItem: {
        height: '100%',
        display: 'flex',
        padding: '5px',
        alignItems: 'center',
        justifyContent: 'center',
    },
});

type Params = {
    token: string;
}

type ArticlePageAPIProps = GetArticleMetadataResponse;
type ArticlePageOwnProps = RouteComponentProps<Params>;

const ArticlePage = asBaseComponent(
    (props: ArticlePageOwnProps & BaseComponentProps & ArticlePageAPIProps) => {
        const [ iframeRef, setIFrameRef ] = useState<HTMLIFrameElement>(null);

        const [ selection, setSelection ] = useState<string>(null);

        useEffect(() => {
            !!iframeRef && iframeRef.addEventListener('load', function() {
                iframeRef.contentWindow.document.addEventListener('selectionchange', function() {
                    setSelection(this.getSelection().toString());
                })
            });
            return () => {
                !!iframeRef && iframeRef.removeEventListener('selectionchange', function() {})
            }
        }, [iframeRef])

        console.log(`selection: ${selection}`);
        const classes = styleClasses();
        return (
            <Grid container>
                <Grid className={classes.navbar} item xs={12}>
                    <Grid className={classes.navbarContainer} container>
                        <Grid className={classes.navbarItem} item xs={6}>
                            <PrimaryButton
                                disabled={!selection}>
                                Lookup Word
                            </PrimaryButton>
                        </Grid>
                        <Grid className={classes.navbarItem} item xs={6}>
                            <WarningButton>
                                Not Working?
                            </WarningButton>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item xs={12}>
                    <iframe
                        ref={setIFrameRef}
                        id="contentIFrame"
                        className={classes.content}
                        src={`/a/${props.articleId}`} />
                </Grid>
            </Grid>
        );
    },
    (
        ownProps: ArticlePageOwnProps,
        onSuccess: (resp: ArticlePageAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getArticleMetadata({
            articleToken: ownProps.match.params.token,
        },
        (resp: GetArticleMetadataResponse) => {
            onSuccess(resp);
        },
        onError);
    },
    true,
);

export default ArticlePage;
