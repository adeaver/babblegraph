import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import Snackbar from '@material-ui/core/Snackbar';

import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Alert from 'common/components/Alert/Alert';

import {
    UserScheduleError,
    getUserSchedule,
    GetUserScheduleResponse,
    updateUserSchedule,
    UpdateUserScheduleResponse,
} from 'ConsumerWeb/api/user/schedule';
import {
    getUserProfile,
    GetUserProfileResponse
} from 'ConsumerWeb/api/useraccounts/useraccounts';

import TimeSelector from './TimeSelector';

const errorMessages: { [k: string]: string } = {
    [UserScheduleError.InvalidUser]:  "Your session may be invalid. Try logging out and logging back in",
    [UserScheduleError.InvalidEmailAddress]: "Incorrect email address. Make sure it’s spelled right or try clicking the manage link from another newsletter email.",
    [UserScheduleError.UnsupportedLanguage]: "Babblegraph currently doesn’t support that language",
    [UserScheduleError.UnsupportedTimezone]: "That timezone is not supported by Babblegraph",
    [UserScheduleError.InvalidTime]: "That time is invalid.",
    [UserScheduleError.InvalidSettings]: "Your desired settings are invalid. Follow hints on the inputs to make sure they’re valid.",
    "other": "Something went wrong processing your request",
}

type Params = {
    token: string;
}

const styleClasses = makeStyles({
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

type SchedulePageProps = RouteComponentProps<Params>

const SchedulePage = (props: SchedulePageProps) => {
    const { token } = props.match.params;

    const [ initialIANATimezone, setInitialIANATimezone ] = useState<string>(Intl.DateTimeFormat().resolvedOptions().timeZone || "America/New_York");
    const [ ianaTimezone, setIANATimezone ] = useState<string>(Intl.DateTimeFormat().resolvedOptions().timeZone || "America/New_York");
    const [ hourIndex, setHourIndex ] = useState<number>(7);
    const [ quarterHourIndex, setQuarterHourIndex ] = useState<number>(0);

    const [ hasSubscription, setHasSubscription ] = useState<boolean>(false);
    const [ emailAddress, setEmailAddress ] = useState<string>(null);

    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ error, setError ] = useState<string>(null);
    const [ success, setSuccess ] = useState<boolean>(false);

    useEffect(() => {
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setHasSubscription(!!resp.subscriptionLevel);
            resp.emailAddress && setEmailAddress(resp.emailAddress);
            getUserSchedule({
                token: token,
                // TODO(multiple-languages): make this dynamic
                languageCode: "es",
            },
            (resp: GetUserScheduleResponse) => {
                setInitialIANATimezone(resp.userIanaTimezone);
                setIANATimezone(resp.userIanaTimezone);
                setHourIndex(resp.hourIndex);
                setQuarterHourIndex(resp.quarterHourIndex);
                setIsLoading(false);
            },
            (err: Error) => {
                setIsLoading(false);
                setError(errorMessages["other"]);
            });
        },
        (err: Error) => {
            setIsLoading(false);
            setError(errorMessages["other"]);
        });
    }, []);

    const handleSubmit = () => {
        if (hasSubscription) {
            console.log("Will update everything");
        } else {
            setIsLoading(true);
            updateUserSchedule({
                emailAddress: emailAddress,
                token: token,
                // TODO(multiple-languages): Make this dynamic
                languageCode: "es",
                hourIndex: hourIndex,
                quarterHourIndex: quarterHourIndex / 15,
                ianaTimezone: ianaTimezone,
            },
            (resp: UpdateUserScheduleResponse) => {
                setIsLoading(false);
                if (resp.error) {
                    setError(errorMessages[resp.error] || errorMessages["other"]);
                } else {
                    setInitialIANATimezone(ianaTimezone);
                    setSuccess(true);
                }
            },
            (err: Error) => {
                setIsLoading(false);
                setError(errorMessages["other"]);
            });
        }
    }

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                <div>
                                    <DisplayCardHeader
                                        title={hasSubscription ? "Newsletter Schedule and Customization" : "Newsletter Schedule"}
                                        backArrowDestination={`/manage/${token}`} />
                                    <Divider />
                                    <TimeSelector
                                        initialIANATimezone={initialIANATimezone}
                                        ianaTimezone={ianaTimezone}
                                        hourIndex={hourIndex}
                                        quarterHourIndex={quarterHourIndex}
                                        handleUpdateIANATimezone={setIANATimezone}
                                        handleUpdateHourIndex={setHourIndex}
                                        handleUpdateQuarterHourIndex={setQuarterHourIndex} />
                                    <ConfirmationForm
                                        emailAddress={emailAddress}
                                        userHasSubscription={hasSubscription}
                                        handleEmailAddressChange={setEmailAddress}
                                        handleSubmit={handleSubmit} />
                                </div>
                            )
                        }
                    </DisplayCard>
                    <Snackbar open={!!error} autoHideDuration={6000} onClose={() => setError(null)}>
                        <Alert severity="error">{error}</Alert>
                    </Snackbar>
                        <Snackbar open={success} autoHideDuration={6000} onClose={() => setSuccess(false)}>
                            <Alert severity="success">Successfully updated your schedule. Changes may take up to 24 hours to take effect!</Alert>
                        </Snackbar>
                </Grid>
            </Grid>
        </Page>
    );
}

type ConfirmationFormProps = {
    emailAddress: string;
    userHasSubscription: boolean;

    handleEmailAddressChange: (emailAddress: string) => void;
    handleSubmit: () => void;
}

const ConfirmationForm = (props: ConfirmationFormProps) => {
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleEmailAddressChange((event.target as HTMLInputElement).value);
    }
    const classes = styleClasses();
    return (
        <form className={classes.confirmationForm} noValidate autoComplete="off">
            <Grid container>
                {
                    props.userHasSubscription ? (
                        <Grid item xs={4} md={5}>
                            &nbsp;
                        </Grid>
                    ) : (
                        <Grid item xs={8} md={10}>
                            <PrimaryTextField
                                id="email"
                                className={classes.emailField}
                                value={props.emailAddress}
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
    )
}

export default SchedulePage;
