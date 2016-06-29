import React from 'react';
import {List, Map} from 'immutable';
import {Header} from './components';
import {Add, Marks} from './containers';
import Radium from 'radium';

var baseStyle = {
  color: "#333",
  fontFamily: "-apple-system, BlinkMacSystemFont, 'Roboto', 'Droid Sans', 'Helvetica Neue', Helvetica, Arial, sans-serif",
  fontSize: "15px",
  lineHeight: "1.5",
  textRendering: "optimizeLegibility",
  maxWidth: "600px",
  margin: "0 auto",
}

export function App(props) {
    return (
        <div style={baseStyle}>
            <Header/>
            <Add/>
            <Marks/>
        </div>
    );
}
App = Radium(App)
