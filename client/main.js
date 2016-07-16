import React from 'react';
import { render } from 'react-dom';
import { App } from './app';

import { createStore, applyMiddleware } from 'redux';
import reducer from './reducer';
import { fetchFeed } from './actions';
import { Provider } from 'react-redux';
import thunkMiddleware from 'redux-thunk';

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

store.dispatch(fetchFeed());

render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('app')
);
