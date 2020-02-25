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
        "type": REQUEST_TYPES.CARD_DETAIL_REQUEST,
        "value": this.uuid
      });
    this.props.backendRequest(request);
  }

  render() {
    if (this.state.CardDetail) {
      let variationLinks = null
      if (this.state.CardDetail.variations.length > 0) {
        variationLinks = (
          <div>
            <Link to={"/card/" + this.uuid}>1</Link>
            {this.state.CardDetail.variations.map((variationUUID, i) =>
              <Link to={"/card/" + variationUUID}>{i + 2}</Link>)};
          </div>
        );
      }

      return (
        <div>
          <div>Name: {this.state.CardDetail.name}</div>
          <div>Artist: {this.state.CardDetail.artist}</div>
          <CardImage uuid={this.uuid} name={this.state.CardDetail.name} />
          {variationLinks}
        </div>
      );
    } else {
      return (<div></div>);
    }
  }
}

export default withRouter(CardDetail);
