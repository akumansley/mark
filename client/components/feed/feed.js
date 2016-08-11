import React from 'react'
import Radium from 'radium'
import Colors from '../../colors'
import { fetchStream, removeMark } from '../../actions'
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
const deleteStyle = {
  padding: '4px 4px 4px 8px',
}

const RawItem = React.createClass({
    render: function() {
      const {removeMark, me, item: i} = this.props;
      function clickDelete(evt) {
        removeMark(i.get('id'));
      }

      let delNode = null;
      if (me.get('feed_id') == i.get('feed_id')) {
        delNode = <div style={deleteStyle}>
            <a onClick={clickDelete}>&times;</a>
          </div>;
      }

      return (
        <div style={itemStyle} key={i.get('id')}>
          <div style={leftStyle}>
            <a href={i.get('url')} style={titleStyle}>{i.get('title')}</a>
            <span style={urlStyle}>
            {i.get('profile').get('name')} - {i.get('short_url')} </span>
          </div>
          {delNode}
        </div>
      );
    }
});

const FeedItem = connect(
  function mapStateToProps(state) {
    return {
      me: state.me
    }
  },
  function mapDispatchToProps(dispatch) {
    return {
      removeMark: (id) => dispatch(removeMark(id)),
    }
  }
)(RawItem);

const PAGE_SIZE = 30;
const TRIGGER_THRESHOLD = 100;

const Component = React.createClass({
  componentWillMount: function () {
    window.addEventListener("scroll", this.handleScroll);
  },

  handleScroll: function (evt) {
    // From http://stackoverflow.com/questions/1145850/how-to-get-height-of-entire-document-with-javascript
    var body = document.body
    var html = document.documentElement;
    var totalHeight = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight,
      html.scrollHeight, html.offsetHeight );

    const innerHeight = window.innerHeight;
    const scrollTop = window.scrollY;

    var totalScrolled = scrollTop + innerHeight;
    if (totalScrolled + TRIGGER_THRESHOLD > totalHeight &&
      !this.props.loading && this.props.hasMore) {
      this.loadMore();
    }
  },


  buildElements: function() {
      return this.props.items.map(item => (
        <FeedItem key={item.get('id')} item={item}/>
      ));
  },

  loadMore: function() {
      this.props.fetchStream(PAGE_SIZE, this.props.items.size);
  },

  render: function() {
    return (
      <div>
        <Add />
        {this.buildElements()}
      </div>
    )
  }
});


const Styled = Radium(Component)

const selectItems = state => state.bookmarks.get('items');
const sortItems = createSelector([selectItems], items => items.valueSeq().sortBy(v => -1 * v.get('created_at')).toList());

function shortUrl(url) {
  const u = new URL(url);
  let r = u.hostname + u.pathname;
  r = r.endsWith("/") ? r.slice(0, -1) : r;
  r = r.endsWith(".html") ? r.slice(0, -5) : r;
  return r;
}

const mixShortUrl = createSelector([sortItems], items => {
  return items.map(i => i.set('short_url', shortUrl(i.get('url')) ));
});

const Connected = connect(
  function mapStateToProps(state) {
    return {
      items: mixShortUrl(state),
      loading: state.bookmarks.get('loading'),
      hasMore: state.bookmarks.get('hasMore'),
    }
  },
  function mapDispatchToProps(dispatch) {
    return {
      fetchStream: (count, offset) => dispatch(fetchStream(count, offset)),
    }

  }
)(Styled);

export const Feed = Connected;
