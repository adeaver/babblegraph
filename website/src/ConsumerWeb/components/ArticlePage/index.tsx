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
import { Heading1, Heading3 } from 'common/typography/Heading';
import {
    PrimaryButton,
    WarningButton,
    ConfirmationButton,
} from 'common/components/Button/Button';
import { loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';
import { openLocationAsNewTab } from 'util/window/Location';
import { getStaticContentURLForPath } from 'util/static/static';

import {
    getArticleMetadata,
    GetArticleMetadataResponse,
    UpdateUserReaderTutorialResponse,
    updateUserReaderTutorial,
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
    tutorialModal: {
        position: 'absolute',
        top: '50%',
        left: '50%',
        width: '50%',
        minWidth: '300px',
        transform: 'translate(-50%, -50%)',
        maxHeight: '500px',
        overflowY: 'scroll',
    },
    tutorialImage: {
        height: 'auto',
        maxWidth: '100%',
        display: 'block',
        margin: 'auto',
    },
    tutorialContainer: {
        display: 'flex',
        justifyContent: 'center',
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
        const [ shouldShowLoginForm, setShouldShowLoginForm ] = useState<boolean>(
            props.userProfile.hasAccount && !props.userProfile.isLoggedIn
        );
        const [ shouldShowTutorial, setShouldShowTutorial ] = useState<boolean>(props.shouldShowTutorial);
        const handleCloseTutorial = () => {
            setShouldShowTutorial(false);
            updateUserReaderTutorial({
                readerToken: props.readerToken,
            },
            (resp: UpdateUserReaderTutorialResponse) => {
                // no-op
            },
            (err: Error) => {
                // no-op
            });
        }

        const [ wordReinforcementToken, setWordReinforcementToken ] = useState<string>(
            !!props.userProfile.nextTokens ? props.userProfile.nextTokens[0] : null
        );
        const [ subscriptionManagementToken, setSubscriptionManagementToken ] = useState<string>(
            !!props.userProfile.nextTokens ? props.userProfile.nextTokens[1] : null
        );

        const [ isModalLoading, setIsModalLoading ] = useState<boolean>(false);
        const [ modalError, setModalError ] = useState<Error>(null);

        const handleLogin = () => {
            setIsModalLoading(true);
            getUserProfileInformation({
                key: RouteEncryptionKey.ArticleReaderKey,
                token: props.readerToken,
                nextKeys: [RouteEncryptionKey.WordReinforcement, RouteEncryptionKey.SubscriptionManagement],
            },
            (resp: GetUserProfileInformationResponse) => {
                setIsModalLoading(false);
                setModalError(null);
                setWordReinforcementToken(resp.userProfile.nextTokens[0]);
                setSubscriptionManagementToken(resp.userProfile.nextTokens[1]);
                setShouldShowLoginForm(false);
            },
            (err: Error) => {
                setIsModalLoading(false);
                setModalError(err);
            });
        }

        const [ iframeRef, setIFrameRef ] = useState<HTMLIFrameElement>(null);
        const [ selection, setSelection ] = useState<string>(null);
        const [ isModalOpen, setIsModalOpen ] = useState<boolean>(false);

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
                        <Grid className={classes.navbarItem} item xs={4}>
                            <PrimaryButton
                                onClick={handleToggleWordSearch(true)}
                                disabled={!selection}>
                                Lookup Word
                            </PrimaryButton>
                        </Grid>
                        <Grid className={classes.navbarItem} item xs={4}>
                            <WarningButton
                                onClick={() => {openLocationAsNewTab(`/out/${props.articleId}`)}}>
                                Not Working?
                            </WarningButton>
                        </Grid>
                        <Grid className={classes.navbarItem} item xs={4}>
                            <PrimaryButton
                                onClick={() => {setShouldShowTutorial(true)}}>
                                Help!
                            </PrimaryButton>
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
                    shouldShowTutorial && (
                        <TutorialModal
                            handleCloseModal={handleCloseTutorial} />
                    )
                }
                {
                    (isModalOpen && !shouldShowTutorial) && (
                        <WordSearchModal
                            shouldShowLoginForm={shouldShowLoginForm}
                            wordReinforcementToken={wordReinforcementToken}
                            subscriptionManagementToken={subscriptionManagementToken}
                            selection={selection}
                            isLoading={isModalLoading}
                            error={modalError}
                            handleCloseModal={handleToggleWordSearch(false)}
                            handleLogin={handleLogin} />
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
    isLoading: boolean;
    selection: string;
    wordReinforcementToken: string | null;
    subscriptionManagementToken: string | null;
    error: Error;

    handleCloseModal: () => void;
    handleLogin: () => void;
}

const WordSearchModal = (props: WordSearchModalProps) => {
    const classes = styleClasses();

    const handleLogin = (_: string) => {
        props.handleLogin();
    }

    let body;
    if (props.isLoading) {
        body = <LoadingSpinner />
    } else if (!!props.error) {
        body = (
            <Heading3 color={TypographyColor.Warning}>
                Something went wrong with your request. Try again later.
            </Heading3>
        )
    } else if (props.shouldShowLoginForm) {
        body = (
            <div>
                <Heading3 color={TypographyColor.Primary}>
                    Login to Babblegraph
                </Heading3>
                <LoginForm
                    onLoginSuccess={handleLogin} />
                <CenteredComponent>
                    <WarningButton
                        className={classes.closeModalButton}
                        onClick={props.handleCloseModal}>
                        Close popup
                    </WarningButton>
                </CenteredComponent>
            </div>
        );
    } else {
        body = (
            <div>
                <WordSearchComponent
                    selection={props.selection}
                    wordReinforcementToken={props.wordReinforcementToken}
                    subscriptionManagementToken={props.subscriptionManagementToken} />
                <CenteredComponent>
                    <WarningButton
                        className={classes.closeModalButton}
                        onClick={props.handleCloseModal}>
                        Close popup
                    </WarningButton>
                </CenteredComponent>
                <Link href={`/manage/${props.wordReinforcementToken}/vocabulary`}>
                    Go to your vocabulary list
                </Link>
            </div>
        );
    }
    return (
        <Modal
            open={true}
            onClose={props.handleCloseModal}>
            <DisplayCard
                className={classes.wordSearchModal}>
                {body}
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

type TutorialModalProps = {
    handleCloseModal: () => void;
}

const TutorialModal = (props: TutorialModalProps) => {
    const classes = styleClasses();
    return (
        <Modal
            open={true}
            onClose={props.handleCloseModal}>
            <DisplayCard
                className={classes.tutorialModal}>
                <Heading1 color={TypographyColor.Primary}>
                    Welcome to the Babblegraph article reader
                </Heading1>
                <Paragraph>
                    Getting started is easy. Just highlight a word you don’t know.
                </Paragraph>
                <img className={classes.tutorialImage} src={getStaticContentURLForPath("assets/reader_tutorial_highlight.jpeg")} />
                <Paragraph>
                    Click on the "Lookup Word" button on the top of the page.
                </Paragraph>
                <CenteredComponent className={classes.tutorialContainer}>
                    <PrimaryButton>
                        Lookup Word
                    </PrimaryButton>
                </CenteredComponent>
                <Paragraph>
                    And then get the definition right here in the article!
                </Paragraph>
                <img className={classes.tutorialImage} src={getStaticContentURLForPath("assets/reader_tutorial_lookup.png")} />
                <CenteredComponent className={classes.tutorialContainer}>
                    <ConfirmationButton
                        onClick={props.handleCloseModal}>
                        Got it
                    </ConfirmationButton>
                </CenteredComponent>
            </DisplayCard>
        </Modal>
    );
}

export default ArticlePage;
