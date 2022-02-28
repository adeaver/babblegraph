export enum PaymentState {
    CreatedUnpaid = 0,
    TrialNoPaymentMethod = 1,
    TrialPaymentMethodAdded = 2,
    Active = 3,
    Errored = 4,
    Terminated = 5,
}

export type PremiumNewsletterSubscription = {
    id: string | undefined;
    paymentState: PaymentState;
    stripePaymentIntentId: string | undefined;
    currentPeriodEnd: string;
    isAutoRenewEnabled: boolean;
}
