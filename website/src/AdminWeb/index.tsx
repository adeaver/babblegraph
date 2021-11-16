import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';

import LoginPage from 'AdminWeb/components/LoginPage/LoginPage';

class App extends React.Component{
    render() {
        return (
            <Router basename="/ops">
                <Switch>
                    <Route component={LoginPage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
