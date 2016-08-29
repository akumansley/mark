import React from 'react';
import { Link, IndexLink } from 'react-router'
import Styles from '../../styles';
import Radium from 'radium';
const m = require('../../assets/m.png')

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


class Component extends React.Component {
  render() {
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
        <p><strong>Here's the address of your Mark server:</strong></p>
        {/* <input type="text" /> */}
        <iframe style={offerIframeStyle} id="offer-iframe"></iframe>
        <p><strong>Copy and paste it here to get started:</strong></p>

        <input key="title" style={Styles.input} placeholder="http://.." type="text"></input>

        <button style={Styles.actionButton}>OK</button>
      </div>
    )
  }
  componentWillMount() {
    requestIframeURL();
  }
}


export const FirstRun = Component;
