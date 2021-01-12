import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';

import HomePage from 'components/HomePage/HomePage';

import SubscriptionManagementDashboardPage from 'components/SubscriptionManagement/SubscriptionManagementDashboardPage';
import InterestSelectionPage from 'components/SubscriptionManagement/InterestSelectionPage';
import DifficultyLevelSettingPage from 'components/SubscriptionManagement/DifficultyLevelSettingPage';

import NotFoundPage from 'components/NotFoundPage/NotFoundPage';

class App extends React.Component{
    render() {
        return (
            <Router>
                <Switch>
                    <Route path="/manage/:token/level" component={DifficultyLevelSettingPage} />
                    <Route path="/manage/:token/interests" component={InterestSelectionPage} />
                    <Route exact path="/manage/:token" component={SubscriptionManagementDashboardPage} />
                    <Route exact path="/" component={HomePage} />
                    <Route component={NotFoundPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
