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

import SubscriptionManagementDashboardPage from 'components/SubscriptionManagement/SubscriptionManagementDashboardPage';
import InterestSelectionPage from 'components/SubscriptionManagement/InterestSelectionPage';
import DifficultyLevelSettingPage from 'components/SubscriptionManagement/DifficultyLevelSettingPage';
import UnsubscribePage from 'components/SubscriptionManagement/UnsubscribePage';
import WordReinforcementPage from 'components/SubscriptionManagement/WordReinforcementPage';

import LoginPage from 'components/UserAccounts/LoginPage';
import SignupPage from 'components/UserAccounts/SignupPage';

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
                    <Route exact path="/manage/:token" component={SubscriptionManagementDashboardPage} />

                    { /* User Account Management */ }
                    <Route path="/login" component={LoginPage} />
                    <Route path="/signup/:token" component={SignupPage} />

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
