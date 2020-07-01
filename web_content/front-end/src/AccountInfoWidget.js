import React from 'react';

class AccountInfoWidget extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      firstName: "",
      lastName: "",
    };

    this.handleLogout = this.handleLogout.bind(this);
  }

  componentDidMount() {
    this.props.userManager.getUser().then(user => {
      this.setState({
        firstName: user.profile.first_name,
        lastName: user.profile.last_name,
      });
    });
  }

  handleLogout(event) {
    this.props.userManager.signoutRedirect();
    event.preventDefault();
  }

  render() {
    return (
      <div id="request_logout">
        Hello, {this.state.firstName} {this.state.lastName}!
        <form onSubmit={this.handleLogout}>
          <input type="submit" value="Log Out" />
        </form>
      </div>
    );
  }
}

export default AccountInfoWidget;
