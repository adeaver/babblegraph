import React, { useState, useEffect, useRef, useCallback } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Modal from '@material-ui/core/Modal';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
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
    UserProfileInformationError,
    UserProfileInformation,
    GetUserProfileInformationResponse,
    getUserProfileInformation,
} from 'ConsumerWeb/api/useraccounts2/useraccounts';
import {
    RouteEncryptionKey,
} from 'ConsumerWeb/api/routes/consts';


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
    wordSearchModal: {
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
    }
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
        const [ isModalOpen, setIsModalOpen ] = useState<boolean>(false);

        const handleToggleWordSearch = (isOpen: boolean) => {
            return () => {
                setIsModalOpen(isOpen);
            }
        }

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
                                onClick={handleToggleWordSearch(true)}
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
                {
                    !!isModalOpen && (
                        <WordSearchModal
                            readerToken={props.readerToken}
                            selection={selection}
                            handleCloseModal={handleToggleWordSearch(false)} />
                    )
                }
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

type WordSearchModalOwnProps = {
    readerToken: string;
    selection: string;

    handleCloseModal: () => void;
}

type WordSearchModalAPIProps = GetUserProfileInformationResponse;

const WordSearchModal = asBaseComponent(
    (props: BaseComponentProps & WordSearchModalAPIProps & WordSearchModalOwnProps) => {
        const classes = styleClasses();
        return (
            <Modal
                open={true}
                onClose={props.handleCloseModal}>
                <DisplayCard
                    className={classes.wordSearchModal}>
                    { props.selection }
                </DisplayCard>
            </Modal>
        );
    },
    (
        ownProps: WordSearchModalOwnProps,
        onSuccess: (resp: WordSearchModalAPIProps) => void,
        onError: (err: Error) => void
    ) => {
        getUserProfileInformation({
            key: RouteEncryptionKey.ArticleReaderKey,
            token: ownProps.readerToken,
            nextKeys: [RouteEncryptionKey.WordReinforcement],
        },
        (resp: GetUserProfileInformationResponse) => {
            onSuccess(resp);
        },
        onError)
    },
    false
)

export default ArticlePage;
