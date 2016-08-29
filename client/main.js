import React from 'react';
import { render } from 'react-dom';
import { App } from './app';
import { Feed } from './components/feed/feed'
import { Me } from './components/me/me'
import { FirstRun } from './components/first-run/first-run'
import Colors from './colors';
import Radium from 'radium';
import {StyleRoot} from 'radium';

import { createStore, applyMiddleware, compose } from 'redux';
import reducer from './reducer';
import { fetchStream, loadProfile } from './actions';
import {actions as meActions} from './resources/me';

import { Provider } from 'react-redux';
import thunkMiddleware from 'redux-thunk';

import { syncHistoryWithStore, routerMiddleware, push } from 'react-router-redux'
import { Router, Route, browserHistory, IndexRoute } from 'react-router'
import promiseMiddleware from 'redux-promise-middleware'


const logger = store => next => action => {
  console.log('dispatching', action)
  let result = next(action)
  console.log('next state', store.getState())
  return result
}



const store = createStore(
  reducer,
  compose(applyMiddleware(thunkMiddleware,
    logger, promiseMiddleware(), routerMiddleware(browserHistory)),
  window.devToolsExtension ? window.devToolsExtension() : f => f)
);

const history = syncHistoryWithStore(browserHistory, store)
store.dispatch(meActions.getMe());

function firstRun() {
  return dispatch => {
    fetch(
      "/api/self",
      { credentials: 'same-origin' }
    ).then(
      res => res.json()
    ).then(
      data => dispatch({type: "GET_SELF_SUCCESS", payload: data}),
      err => {
        dispatch({type: "GET_SELF_FAILED", payload: err});
        dispatch(push('/first-run'));
      }
    )
  }
}
store.dispatch(firstRun());

var baseStyle = {
  color: Colors.primaryText,
  fontFamily: "'Roboto', 'Droid Sans', 'Helvetica Neue', Helvetica, Arial, sans-serif",
  fontSize: "16px",
  lineHeight: "1.5",
  textRendering: "geometricPrecision",
  maxWidth: "650px",
  margin: "0px",
  '@media (min-width: 800px)': {
    margin: "20px 100px",
  },
  '@media (min-width: 1000px)': {
    margin: "50px 200px",
  },
  padding: "16px",
  background: Colors.primaryBackground,
};


render(
  <Provider store={store}>
    <StyleRoot style={baseStyle}>
      <Router history={history}>
        <Route path="/" component={App}>
          <IndexRoute component={Feed} onEnter={() => store.dispatch(fetchStream(30, 0))}/>
          <Route path="me" component={Me} onEnter={() => store.dispatch(meActions.getMe())}/>
        </Route>
        <Route path="first-run" component={FirstRun} />
      </Router>
    </StyleRoot>
  </Provider>,
  document.getElementById('app')
);
