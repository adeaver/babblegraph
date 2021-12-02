import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';

import RegistrationPage from 'AdminWeb/components/RegistrationPage/RegistrationPage';
import LoginPage from 'AdminWeb/components/LoginPage/LoginPage';
import Dashboard from 'AdminWeb/components/Dashboard/Dashboard';
import UserMetricsPage from 'AdminWeb/components/UserMetricsPage/UserMetricsPage';
import PermissionManagerPage from 'AdminWeb/components/PermissionManagerPage/PermissionManagerPage';

class App extends React.Component{
    render() {
        return (
            <Router basename="/ops">
                <Switch>
                    <Route path="/dashboard" component={Dashboard} />
                    <Route path="/user-metrics" component={UserMetricsPage} />
                    <Route path="/permission-manager" component={PermissionManagerPage} />

                    <Route path="/register/:token" component={RegistrationPage} />
                    <Route exact path="/" component={LoginPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
