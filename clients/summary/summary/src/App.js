import React from "react";
import "./App.css";

// TODO: How to I get the errors to show up from the server?
export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      query: "",
      title: "",
      description: "",
      images: [],
      error: ""
    };
  }

  requestData = event => {
    event.preventDefault();
    const { query } = this.state;
    const fetchQuery = "http://localhost:4000/v1/summary/?url=" + query;
    fetch(fetchQuery)
      .then(response => {
        return response.json();
      })
      .catch(err => {
        this.setState({
          query: "",
          title: "",
          description: "",
          images: [],
          error: err.message
        });
      })
      .then(data => {
        if (data) {
          this.setState({
            title: data.title ? data.title : "",
            description: data.description ? data.description : "",
            images: data.images ? data.images : []
          });
        }
      })
      .catch(error => {
        console.log(error);
      });
  };

  handleChange = event => {
    let newVal = event.target.value;
    this.setState({ query: newVal });
  };

  render() {
    const images =
      this.state.images.length !== 0 ? (
        <div className="image-container">
          <p>Images:</p>
          {this.state.images.map(image => {
            return (
              <img
                className="image"
                src={image.url}
                alt={image.alt ? image.alt : ""}
                key={image.url}
              />
            );
          })}
        </div>
      ) : null;
    return (
      <div className="App">
        <form onSubmit={this.requestData}>
          <input
            className="input"
            placeholder="Enter an url"
            aria-label="Enter an url"
            type="search"
            value={this.state.query}
            onChange={this.handleChange}
            required
          />
        </form>

        {this.state.error !== "" ? (
          <p style={{ color: "red" }}>Error: {this.state.error}</p>
        ) : null}

        {this.state.title !== "" ? (
          <div>
            <p>{`The title of the link is:\n ${this.state.title}`}</p>
          </div>
        ) : null}

        {this.state.description !== "" ? (
          <div>
            <p>{`The description of the link is:\n ${this.state.description}`}</p>{" "}
          </div>
        ) : null}

        {images}
      </div>
    );
  }
}
