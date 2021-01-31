declare const window: any;

export const withCaptchaToken = (actionName: string, fn: (token: string) => void) => {
    window.grecaptcha.ready(() => {
        window.grecaptcha.execute("6LfydEEaAAAAAMdgim0VlnzE6TSq01Urg3Qde_ui", { action: actionName })
        .then((token: string) => {
            fn(token);
        });
    });
}

export const loadCaptchaScript = () => {
    const script = document.createElement("script")
    script.src = "https://www.google.com/recaptcha/api.js?render=6LfydEEaAAAAAMdgim0VlnzE6TSq01Urg3Qde_ui"
    document.body.appendChild(script)
}
