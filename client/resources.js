import {createResource} from 'redux-rest-resource';
import Immutable from 'immutable';

const origin = window.location.origin;

export const {types, actions, reducers} = createResource({
  name: 'me',
  url: `${origin}/api/me/`,
  actions: {
    get: {
      credentials: 'same-origin',
      transformResponse: (res) => {
        res.body = Immutable.fromJS(res.body);
        console.log(res);
        return res;
      }
    },
    update: {
      credentials: 'same-origin',
      method: "PUT"
    }
  }
});
