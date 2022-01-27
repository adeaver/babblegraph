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
import BlogListPage from 'AdminWeb/components/Blog/BlogListPage/BlogListPage';
import BlogEditPage from 'AdminWeb/components/Blog/BlogEditPage/BlogEditPage';
import ContentManagerDashboard from 'AdminWeb/components/ContentManager/ContentManagerDashboard';
import TopicListPage from 'AdminWeb/components/ContentManager/topics/TopicListPage';
import TopicManagementPage from 'AdminWeb/components/ContentManager/topics/TopicManagementPage';
import SourcesListPage from 'AdminWeb/components/ContentManager/sources/SourcesListPage';
import SourceManagementPage from 'AdminWeb/components/ContentManager/sources/SourceManagementPage';

class App extends React.Component{
    render() {
        return (
            <Router basename="/ops">
                <Switch>
                    <Route path="/dashboard" component={Dashboard} />

                    <Route path="/permission-manager" component={PermissionManagerPage} />

                    <Route path="/user-metrics" component={UserMetricsPage} />

                    <Route path="/content-manager/topics/:id" component={TopicManagementPage} />
                    <Route path="/content-manager/topics" component={TopicListPage} />
                    <Route path="/content-manager/sources/:id" component={SourceManagementPage} />
                    <Route path="/content-manager/sources" component={SourcesListPage} />
                    <Route path="/content-manager" component={ContentManagerDashboard} />

                    <Route path="/blog-manager/edit/:blogPath" component={BlogEditPage} />
                    <Route path="/blog-manager" component={BlogListPage} />

                    <Route path="/register/:token" component={RegistrationPage} />
                    <Route exact path="/" component={LoginPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
