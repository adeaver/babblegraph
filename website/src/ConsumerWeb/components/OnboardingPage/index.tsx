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

import UserNewsletterPreferencesDisplay from 'ConsumerWeb/components/UserNewsletterPreferencesPage/UserNewsletterPreferencesDisplay';
import InterestSelector from 'ConsumerWeb/components/InterestSelectionPage/InterestSelector';
import TimeSelector from 'ConsumerWeb/components/UserNewsletterSchedulePage/TimeSelector';

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
            case OnboardingStatus.Settings:
                if (!props.subscriptionManagementToken || !props.emailAddress) {
                   props.setError(new Error("Something went wrong"));
                }
                body = (
                    <SettingsSelection
                        emailAddress={props.emailAddress}
                        subscriptionManagementToken={props.subscriptionManagementToken}
                        handleSetStatus={setOnboardingStatus} />
                );
                break;
            case OnboardingStatus.Schedule:
                if (!props.subscriptionManagementToken || !props.emailAddress) {
                   props.setError(new Error("Something went wrong"));
                }
                body = (
                    <ScheduleSelection
                        emailAddress={props.emailAddress}
                        subscriptionManagementToken={props.subscriptionManagementToken}
                        handleSetStatus={setOnboardingStatus} />
                );
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
        props.handleSetStatus(OnboardingStatus.Settings);
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

type SettingsSelectionProps = {
    emailAddress: string;
    subscriptionManagementToken: string;
    handleSetStatus: (newStatus: OnboardingStatus) => void;
}

const SettingsSelection = (props: SettingsSelectionProps) => {
    const advanceStatus = () => {
        props.handleSetStatus(OnboardingStatus.Schedule);
    }
    return (
        <div>
            <Paragraph size={Size.Small}>
                Page 1 of 4
            </Paragraph>
            <Heading1 color={TypographyColor.Primary}>
                It’s your newsletter, let’s make it feel like yours.
            </Heading1>
            <Paragraph>
                Let’s get started with some questions about your basic newsletter experience.
            </Paragraph>
            <UserNewsletterPreferencesDisplay
                languageCode={WordsmithLanguageCode.Spanish}
                subscriptionManagementToken={props.subscriptionManagementToken}
                emailAddress={props.emailAddress}
                postSubmit={advanceStatus}
                omitEmailAddress />
        </div>
    )
}

type ScheduleSelectionProps = {
    emailAddress: string;
    subscriptionManagementToken: string;
    handleSetStatus: (newStatus: OnboardingStatus) => void;
}

const ScheduleSelection = (props: ScheduleSelectionProps) => {
    const advanceStatus = () => {
        props.handleSetStatus(OnboardingStatus.InterestSelection);
    }
    return (
        <div>
            <Paragraph size={Size.Small}>
                Page 2 of 4
            </Paragraph>
            <Heading1 color={TypographyColor.Primary}>
                Is the timing right between us?
            </Heading1>
            <Paragraph>
                Okay, jokes aside. Let’s talk time. You can customize which days and times you want to receive your emails.
            </Paragraph>
            <Paragraph>
                Don’t worry, we don’t share this information with anyone or judge you if your preferred practice time is 3am...
            </Paragraph>
            <TimeSelector
                languageCode={WordsmithLanguageCode.Spanish}
                subscriptionManagementToken={props.subscriptionManagementToken}
                emailAddress={props.emailAddress}
                postSubmit={advanceStatus}
                omitEmailAddress />
        </div>
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
            <Paragraph size={Size.Small}>
                Page 3 of 4
            </Paragraph>
            <Heading1 color={TypographyColor.Primary}>
                Now, what are some topics you’re interested in?
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
