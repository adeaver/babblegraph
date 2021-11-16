import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import Card from '@material-ui/core/Card';
import CircularProgress from '@material-ui/core/CircularProgress';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import RadioGroup from '@material-ui/core/RadioGroup';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryRadio } from 'common/components/Radio/Radio';

import {
    getUserPreferencesForToken,
    GetUserPreferencesForTokenResponse,
    ReadingLevelClassificationForLanguage,
    updateUserPreferencesForToken,
    UpdateUserPreferencesForTokenResponse,
} from 'ConsumerWeb/api/user/difficultyLevel';
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
    radioController: {
        width: '100%',
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
    },
});

const radioFormOptions = [
    {
        displayText: 'Easier',
        value: 'Beginner',
    }, {
        displayText: 'Intermediate',
        value: 'Intermediate',
    }, {
        displayText: 'More Difficult',
        value: 'Advanced',
    }, {
        displayText: 'Most Difficult',
        value: 'Professional',
    }
]

type Params = {
    token: string
}

type DifficultyLevelSettingPageProps = RouteComponentProps<Params>

const DifficultyLevelSettingPage = (props: DifficultyLevelSettingPageProps) => {
    const { token } = props.match.params;

    // User Profile State
    const [ userProfile, setUserProfile ] = useState<GetUserProfileResponse | null>(null);
    const [ isLoadingUserProfile, setIsLoadingUserProfile ] = useState<boolean>(true);

    // Reading Classification State
    const [ readingLevelClassifications, setReadingLevelClassifications ] = useState<Array<ReadingLevelClassificationForLanguage>>([]);
    const [ isLoadingReadingLevelClassification, setIsLoadingReadingLevelClassification ] = useState<boolean>(true);

    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ error, setError ] = useState<Error | null>(null);
    const [ didUpdate, setDidUpdate ] = useState<boolean | null>(null);

    const handleSetReadingLevel = (newReadingLevel: string) => {
        const newClassifications = (readingLevelClassifications || []).map((classification: ReadingLevelClassificationForLanguage) => {
            if (classification.languageCode !== "es") {
                return classification;
            }
            return {
                ...classification,
                readingLevelClassification: newReadingLevel,
            }
        });
        if (!newClassifications.length) {
            newClassifications.push({
                languageCode: "es",
                readingLevelClassification: newReadingLevel,
            });
        }
        setReadingLevelClassifications(newClassifications);
    }
    const getReadingLevel = () => {
        return (readingLevelClassifications || []).reduce((accumulator: string | null, next: ReadingLevelClassificationForLanguage) => {
            if (next.languageCode !== "es") {
                return accumulator;
            }
            return next.readingLevelClassification;
        }, null);
    }
    const handleSubmit = () => {
        setIsLoadingReadingLevelClassification(true);
        setError(null);
        setDidUpdate(null);
        updateUserPreferencesForToken({
            token: token,
            emailAddress: emailAddress || '',
            classificationsByLanguage: readingLevelClassifications,
        },
        (resp: UpdateUserPreferencesForTokenResponse) => {
            setIsLoadingReadingLevelClassification(false);
            setDidUpdate(resp.didUpdate);
        },
        (e: Error) => {
            setIsLoadingReadingLevelClassification(false);
            setError(e);
        });
    }

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
        getUserPreferencesForToken({
            token: token,
        },
        (resp: GetUserPreferencesForTokenResponse) => {
            setReadingLevelClassifications(resp.classificationsByLanguage);
            setIsLoadingReadingLevelClassification(false);
        },
        (e: Error) => {
            setIsLoadingReadingLevelClassification(false);
            setError(e);
        });
    }, []);

    const classes = styleClasses();
    const isLoading = isLoadingReadingLevelClassification || isLoadingUserProfile;
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
                            Select the difficulty you think is appropriate for your reading level. When you’re done, remember to enter your email on the bottom and click ‘Update’ to complete the process.
                        </Paragraph>
                        {
                            isLoading ? (
                                <LoadingScreen />
                            ) : (
                                <ReadingLevelRadioForm
                                    value={getReadingLevel()}
                                    options={radioFormOptions}
                                    userProfile={userProfile}
                                    emailAddress={emailAddress}
                                    handleValueChange={handleSetReadingLevel}
                                    handleEmailAddressChange={setEmailAddress}
                                    handleSubmit={handleSubmit} />
                            )
                        }
                    </Card>
                    <Snackbar open={!!error} autoHideDuration={6000}>
                        <Alert severity="error">Something went wrong processing your request.</Alert>
                    </Snackbar>
                    <Snackbar open={didUpdate} autoHideDuration={6000}>
                        <Alert severity="success">Successfully updated the difficulty level.</Alert>
                    </Snackbar>
                </Grid>
            </Grid>
        </Page>
    );
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
                    Set your difficulty level
                </Paragraph>
            </Grid>
        </Grid>
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

type ReadingLevelRadioFormProps = {
    value: string | null;
    options: ReadingLevelRadioFormOption[];
    emailAddress: string;
    userProfile: GetUserProfileResponse | null;
    handleValueChange: (v: string) => void;
    handleEmailAddressChange: (v: string) => void;
    handleSubmit: () => void;
}

type ReadingLevelRadioFormOption = {
    value: string;
    displayText: string;
}

const ReadingLevelRadioForm = (props: ReadingLevelRadioFormProps) => {
    const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleValueChange((event.target as HTMLInputElement).value);
    };
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleEmailAddressChange((event.target as HTMLInputElement).value);
    };

    const classes = styleClasses();
    return (
        <div>
            <FormControl className={classes.radioController} component="fieldset">
                <RadioGroup aria-label="diffculty-level" name="diffculty-level1" value={props.value} onChange={handleRadioFormChange}>
                    <Grid container>
                    {
                        props.options.map((option: ReadingLevelRadioFormOption, idx: number) => (
                            <ReadingLevelClassificationRadioButton
                                key={`reading-level-option-${idx}`}
                                isSelected={props.value === option.value}
                                value={option.value}
                                displayText={option.displayText} />
                        ))
                    }
                    </Grid>
                </RadioGroup>
            </FormControl>
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
                            disabled={!props.emailAddress || !props.value}>
                            Submit
                        </PrimaryButton>
                    </Grid>
                </Grid>
            </form>
        </div>
    );
}


type ReadingLevelClassificationRadioButtonProps = {
    value: string;
    displayText: string;
    isSelected: boolean;
}

const ReadingLevelClassificationRadioButton = (props: ReadingLevelClassificationRadioButtonProps) => {
    return (
        <Grid item xs={6} md={3}>
            <FormControlLabel value={props.value} control={<PrimaryRadio />} label={props.displayText} />
        </Grid>
    )
}


export default DifficultyLevelSettingPage;
