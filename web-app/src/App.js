import React, { Component } from "react";
import {
  Grid,
  Navbar,
  Button,
  Form,
  FormGroup,
  InputGroup,
  FormControl,
  ListGroup,
  ListGroupItem,
  Glyphicon
} from "react-bootstrap";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      searchTitle: "",
      feeds: null,
      failMsg: ""
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleCreate = this.handleCreate.bind(this);
  }

  handleChange(e) {
    this.setState(
      Object.assign({}, this.state, { searchTitle: e.target.value })
    );
  }

  handleCreate(e) {
    e.preventDefault();

    if (this.state.searchTitle === "") {
      return;
    }

    let body = JSON.stringify({ title: this.state.searchTitle });

    fetch("/api/feeds", {
      method: "POST",
      headers: {
        Accept: "application/json, text/plain, */*",
        "Content-Type": "application/json"
      },
      body: body
    })
      .then(res => {
        if (res.ok) return res.json();
        throw new Error("post api to create feed failed");
      })
      .then(result => {
        this.setState(
          Object.assign({}, this.state, {
            feeds: [result].concat(this.state.feeds)
          })
        );
      })
      .catch(e => {
        console.log("failed to create new feed", e);
      });
  }

  componentWillMount() {
    fetch("/api/feeds")
      .then(res => res.json())
      .then(feeds => {
        this.setState({
          feeds: feeds
        });
      })
      .catch(e => {
        console.warn("fetch feeds failed with error", e);
        this.setState({
          feeds: null,
          failMsg: e.message
        });
      });
  }

  render() {
    return (
      <div>
        <Navbar>
          <Navbar.Header>
            <Navbar.Brand>
              <a href="/">Feedit</a>
            </Navbar.Brand>
          </Navbar.Header>
        </Navbar>
        <Grid>
          <Form onSubmit={this.handleCreate}>
            <FormGroup>
              <InputGroup>
                <FormControl
                  type="text"
                  value={this.state.searchTitle}
                  placeholder="Feed title"
                  onChange={this.handleChange}
                />
                <InputGroup.Button>
                  <Button type="submit" onClick={this.handleCreate}>
                    Create New Feed
                  </Button>
                </InputGroup.Button>
              </InputGroup>
            </FormGroup>
          </Form>
          <FeedList feeds={this.state.feeds} />
        </Grid>
      </div>
    );
  }
}

class Feed extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <ListGroupItem>
        <div>{this.props.title}</div>
        <div>
          <span>[{this.props.email}] </span>
        </div>
        <SubscribeButton url={this.props.url} />
      </ListGroupItem>
    );
  }
}

const FeedEmpty = props => {
  return (
    <div>
      <p> No feed found. Try to add new feed by clicking above button </p>
    </div>
  );
};

const SubscribeButton = props => {
  return (
    <Button bsStyle="link" bsSize="sm" href={props.url}>
      <Glyphicon glyph="list" />
    </Button>
  );
};

const FeedList = props => {
  const feeds = props.feeds;

  if (feeds === null) {
    return <FeedEmpty />;
  }

  const listItems = feeds.map((feed, i) => <Feed {...feed} key={i} />);

  return <ListGroup>{listItems}</ListGroup>;
};

export default App;
