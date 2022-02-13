import { DisplayLanguage } from 'common/model/language/language';

export enum TextBlock {
    UserAlreadyExistsError = 'user-already-exists-error',
    InvalidTokenError = 'invalid-token-error',
    PasswordRequirementsError = 'password-requirements-error',
    PasswordsNoMatchError = 'passwords-no-match-error',
    GenericPasswordError = 'generic-password-error',

    PageTitle = 'page-title',
    IntroParagraph1 = 'intro-paragraph-1',
    IntroParagraph2 = 'intro-paragraph-2',
    IntroParagraph3 = 'intro-paragraph-3',

}

const translations = {
    [DisplayLanguage.English]: {
        [TextBlock.UserAlreadyExistsError]: "There’s already an existing account for that email address",
        [TextBlock.InvalidTokenError]: "The email submitted didn’t match the email address this unique link is for. Make sure you entered the same email address that you received the signup link with.",
        [TextBlock.PasswordRequirementsError]: "The password entered did not match the minimum password requirements",
        [TextBlock.PasswordsNoMatchError]: "The passwords entered did not match.",
        [TextBlock.GenericPasswordError]: "Something went wrong processing your request. Try again, or email hello@babblegraph.com for help.",
        [TextBlock.PageTitle]: "Create an account",
        [TextBlock.IntroParagraph1]: "First step, sign up for a Babblegraph account",
        [TextBlock.IntroParagraph2]: "Why do you need to sign up for an account to access Babblegraph Premium?",
        [TextBlock.IntroParagraph3]: "Security is a big concern when dealing with payment information. Accounts are more secure than managing your Babblegraph subscription.",
    },
    [DisplayLanguage.Spanish]: {
        [TextBlock.UserAlreadyExistsError]: "Ya existe una cuenta para esa dirección de email.",
        [TextBlock.InvalidTokenError]: "La dirección de correo electrónico que proporcionaste no es correcta para este enlace. Por favor, asegúrate de que hubieras proporcionado la dirección correcta.",
        [TextBlock.PasswordRequirementsError]: "La contraseña proporcionada no cumplía con los requisitos mínimos.",
        [TextBlock.PasswordsNoMatchError]: "Las contraseñas proporcionadas no eran las mísmas.",
        [TextBlock.GenericPasswordError]: "Había un problema. Inténtalo de nuevo, o manda un correo electrónico a hello@babblegraph.com por ayuda.",
        [TextBlock.PageTitle]: "Crear una cuenta",
        [TextBlock.IntroParagraph1]: "Primero, es necesario crear una cuenta de Babblegraph",
        [TextBlock.IntroParagraph2]: "¿Por qué es necesario crear una cuenta para usar Babblegraph Premium?",
        [TextBlock.IntroParagraph3]: "La seguridad es muy importante cuando usas los detalles de pagamiento. Una cuenta es más segura.",
    },
}

export const getTextBlocksForLanguage = (language: DisplayLanguage | undefined) => {
    const displayLanguage = !!language ? language : DisplayLanguage.English;
    return translations[displayLanguage];
}
