import React from 'react';
import { render } from 'react-dom';
import { App } from './app';
import { Feed } from './components/feed/feed'
import { Me } from './components/me/me'

import { createStore, applyMiddleware } from 'redux';
import reducer from './reducer';
import { fetchStream } from './actions';
import { Provider } from 'react-redux';
import thunkMiddleware from 'redux-thunk';

import { syncHistoryWithStore } from 'react-router-redux'
import { Router, Route, browserHistory, IndexRoute } from 'react-router'


const logger = store => next => action => {
  console.log('dispatching', action)
  let result = next(action)
  console.log('next state', store.getState())
  return result
}



const store = createStore(
  reducer,
  applyMiddleware(thunkMiddleware, logger)
);

const history = syncHistoryWithStore(browserHistory, store)
store.dispatch(fetchStream());

render(
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={App}>
        <IndexRoute component={Feed} />
        <Route path="me" component={Me} />
      </Route>
    </Router>
  </Provider>,
  document.getElementById('app')
);
