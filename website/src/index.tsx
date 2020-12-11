import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';
import HomePage from 'components/HomePage/HomePage';
import UnsubscribePage from 'components/UnsubscribePage/UnsubscribePage';
import SubscriptionManagementPage from 'components/SubscriptionManagementPage/SubscriptionManagementPage';
import NotFoundPage from 'components/NotFoundPage/NotFoundPage';

class App extends React.Component{
    render() {
        return (
            <Router>
                <Switch>
                    <Route path="/unsubscribe/:token" component={UnsubscribePage} />
                    <Route path="/manage/:token" component={SubscriptionManagementPage} />
                    <Route exact path="/" component={HomePage} />
                    <Route component={NotFoundPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
