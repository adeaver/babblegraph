import React from 'react';

import { Heading1, Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';

import {
    Survey,
    SurveySection,
    SurveyQuestion,
    QuestionType,
} from 'common/data/surveys/typedefs';

const Survey = (props: Survey) => {
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
    return (<h1>{props.questionText}</h1>);
}

export default Survey;
