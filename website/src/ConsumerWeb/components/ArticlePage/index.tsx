import React, { useState, useEffect, useRef, useCallback } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    content: {
        position: 'absolute',
        left: '0',
        top: '80px',
        width: '100%',
        minHeight: '100%',
    },
});

type Params = {
    token: string;
}

type ArticlePageAPIProps = {}
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
                <Grid item xs={12}>
                    <iframe
                        ref={setIFrameRef}
                        id="contentIFrame"
                        className={classes.content}
                        src="/a/abc" />
                </Grid>
            </Grid>
        );
    },
    (
        ownProps: ArticlePageOwnProps,
        onSuccess: (resp: ArticlePageAPIProps) => void,
        onError: (err: Error) => void,
    ) => onSuccess({}),
    true,
);

export default ArticlePage;
