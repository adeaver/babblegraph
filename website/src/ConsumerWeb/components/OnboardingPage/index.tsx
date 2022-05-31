import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';

import {
    OnboardingStatus,
    getOnboardingStatusForUser,
    GetOnboardingStatusForUserResponse,
} from 'ConsumerWeb/api/onboarding';

import InterestSelector from 'ConsumerWeb/components/InterestSelectionPage/InterestSelector';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    submitButton: {
        display: 'block',
        margin: 'auto',
    },
});

type Params = {
    token: string;
}

type OnboardingPageAPIProps = GetOnboardingStatusForUserResponse;

type OnboardingPageOwnProps = RouteComponentProps<Params>;

const OnboardingPage = asBaseComponent(
    (props: OnboardingPageOwnProps & BaseComponentProps & OnboardingPageAPIProps) => {
        const [ onboardingStatus, setOnboardingStatus ] = useState<OnboardingStatus>(props.onboardingStatus);

        let body;
        switch (onboardingStatus) {
            case OnboardingStatus.NotStarted:
                body = <Introduction handleSetStatus={setOnboardingStatus} />
                break;
            case OnboardingStatus.InterestSelection:
                if (!props.subscriptionManagementToken || !props.emailAddress) {
                   props.setError(new Error("Something went wrong"));
                }
                body = (
                    <InterestSelection
                        emailAddress={props.emailAddress}
                        subscriptionManagementToken={props.subscriptionManagementToken}
                        handleSetStatus={setOnboardingStatus} />
                );
                break;
        }

        return (
            <CenteredComponent>
                <DisplayCard>
                    { body }
                </DisplayCard>
            </CenteredComponent>
        );
    },
    (
        ownProps: OnboardingPageOwnProps,
        onSuccess: (resp: OnboardingPageAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getOnboardingStatusForUser({
            onboardingToken: ownProps.match.params.token,
        },
        (resp: GetOnboardingStatusForUserResponse) => {
            onSuccess(resp);
        },
        onError);
    },
    true
);

type IntroductionProps = {
    handleSetStatus: (newStatus: OnboardingStatus) => void;
}

const Introduction = (props: IntroductionProps) => {
    const advanceStatus = () => {
        props.handleSetStatus(OnboardingStatus.InterestSelection);
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={advanceStatus}>
            <Heading1 color={TypographyColor.Primary}>
                Welcome to Babblegraph!
            </Heading1>
            <Paragraph>
                Let’s get your subscription set up. We’ll promise this will be quick.
            </Paragraph>
            <CenteredComponent>
                <PrimaryButton
                    type="submit"
                    className={classes.submitButton}>
                    Get Started
                </PrimaryButton>
            </CenteredComponent>
        </Form>
    )
}

type InterestSelectionProps = {
    emailAddress: string;
    subscriptionManagementToken: string;
    handleSetStatus: (newStatus: OnboardingStatus) => void;
}

const InterestSelection = (props: InterestSelectionProps) => {
    const advanceStatus = () => {
        props.handleSetStatus(OnboardingStatus.Vocabulary);
    }
    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                First, what are some topics you’re interested in?
            </Heading1>
            <Paragraph>
                This will help us tailor your newsletter to your interests. It’s easier to keep up with Spanish if your content is more interesting.
            </Paragraph>
            <Paragraph>
                Don’t worry, we don’t share this information with anyone!
            </Paragraph>
            <InterestSelector
                languageCode={WordsmithLanguageCode.Spanish}
                subscriptionManagementToken={props.subscriptionManagementToken}
                emailAddress={props.emailAddress}
                postSubmit={advanceStatus}
                omitEmailAddress />
        </div>
    )
}

export default OnboardingPage;
