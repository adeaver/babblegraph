import { Survey, QuestionType } from 'common/data/surveys/typedefs';

const HighOpenSurvey: Survey = {
    header: "Help Make Babblegraph Better",
    description: "As the creator and only person working on Babblegraph, it can be hard to figure out what the next best thing to work on is. I use Babblegraph, so I have a list of new features I want to build or existing feature I’d like to improve, but it’s better to know what everyone else wants to see.",
    sections: [
        {
            header: "Rate the following aspects of Babblegraph.",
            description: "",
            questions: [
                {
                    id: "highopen1-interest-slider1",
                    questionText: "The content in every email is interesting.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "highopen1-learning-slider1",
                    questionText: "I am learning new words and grammar.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "highopen1-content-slider1",
                    questionText: "The content that is shown to me matches the interests I selected.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "highopen1-difficulty-slider1",
                    questionText: "The content I receive is too difficult.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "highopen1-amount-content-slider1",
                    questionText: "There is too much content in every email.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "highopen1-format-slider1",
                    questionText: "I enjoy the format of the email.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                },
            ],
        }, {
            header: "Open Ended",
            description: "Here’s your chance to say what’s really on your mind. Maybe a new feature you’d like to see or a current feature that just isn’t working out well for you.",
            questions: [{
                    id: "highopen1-new-feature",
                    questionText: "What’s a new feature you’d like to see?",
                    questionType: QuestionType.FreeFormText,
                    questionBody: {},
                }, {
                    id: "highopen1-freeform",
                    questionText: "Is there anything else you’d like to say?",
                    questionType: QuestionType.FreeFormText,
                    questionBody: {},
                }, {
                    id: "highopen1-contact",
                    questionText: "Can we contact you to follow up on your responses?",
                    questionType: QuestionType.YesOrNo,
                    questionBody: {
                        positiveLabel: "Yes",
                        negativeLabel: "No",
                    },
                },
            ],
        },
    ],
}

export default HighOpenSurvey;
