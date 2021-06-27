import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import { Heading1 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { withCaptchaToken, loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';

import {
    RequestPasswordResetLinkResponse,
    requestPasswordResetLink,
} from 'api/user/forgotPassword';

const styleClasses = makeStyles({
    displayCard: {
        padding: '20px',
        marginTop: '20px',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    textField: {
        width: '100%',
    },
    formGridContainer: {
        alignItems: 'center',
    },
    formGridItem: {
       padding: '5px',
    },
    iconContainer: {
        alignSelf: 'center',
    },
    warningIcon: {
        color: Color.Warning,
    },
});

type ForgotPasswordPageProps = {};

const ForgotPasswordPage = (props: ForgotPasswordPageProps) => {
    const classes = styleClasses();

    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ wasRequestResetLinkSuccessful, setWasRequestResetLinkSuccessful ] = useState<boolean | null>(null);
    const [ error, setError ] = useState<Error>(null);
    const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);

    const handleSubmit = () => {
        setIsLoading(true);
        withCaptchaToken("forgotpassword", (token: string) => {
            requestPasswordResetLink({
                emailAddress: emailAddress,
                captchaToken: token,
            },
            (resp: RequestPasswordResetLinkResponse) => {
                setIsLoading(false);
                setWasRequestResetLinkSuccessful(resp.success);
            },
            (e: Error) => {
                setIsLoading(false);
                setError(e);
            });
        });
    };

    useEffect(() => {
        loadCaptchaScript();
        setHasLoadedCaptcha(true);
    }, []);

    let resultBody = null;
    if (wasRequestResetLinkSuccessful) {
        resultBody = (
            <Paragraph color={TypographyColor.Confirmation}>
                If there is an account associated with the email address entered, you’ll receive an email with a link to reset your password in the next five minutes! Keep in mind that the link in the email is only valid for 15 minutes after receiving it.
            </Paragraph>
        );
    } else if (wasRequestResetLinkSuccessful != null && !wasRequestResetLinkSuccessful || !!error) {
        resultBody = (
            <Paragraph color={TypographyColor.Warning}>
                There was a problem making this request. Try again later or email hello@babblegraph.com for assistance.
            </Paragraph>
        );
    }

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard}>
                        <Heading1 color={TypographyColor.Primary}>
                            Reset your password
                        </Heading1>
                        <Paragraph>
                            If you have a Babblegraph Premium Account, you can use the form below to request a link to reset your password. If you don’t have a premium account and would like one, check out Babblegraph’s BuyMeACoffee page.
                        </Paragraph>
                        { resultBody }
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                <RequestResetLinkForm
                                    emailAddress={emailAddress}
                                    setEmailAddress={setEmailAddress}
                                    canSubmit={hasLoadedCaptcha && !!emailAddress}
                                    handleSubmit={handleSubmit} />
                            )
                        }
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

type RequestResetLinkFormProps = {
    emailAddress: string | null;
    setEmailAddress: (v: string) => void;
    canSubmit: boolean;
    handleSubmit: () => void;
}

const RequestResetLinkForm = (props: RequestResetLinkFormProps) => {
    const classes = styleClasses();

    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.setEmailAddress((event.target as HTMLInputElement).value);
    };
    return (
        <Grid container className={classes.formGridContainer}>
            <Grid item xs={12} md={10} className={classes.formGridItem}>
                <PrimaryTextField
                    className={classes.textField}
                    id="email"
                    label="Email Address"
                    variant="outlined"
                    defaultValue={props.emailAddress}
                    onChange={handleEmailAddressChange} />
            </Grid>
            <Grid item xs={12} md={2} className={classes.formGridItem}>
                <PrimaryButton onClick={props.handleSubmit} disabled={!props.canSubmit}>
                    Request a reset link
                </PrimaryButton>
            </Grid>
        </Grid>
    );
}

export default ForgotPasswordPage;
