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

import BlogHomePage from 'components/Blog/BlogHomePage';
import BlogPostPage from 'components/Blog/BlogPostPage';

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

                    { /* Home Page & About Page */ }
                    <Route path="/about" component={AboutPage} />
                    <Route path="/privacy-policy" component={PrivacyPolicyPage} />
                    <Route exact path="/" component={HomePage} />

                    { /* Blog */ }
                    <Route path="/blog/:path" component={BlogPostPage} />
                    <Route exact path="/blog" component={BlogHomePage} />

                    { /* 404 Page */ }
                    <Route component={NotFoundPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
