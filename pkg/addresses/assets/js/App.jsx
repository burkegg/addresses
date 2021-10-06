class App extends React.Component {
  // a list
  constructor(props) {
    super(props)
    this.state = {
      houses: [],
      version: '',
    }
  }
  componentDidMount() {
    const fetchEndpoint = "api/addresses"
    fetch(fetchEndpoint)
    .then(res => res.json())
    .then(json => {
      console.log('json: \n', json)
    })
    .catch(err => {
      console.log('error:', err)
    })
  }

  render() {
    return (
      <div>not much here yet</div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));
