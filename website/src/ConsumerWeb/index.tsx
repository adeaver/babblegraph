import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';

import HomePage from 'ConsumerWeb/ConsumerWeb/components/HomePage/HomePage';
import AboutPage from 'ConsumerWeb/components/AboutPage/AboutPage';
import PrivacyPolicyPage from 'ConsumerWeb/components/PrivacyPolicyPage/PrivacyPolicyPage';
import PaywallReportPage from 'ConsumerWeb/components/PaywallReportPage/PaywallReportPage';

import SubscriptionManagementDashboardPage from 'ConsumerWeb/components/SubscriptionManagement/SubscriptionManagementDashboardPage';
import InterestSelectionPage from 'ConsumerWeb/components/SubscriptionManagement/InterestSelectionPage';
import UnsubscribePage from 'ConsumerWeb/components/SubscriptionManagement/UnsubscribePage';
import WordReinforcementPage from 'ConsumerWeb/components/SubscriptionManagement/WordReinforcementPage';
import SchedulePage from 'ConsumerWeb/components/SubscriptionManagement/SchedulePage';
import NewsletterPreferencesPage from 'ConsumerWeb/components/SubscriptionManagement/NewsletterPreferencesPage';
import SubscriptionManagementPremiumInformationPage from 'ConsumerWeb/components/SubscriptionManagement/PremiumInformationPage';
import PaymentAndSubscriptionSettingsPage from 'ConsumerWeb/components/SubscriptionManagement/PaymentAndSubscriptionPage';

import LoginPage from 'ConsumerWeb/components/UserAccounts/LoginPage';
import SignupPage from 'ConsumerWeb/components/UserAccounts/SignupPage';
import ForgotPasswordPage from 'ConsumerWeb/components/UserAccounts/ForgotPasswordPage';
import ResetPasswordPage from 'ConsumerWeb/components/UserAccounts/ResetPasswordPage';
import SubscriptionCheckoutPage from 'ConsumerWeb/components/UserAccounts/SubscriptionCheckoutPage';

import NotFoundPage from 'ConsumerWeb/components/NotFoundPage/NotFoundPage';

class App extends React.Component{
    render() {
        return (
            <Router>
                <Switch>
                    { /* Subscription Management */ }
                    <Route path="/manage/:token/unsubscribe" component={UnsubscribePage} />
                    <Route path="/manage/:token/interests" component={InterestSelectionPage} />
                    <Route path="/manage/:token/vocabulary" component={WordReinforcementPage} />
                    <Route path="/manage/:token/schedule" component={SchedulePage} />
                    <Route path="/manage/:token/preferences" component={NewsletterPreferencesPage} />
                    <Route exact path="/manage/:token/premium" component={SubscriptionManagementPremiumInformationPage} />
                    <Route path="/manage/:token/payment-settings" component={PaymentAndSubscriptionSettingsPage} />
                    <Route exact path="/manage/:token" component={SubscriptionManagementDashboardPage} />
                    <Route path="/paywall-thank-you/:token" component={PaywallReportPage} />

                    { /* User Account Management */ }
                    <Route path="/login" component={LoginPage} />
                    <Route path="/signup/:token" component={SignupPage} />
                    <Route path="/checkout/:token" component={SubscriptionCheckoutPage} />
                    <Route path="/forgot-password" component={ForgotPasswordPage} />
                    <Route path="/password-reset/:token" component={ResetPasswordPage} />

                    { /* Home Page & About Page */ }
                    <Route path="/about" component={AboutPage} />
                    <Route path="/privacy-policy" component={PrivacyPolicyPage} />
                    <Route exact path="/" component={HomePage} />

                    { /* 404 Page */ }
                    <Route component={NotFoundPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
