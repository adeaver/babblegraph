import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import RadioGroup from '@material-ui/core/RadioGroup';

import { Heading1, Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { PrimaryRadio } from 'common/components/Radio/Radio';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';

import {
    Survey,
    SurveySection,
    SurveyQuestion,
    QuestionType,
    RadioQuestionBody,
    YesOrNoBody,
} from 'common/data/surveys/typedefs';

const styleClasses = makeStyles({
    submitButton: {
        display: 'block',
        margin: '10px auto',
    },
    radioGroupContainer: {
        width: '100%',
        padding: '0 20px',
        boxSizing: 'border-box',
    },
    freeTextField: {
        width: '100%',
    },
    radioButton: {
        alignContent: 'center',
        width: '20%',
        maxWidth: '20%',
        margin: '0',
    },
    yesNoRadioButton: {
        width: '100%',
        justifyContent: 'center',
    },
});

type SurveyProps = Survey & {
    surveyToken: string;
}

const Survey = (props: SurveyProps) => {
    const classes = styleClasses();
    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                {props.header}
            </Heading1>
            <Paragraph>
                {props.description}
            </Paragraph>
            <div>
                {
                    props.sections.map((s: SurveySection) => {
                        return <Section {...s} />;
                    })
                }
            </div>
            <PrimaryButton className={classes.submitButton}>Submit</PrimaryButton>
        </div>
    );
}

const Section = (props: SurveySection) => {
    return (
        <div>
            <Heading3 color={TypographyColor.Primary}>
                {props.header}
            </Heading3>
            <Paragraph>
                {props.description}
            </Paragraph>
            {
                props.questions.map((question: SurveyQuestion) => {
                    return <Question {...question} />;
                })
            }
        </div>
    );
}

const Question = (props: SurveyQuestion) => {
    if (props.questionType === QuestionType.RadioQuestion) {
        return <RadioQuestion {...props} />;
    } else if (props.questionType === QuestionType.FreeFormText) {
        return <FreeFormTextQuestion {...props} />;
    } else if (props.questionType === QuestionType.YesOrNo) {
        return <YesOrNoQuestion {...props} />;
    } else {
        throw new Error(`unimplemented question type ${props.questionType}`);
    }
}

const RadioQuestion = (props: SurveyQuestion) => {
    const classes = styleClasses();
    const body: RadioQuestionBody = props.questionBody as RadioQuestionBody;

    const [ currentVal, setCurrentVal ] = useState<string | null>(null);
    const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const nextVal = (event.target as HTMLInputElement).value;
        setCurrentVal(nextVal);
    };

    return (
        <div>
            <Paragraph>
                {props.questionText}
            </Paragraph>
            <FormControl className={classes.radioGroupContainer} component="fieldset">
                <RadioGroup aria-label={`${props.id}-radiogroup`} name={`${props.id}-radiogroup`} value={currentVal} onChange={handleRadioFormChange} row>
                        {
                            [0, 1, 2, 3, 4].map((val: number) => {
                                let label: string | undefined = undefined;
                                if (val === 0) {
                                    label = body.scaleMinimumLabel;
                                } else if (val === 4) {
                                    label = body.scaleMaximumLabel;
                                }
                                return (
                                    <FormControlLabel
                                        className={classes.radioButton}
                                        value={`${val+1}`}
                                        control={<PrimaryRadio />}
                                        label={label}
                                        labelPlacement="bottom" />
                                );
                            })
                        }
                </RadioGroup>
            </FormControl>
        </div>
    )
}

const FreeFormTextQuestion = (props: SurveyQuestion) => {
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <PrimaryTextField
                    className={classes.freeTextField}
                    label={props.questionText}
                    rows={5}
                    multiline
                    variant="outlined" />
            </Grid>
        </Grid>
    );
}

const YesOrNoQuestion = (props: SurveyQuestion) => {
    const classes = styleClasses();
    const body: YesOrNoBody = props.questionBody as YesOrNoBody;

    const [ currentVal, setCurrentVal ] = useState<string | null>(null);
    const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const nextVal = (event.target as HTMLInputElement).value;
        setCurrentVal(nextVal);
    };
    return (
        <div>
            <Paragraph>
                {props.questionText}
            </Paragraph>
            <FormControl className={classes.radioGroupContainer} component="fieldset">
                <RadioGroup aria-label={`${props.id}-radiogroup`} name={`${props.id}-radiogroup`} value={currentVal} onChange={handleRadioFormChange} row>
                    <Grid container>
                        <Grid item xs={3}>
                            &nbsp;
                        </Grid>
                        <Grid item xs={3}>
                            <FormControlLabel
                                className={classes.yesNoRadioButton}
                                value={body.positiveLabel}
                                control={<PrimaryRadio />}
                                label={body.positiveLabel} />
                        </Grid>
                        <Grid item xs={3}>
                            <FormControlLabel
                                className={classes.yesNoRadioButton}
                                value={body.negativeLabel}
                                control={<PrimaryRadio />}
                                label={body.negativeLabel} />
                        </Grid>
                    </Grid>
                </RadioGroup>
            </FormControl>
        </div>
    );
}

export default Survey;
