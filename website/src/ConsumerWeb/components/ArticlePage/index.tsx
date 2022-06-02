import React, { useState, useEffect, useRef } from 'react';
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
        const [ selection, setSelection ] = useState<string>(null);

        const iframe = useRef<HTMLIFrameElement>(null);
        useEffect(() => {
            // @ts-ignore
            const contentElement: HTMLIFrameElement = document.getElementById("contentIFrame");
            if (!!contentElement) {
                contentElement.contentWindow.document.addEventListener('selectionchange', function() {
                    // @ts-ignore
                    const selection = this.getSelection().toString();
                    setSelection(selection);
                });
            }
        }, [iframe.current])

        console.log(`selection: ${selection}`);
        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={12}>
                    <iframe
                        ref={iframe}
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
