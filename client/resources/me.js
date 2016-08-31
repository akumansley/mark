import Immutable from 'immutable';

const origin = window.location.origin;

function getMe() {
  return dispatch => {
    fetch("/api/me", {
      credentials: 'same-origin'
    }).then(res => res.json())
    .then(json => dispatch({type:"FETCH_ME_SUCCESS", payload: Immutable.fromJS(json)}),
          err => dispatch({type:"FETCH_ME_FAILED", payload: err}))
  }
}

function updateMe(body) {
  return dispatch => {
    fetch("/api/me", {
      credentials: 'same-origin',
      method: "PUT",
      body: JSON.stringify(body)
    }).then(res => res.json())
    .then(json => dispatch({type:"UPDATE_ME_SUCCESS", payload: Immutable.fromJS(json)}),
          err => dispatch({type:"UPDATE_ME_FAILED", payload: err}))
  }
}

export default {
  getMe,
  updateMe
}
