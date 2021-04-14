import { Survey, QuestionType } from 'common/data/surveys/typedefs';

const LowOpenSurvey: Survey = {
    header: "It looks like Babblegraph isn’t working for you - Let’s fix that!",
    description: "As the creator and only person working on Babblegraph, it’s always great to hear good things about Babblegraph - but it’s most useful to hear about what isn’t working. It seems like Babblegraph isn’t working for you, so you’re a very valuable critic!",
    sections: [
        {
            header: "Rate the following aspects of Babblegraph.",
            description: "",
            questions: [
                {
                    id: "lowopen1-interest-slider1",
                    questionText: "The content in every email is interesting.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "lowopen1-learning-slider1",
                    questionText: "I am learning new words and grammar.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "lowopen1-content-slider1",
                    questionText: "The content that is shown to me matches the interests I selected.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "lowopen1-difficulty-slider1",
                    questionText: "The content I receive is too difficult.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "lowopen1-amount-content-slider1",
                    questionText: "There is too much content in every email.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                }, {
                    id: "lowopen1-format-slider1",
                    questionText: "I enjoy the format of an email.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimumLabel: "Strongly Disagree",
                        scaleMaximumLabel: "Strongly Agree",
                    },
                },
            ],
        }, {
            header: "Do you have any other comments?",
            description: "Here’s your chance to say what’s really on your mind. Maybe a new feature you’d like to see or a current feature that just isn’t working out well for you.",
            questions: [
                {
                    id: "lowopen1-freeform",
                    questionText: "Type anything you want here.",
                    questionType: QuestionType.FreeFormText,
                    questionBody: {},
                },
            ],
        },
    ],
}

export default LowOpenSurvey;
