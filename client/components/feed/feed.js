import React from 'react'
import Radium from 'radium'
import Colors from '../../colors'
import { connect } from 'react-redux'
import { Add } from '../add/add'
import { createSelector } from 'reselect'

// var moreSrc = require('../../assets/more.png')

var itemStyle = {
  display: "flex",
  alignItems: "start",
  flexDirection: "row",
  paddingTop: 6,
  paddingBottom: 18,
}

var titleStyle = {
  lineHeight: "1.2",
  marginBottom: -2,
  display: "block",
  color: Colors.primaryText,
  textDecoration: "none",
};

var urlStyle = {
  fontSize: "13px",
  color: Colors.secondaryText,
  fontWeight: "200",
  overflowWrap: 'break-word',
  wordWrap: 'break-word',
  wordBreak: 'break-word',
};

var leftStyle = {
  flex: 1,
}

var moreStyle = {
  padding: '4px 8px',
  ":hover": {
    boxShadow: "1px 1px 1px 1px #eee" ,
  }
}

const Component = props => {
    const {items} = props;
    return (
      <div>
        <Add />
        <ul>
            {props.items.map(i => {
                return (
                    <div style={itemStyle} key={i.get('id')}>
                      <div style={leftStyle}>
                        <a href={i.get('url')} style={titleStyle}>{i.get('title')}</a>
                        <span style={urlStyle}> {i.get('short_url')}</span>
                      </div>
                      <div>
                      </div>
                    </div>
                )
            })}
        </ul>
      </div>
    )
}

const Styled = Radium(Component)

const selectItems = state => state.bookmarks.get('items');

function shortUrl(url) {
  const u = new URL(url);
  let r = u.hostname + u.pathname;
  r = r.endsWith("/") ? r.slice(0, -1) : r;
  r = r.endsWith(".html") ? r.slice(0, -5) : r;
  return r;
}

const mixShortUrl = createSelector([selectItems], items => {
  return items.map(i => i.set('short_url', shortUrl(i.get('url')) ));
});

const Connected = connect(
  function mapStateToProps(state) {
    return { items: mixShortUrl(state) }
  }
)(Styled);

export const Feed = Connected;
