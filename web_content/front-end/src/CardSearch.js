import React from 'react';
import {Link} from 'react-router-dom';
import {connect} from 'react-redux';

import {SetSymbol} from './CardComponents.js';
import {requestCardSearchResults} from './ReduxActions.js';

const mapStateToProps = (state) => {
  return {
    cards: state.cardSearch.searchResults,
  };
}

class CardSearch extends React.Component {
  constructor(props) {
    super(props);

    this.handleInput = this.handleInput.bind(this);
  }

  handleInput(event) {
    if (event.target.value.length > 1) {
      this.props.dispatch(requestCardSearchResults(event.target.value));
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
        <CardSearchResultsTable cards={this.props.cards} />
      </div>
    );
  }
}

function CardSearchResultsTable(props) {
  // TODO: Only render the first 10 cards right now for performance reasons
  const cardsToRender = props.cards.slice(0, 11)
  const rows = cardsToRender.map((cardInfo) =>
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

const CardSearchWidget = connect(mapStateToProps)(CardSearch);

export default CardSearchWidget;
