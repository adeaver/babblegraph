import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Color from 'common/styles/colors';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Form from 'common/components/Form/Form';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    optionCardContainer: {
        padding: '5px',
    },
    selectedOptionCard: {
        cursor: 'pointer',
        borderColor: Color.Primary,
    },
    notSelectedOptionCard: {
        cursor: 'pointer',
        borderColor: Color.BorderGray,
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
    emailField: {
        width: '100%',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    submitButton: {
        display: 'block',
        margin: 'auto',
    },
});

enum TimeSelection {
    Light = 'light',
    Medium = 'medium',
    Intense = 'intense',
}

type UserNewsletterPreferencesDisplayAPIProps = {}

type UserNewsletterPreferencesDisplayOwnProps = {
    subscriptionManagementToken: string;
    languageCode: WordsmithLanguageCode;
    emailAddress?: string;
    omitEmailAddress?: boolean;
    postSubmit?: () => void;
}

const UserNewsletterPreferencesDisplay = asBaseComponent(
    (props: BaseComponentProps & UserNewsletterPreferencesDisplayAPIProps & UserNewsletterPreferencesDisplayOwnProps) => {
        const [ timeSelection, setTimeSelection ] = useState<TimeSelection>(TimeSelection.Medium);
        const [ arePodcastsEnabled, setArePodcastsEnabled ] = useState<boolean>(true);

        const [ emailAddress, setEmailAddress ] = useState<string>(props.emailAddress);
        const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setEmailAddress((event.target as HTMLInputElement).value);
        };

        const [ isLoading, setIsLoading ] = useState<boolean>(false);

        const handleSubmit = () => {
            !!props.postSubmit && props.postSubmit();
        }

        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={12}>
                    <Heading3 color={TypographyColor.Primary}>
                        How much time do you want to practice?
                    </Heading3>
                </Grid>
                <Grid className={classes.optionCardContainer} item xs={12} md={4}>
                    <DisplayCard
                        className={timeSelection === TimeSelection.Light ? classes.selectedOptionCard : classes.notSelectedOptionCard}
                        onClick={() => {setTimeSelection(TimeSelection.Light)}}>
                        <Paragraph color={timeSelection === TimeSelection.Light ? TypographyColor.Primary : TypographyColor.Gray}>
                            Light
                        </Paragraph>
                        <Paragraph size={Size.Small} color={timeSelection === TimeSelection.Light ? TypographyColor.Primary : TypographyColor.Gray}>
                            Approximately 20 minutes per email
                        </Paragraph>
                    </DisplayCard>
                </Grid>
                <Grid className={classes.optionCardContainer} item xs={12} md={4}>
                    <DisplayCard
                        className={timeSelection === TimeSelection.Medium ? classes.selectedOptionCard : classes.notSelectedOptionCard}
                        onClick={() => {setTimeSelection(TimeSelection.Medium)}}>
                        <Paragraph color={timeSelection === TimeSelection.Medium ? TypographyColor.Primary : TypographyColor.Gray}>
                            Medium
                        </Paragraph>
                        <Paragraph size={Size.Small} color={timeSelection === TimeSelection.Medium ? TypographyColor.Primary : TypographyColor.Gray}>
                            Approximately 45 minutes per email
                        </Paragraph>
                    </DisplayCard>
                </Grid>
                <Grid className={classes.optionCardContainer} item xs={12} md={4}>
                    <DisplayCard
                        className={timeSelection === TimeSelection.Intense ? classes.selectedOptionCard : classes.notSelectedOptionCard}
                        onClick={() => {setTimeSelection(TimeSelection.Intense)}}>
                        <Paragraph color={timeSelection === TimeSelection.Intense ? TypographyColor.Primary : TypographyColor.Gray}>
                            Intense
                        </Paragraph>
                        <Paragraph size={Size.Small} color={timeSelection === TimeSelection.Intense ? TypographyColor.Primary : TypographyColor.Gray}>
                            Approximately 60 minutes per email
                        </Paragraph>
                    </DisplayCard>
                </Grid>
                <Grid item xs={12}>
                    <Paragraph size={Size.Small}>
                        The timing is an estimate of how long it will take. Some articles may take more time, some may take less. If you spend way longer on the emails, it’s probably our fault, not yours.
                    </Paragraph>
                </Grid>
                <Grid item xs={12}>
                    <Heading3 color={TypographyColor.Primary}>
                        Do you want listening practice with your emails?
                    </Heading3>
                    <Paragraph>
                        Babblegraph can send you podcasts with your email. These podcasts are authentic, Spanish-language podcasts made by Spanish speakers for Spanish speakers.
                    </Paragraph>
                </Grid>
                <Grid className={classes.optionCardContainer} item xs={12} md={6}>
                    <DisplayCard
                        className={arePodcastsEnabled ? classes.selectedOptionCard : classes.notSelectedOptionCard}
                        onClick={() => {setArePodcastsEnabled(true)}}>
                        <Paragraph color={arePodcastsEnabled ? TypographyColor.Primary : TypographyColor.Gray}>
                            Yes, please include podcasts
                        </Paragraph>
                    </DisplayCard>
                </Grid>
                <Grid className={classes.optionCardContainer} item xs={12} md={6}>
                    <DisplayCard
                        className={!arePodcastsEnabled ? classes.selectedOptionCard : classes.notSelectedOptionCard}
                        onClick={() => {setArePodcastsEnabled(false)}}>
                        <Paragraph color={!arePodcastsEnabled ? TypographyColor.Primary : TypographyColor.Gray}>
                            No, I think I’m good
                        </Paragraph>
                    </DisplayCard>
                </Grid>
                <Grid item xs={12}>
                    <Form
                        className={classes.confirmationForm}
                        handleSubmit={handleSubmit}>
                        <Grid container>
                        {
                            !!props.omitEmailAddress ? (
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
                                    type="submit"
                                    className={classes.submitButton}
                                    disabled={!emailAddress && !props.omitEmailAddress}>
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                </Grid>
            </Grid>
        );
    },
    (
        ownProps: UserNewsletterPreferencesDisplayOwnProps,
        onSuccess: (resp: UserNewsletterPreferencesDisplayAPIProps) => void,
        onError: (err: Error) => void,
    ) => onSuccess({}),
    false
);

// Okay, so here's what needs to happen here:
// 1) "How much practice would you like to receive?"
// - Light (about 20 minutes) - 4 articles + < 15 minutes of podcasts
// - Medium (about 45 minutes) - 8 articles + 15-30 minutes of podcasts
// - Heavy (about 60 minutes) - 12 articles + < 15-30 minutes of podcasts

// 2) Would you like listening practice or just reading practice?
// - Include just articles
// - Include both

// 3) Are you okay with potentially explicit podcasts? (i.e. bad words or potentially triggering content?)
// - I'm fine with it
// - I'd really rather not

export default UserNewsletterPreferencesDisplay;
