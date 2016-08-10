import { List, Map } from 'immutable';
import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux'


const initBookmarks = Map({
  items: Map({}),
  loading: false,
  error: null,
  hasMore: true
});


function bookmarks(state=initBookmarks, action) {
  switch (action.type) {
    case 'FETCH_STREAM':
      return state.set('loading', true);
    case 'FETCH_STREAM_SUCCESS':
      const notLoading = state.set('loading', false);
      const mergedItems = notLoading.update('items', items => (
                    items.withMutations(m => {
                      for (let i of action.payload.values()) {
                        m.set(i.get('id'), i)
                      }
                    })
                  ));
      const hasMore = mergedItems.set('hasMore', action.payload.size? true: false);
      return hasMore;
    case 'FETCH_STREAM_FAILED':
      return state.set('loading', false)
                  .set('error', action.payload);
    case 'POST_MARK':
      return state.set('loading', true);
    case 'ADD_MARK_SUCCESS':
      return state.set('loading', false);
    case 'ADD_MARK_FAILED':
      return state.set('loading', false)
                  .set('error', action.payload);
    case 'REMOVE_MARK_SUCCESS':
      return state.update('items', items => (items.delete(action.payload)));
    default:
      return state;
  }
}

function showTitle(state=false, action) {
  switch (action.type) {
    case 'UPDATE_URL':
      return action.payload !== "";
    case 'ADD_MARK_SUCCESS':
      return false;
    default:
      return state;
  }
}

function url(state="", action) {
  switch(action.type) {
    case 'UPDATE_URL':
      return action.payload;
    case 'ADD_MARK_SUCCESS':
      return "";
    case 'ADD_MARK_FAILED':
      return state;
    default:
      return state
  }
}

// Probably want to track a bit saying whether the user's edited the current title
function title(state="", action) {
  switch(action.type) {
    case 'LOAD_TITLE_SUCCESS':
      return action.payload;
    case 'LOAD_TITLE_FAILED':
      return state;
    case 'ADD_MARK_SUCCESS':
      return "";
    case 'UPDATE_TITLE':
      return action.payload;
    case 'UPDATE_URL':
      if (action.payload === "") {
        return ""
      }
      return state
    default:
      return state;
  }
}

function me(state=Map({}), action) {
  switch (action.type) {
    case 'LOAD_PROFILE_SUCCESS':
      return action.payload;
    case 'UPDATE_PROFILE_SUCCESS':
      return action.payload;
    case 'UPDATE_NAME':
      return state.set('name', action.payload)
    default:
      return state
  }
}

const reducer = combineReducers({
  bookmarks,
  showTitle,
  url,
  title,
  me,
  routing: routerReducer
});

export default reducer;
