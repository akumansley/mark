import  { connect } from 'react-redux';
import * as components from './components';
import { fetchFeed } from './actions';

export const Feed = connect(
  function mapStateToProps(state) {
    return { items: state.bookmarks.get('items') }
  }
)(components.Feed);

export const Add = connect(null,
  function mapDispatchToProps(dispatch) {
    return { addMark: url => {dispatch(addMark(url));} }
  }
)(components.Add);
