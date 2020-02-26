import React from 'react';
import {Link} from "react-router-dom";

import {SetSymbol} from './CardComponents.js';
import {REQUEST_TYPES} from './Constants.js';

class CardSearch extends React.Component {
  constructor(props) {
    super(props);

    this.state = {cards: []}
    
    this.handleInput = this.handleInput.bind(this);
  }

  handleInput(event) {
    if (event.target.value.length > 1) {
      const request = JSON.stringify(
        {
          "type": this.props.apiTypesMap.CardSearchRequest,
          "value": event.target.value
        });
      this.props.backendRequest(request)
    } else {
      this.setState({cards: []});
    }
  }

  receiveNewCards(newCards) {
    this.setState({cards: newCards});
  }

  render() {
    return (
      <div id="card_search">
        <form>
          <label>
            Search for a Card:
            <input type="text" onInput={this.handleInput} />
          </label>
        </form>
        <CardSearchResultsTable cards={this.state.cards} />
      </div>
    );
  }
}

function CardSearchResultsTable(props) {
  const rows = props.cards.map((cardInfo) =>
    <CardSearchResultRow
      key={cardInfo.uuid}
      cardName={cardInfo.name}
      cardUUID={cardInfo.uuid}
      setName={cardInfo.setName}
      setCode={cardInfo.setKeyruneCode} />
    );

  return (
    <div id="card_results">
      <table>
        <thead>
          <tr>
            <th>Card Name</th>
            <th>Set</th>
          </tr>
        </thead>
        <tbody>
          {rows}
        </tbody>
      </table>
    </div>
  );
}

function CardSearchResultRow(props) {
    return (
      <tr>
        <td>
          <Link to={"/card/" + props.cardUUID}>
            {props.cardName}
          </Link>
        </td>
        <td><SetSymbol setName={props.setName} setCode={props.setCode} /></td>
      </tr>
    );
}

export {CardSearch, CardSearchResultsTable, CardSearchResultRow};
