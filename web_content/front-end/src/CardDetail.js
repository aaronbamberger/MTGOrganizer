import React from 'react';
import { withRouter, Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { CardImage } from './CardComponents.js';
import { requestCardDetail } from './ReduxActions.js';

const mapStateToProps = (state) => {
  return {
    cardUUID: state.cardDetail.uuid,
    cardDetail: state.cardDetail.cardDetail,
  };
}

class CardDetail extends React.Component {
  componentDidMount() {
    this.sendRequestForUpdatedCardInfo();
  }

  componentDidUpdate() {
    if (this.props.cardDetail) {
      if (this.props.cardUUID !== this.props.match.params.uuid) {
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
      variations.sort();
      let variationLinks = null
      if (variations.length > 1) {
        variationLinks = (
          <div>
            Variations:
            {variations.map((variationUUID, i) =>
              <Link key={variationUUID} to={"/card/" + variationUUID}>
                <CardImage
                  uuid={variationUUID}
                  sizePercent={5}
                  isDisabled={this.uuid === variationUUID} />
              </Link>)}
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
