import React from 'react'
import Radium from 'radium'
import Colors from '../../colors'
import { fetchStream } from '../../actions'
import { connect } from 'react-redux'
import { Add } from '../add/add'
import { createSelector } from 'reselect'
import Infinite from 'react-infinite'

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

const FeedItem = React.createClass({
    render: function() {
      const i = this.props.item;
      return (
        <div style={itemStyle} key={i.get('id')}>
          <div style={leftStyle}>
            <a href={i.get('url')} style={titleStyle}>{i.get('title')}</a>
            <span style={urlStyle}>
            {i.get('profile').get('name')} - {i.get('short_url')} </span>
          </div>
        </div>
      );
    }
});

const PAGE_SIZE = 30;

const Component = React.createClass({
    getInitialState: function() {
        return {
            elements: this.buildElements(),
        }
    },

    buildElements: function(start, end) {
        return this.props.items.map(item => (
          <FeedItem key={item.get('id')} item={item}/>
        ));
    },

    handleInfiniteLoad: function() {
        this.props.fetchStream(PAGE_SIZE, this.props.items.size);
    },

    componentWillReceiveProps: function (newProps) {
      this.setState({
        elements: this.buildElements()
      });
    },

    render: function() {
      return (
        <div>
          <Add />
          <Infinite elementHeight={40}
            infiniteLoadBeginEdgeOffset={200}
            onInfiniteLoad={this.handleInfiniteLoad}
            isInfiniteLoading={this.props.loading}
            useWindowAsScrollContainer={true}>
              {this.state.elements}
          </Infinite>
        </div>
      )
    }
});


const Styled = Radium(Component)

const selectItems = state => state.bookmarks.get('items').valueSeq().sortBy(v => v.created_at);

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
    return {
      items: mixShortUrl(state),
      loading: state.bookmarks.get('loading'),
    }
  },
  function mapDispatchToProps(dispatch) {
    return {
      fetchStream: (count, offset) => dispatch(fetchStream(count, offset))
    }

  }
)(Styled);

export const Feed = Connected;
