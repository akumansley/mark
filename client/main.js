import React from 'react';
import { render } from 'react-dom';
import { App } from './app';
import { Feed } from './components/feed/feed'
import { Me } from './components/me/me'

import { createStore, applyMiddleware, compose } from 'redux';
import reducer from './reducer';
import { fetchStream, loadProfile } from './actions';
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
  compose(applyMiddleware(thunkMiddleware, logger),
  window.devToolsExtension ? window.devToolsExtension() : f => f)
);

const history = syncHistoryWithStore(browserHistory, store)

render(
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={App}>
        <IndexRoute component={Feed} onEnter={() => store.dispatch(fetchStream(30, 0))}/>
        <Route path="me" component={Me} onEnter={() => store.dispatch(loadProfile())}/>
      </Route>
    </Router>
  </Provider>,
  document.getElementById('app')
);
