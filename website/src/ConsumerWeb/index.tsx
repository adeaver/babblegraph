import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';

import HomePage from 'ConsumerWeb/components/HomePage/HomePage';
import AboutPage from 'ConsumerWeb/components/AboutPage/AboutPage';
import PrivacyPolicyPage from 'ConsumerWeb/components/PrivacyPolicyPage/PrivacyPolicyPage';
import PricingPage from 'ConsumerWeb/components/PricingPage/PricingPage';

import AdvertisingPolicyPage from 'ConsumerWeb/components/AdvertisingPolicyPage/AdvertisingPolicyPage';

import BlogListPage from 'ConsumerWeb/components/BlogListPage/BlogListPage';
import BlogPostPage from 'ConsumerWeb/components/BlogPostPage/BlogPostPage';

import ArticlePage from 'ConsumerWeb/components/ArticlePage';

import SubscriptionManagementHomePage from 'ConsumerWeb/components/SubscriptionManagement/SubscriptionManagementHomePage';
import InterestSelectionPage from 'ConsumerWeb/components/SubscriptionManagement/InterestSelectionPage';
import UnsubscribePage from 'ConsumerWeb/components/UnsubscribePage/UnsubscribePage';
import WordReinforcementPage from 'ConsumerWeb/components/WordReinforcementPage';
import UserNewsletterPreferencesPage from 'ConsumerWeb/components/UserNewsletterPreferencesPage/UserNewsletterPreferencesPage';

import LoginPage from 'ConsumerWeb/components/UserAccounts/LoginPage';
import CreateUserAccountPage from 'ConsumerWeb/components/CreateUserAccountPage/CreateUserAccountPage';
import ForgotPasswordPage from 'ConsumerWeb/components/UserAccounts/ForgotPasswordPage';
import ResetPasswordPage from 'ConsumerWeb/components/UserAccounts/ResetPasswordPage';

import PremiumNewsletterSubscriptionCheckoutPage from 'ConsumerWeb/components/PremiumNewsletterSubscriptionCheckoutPage/PremiumNewsletterSubscriptionCheckoutPage';
import PremiumInformationPage from 'ConsumerWeb/components/PremiumInformationPage/PremiumInformationPage';
import PremiumNewsletterSubscriptionManagementPage from 'ConsumerWeb/components/PremiumNewsletterSubscriptionManagementPage/PremiumNewsletterSubscriptionManagementPage';

import PodcastPlayerPage from 'ConsumerWeb/components/PodcastPlayerPage/PodcastPlayerPage';

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
                    <Route path="/manage/:token/preferences" component={UserNewsletterPreferencesPage} />
                    <Route exact path="/manage/:token/premium" component={PremiumInformationPage} />
                    <Route path="/manage/:token/payment-settings" component={PremiumNewsletterSubscriptionManagementPage} />
                    <Route exact path="/manage/:token" component={SubscriptionManagementHomePage} />

                    { /* Content */ }
                    <Route path="/article/:token" component={ArticlePage} />
                    <Route path="/podcast/:userPodcastID" component={PodcastPlayerPage} />

                    { /* User Account Management */ }
                    <Route path="/login" component={LoginPage} />
                    <Route path="/signup/:token" component={CreateUserAccountPage} />
                    <Route path="/checkout/:token" component={PremiumNewsletterSubscriptionCheckoutPage} />
                    <Route path="/forgot-password" component={ForgotPasswordPage} />
                    <Route path="/password-reset/:token" component={ResetPasswordPage} />

                    { /* Blog */ }
                    <Route path="/blog/:blogPath" component={BlogPostPage} />
                    <Route exact path="/blog" component={BlogListPage} />

                    { /* Home Page & About Page */ }
                    <Route path="/about" component={AboutPage} />
                    <Route path="/privacy-policy" component={PrivacyPolicyPage} />
                    <Route path="/pricing" component={PricingPage} />
                    <Route path="/advertising-policy" component={AdvertisingPolicyPage} />
                    <Route exact path="/" component={HomePage} />

                    { /* 404 Page */ }
                    <Route component={NotFoundPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
