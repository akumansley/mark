import { List, Map } from 'immutable';

const init = List([]);

export default function(marks=init, action) {
  switch(action.type) {
    case 'ADD_MARK':
      return marks.push(Map(action.payload));
    default:
      return marks
  }
}
