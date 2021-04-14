import { Survey, QuestionType } from 'common/data/surveys/typedefs';

const LowOpenSurvey: Survey = {
    header: "It looks like Babblegraph isn’t working for you - Let’s fix that!",
    description: "As the creator and only person working on Babblegraph, it’s always great to hear good things about Babblegraph - but it’s most useful to hear about what isn’t working. It seems like Babblegraph isn’t working for you, so you’re a very valuable critic!",
    sections: [
        {
            header: "Rate the following aspects of Babblegraph.",
            description: "A score of 5 means you agree, and a score of 1 means you disagree.",
            questions: [
                {
                    id: "lowopen1-interest-slider1",
                    questionText: "The content in every email is interesting.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimum: 1,
                        scaleMaximum: 5,
                    },
                }, {
                    id: "lowopen1-learning-slider1",
                    questionText: "I am learning new words and grammar.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimum: 1,
                        scaleMaximum: 5,
                    },
                }, {
                    id: "lowopen1-content-slider1",
                    questionText: "The content that is shown to me matches the interests I selected.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimum: 1,
                        scaleMaximum: 5,
                    },
                }, {
                    id: "lowopen1-difficulty-slider1",
                    questionText: "The content I receive is too difficult.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimum: 1,
                        scaleMaximum: 5,
                    },
                }, {
                    id: "lowopen1-difficulty-slider1",
                    questionText: "There is too much content in every email.",
                    questionType: QuestionType.RadioQuestion,
                    questionBody: {
                        scaleMinimum: 1,
                        scaleMaximum: 5,
                    },
                },
            ],
        },
    ],
}

export default LowOpenSurvey;
