import { DisplayLanguage } from 'common/model/language/language';

export enum ClientError {
    InvalidToken = "invalid-token",
    NoAuth = "no-auth",
    IncorrectKey = "incorrect-key",
    InvalidLanguageCode = "invalid-language",
    InvalidEmailAddress = "invalid-email-address",

    DefaultError = 'unknown-error',
}

export const asReadable = (error: ClientError, displayLanguage: DisplayLanguage) => {
    const errorMessagesForLanguage = errorMessages[displayLanguage] || errorMessages[DisplayLanguage.English];
    return errorMessagesForLanguage[error] || errorMessagesForLanguage[ClientError.DefaultError]
}

const errorMessages = {
    [DisplayLanguage.Spanish]: {
        [ClientError.InvalidToken]: "Este enlace se ha caducado. Haz clic en el enlace que está en la parte inferior de su boletín.",
        [ClientError.InvalidLanguageCode]: "Babblegraph no tiene ese idioma, pruébalo luego",
        [ClientError.InvalidEmailAddress]: "La dirección de correo electrónico que proporcionaste no es correcta para este enlace. Por favor, asegúrate de que hubieras proporcionado la dirección correcta.",
        [ClientError.IncorrectKey]: "Esa acción no se permite.",
        [ClientError.NoAuth]: "Esa acción no se permite.",
        [ClientError.DefaultError]: "Parece que hay un problema con la entrega a su dirección. Póngase en contacto con el soporte técnico para solucionar el problema. Manda un email a hello@babblegraph.com.",
    },
    [DisplayLanguage.English]: {
        [ClientError.InvalidToken]: "This link has expired. Try clicking on the link at the bottom of your newsletter",
        [ClientError.InvalidLanguageCode]: "That language is not yet supported by Babblegraph, try again later",
        [ClientError.InvalidEmailAddress]: "The email address that you provided was not correct",
        [ClientError.IncorrectKey]: "That action is forbidden.",
        [ClientError.NoAuth]: "That action is forbidden.",
        [ClientError.DefaultError]: "Something went wrong processing that request. Try again later or email hello@babblegraph.com for support.",
    },
};
