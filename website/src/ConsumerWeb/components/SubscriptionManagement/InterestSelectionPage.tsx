import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Divider from '@material-ui/core/Divider';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import {
    contentTopicDisplayMappings,
    ContentTopicDisplayMapping,
    getUserContentTopicsForToken,
    GetUserContentTopicsForTokenResponse,

    updateUserContentTopicsForToken,
    UpdateUserContentTopicsForTokenResponse,
} from 'ConsumerWeb/api/user/contentTopics';
import {
    getUserProfile,
    GetUserProfileResponse
} from 'ConsumerWeb/api/useraccounts/useraccounts';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
    },
    contentHeaderBackArrow: {
        alignSelf: 'center',
        cursor: 'pointer',
    },
    loadingSpinner: {
        color: Color.Primary,
        display: 'block',
        margin: 'auto',
    },
    submitButton: {
        display: 'block',
        margin: 'auto',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    emailField: {
        width: '100%',
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
});

type Params = {
    token: string
}

type InterestSelectionPageProps = RouteComponentProps<Params>

const InterestSelectionPage = (props: InterestSelectionPageProps) => {
    const classes = styleClasses();
    const { token } = props.match.params;

    // User Profile State
    const [ userProfile, setUserProfile ] = useState<GetUserProfileResponse | null>(null);
    const [ isLoadingUserProfile, setIsLoadingUserProfile ] = useState<boolean>(true);

    // Content Topics State
    const [ selectedContentTopics, setSelectedContentTopics ] = useState<Object>({});
    const [ isLoadingContentTopics, setIsLoadingContentTopics ] = useState<boolean>(true);

    const [ emailAddress, setEmailAddress ] = useState<string>('');
    const [ didUpdate, setDidUpdate ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    useEffect(() => {
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setIsLoadingUserProfile(false);
            if (!!resp.emailAddress) {
                setEmailAddress(resp.emailAddress);
                setUserProfile(resp);
            }
        },
        (e: Error) => {
            setIsLoadingUserProfile(false);
            setError(e);
        });
        getUserContentTopicsForToken({
            token: token,
        },
        (resp: GetUserContentTopicsForTokenResponse) => {
            setIsLoadingContentTopics(false);
            setSelectedContentTopics((resp.contentTopics || []).reduce((accumulator: Object, next: string) => ({
                ...accumulator,
                [next]: true,
            }), {}));
        },
        (e: Error) => {
            setIsLoadingContentTopics(false);
            setError(e);
        });
    }, []);

    const handleSelectContentTopicMapping = (apiValue: string[]) => {
        setSelectedContentTopics(apiValue.reduce((accumulator: Object, next: string) => ({
            ...accumulator,
            [next]: !accumulator[next],
        }), selectedContentTopics));
    };
    const handleSubmit = () => {
        setIsLoadingContentTopics(true);
        setDidUpdate(false);
        setError(null);
        updateUserContentTopicsForToken({
            token: token,
            emailAddress: emailAddress,
            contentTopics: Object.keys(selectedContentTopics).filter((key: string) => !!selectedContentTopics[key]),
        },
        (resp: UpdateUserContentTopicsForTokenResponse) => {
            setIsLoadingContentTopics(false);
            setDidUpdate(true);
        },
        (e: Error) => {
            setIsLoadingContentTopics(false);
            setError(e);
        });
    }

    const isLoading = isLoadingContentTopics || isLoadingUserProfile;
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <ContentHeader token={token} />
                        <Divider />
                        <Paragraph size={Size.Medium} align={Alignment.Left}>
                            Click on the topics that interest you to receive emails with content on that topic. You can select as many as you’d like. When you’re done, enter your email at the bottom and click ‘Update’ to complete the process. Not every email will contain content with all the topics you’ve picked.
                        </Paragraph>
                        {
                            isLoading ? (
                                <LoadingScreen />
                            ) : (
                                <ContentTopicSelectionForm
                                    contentTopicDisplayMappings={contentTopicDisplayMappings.map((m: ContentTopicDisplayMapping) => ({
                                            ...m,
                                            isChecked: m.apiValue.reduce((val: boolean, next: string) => (
                                                val || !!selectedContentTopics[next]
                                            ), false),
                                        })
                                    )}
                                    userProfile={userProfile}
                                    emailAddress={emailAddress}
                                    handleEmailAddressChange={setEmailAddress}
                                    handleSelectContentTopicMapping={handleSelectContentTopicMapping}
                                    handleSubmit={handleSubmit} />
                            )
                        }
                        <Snackbar open={!!error} autoHideDuration={6000}>
                            <Alert severity="error">Something went wrong processing your request.</Alert>
                        </Snackbar>
                        <Snackbar open={didUpdate} autoHideDuration={6000}>
                            <Alert severity="success">Successfully updated your email topic interests.</Alert>
                        </Snackbar>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

