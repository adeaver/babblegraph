import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';
import HomePage from 'components/HomePage/HomePage';

class App extends React.Component{
    render() {
        return (
            <Router>
                <Switch>
                    <Route path="/">
                        <HomePage />
                    </Route>
                </Switch>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
