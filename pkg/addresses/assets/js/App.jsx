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

    let search = { Term: "MA" }

    let initReq = {
      method: "POST",
      body: JSON.stringify(search)
    }

    fetch(fetchEndpoint, initReq)
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
      <div>b not much here yet</div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));
