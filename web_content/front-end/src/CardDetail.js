import React from 'react';
import { withRouter, NavLink } from 'react-router-dom';
import { connect } from 'react-redux';

import { CardImage } from './CardComponents.js';
import { requestCardDetail } from './ReduxActions.js';

const mapStateToProps = (state) => {
  return {
    backendReady: state.backend.ready,
    cardUUID: state.cardDetail.uuid,
    cardDetail: state.cardDetail.cardDetail,
  };
}

class CardDetail extends React.Component {
  componentDidMount() {
    if (this.props.backendReady) {
      this.sendRequestForUpdatedCardInfo();
    }
  }

  componentDidUpdate() {
    if (this.props.cardDetail) {
      if (this.props.cardUUID !== this.props.match.params.uuid) {
        console.log("Sent request due to uuid mismatch");
        this.sendRequestForUpdatedCardInfo();
      }
    } else {
      // If we haven't received a card detail object yet, check to see
      // if the backend just became ready, and if so, sent a request for the
      // detail info
      if (this.props.backendReady) {
        console.log("Sent request due to backend ready");
        this.sendRequestForUpdatedCardInfo();
      }
    }
  }

  sendRequestForUpdatedCardInfo() {
    const uuid = this.props.match.params.uuid;
    this.props.dispatch(requestCardDetail(uuid));
  }

  render() {
    if (this.props.cardDetail) {
      let variations = [ this.props.cardUUID ].concat(this.props.cardDetail.variations);
      console.log("Card UUID: " + this.props.cardUUID);
      console.log("Card Variations: ", this.props.cardDetail.variations);
      console.log("Variations List: ", variations);
      variations.sort();
      let variationLinks = null
      if (variations.length > 1) {
        variationLinks = (
          <div>
            Variations:
            {variations.map((variationUUID, i) =>
              <NavLink
                  key={variationUUID}
                  to={"/card/" + variationUUID}
                  activeStyle={{opacity: "50%", pointerEvents: "none"}}>
                <CardImage
                  uuid={variationUUID}
                  sizePercent={5} />
              </NavLink>)}
              &nbsp;
          </div>
        );
      }

      let legalities = null
      if (Object.keys(this.props.cardDetail.legalities).length > 0) {
        legalities = (
          <SortedKeyValueTable
            keyHeader="Game Format"
            valueHeader="Legality"
            data={this.props.cardDetail.legalities} />
        );
      }

      let leadershipSkills = null
      if (Object.keys(this.props.cardDetail.leadershipSkills).length > 0) {
        leadershipSkills = (
          <SortedKeyValueTable
            keyHeader="Game Format"
            valueHeader="Leader Legal?"
            data={this.props.cardDetail.leadershipSkills} />
        );
      }

      return (
        <div>
          <div>Name: {this.props.cardDetail.name}</div>
          <div>Artist: {this.props.cardDetail.artist}</div>
          <CardImage
            uuid={this.props.cardUUID}
            name={this.props.cardDetail.name}
            sizePercent={20} />
          {variationLinks}
          {legalities}
          {leadershipSkills}
        </div>
      );
    } else {
      return (<div></div>);
    }
  }
}

function SortedKeyValueTable(props) {
  const keys = Object.keys(props.data);
  keys.sort();

  const rows = keys.map((key) =>
    <KeyValueTableRow
      key={key}
      dataKey={key}
      dataValue={props.data[key]} />
  );

   return (
    <div id={props.tableName}>
      <table>
        <thead>
          <tr>
            <th>{props.keyHeader}</th>
            <th>{props.valueHeader}</th>
          </tr>
        </thead>
        <tbody>
          {rows}
        </tbody>
      </table>
    </div>
  );
}

function KeyValueTableRow(props) {
  return (
    <tr>
      <td>{props.dataKey}</td>
      <td>{props.dataValue.toString()}</td>
    </tr>
  );
}

export default connect(mapStateToProps)(withRouter(CardDetail));
