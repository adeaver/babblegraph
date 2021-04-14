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
    FreeFormText = 'FreeFormText',
    YesOrNo = 'YesOrNo',
}

export type SurveyQuestion = {
    id: string;
    questionText: string;
    questionType: QuestionType;
    questionBody: QuestionBody;
}

export type QuestionBody = RadioQuestionBody | FreeFormTextBody | YesOrNoBody;

export type RadioQuestionBody = {
    scaleMinimumLabel: string;
    scaleMaximumLabel: string;
}

export type FreeFormTextBody = {};

export type YesOrNoBody = {
    positiveLabel: string;
    negativeLabel: string;
}
