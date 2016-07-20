import React from 'react';
import {List, Map} from 'immutable';
import {Header} from './components/header/header';
import Radium from 'radium';
import {StyleRoot} from 'radium';
import Colors from './colors';

var baseStyle = {
  color: Colors.primaryText,
  fontFamily: "'Roboto', 'Droid Sans', 'Helvetica Neue', Helvetica, Arial, sans-serif",
  fontSize: "16px",
  lineHeight: "1.5",
  textRendering: "geometricPrecision",
  maxWidth: "650px",
  margin: "0",
  '@media (min-width: 800px)': {
    margin: "0 100px",
  },
  '@media (min-width: 1000px)': {
    margin: "0 200px",
  },
  borderLeft: "1px solid #eee",
  borderRight: "1px solid #eee",
  padding: "24px",
  background: Colors.primaryBackground,
}

export function App(props) {
    return (
        <StyleRoot style={baseStyle}>
            <Header/>
            {props.children}
        </StyleRoot>
    );
}
App = Radium(App)
