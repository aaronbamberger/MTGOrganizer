import React from 'react';

import {SetSymbol} from './CardComponents.js'

class CardSearch extends React.Component {
  constructor(props) {
    super(props);
    
    this.handleInput = this.handleInput.bind(this);
  }

  handleInput(event) {
    if (event.target.value.length > 1) {
      const request = JSON.stringify({"type": 0, "value": event.target.value});
      this.props.backendRequest(request)
    } else {
      this.setState({"cards": []});
    }
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
        <CardSearchResultsTable cards={this.props.cards} />
      </div>
    );
  }
}

function CardSearchResultsTable(props) {
  const rows = props.cards.map((cardInfo) =>
    <CardSearchResultRow
      cardName={cardInfo.name}
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
        <td>{props.cardName}</td>
        <td><SetSymbol setName={props.setName} setCode={props.setCode} /></td>
      </tr>
    );
}

export {CardSearch, CardSearchResultsTable, CardSearchResultRow};
