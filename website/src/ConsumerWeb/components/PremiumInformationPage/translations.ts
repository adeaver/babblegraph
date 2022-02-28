import { DisplayLanguage } from 'common/model/language/language';

export enum TextBlock {
    Subheading = 'subheading',

    SchedulingTitle = 'scheduling',
    SchedulingDescription = 'scheduling-description',

    CustomizationTitle = 'customization',
    CustomizationDescription = 'customization-description',

    Price = 'price',
    TrialText = 'trial-text',

    CallToAction = 'call-to-action',
    AccountDisclaimer = 'account-disclaimer',
}

const translations = {
    [DisplayLanguage.English]: {
        [TextBlock.Subheading]: "Get access to new features that improve your learning and give you more control over how Babblegraph fits into your journey!",
        [TextBlock.SchedulingTitle]: "Newsletter Scheduling",
        [TextBlock.SchedulingDescription]: "Choose which days you receive newsletters from Babblegraph and which days you don’t.",
        [TextBlock.CustomizationTitle]: "Newsletter Customization",
        [TextBlock.CustomizationDescription]: "Want fewer articles in your newsletter? Want guarantee that you’ll see a certain topic in every email? Premium Subscribers can do just that with the newsletter customization tool!",
        [TextBlock.Price]: "US$29/year",
        [TextBlock.TrialText]: "The best part is that you can try it for free for 14 days!",
        [TextBlock.CallToAction]: "Start your Babblegraph Premium Subscription",
        [TextBlock.AccountDisclaimer]: "You already have a premium subscription!",
    },
    [DisplayLanguage.Spanish]: {
        [TextBlock.Subheading]: "¡Obtenga acceso a nuevas herramientas para mejorar tu aprendizaje y tener más control de la manera en que Babblegraph es parte de tu aprendizaje!",
        [TextBlock.SchedulingTitle]: "Más control sobre el horario de boletines",
        [TextBlock.SchedulingDescription]: "Escoge cuales días en que recibes boletines de Babblegraph y cuales días en que no los recibes.",
        [TextBlock.CustomizationTitle]: "Más control sobre la personalización de boletines",
        [TextBlock.CustomizationDescription]: "¿Quieres menos articúlos en los boletines? ¿Quieres aseguar que un tema aparece en los boletines? Los suscriptores de Babblegraph Premium tienen este típo de control sobre sus boletines.",
        [TextBlock.Price]: "US$29/año",
        [TextBlock.TrialText]: "¡Lo mejor es que puedes probarlo gratis por 14 días!",
        [TextBlock.CallToAction]: "Comienza tu suscripción a Babblegraph Premium",
        [TextBlock.AccountDisclaimer]: "Ya tienes una suscripción a Babblegraph Premium",
    },
}

export const getTextBlocksForLanguage = (language: DisplayLanguage | undefined) => {
    const displayLanguage = !!language ? language : DisplayLanguage.English;
    return translations[displayLanguage];
}
