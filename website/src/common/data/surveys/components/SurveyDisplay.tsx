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

import {
    Survey,
    SurveySection,
    SurveyQuestion,
    QuestionType,
    RadioQuestionBody,
} from 'common/data/surveys/typedefs';

const styleClasses = makeStyles({
    radioGroupContainer: {
        width: '100%',
    },
});

type SurveyProps = Survey & {
    surveyToken: string;
}

const Survey = (props: SurveyProps) => {
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
    } else {
        throw new Error('unimplemented question body');
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
            <Grid container>
                <Grid item xs={12}>
                    <Paragraph>
                        {props.questionText}
                    </Paragraph>
                </Grid>
            </Grid>
            <FormControl className={classes.radioGroupContainer} component="fieldset">
                <RadioGroup aria-label={`${props.id}-radiogroup`} name={`${props.id}-radiogroup`} value={currentVal} onChange={handleRadioFormChange}>
                    <Grid container>
                        <Grid item xs={1}>
                            &nbsp;
                        </Grid>
                        {
                            [0, 1, 2, 3, 4].map((val: number) => {
                                let label: string | undefined = undefined;
                                if (val === 0) {
                                    label = body.scaleMinimumLabel;
                                } else if (val === 4) {
                                    label = body.scaleMaximumLabel;
                                }
                                return (
                                    <Grid key={`${props.id}-${val}`} item xs={2}>
                                        <FormControlLabel value={`${val+1}`} control={<PrimaryRadio />} label={label} labelPlacement="bottom" />
                                    </Grid>
                                );
                            })
                        }
                        <Grid item xs={1}>
                            &nbsp;
                        </Grid>
                    </Grid>
                </RadioGroup>
            </FormControl>
        </div>
    )
}

export default Survey;
