import  { connect } from 'react-redux';
import * as components from './components';
import { addMark } from './actions';

export const Marks = connect(
  function mapStateToProps(state) {
    return { items: state }
  }
)(components.Marks);

export const Add = connect(null,
  function mapDispatchToProps(dispatch) {
    return { addMark: url => {dispatch(addMark(url));} }
  }
)(components.Add);
