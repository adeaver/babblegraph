import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';

import HomePage from 'components/HomePage/HomePage';
import AboutPage from 'components/AboutPage/AboutPage';
import PrivacyPolicyPage from 'components/PrivacyPolicyPage/PrivacyPolicyPage';
import PaywallReportPage from 'components/PaywallReportPage/PaywallReportPage';

import SubscriptionManagementDashboardPage from 'components/SubscriptionManagement/SubscriptionManagementDashboardPage';
import InterestSelectionPage from 'components/SubscriptionManagement/InterestSelectionPage';
import DifficultyLevelSettingPage from 'components/SubscriptionManagement/DifficultyLevelSettingPage';
import UnsubscribePage from 'components/SubscriptionManagement/UnsubscribePage';
import WordReinforcementPage from 'components/SubscriptionManagement/WordReinforcementPage';
import SchedulePage from 'components/SubscriptionManagement/SchedulePage';
import NewsletterPreferencesPage from 'components/SubscriptionManagement/NewsletterPreferencesPage';
import SubscriptionManagementPremiumInformationPage from 'components/SubscriptionManagement/PremiumInformationPage';

import LoginPage from 'components/UserAccounts/LoginPage';
import SignupPage from 'components/UserAccounts/SignupPage';
import ForgotPasswordPage from 'components/UserAccounts/ForgotPasswordPage';
import ResetPasswordPage from 'components/UserAccounts/ResetPasswordPage';
import SubscriptionCheckoutPage from 'components/UserAccounts/SubscriptionCheckoutPage';

import NotFoundPage from 'components/NotFoundPage/NotFoundPage';

class App extends React.Component{
    render() {
        return (
            <Router>
                <Switch>
                    { /* Subscription Management */ }
                    <Route path="/manage/:token/unsubscribe" component={UnsubscribePage} />
                    <Route path="/manage/:token/level" component={DifficultyLevelSettingPage} />
                    <Route path="/manage/:token/interests" component={InterestSelectionPage} />
                    <Route path="/manage/:token/vocabulary" component={WordReinforcementPage} />
                    <Route path="/manage/:token/schedule" component={SchedulePage} />
                    <Route path="/manage/:token/preferences" component={NewsletterPreferencesPage} />
                    <Route exact path="/manage/:token/premium" component={SubscriptionManagementPremiumInformationPage} />
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
