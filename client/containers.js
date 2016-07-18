import  { connect } from 'react-redux';
import * as components from './components';
import { fetchFeed, addMark, showTitle, hideTitle, updateUrl } from './actions';

export const Feed = connect(
  function mapStateToProps(state) {
    return { items: state.bookmarks.get('items') }
  }
)(components.Feed);

export const Add = connect(
  function mapStateToProps(state) {
    return {
      shouldShowTitle: state.showTitle,
      url: state.url
    };
  },
  function mapDispatchToProps(dispatch) {
    return {
      addMark: (url, title) => {dispatch(addMark(url, title))},
      showTitle: () => {dispatch(showTitle())},
      hideTitle: () => {dispatch(hideTitle())},
      updateUrl: url => {dispatch(updateUrl(url))},
    }
  }
)(components.Add);
