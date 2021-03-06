import React from "react";
import "./App.css";

export default class Channel extends React.Component {
  constructor(props) {
    super(props);
    this.state = { text: "", editMessage: null };
  }

  startEdit = (event, i) => {
    const message = this.props.messages[i];
    this.setState({ editMessage: message, text: message.body });
    this.input.focus();
  };

  handleEdit = () => {
    const newMessage = this.state.text;
    const message = this.state.editMessage;

    fetch("https://api.alexst.me/v1/messages/" + message._id, {
      method: "PATCH",
      headers: {
        Authorization: localStorage.getItem("Auth")
      },
      body: JSON.stringify({
        message: newMessage
      })
    })
      .then(response => {
        return response.json();
      })
      .then(newMessage => {
        console.log(newMessage);
      })
      .catch(err => {
        console.log(err);
      });
  };

  deleteItem = (event, i) => {
    const message = this.props.messages[i];
    fetch("https://api.alexst.me/v1/messages/" + message._id, {
      method: "DELETE",
      headers: {
        Authorization: localStorage.getItem("Auth")
      }
    })
      .then(resp => {
        return resp.text();
      })
      .catch(err => {
        console.log(err);
      })
      .then(text => {
        console.log(text);
      });
  };

  handleChange = event => {
    let field = event.target.name;
    let value = event.target.value;

    let changes = {};
    changes[field] = value;
    this.setState(changes);
  };

  handleAdd = () => {
    const message = this.state.text;
    fetch(`https://api.alexst.me/v1/channels/${this.props.channel._id}`, {
      method: "POST",
      headers: {
        "Content-Type": " application/json",
        Authorization: localStorage.getItem("Auth")
      },
      body: JSON.stringify({
        body: message
      })
    })
      .then(response => {
        return response.json();
      })
      .then(m => {
        // console.log(m);
        this.setState({ text: "" });
      })
      .catch(err => {
        alert(err);
        return;
      });
  };

  render() {
    const messages = this.props.messages.map((m, i) => {
      return (
        <div key={i} style={{ display: "flex", flexDirection: "column" }}>
          <div
            style={{
              display: "flex",
              flexDirection: "row"
            }}
          >
            <p style={{ fontWeight: "bold", marginBottom: 0 }}>
              {m.creator.userName}
            </p>
            <p style={{ paddingLeft: 10, paddingRight: 10 }}>-</p>
            <p>{m.body}</p>
          </div>
          {m.creator.id === this.props.user.id ? (
            <div style={{ display: "flex", flexDirection: "row" }}>
              <div onClick={event => this.deleteItem(event, i)}>
                <p style={{ color: "red" }}>Delete</p>
              </div>
              <div onClick={event => this.startEdit(event, i)}>
                <p style={{ color: "blue", marginLeft: 10 }}>Edit</p>
              </div>
            </div>
          ) : null}
        </div>
      );
    });
    return (
      <div style={styles.channelContainer}>
        <p style={{ fontWeight: "bold", fontSize: 24 }}>
          Channel: {this.props.channel ? this.props.channel.name : null}
        </p>
        <p style={{ textDecoration: "underline" }}>
          Messages Below (Newest at the top. I'm sorry)
        </p>
        <div style={styles.messageContainer}>{messages}</div>
        {/* Add */}
        <div
          style={{ display: "flex", flexDirection: "row" }}
          className="form-group"
        >
          <input
            ref={input => {
              this.input = input;
            }}
            className="form-control"
            id="text"
            name="text"
            placeholder={"Type something"}
            value={this.state.text}
            onChange={this.handleChange}
          />
          <button
            className="btn btn-primary mr-2"
            onClick={this.state.editMessage ? this.handleEdit : this.handleAdd}
          >
            {this.state.editMessage ? "Edit" : "Add"}
          </button>
        </div>
      </div>
    );
  }
}

const styles = {
  channelContainer: {
    display: "flex",
    flexDirection: "column",
    marginLeft: 80,
    width: 600
  },
  messageContainer: {
    height: 600,
    overflowY: "scroll",
    backgroundColor: "gray"
  }
};
