import React from 'react';
import {List, Map} from 'immutable';
import {Header} from './components';
import {Add, Feed} from './containers';
import Radium from 'radium';
import Colors from './colors';

var baseStyle = {
  color: Colors.primaryText,
  fontFamily: "'Roboto', 'Droid Sans', 'Helvetica Neue', Helvetica, Arial, sans-serif",
  fontSize: "14px",
  lineHeight: "1.5",
  textRendering: "optimizeLegibility",
  maxWidth: "400px",
  margin: "0 auto",
}

export function App(props) {
    return (
        <div style={baseStyle}>
            <Header/>
            <Add/>
            <Feed/>
        </div>
    );
}
App = Radium(App)
