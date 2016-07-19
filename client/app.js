import React from 'react';
import {List, Map} from 'immutable';
import {Header} from './components/header/header';
import Radium from 'radium';
import Colors from './colors';

var baseStyle = {
  color: Colors.primaryText,
  fontFamily: "'Roboto', 'Droid Sans', 'Helvetica Neue', Helvetica, Arial, sans-serif",
  fontSize: "16px",
  lineHeight: "1.5",
  textRendering: "geometricPrecision",
  maxWidth: "600px",
  margin: "0 auto",
  background: Colors.background,
}

export function App(props) {
    return (
        <div style={baseStyle}>
          <Header/>
          {props.children}
        </div>
    );
}
App = Radium(App)
