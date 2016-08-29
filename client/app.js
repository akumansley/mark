import React from 'react';
import {List, Map} from 'immutable';
import {Header} from './components/header/header';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';

export function App(props) {
    return (
      <div>
        <Header/>
        <ReactCSSTransitionGroup
            component="div"
            transitionName="fade"
            transitionAppear={true}
            transitionEnter={true}
            transitionAppearTimeout={200}
            transitionEnterTimeout={200}
            transitionLeaveTimeout={200}>
            {React.cloneElement(props.children, {
              key: location.pathname
            })}
        </ReactCSSTransitionGroup>
      </div>
    );
}
