import fetch from 'isomorphic-fetch'
import Immutable from 'immutable';
import { createAction } from 'redux-actions';
import { isWebUri } from 'valid-url';
import _ from 'lodash';

export const REQUEST_STREAM = 'REQUEST_STREAM';
function requestStream() {
  return {
    type: REQUEST_STREAM,
  }
}

export const FETCH_STREAM_SUCCESS = 'FETCH_STREAM_SUCCESS';
function fetchStreamSuccess(items) {
  return {
    type: FETCH_STREAM_SUCCESS,
    payload: items,
  }
}

export const FETCH_STREAM_FAILED = 'FETCH_STREAM_FAILED';
function fetchStreamFailed(error) {
  return {
    type: FETCH_STREAM_FAILED,
    payload: error,
    error: true
  }
}

// this is a thunk
export function fetchStream(count, offset, feedId) {
  return function (dispatch) {
    // start the request
    dispatch(requestStream());

    let qs = "?count=" + encodeURIComponent(count) +
             "&offset=" + encodeURIComponent(offset);

    if (feedId) {
      qs = qs + "&feedId=" + encodeURIComponent(feedId);
    } else {
      feedId = "me";
    }

    return fetch('/api/stream' + qs, {credentials: 'same-origin'})
            .then(res => {
              if (res.status >= 400) {
                throw new Error(res.status);
              }
              return res.json();
            })
            .then(json => dispatch(fetchStreamSuccess(
              Immutable.Map({
                feedId: feedId,
                items: Immutable.fromJS(json),
              })
            )))
            .catch(err => dispatch(fetchStreamFailed(err)));
  }
}

// experimentin' here
const ADD_MARK_SUCCESS = "ADD_MARK_SUCCESS"
const ADD_MARK_FAILED = "ADD_MARK_FAILED"
const POST_MARK = "POST_MARK"

const postMark = createAction(POST_MARK);
const addMarkSuccess = createAction(ADD_MARK_SUCCESS);
const addMarkFailed = createAction(ADD_MARK_FAILED);

export function addMark(url, title) {
  return dispatch => {
    dispatch(postMark());
    if (!isWebUri(url)) {
      dispatch(addMarkFailed(new Error("Invalid URL")))
      return;
    }

    return fetch('/api/bookmark', { method: "POST",
      credentials: 'same-origin',
      body: JSON.stringify({
        url: url,
        title: title
      })
    }).then(res => {
      if (res.status >= 400) {
        throw new Error(res.status);
      }
      return res.json();
    }).then(json => {
      dispatch(fetchStream(30, 0));
      dispatch(addMarkSuccess());
    }).catch(err => dispatch(addMarkFailed(err)));
  }
}


const REMOVE_MARK_SUCCESS = "REMOVE_MARK_SUCCESS"
const REMOVE_MARK_FAILED = "REMOVE_MARK_FAILED"
const REMOVE_MARK = "REMOVE_MARK"

const removeMarkSuccess = createAction(REMOVE_MARK_SUCCESS);
const removeMarkFailed = createAction(REMOVE_MARK_FAILED);

export function removeMark(id) {
  return dispatch => {
    if (!id) {
      dispatch(removeMarkFailed(new Error("Cannot delete null id")));
      return
    }
    return fetch('/api/bookmark/' + id, { method: "DELETE",  credentials: 'same-origin'})
    .then(res => {
      if (res.status >= 400) {
        throw new Error(res.status);
      }
      return res.text();
    }).then(json => {
      dispatch(removeMarkSuccess(id));
    }).catch(err => dispatch(removeMarkFailed(err)));
  }
}


const UPDATE_URL = "UPDATE_URL";
const LOAD_TITLE = "LOAD_TITLE";
const LOAD_TITLE_SUCCESS = "LOAD_TITLE_SUCCESS";
const LOAD_TITLE_FAILED = "LOAD_TITLE_FAILED";

const UPDATE_TITLE = "UPDATE_TITLE";

const loadTitleSuccess = createAction(LOAD_TITLE_SUCCESS);
const loadTitleFailed = createAction(LOAD_TITLE_FAILED);

export const updateTitle = createAction(UPDATE_TITLE);

function loadTitleRaw(url) {
  return dispatch => {
    let qs = "?url=" + encodeURIComponent(url)
    return fetch('/views/title' + qs, {credentials: 'same-origin'})
    .then(res => {
      if (res.status >= 400) {
        throw new Error(res.status);
      }
      return res.text();
    }).then(title => dispatch(loadTitleSuccess(title)))
    .catch(err => dispatch(loadTitleFailed(err)));
  }
}

const loadTitle = _.throttle(loadTitleRaw, 300);

export function updateUrl(url) {
  return dispatch =>  {
    dispatch({
      type: UPDATE_URL,
      payload: url
    });
    if (isWebUri(url)) {
      dispatch(loadTitle(url));
    }
  }
}
