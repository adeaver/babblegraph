import { DisplayLanguage } from 'common/model/language/language';

export enum TextBlock {
    ErrorMessageInvalidEmailAddress = 'error-message-invalid-email-address',
    ErrorMessageRateLimited = 'error-message-rate-limited',
    ErrorMessageIncorrectStatus = 'error-message-incorrect-status',
    ErrorMessageLowScore = 'error-message-low-score',
    ErrorMessageDefault = 'error-message-default',

    EmailAddressInputLabel = 'email-address-input-label',
    SignupFormButtonText = 'signup-form-button-text',

    VerificationLocationConfirmation = 'verification-location-confirmation',
    VerificationInstructions = 'verification-instructions',
    VerificationWarningDisclaimer = 'verification-warning-disclaimer',
    VerificationResendButtonText = 'verification-resend-button-text',
}

const translations = {
    [DisplayLanguage.English]: {
        [TextBlock.ErrorMessageInvalidEmailAddress]: "Hmm, the email address you gave doesn’t appear to be valid. Check to make sure that you spelled everything right.",
        [TextBlock.ErrorMessageRateLimited]: "It looks like we’re having some trouble reaching you. Contact our support so we can get you on the list!",
        [TextBlock.ErrorMessageIncorrectStatus]: "It looks like you’re already signed up for Babblegraph!",
        [TextBlock.ErrorMessageLowScore]: "We’re having some trouble verifying your request. Contact us at hello@babblegraph.com to finish signing up.",
        [TextBlock.ErrorMessageDefault]: "Something went wrong. Contact our support so we can get you on the list!",
        [TextBlock.EmailAddressInputLabel]: "Email address",
        [TextBlock.SignupFormButtonText]: "Try it for free!",
        [TextBlock.VerificationLocationConfirmation]: "Check your email for a verification email from hello@babblegraph.com. We sent it to",
        [TextBlock.VerificationInstructions]: "You’ll need to click the button in the verification email that was just sent to you in order to start receiving emails from Babblegraph.",
        [TextBlock.VerificationWarningDisclaimer]: "It can take up to 5 minutes for the email to make its way to your inbox.",
        [TextBlock.VerificationResendButtonText]: "Resend the verification email",
    },
    [DisplayLanguage.Spanish]: {
        [TextBlock.ErrorMessageInvalidEmailAddress]: "La dirección de correo electrónico que proporcionó no es válida. Asegúrate de que esté escrito correctamente.",
        [TextBlock.ErrorMessageRateLimited]: "Parece que hay un problema con la entrega a su dirección. Póngase en contacto con el soporte técnico para solucionar el problema. Manda un email a hello@babblegraph.com.",
        [TextBlock.ErrorMessageIncorrectStatus]: "Ya estás en la lista de Babblegraph!",
        [TextBlock.ErrorMessageLowScore]: "Parece que hay un problema con la entrega a su dirección. Póngase en contacto con el soporte técnico para solucionar el problema. Manda un email a hello@babblegraph.com",
        [TextBlock.ErrorMessageDefault]: "Parece que hay un problema con la entrega a su dirección. Póngase en contacto con el soporte técnico para solucionar el problema. Manda un email a hello@babblegraph.com.",
        [TextBlock.EmailAddressInputLabel]: "Dirección de email",
        [TextBlock.SignupFormButtonText]: "¡Regístrese gratis!",
        [TextBlock.VerificationLocationConfirmation]: "Busque en su inbox un email de verificación de hello@babblegraph.com. Lo enviamos a",
        [TextBlock.VerificationInstructions]: "Necesita hacer clic en el botón en este email para empezar a recibir los emails diarios de Babblegraph.",
        [TextBlock.VerificationWarningDisclaimer]: "Pueden pasar hasta cinco minutos antes de recibir este email.",
        [TextBlock.VerificationResendButtonText]: "Enviar de nuevo el email de verificación",
    },
}

export const getTextBlocksForLanguage = (language: DisplayLanguage | undefined) => {
    const displayLanguage = !!language ? language : DisplayLanguage.English;
    return translations[displayLanguage];
}
