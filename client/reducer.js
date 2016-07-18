import { List, Map } from 'immutable';
import { combineReducers } from 'redux';

const initBookmarks = Map({
  items: List([]),
  loading: false,
  error: null
});


function bookmarks(state=initBookmarks, action) {
  switch (action.type) {
    case 'FETCH_FEED':
      return state.set('loading', true);
    case 'FETCH_FEED_SUCCESS':
      return state.set('loading', false)
                  .set('items', action.payload);
    case 'FETCH_FEED_FAILED':
      return state.set('loading', false)
                  .set('error', action.payload);
    case 'POST_MARK':
      return state.set('loading', true);
    case 'ADD_MARK_SUCCESS':
      return state.set('loading', false)
                  .update('items', i => i.push(action.payload));
    case 'ADD_MARK_FAILED':
      return state.set('loading', false)
                  .set('error', action.payload);
    default:
      return state;
  }
}

function showTitle(state=false, action) {
  switch (action.type) {
    case 'SHOW_TITLE':
      return true;
    case 'HIDE_TITLE':
      return false;
    default:
      return state;
  }
}

const reducer = combineReducers({
  bookmarks,
  showTitle
});

export default reducer;
