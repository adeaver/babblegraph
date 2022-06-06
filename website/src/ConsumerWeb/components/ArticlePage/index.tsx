import React, { useState, useEffect, useRef, useCallback } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Modal from '@material-ui/core/Modal';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Color from 'common/styles/colors';
import Paragraph from 'common/typography/Paragraph';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Link from 'common/components/Link/Link';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading3 } from 'common/typography/Heading';
import {
    PrimaryButton,
    WarningButton,
    ConfirmationButton,
} from 'common/components/Button/Button';
import { loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';
import { openLocationAsNewTab } from 'util/window/Location';

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
import { UserVocabularyEntry } from 'ConsumerWeb/api/user/userVocabulary';
import {
    RouteEncryptionKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import LoginForm from 'ConsumerWeb/components/UserAccounts/LoginForm';
import WordSearchDisplay from 'ConsumerWeb/components/WordReinforcementPage/WordSearchDisplay';
import {
    withUserVocabulary,
    InjectedUserVocabularyComponentProps,
} from 'ConsumerWeb/components/WordReinforcementPage/withUserVocabulary';

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
    storyLoadingSpinner: {
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        zIndex: 1,
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
        width: '50%',
        minWidth: '300px',
        transform: 'translate(-50%, -50%)',
        maxHeight: '500px',
        overflowY: 'scroll',
    },
    closeModalButton: {
        width: '100%',
    },
});

type Params = {
    token: string;
}

type ArticlePageAPIProps = GetArticleMetadataResponse & GetUserProfileInformationResponse;
type ArticlePageOwnProps = RouteComponentProps<Params>;

const ArticlePage = asBaseComponent(
    (props: ArticlePageOwnProps & BaseComponentProps & ArticlePageAPIProps) => {
        const [ iframeRef, setIFrameRef ] = useState<HTMLIFrameElement>(null);
        const [ selection, setSelection ] = useState<string>(null);
        const [ isModalOpen, setIsModalOpen ] = useState<boolean>(false);

        const [ shouldShowLoginForm, setShouldShowLoginForm ] = useState<boolean>(
            props.userProfile.hasAccount && !props.userProfile.isLoggedIn
        );

        const handleToggleWordSearch = (isOpen: boolean) => {
            return () => {
                setIsModalOpen(isOpen);
            }
        }

        const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);
        const [ hasIFrameLoaded, setHasIFrameLoaded ] = useState<boolean>(true);

        useEffect(() => {
            if (!hasLoadedCaptcha) {
                loadCaptchaScript();
                setHasLoadedCaptcha(true);
            }
            !!iframeRef && iframeRef.addEventListener('load', function() {
                setHasIFrameLoaded(false);

                const iframeAnchors = iframeRef.contentWindow.document.getElementsByTagName("a");
                Object.values(iframeAnchors).forEach((anchor: HTMLAnchorElement) => {
                    anchor.onclick = function(e: MouseEvent) {
                        e.preventDefault();
                        window.open(
                            (e.target as HTMLAnchorElement).href,
                            '_blank'
                        );
                    }
                });

                iframeRef.contentWindow.document.addEventListener('selectionchange', function() {
                    setSelection(this.getSelection().toString());
                })
            });
            return () => {
                !!iframeRef && iframeRef.removeEventListener('selectionchange', function() {})
            }
        }, [iframeRef])

        const classes = styleClasses();
        if (!hasLoadedCaptcha) {
            return <LoadingSpinner />;
        }
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
                            <WarningButton
                                onClick={() => {openLocationAsNewTab(`/out/${props.articleId}`)}}>
                                Not Working?
                            </WarningButton>
                        </Grid>
                    </Grid>
                </Grid>
                {
                    hasIFrameLoaded && <LoadingSpinner className={classes.storyLoadingSpinner} />
                }
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
                            shouldShowLoginForm={shouldShowLoginForm}
                            wordReinforcementToken={props.userProfile.nextTokens[0]}
                            subscriptionManagementToken={props.userProfile.nextTokens[1]}
                            selection={selection}
                            handleCloseModal={handleToggleWordSearch(false)}
                            handleToggleLoginForm={setShouldShowLoginForm} />
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
            getUserProfileInformation({
                key: RouteEncryptionKey.ArticleReaderKey,
                token: resp.readerToken,
                nextKeys: [RouteEncryptionKey.WordReinforcement, RouteEncryptionKey.SubscriptionManagement],
            },
            (resp2: GetUserProfileInformationResponse) => {
                onSuccess({
                    ...resp,
                    ...resp2
                });
            },
            onError)
        },
        onError);
    },
    true,
);

type WordSearchModalProps = {
    shouldShowLoginForm: boolean;
    selection: string;
    wordReinforcementToken: string;
    subscriptionManagementToken: string;

    handleCloseModal: () => void;
    handleToggleLoginForm: (v: boolean) => void;
}

const WordSearchModal = (props: WordSearchModalProps) => {
    const handleLogin = (_: string) => {
        props.handleToggleLoginForm(false);
    }
    const classes = styleClasses();
    return (
        <Modal
            open={true}
            onClose={props.handleCloseModal}>
            <DisplayCard
                className={classes.wordSearchModal}>
                {
                    props.shouldShowLoginForm ? (
                        <div>
                            <Heading3 color={TypographyColor.Primary}>
                                Login to Babblegraph
                            </Heading3>
                            <LoginForm
                                onLoginSuccess={handleLogin} />
                        </div>
                    ) : (
                        <WordSearchComponent
                            selection={props.selection}
                            wordReinforcementToken={props.wordReinforcementToken}
                            subscriptionManagementToken={props.subscriptionManagementToken} />
                    )
                }
                <CenteredComponent>
                    <WarningButton
                        className={classes.closeModalButton}
                        onClick={props.handleCloseModal}>
                        Close popup
                    </WarningButton>
                </CenteredComponent>
                {
                    !props.shouldShowLoginForm && (
                        <Link href={`/manage/${props.wordReinforcementToken}/vocabulary`}>
                            Go to your vocabulary list
                        </Link>
                    )
                }
            </DisplayCard>
        </Modal>
    );
}

type WordSearchComponentProps = {
    selection: string;
    wordReinforcementToken: string;
    subscriptionManagementToken: string;
}

const WordSearchComponent = withUserVocabulary(
    (props: WordSearchComponentProps & InjectedUserVocabularyComponentProps) => (
        <WordSearchDisplay
            searchTerms={props.selection.trim().split(/ +/g)}
            wordReinforcementToken={props.wordReinforcementToken}
            subscriptionManagementToken={props.subscriptionManagementToken}
            userVocabularyEntries={props.userVocabularyEntries}
            handleAddNewUserVocabularyEntry={props.handleAddNewVocabularyEntry} />
    )
)

export default ArticlePage;
