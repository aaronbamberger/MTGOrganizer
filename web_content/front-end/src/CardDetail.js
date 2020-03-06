import React from 'react';
import {withRouter, Link} from 'react-router-dom';

import {REQUEST_TYPES} from './Constants.js';
import {CardImage} from './CardComponents.js';


class CardDetail extends React.Component {
  constructor(props) {
    super(props);

    this.receiveCardDetail = this.receiveCardDetail.bind(this);

    this.uuid = props.match.params.uuid;
    this.state = {CardDetail: null};
  }

  componentDidMount() {
    this.sendRequestForUpdatedCardInfo();
  }

  componentDidUpdate() {
    if (this.uuid != this.props.match.params.uuid) {
      this.uuid = this.props.match.params.uuid;
      this.sendRequestForUpdatedCardInfo();
    }
  }

  receiveCardDetail(newCardDetail) {
    this.setState({CardDetail: newCardDetail});
  }

  sendRequestForUpdatedCardInfo() {
    const request = JSON.stringify(
      {
        "type": this.props.apiTypesMap.CardDetailRequest,
        "value": this.uuid
      });
    this.props.backendRequest(request);
  }

  render() {
    if (this.state.CardDetail) {
      let variations = [ this.uuid ].concat(this.state.CardDetail.variations);
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
      if (Object.keys(this.state.CardDetail.legalities).length > 0) {
        legalities = (
          <SortedKeyValueTable
            keyHeader="Game Format"
            valueHeader="Legality"
            data={this.state.CardDetail.legalities} />
        );
      }

      let leadershipSkills = null
      if (Object.keys(this.state.CardDetail.leadershipSkills).length > 0) {
        leadershipSkills = (
          <SortedKeyValueTable
            keyHeader="Game Format"
            valueHeader="Leader Legal?"
            data={this.state.CardDetail.leadershipSkills} />
        );
      }

      return (
        <div>
          <div>Name: {this.state.CardDetail.name}</div>
          <div>Artist: {this.state.CardDetail.artist}</div>
          <CardImage
            uuid={this.uuid}
            name={this.state.CardDetail.name}
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

export default withRouter(CardDetail);
