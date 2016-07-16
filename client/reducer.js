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
                  .set('items', action.items);
    case 'FETCH_FEED_FAILED':
      return state.set('loading', false)
                  .set('error', action.error);
    case 'ADD_BOOKMARK':
      return state.push(Map(action.payload));
    case 'ADD_BOOKMARK_SUCCESS':
      return state.push(Map(action.payload));
    case 'ADD_BOOKMARK_FAILED':
      return state.push(Map(action.payload));
    default:
      return state;
  }
}

const reducer = combineReducers({
  bookmarks,
});

export default reducer;
