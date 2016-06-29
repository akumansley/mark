import React from 'react';
import { render } from 'react-dom';
import { App } from './app';

import { createStore } from 'redux';
import reducer from './reducer';
import { Provider } from 'react-redux';

const store = createStore(reducer);

render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('app')
);
