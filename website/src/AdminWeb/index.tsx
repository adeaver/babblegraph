import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link
} from 'react-router-dom';

class App extends React.Component{
    render() {
        return (
            <Router basename="/ops">
                <Switch>
                    <Route path="/hello" component={Temporary} />
                </Switch>
            </Router>
        );
    }
}

const Temporary = () => (
    <p>Hello</p>
)

ReactDOM.render(<App />, document.getElementById('content'));
