import React from 'react';
import { Link, IndexLink } from 'react-router'
import Styles from '../../styles';
import Radium from 'radium';
const m = require('../../assets/m.png')
import { push } from 'react-router-redux'
import {connect} from 'react-redux';


const imgStyle = {
  width: 30,
  height: 30,
  margin: "0 auto 30px auto",
  display: "block"
}

const offerIframeStyle = {
  width: "100%",
  height: "55px",
  margin: 0,
  border: 0,
}

function requestIframeURL() {
   var template = "$API_HOST#$API_TOKEN";
   console.log("req");
   window.parent.postMessage({renderTemplate: {
     rpcId: "0",
     template: template,
     clipboardButton: 'right'
   }}, "*");
 }

 var copyIframeURLToElement = function(event) {
  if (event.data.rpcId === "0") {
    if (event.data.error) {
      console.log("ERROR: " + event.data.error);
    } else {
      var el = document.getElementById("offer-iframe");
      el.setAttribute("src", event.data.uri);
    }
  }
};

const inputStyles = Object.assign({
  width: "100%",
}, Styles.input)

function putSelf(url) {
  return dispatch => {
    fetch(
      "/api/self",
      {
        credentials: 'same-origin',
        method: 'PUT',
        body: JSON.stringify({
          url: url,
        }),
      }
    ).then(
      res => res.json()
    ).then(
      data => {
        dispatch({type: "PUT_SELF_SUCCESS", payload: data})
        dispatch(push('/'));
      },
      err => {
        dispatch({type: "PUT_SELF_FAILED", payload: err});
      }
    )
  }
}


class Component extends React.Component {
  render() {
    const {putSelf} = this.props;

    const submit = evt => {
      const url = this.urlInput.value;
      putSelf(url);
    }

    return (
      <div>
        <img src={m} style={imgStyle}/>
        <p>Mark's a network of people sharing links with one another.
        It's a bit different than you're used to. Here's what you need to know.</p>
        <ul>
          <li>Mark is distributed, so nobody owns or controls it</li>
          <li>Mark can't be bought, closed or shut down, or show ads</li>
          <li>Deleting a bookmark hides it, but people can still find it if they go looking</li>
        </ul>
        <p><strong>This is the address of your Mark server</strong></p>
        <iframe style={offerIframeStyle} id="offer-iframe"></iframe>
        <p><strong>Copy and paste it here to get started:</strong></p>

        <input key="url"
          ref={el => this.urlInput = el}
          style={inputStyles}
          placeholder="http://.."
          type="text">
        </input>

        <button onClick={submit} style={Styles.actionButton}>OK</button>
      </div>
    )
  }
  componentWillMount() {
    requestIframeURL();
    window.addEventListener("message", copyIframeURLToElement);
  }
}

export const FirstRun = connect(null, function mapDispatchToProps(dispatch) {
    return {
        putSelf: (url) => {
            dispatch(putSelf(url))
        },
    }
})(Component);
