import  { connect } from 'react-redux';
import * as components from './components';
import { fetchFeed, addMark } from './actions';

export const Feed = connect(
  function mapStateToProps(state) {
    return { items: state.bookmarks.get('items') }
  }
)(components.Feed);

export const Add = connect(null,
  function mapDispatchToProps(dispatch) {
    return { addMark: url => {dispatch(addMark(url, "TODO: Implement"));} }
  }
)(components.Add);
