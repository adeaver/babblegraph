import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Paragraph from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { setLocation } from 'util/window/Location';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Color from 'common/styles/colors';
import Form from 'common/components/Form/Form';

import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';

import {
    createUser,
    CreateUserError,
    CreateUserResponse,
} from 'ConsumerWeb/api/useraccounts2/useraccounts';

const styleClasses = makeStyles({
    formContainer: {
        display: 'flex',
        flexDirection: 'column',
    },
    createUserFormTextField: {
        margin: '10px 0',
    },
    loadingSpinner: {
        color: Color.TextGray,
        display: 'block',
        margin: '0 15px',
        // Hello, Jankiness, my old friend
        height: '0.875rem !important',
        width: '0.875rem !important',
    },
});

const createUserErrorMessages = {
   [CreateUserError.AlreadyExists]: "There’s already an existing account for that email address",
   [CreateUserError.InvalidToken]: "The email submitted didn’t match the email address this unique link is for. Make sure you entered the same email address that you received the signup link with.",
   [CreateUserError.PasswordRequirements]: "The password entered did not match the minimum password requirements",
   [CreateUserError.PasswordsNoMatch]: "The passwords entered did not match.",
   "default": "Something went wrong processing your request. Try again, or email hello@babblegraph.com for help.",
}

type Params = {
    token: string;
}

type CreateUserAccountPageOwnProps = RouteComponentProps<Params>;

const CreateUserAccountPage = withUserProfileInformation<CreateUserAccountPageOwnProps>(
    RouteEncryptionKey.CreateUser,
    [RouteEncryptionKey.SubscriptionManagement, RouteEncryptionKey.PremiumSubscriptionCheckout],
    (ownProps: CreateUserAccountPageOwnProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.CheckoutPage,
    (props: CreateUserAccountPageOwnProps & UserProfileComponentProps) => {
        const [ subscriptionManagementToken, premiumSubscriptionCheckoutToken ] = props.userProfile.nextTokens;
        if (!!props.userProfile.subscriptionLevel) {
            setLocation(`/manage/${subscriptionManagementToken}`);
            return <div />;
        } else if (props.userProfile.hasAccount) {
            setLocation(`/checkout/${premiumSubscriptionCheckoutToken}`);
            return <div />;
        }

        const [ emailAddress, setEmailAddress ] = useState<string>(null);
        const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setEmailAddress((event.target as HTMLInputElement).value);
        }
        const [ password, setPassword ] = useState<string>(null);
        const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setPassword((event.target as HTMLInputElement).value);
        }
        const [ confirmPassword, setConfirmPassword ] = useState<string>(null);
        const handleConfirmPasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setConfirmPassword((event.target as HTMLInputElement).value);
        }

        const [ isLoading, setIsLoading ] = useState<boolean>(false);
        const [ errorMessage, setErrorMessage ] = useState<string>(null);


        const handleSubmit = () => {
            setIsLoading(true);
            createUser({
                createUserToken: props.match.params.token,
                emailAddress: emailAddress,
                password: password,
                confirmPassword: confirmPassword,
            },
            (resp: CreateUserResponse) => {
                setIsLoading(false);
                if (!resp.createUserError) {
                    setLocation(`/checkout/${premiumSubscriptionCheckoutToken}`);
                    return;
                }
                setErrorMessage(createUserErrorMessages[resp.createUserError] || createUserErrorMessages["default"]);
            },
            (err: Error) => {
                setIsLoading(false);
                setErrorMessage(createUserErrorMessages["default"]);
            });
        }

        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <DisplayCardHeader
                            title="Create an account"
                            backArrowDestination={`/manage/${subscriptionManagementToken}`} />
                            <Paragraph align={Alignment.Left}>
                                First step, sign up for a Babblegraph account
                            </Paragraph>
                            <Paragraph align={Alignment.Left}>
                                Why do you need to sign up for an account to access Babblegraph Premium?
                            </Paragraph>
                            <Paragraph align={Alignment.Left}>
                                Security is a big concern when dealing with payment information. Accounts are more secure than managing your Babblegraph subscription.
                            </Paragraph>
                            <Form
                                className={classes.formContainer}
                                handleSubmit={handleSubmit}>
                                <PrimaryTextField
                                    className={classes.createUserFormTextField}
                                    id="email"
                                    label="Confirm Your Email Address"
                                    variant="outlined"
                                    defaultValue={emailAddress}
                                    disabled={isLoading}
                                    onChange={handleEmailAddressChange} />
                                <Paragraph align={Alignment.Left}>
                                    Password Requirements:
                                </Paragraph>
                                <ul>
                                    <PasswordConstraint isConstraintMet={password && password.length > 8}>
                                        At least 8 characters
                                    </PasswordConstraint>
                                    <PasswordConstraint isConstraintMet={password && password.length < 32} >
                                        No more than 32 characters
                                    </PasswordConstraint>
                                    <PasswordConstraint isConstraintMet={false}>
                                        At least three of the following:
                                    </PasswordConstraint>
                                        <ul>
                                            <PasswordConstraint isConstraintMet={password && !!password.match(/[a-z]/)}>
                                                Lower Case Latin Letter (a-z)
                                            </PasswordConstraint>
                                            <PasswordConstraint isConstraintMet={password && !!password.match(/[A-Z]/)}>
                                                Upper Case Latin Letter (A-Z)
                                            </PasswordConstraint>
                                            <PasswordConstraint isConstraintMet={password && !!password.match(/[0-9]/)}>
                                                Number (0-9)
                                            </PasswordConstraint>
                                            <PasswordConstraint isConstraintMet={password && !!password.match(/[^0-9a-zA-Z]/)}>
                                                Special Character (such as !@#$%^&*)
                                            </PasswordConstraint>
                                        </ul>
                                </ul>
                                <PrimaryTextField
                                    className={classes.createUserFormTextField}
                                    id="password"
                                    label="Password"
                                    type="password"
                                    variant="outlined"
                                    defaultValue={password}
                                    disabled={isLoading}
                                    onChange={handlePasswordChange} />
                                <PrimaryTextField
                                    className={classes.createUserFormTextField}
                                    id="confirm-password"
                                    label="Confirm Password"
                                    type="password"
                                    variant="outlined"
                                    disabled={isLoading}
                                    defaultValue={confirmPassword}
                                    onChange={handleConfirmPasswordChange} />
                                <PrimaryButton
                                    className={classes.createUserFormSubmitButton}
                                    type='submit'
                                    disabled={!emailAddress || !password || !confirmPassword || isLoading}>
                                    {
                                        isLoading && (
                                            <CircularProgress className={classes.loadingSpinner} />
                                        )
                                    }
                                    Sign Up
                                </PrimaryButton>
                            </Form>
                            <Snackbar open={!!errorMessage} autoHideDuration={6000} onClose={() => {setErrorMessage(null)}}>
                                <Alert severity="error">{errorMessage}</Alert>
                            </Snackbar>
                    </DisplayCard>
                </Grid>
            </Grid>
        );
    }
);

type PasswordConstraintProps = {
    isConstraintMet: boolean;
    children: string | JSX.Element;
}

const PasswordConstraint = (props: PasswordConstraintProps) => {
    return (
        <li>
            <Paragraph
                color={props.isConstraintMet ? TypographyColor.Confirmation : TypographyColor.Gray}
                align={Alignment.Left}>
                {props.children}
            </Paragraph>
        </li>
    )
}


export default CreateUserAccountPage;
