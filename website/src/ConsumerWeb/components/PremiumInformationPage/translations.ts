import { DisplayLanguage } from 'common/model/language/language';

export enum TextBlock {
    Subheading = 'subheading',

    PodcastsTitle = 'podcasts',
    PodcastsDescription = 'podcasts-description',

    AdvertisementsTitle = 'advertisements',
    AdvertisementsDescription = 'advertisements-description',

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
        [TextBlock.PodcastsTitle]: "Podcasts",
        [TextBlock.PodcastsDescription]: "Receive podcasts in your newsletter in addition to news articles",
        [TextBlock.AdvertisementsTitle]: "No advertisements in your newsletter",
        [TextBlock.AdvertisementsDescription]: "As a Babblegraph Premium Subscriber, you will not receive any advertisements in your newsletter.",
        [TextBlock.CustomizationTitle]: "Better Newsletter Customization and Scheduling",
        [TextBlock.CustomizationDescription]: "Babblegraph Premium comes with more powerful scheduling and customization tools so you can adapt your daily newsletter to your specific needs.",
        [TextBlock.Price]: "US$29/year",
        [TextBlock.TrialText]: "The best part is that you can try it for free for 14 days!",
        [TextBlock.CallToAction]: "Start your Babblegraph Premium Subscription",
        [TextBlock.AccountDisclaimer]: "You already have a premium subscription!",
    },
    [DisplayLanguage.Spanish]: {
        [TextBlock.Subheading]: "¡Obtenga acceso a nuevas herramientas para mejorar tu aprendizaje y tener más control de la manera en que Babblegraph es parte de tu aprendizaje!",
        [TextBlock.AdvertisementsTitle]: "No habrá anuncios en su boletín",
        [TextBlock.AdvertisementsDescription]: "Como suscriptor Premium de Babblegraph, no recibirás ningún anuncio en tu boletines. ",
        [TextBlock.PodcastsTitle]: "Podcasts",
        [TextBlock.PodcastsDescription]: "Recibe podcasts en los boletines además de artículos de noticias.",
        [TextBlock.CustomizationTitle]: "Mejor personalización y programación de boletines",
        [TextBlock.CustomizationDescription]: "Babblegraph Premium viene con herramientas de programación y personalización más potentes para que pueda adaptar su boletín diario a sus necesidades específicas.",
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
