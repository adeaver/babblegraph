import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { TypographyColor } from 'common/typography/common';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';

import {
    UserBillingInformation,

    GetBillingInformationForEmailAddressResponse,
    getBillingInformationForEmailAddress,
} from 'AdminWeb/api/billing/billing';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
    },
});

const BillingManagementPage = () => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const [ userBillingInformation, setUserBillingInformation ] = useState<UserBillingInformation>(null);

    const [ emailAddress, setEmailAddress ] = useState<string>(null);
    const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmailAddress((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        setIsLoading(true);
        getBillingInformationForEmailAddress({
            emailAddress: emailAddress,
        },
        (resp: GetBillingInformationForEmailAddressResponse) => {
            setIsLoading(false);
            setUserBillingInformation(resp.billingInformation);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }

    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Warning}>
                An error occurred
            </Heading3>
        );
    } else if (userBillingInformation !== null) {
        body = (
            <UserBillingInformationDisplay
                userBillingInformation={userBillingInformation} />
        );
    }

    const classes = styleClasses();
    return (
        <Page>
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Billing Manager"
                        backArrowDestination="/ops/dashboard" />
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                            <Grid item xs={12}>
                                <PrimaryTextField
                                    id="email-address"
                                    className={classes.formComponent}
                                    label="Email Address"
                                    variant="outlined"
                                    defaultValue={emailAddress}
                                    onChange={handleEmailAddressChange} />
                            </Grid>
                            <Grid item xs={6}>
                                <PrimaryButton
                                    className={classes.formComponent}
                                    disabled={!emailAddress}
                                    type="submit">
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                </DisplayCard>
            </CenteredComponent>
            {body}
        </Page>
    );
}

type UserBillingInformationDisplayProps = {
    userBillingInformation: UserBillingInformation | undefined;
}

const UserBillingInformationDisplay = (props: UserBillingInformationDisplayProps) => {
    return (
        <CenteredComponent>
            <DisplayCard>
                <Heading3 color={!!props.userBillingInformation.userAccountStatus ? TypographyColor.Primary : TypographyColor.Gray}>
                    Subscription Status: {props.userBillingInformation.userAccountStatus || "inactive"}
                </Heading3>
                <Paragraph>
                    Account Type: {props.userBillingInformation.externalIdType}
                </Paragraph>
            </DisplayCard>
        </CenteredComponent>
    );
}

export default BillingManagementPage;
