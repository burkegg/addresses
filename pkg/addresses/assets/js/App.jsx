class App extends React.Component {
  // a list
  constructor(props) {
    super(props)
    this.state = {
      houses: [],
      version: '',
    }
  }

  render() {
    return (
      <div>not much here yet</div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));
