import  { connect } from 'react-redux';
import * as components from './components';
import { fetchFeed, addMark, showTitle, hideTitle } from './actions';

export const Feed = connect(
  function mapStateToProps(state) {
    return { items: state.bookmarks.get('items') }
  }
)(components.Feed);

export const Add = connect(
  function mapStateToProps(state) {
    return { shouldShowTitle: state.showTitle };
  },
  function mapDispatchToProps(dispatch) {
    return {
      addMark: url => {dispatch(addMark(url, "TODO: Implement"))},
      showTitle: () => {dispatch(showTitle())},
      hideTitle: () => {dispatch(hideTitle())},
    }
  }
)(components.Add);
