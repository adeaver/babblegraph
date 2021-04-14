export type Survey = {
    header: string;
    description: string;
    sections: Array<SurveySection>;
}

export type SurveySection = {
    header: string;
    description: string;
    questions: Array<SurveyQuestion>;
}

export enum QuestionType {
    RadioQuestion = 'RadioQuestion',
}

export type SurveyQuestion = {
    id: string;
    questionText: string;
    questionType: QuestionType;
    questionBody: QuestionBody;
}

export type QuestionBody = RadioQuestionBody

export type RadioQuestionBody = {
    scaleMinimumLabel: string;
    scaleMaximumLabel: string;
}