const LoadingScreen = () => {
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <CircularProgress className={classes.loadingSpinner} />
                <Paragraph size={Size.Medium} align={Alignment.Center}>
                    Loading, please wait.
                </Paragraph>
            </Grid>
        </Grid>
    )
}

type ContentHeaderProps = {
    token: string;
}

const ContentHeader = (props: ContentHeaderProps) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <Grid container>
            <Grid className={classes.contentHeaderBackArrow} onClick={() => history.push(`/manage/${props.token}`)} item xs={1}>
                <ArrowBackIcon color='action' />
            </Grid>
            <Grid item xs={11}>
                <Paragraph size={Size.Large} color={TypographyColor.Primary} align={Alignment.Left}>
                    Manage Your Interests
                </Paragraph>
            </Grid>
        </Grid>
    );
}

type ContentTopicDisplayMappingWithChecked = {
    isChecked: boolean;
} & ContentTopicDisplayMapping;

type ContentTopicSelectionFormProps = {
    contentTopicDisplayMappings: Array<ContentTopicDisplayMappingWithChecked>;
    emailAddress: string;
    userProfile: GetUserProfileResponse | null;
    handleSelectContentTopicMapping: (v: string[]) => void;
    handleEmailAddressChange: (v: string) => void;
    handleSubmit: () => void;
}

const ContentTopicSelectionForm = (props: ContentTopicSelectionFormProps) => {
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleEmailAddressChange((event.target as HTMLInputElement).value);
    };

    const classes = styleClasses();
    return (
        <FormGroup>
            <Grid container>
                {
                    props.contentTopicDisplayMappings.map((mapping: ContentTopicDisplayMappingWithChecked, idx: number) => (
                        <Grid item key={`contentTopicGridItem-${idx}`} xs={6} md={4}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        checked={mapping.isChecked}
                                        onChange={() => { props.handleSelectContentTopicMapping(mapping.apiValue) }}
                                        name={`checkbox-${mapping.displayText}`} />
                                }
                                label={mapping.displayText} />
                        </Grid>
                    ))
                }
            </Grid>
            <Divider />
            <form className={classes.confirmationForm} noValidate autoComplete="off">
                <Grid container>
                    {
                        !!props.userProfile ? (
                            <Grid item xs={4} md={5}>
                                &nbsp;
                            </Grid>
                        ) : (
                            <Grid item xs={8} md={10}>
                                <PrimaryTextField
                                    id="email"
                                    className={classes.emailField}
                                    label="Email Address"
                                    variant="outlined"
                                    onChange={handleEmailAddressChange} />
                            </Grid>
                        )
                    }
                    <Grid item xs={4} md={2} className={classes.submitButtonContainer}>
                        <PrimaryButton
                            className={classes.submitButton}
                            onClick={props.handleSubmit}
                            disabled={!props.emailAddress}>
                            Submit
                        </PrimaryButton>
                    </Grid>
                </Grid>
            </form>
        </FormGroup>
    );
}

export default InterestSelectionPage;
