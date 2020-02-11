import React from "react";
import "./App.css";

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      user: { FirstName: "Alex", lastname: "Tan" },
      showUpdate: false
    };
  }

  handleSignIn = (email, password) => {
    // TODO
  };

  handleSignUp = (
    firstname,
    lastname,
    email,
    password,
    passwordConf,
    username
  ) => {
    // TODO:
  };

  handleSignOut = () => {
    // TODO
  };

  handleUpdateChange = () => {
    this.setState({ showUpdate: !this.state.showUpdate });
  };

  handleUpdate = (fNameChange, lNameChange) => {
    // TODO
    this.setState({ showUpdate: false });
  };

  render() {
    return (
      <div className="App">
        <p>Welcome to my application</p>
        {!this.state.user ? (
          <SignUpForm
            handleSignIn={this.handleSignIn}
            handleSignUp={this.handleSignUp}
          />
        ) : (
          <div>
            <p>Congratulations on signing in!</p>
            <p>
              Your name is {this.state.user.FirstName}{" "}
              {this.state.user.lastname}
            </p>

            {!this.state.showUpdate ? (
              <div>
                <p>Best I can do for you are these two options</p>
                <button
                  className="btn btn-primary mr-2"
                  onClick={this.handleSignOut}
                >
                  Sign Out
                </button>
                <button
                  className="btn btn-primary mr-2"
                  onClick={this.handleUpdateChange}
                >
                  Update Profile
                </button>
              </div>
            ) : (
              <UpdateForm
                cancelUpdate={this.handleUpdateChange}
                handleUpdate={(fNameChange, lNameChange) =>
                  this.handleUpdate(fNameChange, lNameChange)
                }
              />
            )}
          </div>
        )}
      </div>
    );
  }
}

class UpdateForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = { firstNameChange: "", lastNameChange: "" };
  }

  handleChange = event => {
    let field = event.target.name;
    let value = event.target.value;

    let changes = {};
    changes[field] = value;
    this.setState(changes);
  };

  handleUpdate = event => {
    event.preventDefault();
    this.props.handleUpdate(
      this.state.firstNameChange,
      this.state.lastNameChange
    );
  };

  cancelUpdate = event => {
    event.preventDefault();
    this.props.cancelUpdate();
  };

  render() {
    return (
      <div>
        <p>User Updates:</p>
        <div className="form-group">
          <label htmlFor="firstNameChange">First Name</label>
          <input
            className="form-control"
            id="firstNameChange"
            name="firstNameChange"
            onChange={this.handleChange}
          />
        </div>

        <div className="form-group">
          <label htmlFor="lastNameChange">First Name</label>
          <input
            className="form-control"
            id="lastNameChange"
            name="lastNameChange"
            onChange={this.handleChange}
          />
        </div>

        <div className="form-group">
          <button className="btn btn-primary mr-2" onClick={this.handleUpdate}>
            Save Updates
          </button>
          <button className="btn btn-primary mr-2" onClick={this.cancelUpdate}>
            Cancel
          </button>
        </div>
      </div>
    );
  }
}

class SignUpForm extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      signIn: false,
      firstname: "",
      lastname: "",
      email: "",
      password: "",
      passwordConf: "",
      username: ""
    };
  }

  //update state for specific field
  handleChange = event => {
    let field = event.target.name;
    let value = event.target.value;

    let changes = {};
    changes[field] = value;
    this.setState(changes);
  };

  //handle signUp button
  handleSignUp = event => {
    event.preventDefault();
    this.props.handleSignUp(
      this.state.firstname,
      this.state.lastname,
      this.state.email,
      this.state.password,
      this.state.passwordConf,
      this.state.username
    );
  };

  //handle signIn button
  handleSignIn = event => {
    event.preventDefault();
    this.props.handleSignIn(this.state.email, this.state.password);
  };

  changeSignIn = event => {
    event.preventDefault();
    this.setState({ signIn: !this.state.signIn });
  };

  render() {
    return (
      <form>
        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="firstname">First Name</label>
            <input
              className="form-control"
              id="firstname"
              name="firstname"
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="lastname">Last Name</label>
            <input
              className="form-control"
              id="lastname"
              name="lastname"
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            className="form-control"
            id="email"
            type="email"
            name="email"
            onChange={this.handleChange}
          />
        </div>

        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            className="form-control"
            id="password"
            type="password"
            name="password"
            onChange={this.handleChange}
          />
        </div>

        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="passwordConf">Confirm Password</label>
            <input
              className="form-control"
              id="passwordConf"
              type="password"
              name="passwordConf"
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        {this.state.signIn ? (
          <div className="form-group">
            <label htmlFor="username">Username</label>
            <input
              className="form-control"
              id="username"
              name="username"
              onChange={this.handleChange}
            />
          </div>
        ) : null}

        {this.state.signIn ? (
          <div className="form-group">
            <button
              className="btn btn-primary mr-2"
              onClick={this.handleSignUp}
            >
              SignUp
            </button>
            <button
              className="btn btn-primary mr-2"
              onClick={this.changeSignIn}
            >
              Back to Signin
            </button>
          </div>
        ) : (
          <div className="form-group">
            <button
              className="btn btn-primary mr-2"
              onClick={this.handleSignIn}
            >
              Login
            </button>
            <button
              className="btn btn-primary mr-2"
              onClick={this.changeSignIn}
            >
              SignUp
            </button>
          </div>
        )}
      </form>
    );
  }
}
