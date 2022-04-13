import { DisplayLanguage } from 'common/model/language/language';

export enum TextBlock {
    Subheading = 'subheading',

    PhrasesTitle = 'phrases',
    PhrasesDescription = 'phrases-description',

    PodcastsTitle = 'podcasts',
    PodcastsDescription = 'podcasts-description',

    AdvertisementsTitle = 'advertisements',
    AdvertisementsDescription = 'advertisements-description',

    IndependentCreatorTitle = 'independent-creator',
    IndependentCreatorDescription = 'independent-creator-description',

    AndMoreTitle = 'and-more',
    AndMoreDescription = 'and-more-description',

    Price = 'price',
    TrialText = 'trial-text',

    CallToAction = 'call-to-action',
    AccountDisclaimer = 'account-disclaimer',
}

const translations = {
    [DisplayLanguage.English]: {
        [TextBlock.Subheading]: "Get access to new features that improve your learning and give you more control over how Babblegraph fits into your journey!",
        [TextBlock.PhrasesTitle]: "Add phrases like “tener calor” to your vocabulary list",
        [TextBlock.PhrasesDescription]: "Go beyond simple vocabulary words and start reinforcing entire phrases with your newsletter.",
        [TextBlock.PodcastsTitle]: "Podcasts",
        [TextBlock.PodcastsDescription]: "Receive podcasts in your newsletter in addition to news articles",
        [TextBlock.AdvertisementsTitle]: "No advertisements in your newsletter",
        [TextBlock.AdvertisementsDescription]: "As a Babblegraph Premium Subscriber, you will not receive any advertisements in your newsletter.",
        [TextBlock.IndependentCreatorTitle]: "Support an independent creator",
        [TextBlock.IndependentCreatorDescription]: "Believe it or not, Babblegraph is built, maintained, and marketed by just one person! If you enjoy it, consider using Premium to support me!",
        [TextBlock.AndMoreTitle]: "...And more coming!",
        [TextBlock.AndMoreDescription]: "I am constantly adding new features, and plan to add more premium only features in the future. Get access to those as I release them by signing up for Premium today",
        [TextBlock.Price]: "US$29/year",
        [TextBlock.TrialText]: "The best part is that you can try it for free for 14 days!",
        [TextBlock.CallToAction]: "Start your Babblegraph Premium Subscription",
        [TextBlock.AccountDisclaimer]: "You already have a premium subscription!",
    },
    [DisplayLanguage.Spanish]: {
        [TextBlock.Subheading]: "¡Obtén acceso a nuevas herramientas para mejorar tu aprendizaje y tener más control de la manera en que Babblegraph es parte de tu aprendizaje!",
        [TextBlock.PhrasesTitle]: "Agrega frases, como “tener calor” a tu lista de vocabulario",
        [TextBlock.PhrasesDescription]: "Vaya más allá de las simples palabras de vocabulario y comienza a reforzar frases completas con tu boletín informativo.",
        [TextBlock.AdvertisementsTitle]: "No habrá anuncios en su boletín",
        [TextBlock.AdvertisementsDescription]: "Como suscriptor Premium de Babblegraph, no recibirás ningún anuncio en tu boletines. ",
        [TextBlock.PodcastsTitle]: "Podcasts",
        [TextBlock.PodcastsDescription]: "Recibe podcasts en los boletines además de artículos de noticias.",
        [TextBlock.IndependentCreatorTitle]: "Apoya a un creador independiente",
        [TextBlock.IndependentCreatorDescription]: "Babblegraph es construido, mantenido y comercializado por una sola persona. Si lo disfrutas, ¡considera usar Premium para apoyarme! ",
        [TextBlock.AndMoreTitle]: "... y más viene",
        [TextBlock.AndMoreDescription]: "Constantemente agrego nuevas funciones y planeo agregar más funciones exclusivas premium en el futuro. Obtén acceso a ellos a medida que los publique por registrarte en Premium hoy",
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
