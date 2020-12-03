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

class App extends React.Component{
    render() {
        return (
            <Router>
                <Switch>
                    <Route path="/unsubscribe/:userID" component={UnsubscribePage} />
                    <Route exact path="/" component={HomePage} />
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
