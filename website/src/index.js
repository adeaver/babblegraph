import React from 'react';
import ReactDOM from 'react-dom';

class App extends React.Component{
    render() {
        return (
            <div>If you're seeing this, everything is doubly fine.</div>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('content'));
