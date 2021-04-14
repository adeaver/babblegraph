import { Survey, QuestionType } from 'common/data/surveys/typedefs';

const UnsubscribeSurvey: Survey = {
    header: "It’s sad to see you go!",
    description: "It’s sad to see that Babblegraph isn’t working for you. Your feedback is invaluable in helping to make Babblegraph better for others. Please share why you’re leaving. Your unsubscribe request has already been confirmed, so if you don’t fill out this survey, you will not receive any emails.",
    sections: [
        {
            header: "Unsubscribe Feedback",
            description: "Please share why you’re unsubscribing. Feedback like this helps to know what key points of Babblegraph need to be worked on the most.",
            questions: [{
                    id: "unsubscribe1-freeform",
                    questionText: "Type your feedback here.",
                    questionType: QuestionType.FreeFormText,
                    questionBody: {},
                },
            ],
        },
    ],
}

export default UnsubscribeSurvey;
