import fetch from 'isomorphic-fetch'
import Immutable from 'immutable';

const uid = () => Math.random().toString(34).slice(2);

export const REQUEST_FEED = 'REQUEST_FEED';
function requestFeed() {
  return {
    type: REQUEST_FEED,
  }
}

export const FETCH_FEED_SUCCESS = 'FETCH_FEED_SUCCESS';
function fetchFeedSuccess(items) {
  return {
    type: FETCH_FEED_SUCCESS,
    items: items,
  }
}

export const FETCH_FEED_FAILED = 'FETCH_FEED_FAILED';
function fetchFeedFailed(error) {
  return {
    type: FETCH_FEED_FAILED,
    error: error
  }
}

// this is a thunk
export function fetchFeed() {
  return function (dispatch) {
    // start the request
    dispatch(requestFeed());

    return fetch('/api/feed')
            .then(res => {
              if (res.status >= 400) {
                throw new Error(res.status);
              }
              return res.json();
            })
            .then(json => dispatch(fetchFeedSuccess(Immutable.fromJS(json))))
            .catch(err => dispatch(fetchFeedFailed(err)));
  }
}


export function addMark(url) {
  return {
    type: 'ADD_MARK',
    payload: {
      id: uid(),
      url: url
    }
  };
}
